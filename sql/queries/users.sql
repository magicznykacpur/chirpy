-- name: Createuser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid (), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserEmailAndPassword :exec
UPDATE users SET updated_at = NOW(), email = $1, hashed_password = $2 WHERE id = $3;

-- name: UpdateIsChirpyRed :exec
UPDATE users SET updated_at = NOW(), is_chirpy_red = TRUE WHERE id = $1;