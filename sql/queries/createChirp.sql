-- name: CreateChirp :one
insert into chirps(id, created_at,updated_at,body,user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    Now(),
    $1,
    $2
)
RETURNING *;