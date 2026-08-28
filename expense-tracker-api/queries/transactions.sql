-- name: CreateTransaction :exec
INSERT INTO transactions (family_id, user_id, transaction_type, category_id, payment_method_id, upi_app_id, amount, transaction_date, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW());

-- name: GetTransactionByID :one
SELECT t.id, t.family_id, t.user_id, t.transaction_type, tc.name as category, pm.name as payment_method, t.upi_app_id, t.amount, t.transaction_date, t.description, t.created_at, t.updated_at, t.deleted_at
FROM transactions t
LEFT JOIN transaction_categories tc ON t.category_id = tc.id
LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
WHERE t.id = $1 AND t.deleted_at IS NULL;

-- name: ListTransactionsByFamily :many
SELECT t.id, t.family_id, t.user_id, t.transaction_type, tc.name as category, pm.name as payment_method, t.upi_app_id, t.amount, t.transaction_date, t.description, t.created_at, t.updated_at, t.deleted_at
FROM transactions t
LEFT JOIN transaction_categories tc ON t.category_id = tc.id
LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
WHERE t.family_id = $1 AND t.deleted_at IS NULL
ORDER BY t.transaction_date DESC;

-- name: ListTransactionsByUser :many
SELECT t.id, t.family_id, t.user_id, t.transaction_type, tc.name as category, pm.name as payment_method, t.upi_app_id, t.amount, t.transaction_date, t.description, t.created_at, t.updated_at, t.deleted_at
FROM transactions t
LEFT JOIN transaction_categories tc ON t.category_id = tc.id
LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
WHERE t.user_id = $1 AND ($2::UUID IS NULL OR t.family_id = $2) AND t.deleted_at IS NULL
ORDER BY t.transaction_date DESC;

-- name: SoftDeleteTransaction :exec
UPDATE transactions
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;
