-- name: CreateTask :one
INSERT INTO tasks (title, description, status)
VALUES ($1, $2, $3)
    RETURNING id, title, description, status, created_at, updated_at, deleted_at;

-- name: GetTaskByID :one
SELECT *
FROM tasks
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateTask :one
UPDATE tasks
SET title       = $2,
    description = $3,
    status      = $4,
    updated_at  = NOW()
WHERE id = $1
  AND deleted_at IS NULL
    RETURNING *;

-- name: DeleteTask :exec
UPDATE tasks
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetTasks :many
SELECT *
FROM tasks
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListTasks :many
SELECT *
FROM tasks
WHERE deleted_at IS NULL
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;
