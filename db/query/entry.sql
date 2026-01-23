-- name: CreateEntry :one
INSERT INTO entries (
    wallet_id,
    payment_id,
    amount
) VALUES (
             $1, $2, $3
         )
    RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE wallet_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;

