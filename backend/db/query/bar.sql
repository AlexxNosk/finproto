-- name: CreateBar :one
INSERT INTO bars (
  timestamp,
  open,
  high,
  low,
  close,
  volume
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetBar :one
SELECT * FROM bars
WHERE id = $1 LIMIT 1;

-- name: ListBars :many
SELECT * FROM bars
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateBars :one
UPDATE bars
  set Open = $2,
  High = $3,
  Low = $4,
  Close = $5,
  Volume  = $6
WHERE id = $1
RETURNING *;

-- name: DeleteAllBars :exec
DELETE FROM bars;

-- name: TruncateBars :exec
TRUNCATE TABLE bars RESTART IDENTITY;