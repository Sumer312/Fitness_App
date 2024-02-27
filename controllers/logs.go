package controllers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sumer312/Health-App-Backend/views/pages"
)

func (apiCfg *Api) LogsRender(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Not logged in" }`)
		w.WriteHeader(400)
		return
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Not logged in" }`)
		w.WriteHeader(400)
		return
	}
	user_daily, err := apiCfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
	list := make([]pages.DailyLogs, 0)
	for _, ele := range user_daily {
		cur := pages.DailyLogs{
			Id:        ele.ID,
			CreatedAt: ele.CreatedAt,
			Calories:  float32(ele.Calories),
			Carbs:     float32(ele.Carbohydrates),
			Protein:   float32(ele.Protein),
			Fat:       float32(ele.Fat),
			Fiber:     float32(ele.Fiber),
		}
		list = append(list, cur)
	}
  var isEmpty bool
  if len(list) == 0 {
    isEmpty = true
  } else {
    isEmpty = false
  }
	pages.Logs(list, isEmpty).Render(r.Context(), w)
}
