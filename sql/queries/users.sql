-- sql (PostgreSQL)
-- name: CreateUser :one
INSERT INTO public.users (id, created_at, updated_at, email,hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;