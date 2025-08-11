-- name: CreateHello :one
INSERT INTO hello (message)
VALUES ($1)
RETURNING id, message;

-- name: GetHello :one
SELECT id, message
FROM hello
WHERE id = $1;

-- name: ListHellos :many
SELECT id, message
FROM hello
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateHello :exec
UPDATE hello
SET message = $1
WHERE id = $2;

-- name: DeleteHello :exec
DELETE FROM hello
WHERE id = $1;

