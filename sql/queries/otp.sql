-- name: CreateOTP :one
insert into otp (id, mobile, otp) values ($1, $2, $3) RETURNING *;

-- name: DeleteOTP :exec
delete from otp where id = $1;

-- name: GetOTP :one
select * from otp where id = $1;