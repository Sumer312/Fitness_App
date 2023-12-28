package main

import (
	"database/sql"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sumer312/Health-App-Backend/controllers"
	"github.com/sumer312/Health-App-Backend/internal/database"
	"github.com/sumer312/Health-App-Backend/views/pages"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load()
	port := "5000"
	router := chi.NewRouter()
	fmt.Println("using chi")
	authRouter := chi.NewRouter()
	dbConnString := os.Getenv("DB_URL")
	log.Println(dbConnString)
	conn, connerr := sql.Open("postgres",
		dbConnString,
	)
	if connerr != nil {
		log.Fatalln("error conncting", connerr)
	}
	db := database.New(conn)
	apiCfg := controllers.Api{DB: db}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	router.Handle("/login_page", templ.Handler(pages.Login()))
	router.Handle("/user_input_page", templ.Handler(pages.UserInput()))
	authRouter.Post("/login", apiCfg.LoginHandler)
	authRouter.Post("/signup", apiCfg.SignupHandler)
	router.Post("/calorie-tracker", validateJWT(apiCfg.CalorieInputHandler))
	router.Post("/user-input", validateJWT(apiCfg.InputHandler))
	router.HandleFunc("/profile", validateJWT(apiCfg.Profile))

	router.Mount("/auth", authRouter)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("Server starting on port:" + port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("OOPs something went wrong")
	}
}
