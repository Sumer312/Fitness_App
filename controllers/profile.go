package controllers

import (
	"fmt"
	"net/http"

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
	obj.CreatedAt = userInput.CreatedAt
	total := (calories - defict + surplus) / 7716.0
	if userInput.Program == program_fatLoss {
		obj.WeightProgress = float32(total/int(userInput.Deficit.Int32)) * 100
		obj.TimeFrameProgress = float32(2)
	} else if userInput.Program == program_muscleGain {
	} else {
	}
}

func (apiCfg *Api) ProfileHandler(w http.ResponseWriter, r *http.Request) {
}
