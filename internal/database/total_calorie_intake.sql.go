// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: total_calorie_intake.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTotalCalorieIntake = `-- name: CreateTotalCalorieIntake :one
INSERT INTO total_calorie_intake(id, created_at, user_id, calories, total_deficit, total_surplus)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, calories, total_deficit, total_surplus, user_id
`

type CreateTotalCalorieIntakeParams struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UserID       uuid.UUID
	Calories     float64
	TotalDeficit float64
	TotalSurplus float64
}

func (q *Queries) CreateTotalCalorieIntake(ctx context.Context, arg CreateTotalCalorieIntakeParams) (TotalCalorieIntake, error) {
	row := q.db.QueryRowContext(ctx, createTotalCalorieIntake,
		arg.ID,
		arg.CreatedAt,
		arg.UserID,
		arg.Calories,
		arg.TotalDeficit,
		arg.TotalSurplus,
	)
	var i TotalCalorieIntake
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Calories,
		&i.TotalDeficit,
		&i.TotalSurplus,
		&i.UserID,
	)
	return i, err
}

const deleteRedundantRows = `-- name: DeleteRedundantRows :exec
DELETE FROM total_calorie_intake WHERE calories = 0
`

func (q *Queries) DeleteRedundantRows(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteRedundantRows)
	return err
}

const deleteUserRecord = `-- name: DeleteUserRecord :exec
DELETE FROM total_calorie_intake WHERE user_id = $1
`

func (q *Queries) DeleteUserRecord(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUserRecord, userID)
	return err
}

const getMostRecentUserKcalByUserId = `-- name: GetMostRecentUserKcalByUserId :one
SELECT id, created_at, calories, total_deficit, total_surplus, user_id FROM total_calorie_intake WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1
`

func (q *Queries) GetMostRecentUserKcalByUserId(ctx context.Context, userID uuid.UUID) (TotalCalorieIntake, error) {
	row := q.db.QueryRowContext(ctx, getMostRecentUserKcalByUserId, userID)
	var i TotalCalorieIntake
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Calories,
		&i.TotalDeficit,
		&i.TotalSurplus,
		&i.UserID,
	)
	return i, err
}

const getTotalCalorieIntakeByUserId = `-- name: GetTotalCalorieIntakeByUserId :many
SELECT id, created_at, calories, total_deficit, total_surplus, user_id FROM total_calorie_intake WHERE user_id = $1
`

func (q *Queries) GetTotalCalorieIntakeByUserId(ctx context.Context, userID uuid.UUID) ([]TotalCalorieIntake, error) {
	rows, err := q.db.QueryContext(ctx, getTotalCalorieIntakeByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TotalCalorieIntake
	for rows.Next() {
		var i TotalCalorieIntake
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Calories,
			&i.TotalDeficit,
			&i.TotalSurplus,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
