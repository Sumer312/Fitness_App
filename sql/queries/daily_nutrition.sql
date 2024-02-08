-- name: CreateDailyNutrition :one
INSERT INTO daily_calorie_intake(id, created_at, updated_at, user_id, calories, carbohydrates, protien, fat, fiber)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetDailyNutritionOfUserByUserId :many
SELECT * FROM daily_calorie_intake WHERE user_id = $1;

-- name: DeleteDailyNutritionOfUserByUserId :exec
DELETE FROM daily_calorie_intake WHERE user_id = $1;
