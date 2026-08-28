-- name: ListPaymentMethods :many
SELECT id, name, description
FROM payment_methods
ORDER BY name;

-- name: GetPaymentMethodByName :one
SELECT id, name, description
FROM payment_methods
WHERE name = $1;

-- name: GetPaymentMethodByID :one
SELECT id, name, description
FROM payment_methods
WHERE id = $1;
