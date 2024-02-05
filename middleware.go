package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/controllers"
	"log"
	"net/http"
	"os"
	"time"
)

func refresh(w http.ResponseWriter, refreshToken string) (bool, error) {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
			fmt.Println("something is wrong in parsing")
		}
		return SECRET, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not authorized"))
		fmt.Println("something is wrong in parsing ", err)
	}
	if token.Valid {
		subject, err := token.Claims.GetSubject()
		if err != nil {
			log.Printf("lint 32 %s", err)
			return false, err
		}
		newAccessToken, err := controllers.CreateJWT(time.Minute*5, subject)
		if err != nil {
			log.Fatal("error creating new access token")
			return false, err
		}
		fmt.Println(newAccessToken)
		cookie := http.Cookie{Name: "access-token", Value: newAccessToken, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
		http.SetCookie(w, &cookie)
	} else {
		log.Fatal("Login again")
		return false, errors.New("token not valid")
	}
	return true, nil
}

func validateJWT(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		cookies := r.Cookies()
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}
		if accessToken == nil || err != nil {
			w.Header().Add("HX-Trigger", `{ "errorToast" : "401 Unauthorized" }`)
			w.Header().Add("Hx-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
		}
		if accessToken != nil {
			token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.Header().Add("HX-Trigger", `{ "errorToast" : "401 Unauthorized" }`)
					w.Header().Add("Hx-Redirect", "http://localhost:5000/view/login")
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Println("something is wrong in parsing")
				}
				return SECRET, nil
			})
			if err != nil {
				fmt.Println(err)
				refresh_token, err := r.Cookie("refresh-token")
				if err != nil {
					log.Printf("line 77 %s", err)
				}
				if refresh_token != nil {
					status_ok, err := refresh(w, refresh_token.Value)
					if status_ok == false || err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(err.Error()))
						fmt.Println(err)
					} else {
						next(w, r)
					}
				}
			}
			if token.Valid {
				next(w, r)
			}
		} else {
			w.Header().Add("HX-Trigger", `{ "errorToast" : "401 Unauthorized" }`)
			w.Header().Add("Hx-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
