-- name: CreateTransactionHistory :one
INSERT INTO transaction_history (
  account_id,
  amount,
  transaction_type,
  transfer_history_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTransactionHistory :one
SELECT * FROM transaction_history
WHERE id = $1 LIMIT 1;

-- name: ListTransactionHistories :many
SELECT * FROM transaction_history
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;