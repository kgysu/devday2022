-- name: AddProduct :one
INSERT INTO products (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;
