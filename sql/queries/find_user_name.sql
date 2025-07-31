-- name: FindUserName :one
SELECT name FROM users
WHERE id = $1;