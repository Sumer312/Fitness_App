package partials

import (
	"github.com/joho/godotenv"
	"os"
)

var err error = godotenv.Load()
var base_url string = os.Getenv("BASE_URL")
var DrawerAuthFlag bool = false
