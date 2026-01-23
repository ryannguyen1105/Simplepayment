-- name: CreateWallet :one
INSERT INTO wallets (
    user_id, balance
) VALUES (
             $1, $2
         )
RETURNING *;

-- name: GetWalletByID :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletByIDForUpdate :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListWallets :many
SELECT id, user_id, balance, created_at
FROM wallets
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateWalletBalance :one
UPDATE wallets
SET balance = $2
WHERE id = $1
    RETURNING *;

-- name: AddWalletBalance :one
UPDATE wallets
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
    RETURNING *;