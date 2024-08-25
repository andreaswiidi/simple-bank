-- name: CreateTransfer :one
INSERT INTO transfers_history (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetDetilTransfer :one
SELECT * FROM transfers_history
WHERE id = $1 LIMIT 1;

-- name: ListTransfersHistoryAccount :many
SELECT * FROM transfers_history
WHERE 
    from_account_id = $1 OR
    to_account_id = $1 OR
    (to_account_id = $1 AND (from_account_id = $2))
ORDER BY id
LIMIT $3
OFFSET $4;