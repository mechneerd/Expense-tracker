-- name: ListUPIApps :many
SELECT id, name, package_name, created_at, updated_at
FROM upi_apps
ORDER BY name;

-- name: GetUPIAppByName :one
SELECT id, name, package_name, created_at, updated_at
FROM upi_apps
WHERE name = $1;

-- name: GetUPIAppByID :one
SELECT id, name, package_name, created_at, updated_at
FROM upi_apps
WHERE id = $1;
