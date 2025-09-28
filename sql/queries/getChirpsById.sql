-- name: GetChirpById :one
select *
from chirps
where id = $1;