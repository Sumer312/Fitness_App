package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/controllers"
	"github.com/sumer312/Health-App-Backend/views/pages"
	"github.com/sumer312/Health-App-Backend/views/partials"
)

func refresh(w http.ResponseWriter, refreshToken string) error {
	godotenv.Load()
	var SECRET []byte = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println("malformed refresh token")
		}
		return SECRET, nil
	})
	if err != nil {
		fmt.Println("something is wrong in parsing refresh token", err)
		return err
	}
	if token.Valid {
		subject, err := token.Claims.GetSubject()
		if err != nil {
			fmt.Println(err)
			return err
		}
		newAccessToken, err := controllers.CreateJWT(time.Minute*5, subject)
		if err != nil {
			fmt.Println("error createing access token", err)
			return err
		}
		cookie := http.Cookie{Name: "access-token", Value: newAccessToken, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Println("login again")
		return errors.New("token not valid")
	}
	return nil
}

func controllerMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	godotenv.Load()
	var SECRET []byte = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		if err != nil {
			fmt.Println(err)
			w.Header().Add("HX-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				fmt.Println("malformed access token")
			}
			return SECRET, nil
		})
		if err != nil {
			fmt.Println(err)
			refresh_token, err := r.Cookie("refresh-token")
			if err != nil {
				fmt.Println(err)
				w.Header().Add("HX-Redirect", "http://localhost:5000/view/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			err = refresh(w, refresh_token.Value)
			if err != nil {
				w.Header().Add("HX-Redirect", "http://localhost:5000/view/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			partials.DrawerAuthFlag = true
			next(w, r)
		}
		if token.Valid {
			partials.DrawerAuthFlag = true
			next(w, r)
		} else {
			w.Header().Add("HX-Redirect", "http://localhost:5000/view/login")
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	})
}

func viewRenderInControllerMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	godotenv.Load()
	var SECRET []byte = []byte(os.Getenv("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access-token")
		if err != nil {
			fmt.Println(err)
			pages.Login().Render(r.Context(), w)
			return
		}
		token, err := jwt.Parse(accessToken.Value, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				fmt.Println("malformed access token")
			}
			return SECRET, nil
		})
		if err != nil {
			fmt.Println(err)
			refresh_token, err := r.Cookie("refresh-token")
			if err != nil {
				fmt.Println(err)
				pages.Login().Render(r.Context(), w)
				return
			}
			err = refresh(w, refresh_token.Value)
			if err != nil {
				pages.Login().Render(r.Context(), w)
				return
			}
			partials.DrawerAuthFlag = true
			next(w, r)
			return
		}
		if token.Valid {
			partials.DrawerAuthFlag = true
			next(w, r)
		} else {
			pages.Login().Render(r.Context(), w)
		}
		return
	})
}
