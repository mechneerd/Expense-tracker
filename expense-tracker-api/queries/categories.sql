-- name: ListCategoriesByType :many
SELECT id, type, name, created_at, updated_at
FROM transaction_categories
WHERE type = $1
ORDER BY name;

-- name: GetCategoryByName :one
SELECT id, type, name, created_at, updated_at
FROM transaction_categories
WHERE name = $1;

-- name: CreateCategory :exec
INSERT INTO transaction_categories (type, name, description, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW());
