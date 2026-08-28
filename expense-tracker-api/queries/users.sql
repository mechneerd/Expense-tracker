-- name: GetUserByEmail :one
SELECT id, google_id, email, email_verified_at, first_name, last_name, phone, avatar_url, status, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, google_id, email, email_verified_at, first_name, last_name, phone, avatar_url, status, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (google_id, email, first_name, last_name, avatar_url, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW());

-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, phone = $3, avatar_url = $4, status = $5, updated_at = NOW()
WHERE id = $6;

-- name: VerifyUserEmail :exec
UPDATE users
SET email_verified_at = NOW(), updated_at = NOW()
WHERE id = $1;
