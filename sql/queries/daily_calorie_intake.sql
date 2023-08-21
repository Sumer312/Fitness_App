-- name: CreateDailyCalorieIntake :one
INSERT INTO daily_calorie_intake(id, created_at, updated_at, user_id, calories)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDailyCalories :many
SELECT * FROM daily_calorie_intake WHERE user_id = $1;

-- name: DeleteDailyCalories :exec
DELETE FROM daily_calorie_intake WHERE user_id = $1;
