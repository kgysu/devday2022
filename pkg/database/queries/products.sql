-- name: AddProduct :one
INSERT INTO products (
  name,
  kind
) VALUES (
  $1,
  $2
)
RETURNING *;


-- name: GetProducts :many
SELECT * FROM products ORDER BY create_time;


-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;


-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;
