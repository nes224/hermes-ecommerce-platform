-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    first_name,
    last_name,
    phone_number,
    role
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, email, first_name, last_name, phone_number, role, is_active, is_verified, last_login_at, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, first_name, last_name, phone_number, role, is_active, is_verified, last_login_at, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, first_name, last_name, phone_number, role, is_active, is_verified, last_login_at, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    role = COALESCE(sqlc.narg('role'), role),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    is_verified = COALESCE(sqlc.narg('is_verified'), is_verified),
    last_login_at = COALESCE(sqlc.narg('last_login_at'), last_login_at)
WHERE id = $1
RETURNING id, email, first_name, last_name, phone_number, role, is_active, is_verified, last_login_at, created_at, updated_at;

-- name: ListUsers :many
SELECT 
    id, email, first_name, last_name, phone_number, role, is_active, is_verified, last_login_at, created_at, updated_at
FROM users
WHERE
    (
        sqlc.narg('search')::text IS NULL OR 
        email ILIKE '%' || sqlc.narg('search')::text || '%' OR
        first_name ILIKE '%' || sqlc.narg('search')::text || '%' OR
        last_name ILIKE '%' || sqlc.narg('search')::text || '%' OR
        phone_number ILIKE '%' || sqlc.narg('search')::text || '%'
    )
    AND (sqlc.narg('role')::user_role IS NULL OR role = sqlc.narg('role'))
    AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
    AND (sqlc.narg('is_verified')::boolean IS NULL OR is_verified = sqlc.narg('is_verified'))
ORDER BY
    CASE WHEN sqlc.arg('order_by')::text = 'created_at' AND sqlc.arg('is_desc')::boolean = true THEN created_at END DESC,
    CASE WHEN sqlc.arg('order_by')::text = 'created_at' AND sqlc.arg('is_desc')::boolean = false THEN created_at END ASC,
    CASE WHEN sqlc.arg('order_by')::text = 'email' AND sqlc.arg('is_desc')::boolean = true THEN email END DESC,
    CASE WHEN sqlc.arg('order_by')::text = 'email' AND sqlc.arg('is_desc')::boolean = false THEN email END ASC,
    CASE WHEN sqlc.arg('order_by')::text = 'first_name' AND sqlc.arg('is_desc')::boolean = true THEN first_name END DESC,
    CASE WHEN sqlc.arg('order_by')::text = 'first_name' AND sqlc.arg('is_desc')::boolean = false THEN first_name END ASC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE 
    (
        sqlc.narg('search')::text IS NULL OR 
        email ILIKE '%' || sqlc.narg('search')::text || '%' OR
        first_name ILIKE '%' || sqlc.narg('search')::text || '%' OR
        last_name ILIKE '%' || sqlc.narg('search')::text || '%' OR
        phone_number ILIKE '%' || sqlc.narg('search')::text || '%'
    )
    AND (sqlc.narg('role')::user_role IS NULL OR role = sqlc.narg('role'))
    AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
    AND (sqlc.narg('is_verified')::boolean IS NULL OR is_verified = sqlc.narg('is_verified'));

-- name: GetUserForAuth :one
SELECT id, email, password_hash, role, is_active
FROM users
WHERE email = $1 LIMIT 1;