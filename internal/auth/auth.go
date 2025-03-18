package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Auth struct {
}

const (
	key    = "MyKey123"
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("err when load env ", err)
	}
	googleClId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClSr := os.Getenv("GOOGLE_CLIENT_SECRET")
	store := sessions.NewCookieStore([]byte(key))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd
	gothic.Store = store
	goth.UseProviders(
		google.New(googleClId, googleClSr, "http://localhost:3000/auth/google/callback"),
	)

}
