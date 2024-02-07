package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *Api) Profile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		UserId uuid.UUID `json:"user_id"`
	}
	params := &parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("line 20", err)
	}
	userInput, err := apiCfg.DB.GetUserInputById(r.Context(), params.UserId)
	if err != nil {
		log.Fatal("line 24", err)
	}
	userCalories, err := apiCfg.DB.GetTotalCalories(r.Context(), params.UserId)
	if err != nil {
		log.Fatal("line 28", err)
	}
	totalCalories := 0
	for i := 0; i < len(userCalories); i++ {
		totalCalories += int(userCalories[i].Calories)
		fmt.Println(userCalories[i].Calories)
	}
	weight_diff := float32(userInput.Weight) - float32(userInput.DesiredWeight.Int32)
	weight_in_deficit := (float32(userInput.CurrKcal)*float32(len(userCalories)) - float32(totalCalories)) / 7716.0
	var percentage float32 = weight_in_deficit * 100 / weight_diff
	fmt.Println(percentage)
	fmt.Println(totalCalories)
}
