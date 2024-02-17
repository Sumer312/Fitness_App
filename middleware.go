package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/controllers"
	"github.com/sumer312/Health-App-Backend/views/pages"
	"github.com/sumer312/Health-App-Backend/views/partials"
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
			log.Printf("%s\n", err)
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

func controllerMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		if accessToken == nil || err != nil {
			w.Header().Add("Hx-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		if accessToken != nil {
			token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
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
					log.Printf("%s\n", err)
				}
				if refresh_token != nil {
					status_ok, err := refresh(w, refresh_token.Value)
					if status_ok == false || err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(err.Error()))
						fmt.Println(err)
						return
					} else {
						partials.DrawerAuthFlag = true
						next(w, r)
					}
				}
			}
			if token.Valid {
				partials.DrawerAuthFlag = true
				next(w, r)
			}
		} else {
			w.Header().Add("Hx-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

func viewRenderInControllerMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		if accessToken == nil || err != nil {
			pages.Login().Render(r.Context(), w)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}
		if accessToken != nil {
			token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					pages.Login().Render(r.Context(), w)
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Println("something is wrong in parsing")
				}
				return SECRET, nil
			})
			if err != nil {
				fmt.Println(err)
				refresh_token, err := r.Cookie("refresh-token")
				if err != nil {
					log.Printf("%s\n", err)
				}
				if refresh_token != nil {
					status_ok, err := refresh(w, refresh_token.Value)
					if status_ok == false || err != nil {
						pages.Login().Render(r.Context(), w)
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte(err.Error()))
						fmt.Println(err)
						return
					} else {
						partials.DrawerAuthFlag = true
						next(w, r)
					}
				}
			}
			if token.Valid {
				partials.DrawerAuthFlag = true
				next(w, r)
			}
		} else {
			pages.Login().Render(r.Context(), w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}
