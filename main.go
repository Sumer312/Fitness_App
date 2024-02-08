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

func swap(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	w.WriteHeader(200)
}

func main() {
	godotenv.Load()
	port := "5000"
	router := chi.NewRouter()
	fmt.Println("using chi")
	viewRouter := chi.NewRouter()
	serverRouter := chi.NewRouter()
	dbConnString := os.Getenv("DB_URL")
	conn, connerr := sql.Open("postgres",
		dbConnString,
	)
	log.Println(dbConnString)
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

	viewRouter.Handle("/login", templ.Handler(pages.Login()))
	viewRouter.Handle("/signup", templ.Handler(pages.Signup()))
	viewRouter.Handle("/user-input/fatloss", templ.Handler(pages.UserInputFatloss()))
	viewRouter.Handle("/user-input/muscle", templ.Handler(pages.UserInputMuscle()))
	viewRouter.Handle("/user-input/maintain", templ.Handler(pages.UserInputMaintain()))
	viewRouter.Handle("/kcal-calc", templ.Handler(pages.KcalCalc()))
	viewRouter.Handle("/logs", templ.Handler(pages.Logs()))
	viewRouter.Handle("/daily-input", templ.Handler(pages.DailyInput("70", "20", "50", "10", "40")))

	serverRouter.Post("/login", apiCfg.LoginHandler)
	serverRouter.Post("/signup", apiCfg.SignupHandler)
	serverRouter.Post("/user-input", validateJWT(apiCfg.InputHandler))
	serverRouter.Post("/calorie-tracker", validateJWT(apiCfg.CalorieInputHandler))
	serverRouter.Post("/nutrition-api-request", apiCfg.ApiRequest)
	serverRouter.HandleFunc("/profile", validateJWT(apiCfg.Profile))
	serverRouter.Post("/logs", apiCfg.Logs)

	router.Handle("/", templ.Handler(pages.Home()))
	router.Handle("/swap", http.HandlerFunc(swap))
	router.Mount("/view", viewRouter)
	router.Mount("/server", serverRouter)

	router.Handle("/*", templ.Handler(pages.NotFound()))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("Server starting on port " + port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("OOPs something went wrong")
	}
}
