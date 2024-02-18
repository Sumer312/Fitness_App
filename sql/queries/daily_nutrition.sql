-- name: CreateDailyNutrition :one
INSERT INTO daily_nutrition_intake(id, created_at, user_id, calories, carbohydrates, protien, fat, fiber)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDailyNutritionOfUserByUserId :many
SELECT * FROM daily_nutrition_intake WHERE user_id = $1;

-- name: DeleteDailyNutritionOfUserByUserId :exec
DELETE FROM daily_nutrition_intake WHERE user_id = $1;

-- name: DeleteRowFromDailyNutritionTableById :exec
DELETE FROM daily_nutrition_intake WHERE id = $1;
