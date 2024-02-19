package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/views/pages"
)

func (apiCfg *Api) ProfileRender(w http.ResponseWriter, r *http.Request) {
	temp, err := r.Cookie("user-id")
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{"errorToast" : "not logged in"}`)
		w.WriteHeader(400)
		return
	}
	userId, err := uuid.Parse(temp.Value)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("HX-Trigger", `{"errorToast" : "not logged in"}`)
		w.WriteHeader(400)
		return
	}
	userInput, err := apiCfg.DB.GetUserInputById(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	totalCalories, err := apiCfg.DB.GetTotalCalorieIntakeByUserId(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
	}
	var obj pages.TrackProgress
	calories, defict, surplus := 0, 0, 0
	for _, i := range totalCalories {
		calories += int(i.Calories)
		defict += int(i.TotalDeficit)
		surplus += int(i.TotalSurplus)
	}
	obj.Id = userId
	obj.CreatedAt = userInput.CreatedAt
	obj.Program = userInput.Program
	var current_deficit float32 = (float32(defict) - float32(surplus)) / 7716.0

	if userInput.Program == program_fatLoss {
		obj.WeightProgress = current_deficit / float32(userInput.Weight-userInput.DesiredWeight.Int32) * 100
		createdTimePlusTimeFrame := obj.CreatedAt.Add(time.Duration(userInput.TimeFrame.Int32) * 7 * 24 * time.Hour)
		inversePercentage := float32(createdTimePlusTimeFrame.Local().Sub(time.Now().Local()).Hours()/(7*24)) / float32(userInput.TimeFrame.Int32) * 100
		obj.TimeFrameProgress = 100 - inversePercentage
	} else if userInput.Program == program_muscleGain {
		createdTimePlusTimeFrame := obj.CreatedAt.Add(time.Duration(userInput.TimeFrame.Int32) * 7 * 24 * time.Hour)
		inversePercentage := float32(createdTimePlusTimeFrame.Local().Sub(time.Now().Local()).Hours()/(7*24)) / float32(userInput.TimeFrame.Int32) * 100
		obj.TimeFrameProgress = 100 - inversePercentage
	} else {
		obj.WeightProgress = current_deficit
		obj.TimeFrameProgress = -1
	}
	pages.Profile(obj).Render(r.Context(), w)
}

func (apiCfg *Api) ProfileHandler(w http.ResponseWriter, r *http.Request) {
}
