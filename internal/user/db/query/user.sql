-- name: CreateUser :one
INSERT INTO users (username, email)
    VALUES ($1, $2)
RETURNING
    *;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    "username"
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    username = $1,
    email = $2
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
