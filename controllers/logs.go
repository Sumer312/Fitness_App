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
		log.Fatalln(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
	}
	user_daily, err := apiCfg.DB.GetDailyNutritionOfUserByUserId(r.Context(), userId)
	list := make([]pages.DailyLogs, 0)
	for _, ele := range user_daily {
		cur := pages.DailyLogs{
			Id:        ele.ID,
			CreatedAt: ele.CreatedAt,
			Calories:  float32(ele.Calories),
			Carbs:     float32(ele.Carbohydrates),
			Protien:   float32(ele.Protien),
			Fat:       float32(ele.Fat),
			Fiber:     float32(ele.Fiber),
		}
		list = append(list, cur)
	}
	pages.Logs(list).Render(r.Context(), w)
}
