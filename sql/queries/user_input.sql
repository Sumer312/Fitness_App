-- name: CreateUserInput :one
INSERT INTO user_input (id, user_id, created_at, updated_at, height, weight, desired_weight, time_frame, bmi , program, curr_kcal, deficit)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetUserInput :one
SELECT * FROM user_input where user_id = $1;
