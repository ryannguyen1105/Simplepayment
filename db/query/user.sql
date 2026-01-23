-- name: CreateUser :one
INSERT INTO users (
    email,
    username,
    hashed_password
) VALUES (
             $1, $2,$3
         )
    RETURNING id, email, username, created_at;

-- name: GetUserByID :one
SELECT id, email, username, created_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, username, hashed_password, created_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, username, created_at
FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserEmail :one
UPDATE users
set email = $2
WHERE id = $1
    RETURNING id, email, username, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
