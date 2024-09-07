-- name: CreateTask :one
INSERT INTO tasks (type, value, state, creation_time, last_update_time)
VALUES ($1, $2, 'received', NOW(), NOW())
RETURNING id;

-- name: UpdateTaskState :exec
UPDATE tasks
SET state = $2, last_update_time = NOW()
WHERE id = $1;
