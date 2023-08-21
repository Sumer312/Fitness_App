package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/internal/database"
)

func (apiCfg *apiConfig) input_handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		UserId         uuid.UUID `json:"user_id"`
		Height         int       `json:"height"`
		Weight         int       `json:"weitht"`
		Desired_Weight *int      `json:"desired_weight"`
		TimeFrame      *int      `json:"time-frame"`
		Bmi            int       `json:"bmi"`
		Program        string    `json:"program"`
		Curr_Kcal      int       `json:"curr_kcal"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("line 28", err)
	}
	if params.Desired_Weight != nil && params.TimeFrame != nil && *params.TimeFrame != 0 {
		deficit := deficit_calc(params.Weight, *params.Desired_Weight, *params.TimeFrame)
		apiCfg.DB.CreateUserInput(r.Context(), database.CreateUserInputParams{
			ID:            uuid.New(),
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
			UserID:        params.UserId,
			Height:        int32(params.Height),
			Weight:        int32(params.Weight),
			DesiredWeight: sql.NullInt32{Int32: int32(*params.Desired_Weight), Valid: true},
			TimeFrame:     sql.NullInt32{Int32: int32(*params.TimeFrame), Valid: true},
			CurrKcal:      int32(params.Curr_Kcal),
			Bmi:           int32(params.Bmi),
			Program:       params.Program,
			Deficit:       deficit,
		})
	} else {
		apiCfg.DB.CreateUserInput(r.Context(), database.CreateUserInputParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    params.UserId,
			Height:    int32(params.Height),
			Weight:    int32(params.Weight),
			CurrKcal:  int32(params.Curr_Kcal),
			Bmi:       int32(params.Bmi),
			Program:   params.Program,
		})
	}
}
