-- name: CreatePayment :one
INSERT INTO payments (
    from_user_id,
    to_user_id,
    amount,
    status
) VALUES (
             $1, $2, $3, 'pending'
         )
    RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1;

-- name: ListPayments :many
SELECT * FROM payments
WHERE from_user_id = $1 OR
      to_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: CancelPayment :one
UPDATE payments
SET status = 'canceled'
WHERE id = $1
    AND status = 'pending'
    RETURNING *;



