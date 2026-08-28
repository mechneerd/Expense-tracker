-- name: GetFamilyByID :one
SELECT id, name, unique_code, created_by, status, created_at, updated_at
FROM families
WHERE id = $1;

-- name: GetFamilyByUniqueCode :one
SELECT id, name, unique_code, created_by, status, created_at, updated_at
FROM families
WHERE unique_code = $1;

-- name: CreateFamily :exec
INSERT INTO families (name, unique_code, created_by, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());

-- name: ListFamilyMembers :many
SELECT fm.id, fm.family_id, fm.user_id, fm.family_role, fm.status, fm.joined_at, fm.created_at, fm.updated_at
FROM family_members fm
WHERE fm.family_id = $1
ORDER BY fm.created_at;

-- name: AddFamilyMember :exec
INSERT INTO family_members (family_id, user_id, family_role, status, joined_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW());

-- name: GetFamilyMember :one
SELECT id, family_id, user_id, family_role, status, joined_at, created_at, updated_at
FROM family_members
WHERE family_id = $1 AND user_id = $2;

-- name: UpdateFamilyMemberStatus :exec
UPDATE family_members
SET status = $1, updated_at = NOW()
WHERE family_id = $2 AND user_id = $3;
