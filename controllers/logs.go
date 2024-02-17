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
  list := make([]pages.DailyLogs, 0);
	for i := 0; i < len(user_daily); i++ {
    cur := pages.DailyLogs{
			Id:        user_daily[i].ID,
			CreatedAt: user_daily[i].CreatedAt,
			Calories:  float32(user_daily[i].Calories),
			Carbs:     float32(user_daily[i].Carbohydrates),
			Protien:   float32(user_daily[i].Protien),
			Fat:       float32(user_daily[i].Fat),
			Fiber:     float32(user_daily[i].Fiber),
		}
    list = append(list, cur)
	}
	pages.Logs(list).Render(r.Context(), w)
}
