package controllers

import (
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type Api struct {
	DB *database.Queries
}

func CreateJWT(expiresIn time.Duration, subject string) (string, error) {
	godotenv.Load()
	var SECRET = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["exp"] = time.Now().Add(expiresIn).Unix()
	claim["sub"] = subject
	tokenStr, err := token.SignedString(SECRET)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenStr, nil
}

func (apiCfg *Api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
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
		log.Println(err)
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
		log.Fatalln(err)
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln(err)
	}
	access_cookie := http.Cookie{Name: "access-token", Path: "/", Value: accessToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	refresh_cookie := http.Cookie{Name: "refresh-token", Path: "/", Value: refreshToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	user_id := http.Cookie{Name: "user-id", Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &user_id)
	/* w.Header().Add("HX-Redirect", "http://localhost:5000") */
	w.WriteHeader(200)
}

func (apiCfg *Api) SignupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

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
		log.Fatal("error storing password")
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
		log.Fatalln("access token signup error ", err)
	}
	refreshToken, err := CreateJWT(time.Hour*840, user.Name)
	if err != nil {
		log.Fatalln("refresh token signup err ", err)
	}
	access_cookie := http.Cookie{Name: "access-token", Path: "/", Value: accessToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	refresh_cookie := http.Cookie{Name: "refresh-token", Path: "/", Value: refreshToken, HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	user_id := http.Cookie{Name: "user_id", Path: "/", Value: user.ID.String(), HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode}
	/* w.Header().Add("HX-Redirect", "http://localhost:5000") */
	w.WriteHeader(200)
	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)
	http.SetCookie(w, &user_id)
}
