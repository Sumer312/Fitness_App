package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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
		newAccessToken, err := createJWT(time.Minute*5, subject)
		if err != nil {
			log.Fatal("error creating new access token")
			return false, err
		}
		fmt.Println(newAccessToken)
		cookie := http.Cookie{Name: "access-token", Value: newAccessToken, HttpOnly: true}
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
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no access token found"))
			fmt.Println(err)
		}
		if accessToken != nil {
			token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
					fmt.Println("something is wrong in parsing")
				}
				return SECRET, nil
			})
			if err != nil {
				fmt.Println(err)
				fmt.Println(r.Header["Authorization"])
				if r.Header["Authorization"] != nil {
					auth_header := r.Header.Get("Authorization")
					split_arr := strings.Split(auth_header, " ")
					refreshToken := split_arr[1]
					fmt.Print(refreshToken)
					status_ok, err := refresh(w, refreshToken)
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
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}
