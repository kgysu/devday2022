-- name: AddUser :one
INSERT INTO users (
  role,
  name
) VALUES (
  $1, 
  $2
)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;


-- name: GetUsers :many
SELECT * FROM users;
