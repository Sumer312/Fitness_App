-- name: CreateTotalCalorieIntake :one
INSERT INTO total_calorie_intake(id, created_at, user_id, calories, total_deficit, total_surplus)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
 
-- name: GetTotalCalories :many
SELECT * FROM total_calorie_intake WHERE user_id = $1;
 
-- name: GetMostRecentUserKcal :one
SELECT * FROM total_calorie_intake WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: DeleteRedundantRows :exec
DELETE FROM total_calorie_intake WHERE calories = 0;

-- name: DeleteUserRecord :exec
DELETE FROM total_calorie_intake WHERE user_id = $1;
