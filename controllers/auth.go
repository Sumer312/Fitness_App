package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"github.com/sumer312/Health-App-Backend/views/partials"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWT(expiresIn time.Duration, subject string) (string, error) {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["exp"] = time.Now().Add(expiresIn).Unix()
	claim["sub"] = subject
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenStr, nil
}

func (apiCfg *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	base_url := os.Getenv("BASE_URL")
	r.ParseForm()
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	if len(email) == 0 || len(password) == 0 {
		w.Header().Add("HX-Trigger", `{ "warnToast" : "Fields should not be empty" }`)
		w.WriteHeader(400)
		return
	}
	user, err := apiCfg.DB.GetUserByEmail(
		r.Context(),
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Add("HX-Trigger", `{ "errorToast" : "No such user" }`)
			w.WriteHeader(400)
			return
		}
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Cannot connect to database" }`)
		w.WriteHeader(500)
		return
	}
	passwordCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordCheck != nil {
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Wrong password" }`)
		w.WriteHeader(401)
		return
	}
	accessToken, err := CreateJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Println(err)
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Println(err)
	}
	access_cookie := http.Cookie{Name: access_token_cookie_name, Path: "/", Value: accessToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	refresh_cookie := http.Cookie{Name: refresh_token_cookie_name, Path: "/", Value: refreshToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	user_id := http.Cookie{Name: user_id_cookie_name, Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	TempChan := make(chan bool)
	go func() {
		TempChan <- func(tempW http.ResponseWriter, ac http.Cookie, rc http.Cookie, uid http.Cookie) bool {
			http.SetCookie(tempW, &ac)
			http.SetCookie(tempW, &rc)
			http.SetCookie(tempW, &uid)
			return true
		}(w, access_cookie, refresh_cookie, user_id)
	}()
	flag := <-TempChan
	if flag {
		partials.DrawerAuthFlag = true
		w.Header().Add("HX-Redirect", base_url)
		w.WriteHeader(200)
	}
}

func (apiCfg *Api) SignupHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	base_url := os.Getenv("BASE_URL")
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	_, err := apiCfg.DB.GetUserByEmail(r.Context(), sql.NullString{String: email, Valid: true})
	if err != sql.ErrNoRows {
		w.Header().Add("HX-Trigger", `{ "errorToast" : "User already exists" }`)
		w.WriteHeader(500)
		return
	}

	if len(name) == 0 || len(email) == 0 || len(password) == 0 || len(confirmPassword) == 0 {
		w.Header().Add("HX-Trigger", `{ "warnToast" : "Fields should not be empty" }`)
		w.WriteHeader(400)
		return
	}

	if confirmPassword != password {
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Passwords do not match" }`)
		w.WriteHeader(401)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error storing password")
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Error storing password" }`)
		w.WriteHeader(500)
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Email:     sql.NullString{String: email, Valid: true},
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Println(err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Cannot connect to database" }`)
		w.WriteHeader(500)
		return
	}
	accessToken, err := CreateJWT(time.Minute*5, user.Name)
	if err != nil {
		log.Println("access token signup error ", err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Server error" }`)
		w.WriteHeader(500)
		return
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Println("refresh token signup err ", err)
		w.Header().Add("HX-Trigger", `{ "errorToast" : "Server error" }`)
		w.WriteHeader(500)
		return
	}
	access_cookie := http.Cookie{Name: access_token_cookie_name, Path: "/", Value: accessToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	refresh_cookie := http.Cookie{Name: refresh_token_cookie_name, Path: "/", Value: refreshToken, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	user_id := http.Cookie{Name: user_id_cookie_name, Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}
	TempChan := make(chan bool)
	go func() {
		TempChan <- func(tempW http.ResponseWriter, ac http.Cookie, rc http.Cookie, uid http.Cookie) bool {
			http.SetCookie(tempW, &ac)
			http.SetCookie(tempW, &rc)
			http.SetCookie(tempW, &uid)
			return true
		}(w, access_cookie, refresh_cookie, user_id)
	}()
	flag := <-TempChan
	if flag {
		partials.DrawerAuthFlag = true
		w.Header().Add("HX-Redirect", base_url)
		w.WriteHeader(200)
	}
}

func (apiCfg *Api) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	base_url := os.Getenv("BASE_URL")
	access_cookie := http.Cookie{Name: access_token_cookie_name, Path: "/", Value: "", MaxAge: 0}
	refresh_cookie := http.Cookie{Name: refresh_token_cookie_name, Path: "/", Value: "", MaxAge: 0}
	user_id := http.Cookie{Name: user_id_cookie_name, Path: "/", Value: "", MaxAge: 0}
	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &user_id)
	partials.DrawerAuthFlag = false
	w.Header().Add("HX-Redirect", base_url)
	w.WriteHeader(200)
	return
}

func (apiCfg *Api)OAuth2() {
	godotenv.Load()
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback")
	)
}
