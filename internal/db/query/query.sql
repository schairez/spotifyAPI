-- name: CreateUser :one 
INSERT INTO users(
        id,
        name,
        email,
        image_url,
        country,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: GetUserByID :one 
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: GetAllUsers :many
SELECT *
FROM users;
-- name: UpdateUserImageURL :one
UPDATE users
SET image_url = $2
WHERE id = $1
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;