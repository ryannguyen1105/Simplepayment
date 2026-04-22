-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
             $1, $2, $3, $4
         )
    RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 AND deleted = false LIMIT 1;

-- name: DeleteUser :exec
UPDATE users
SET deleted = true
WHERE username = $1;