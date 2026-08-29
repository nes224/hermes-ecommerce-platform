package main

import (
	"context"
	"fmt"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/handler"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository"
	"hermes-ecommerce-platform/apps/user-service/internal/adapters/repository/db"
	"hermes-ecommerce-platform/apps/user-service/internal/core/services"
	"hermes-ecommerce-platform/pkg/config"
	"hermes-ecommerce-platform/pkg/token"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ Cannot load config: %v", err)
	}
	log.Printf("🚀 Starting user-service in [%s] mode...", cfg.Environment)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	accessTokenDuration := cfg.AccessTokenDuration
	refreshTokenDuration := 24 * 7 * time.Hour

	migrationURL := "file://db/migrations" 
	runDBMigration(migrationURL, cfg.DBUrl)

	dbPool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Cannot connect to Postgres: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("❌ Cannot ping Postgres: %v", err)
	}
	log.Println("✅ Connected to Postgres successfully")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("X Cannot connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Connected to Redis successfully")

	tokenMaker, err := token.NewJWTMaker(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("Cannot create token maker: %v", err)
	}

	userRepo := db.New(dbPool)
	sessionRepo := repository.NewRedisRepository(redisClient)
	authService := services.NewAuthService(
		userRepo,
		&sessionRepo,
		tokenMaker,
		accessTokenDuration,
		refreshTokenDuration,
	)
	authHandler := handler.NewAuthHandler(authService)
	router := handler.SetupRouter(cfg.Environment, authHandler)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		log.Printf("HTTP Server is listening on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")

}

func runDBMigration(migrationURL string, dbURL string) {
	migration, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		log.Fatalf("Cannot create new migrate instance: %v", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Failed to run migrate up: %v", err)
	}
	log.Println("✅ DB migrated successfully")
}
