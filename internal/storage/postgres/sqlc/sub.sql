-- name: CreateSubscription :one
INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetSubscriptionByID :one
SELECT *
FROM subscriptions
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET service_name = $2,
    price        = $3,
    start_date   = $4,
    end_date     = $5,
    updated_at   = NOW()
WHERE id = $1
  AND deleted_at IS NULL RETURNING *;

-- name: DeleteSubscription :exec
UPDATE subscriptions
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetSubscriptions :many
SELECT *
FROM subscriptions
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListSubscriptions :many
SELECT *
FROM subscriptions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalSubscription :many
SELECT *
FROM subscriptions
WHERE deleted_at IS NULL
  AND start_date >= @start_date_from
  AND (@start_date_to::timestamp IS NULL OR start_date <= @start_date_to)
  AND user_id = @user_id
  AND service_name = @service_name;


-- name: GetTotalSubscriptionCost :one
SELECT COALESCE(SUM(price), 0)::BIGINT AS total
FROM subscriptions
WHERE deleted_at IS NULL
  AND start_date >= @start_date_from
  AND (@start_date_to::timestamp IS NULL OR start_date <= @start_date_to)
  AND user_id = @user_id
  AND service_name = @service_name;
