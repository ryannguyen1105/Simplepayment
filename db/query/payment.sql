-- name: CreatePayment :one
INSERT INTO payments (
    from_wallet_id,
    to_wallet_id,
    amount,
    status
) VALUES (
             $1, $2, $3, $4
         )
    RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
WHERE from_wallet_id = $1 OR
      to_wallet_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: CancelPayment :one
UPDATE payments
SET status = 'canceled'
WHERE id = $1
RETURNING *;



