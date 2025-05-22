-- name: RegisterUser :exec
INSERT INTO user_table(username, public_key) 
VALUES (sqlc.arg(username), sqlc.arg(public_key));

-- name: GetPublicKey :one 
SELECT public_key
FROM user_table 
WHERE username = sqlc.arg(username);

