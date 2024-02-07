package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *Api) Logs(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("user-id")
	if err != nil {
		log.Fatalln(err)
	}
	userId, err := uuid.Parse(cookieVal.Value)
	if err != nil {
		log.Println(err)
	}
	user, err := apiCfg.DB.GetUserInputById(r.Context(), userId)
	fmt.Println(user)
}
