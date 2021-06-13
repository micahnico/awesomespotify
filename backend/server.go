package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/micahnico/awesomespotify/backend/authenticate"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/micahnico/awesomespotify/backend/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(loadSpotifyClientHandler())

	r.Route("/api", func(r chi.Router) {
		r.Get("/login", routes.Login)
		r.Get("/test", routes.TestRoute)
		r.Get("/callback", authenticate.CompleteAuth)
		r.Get("/lyrics/find", routes.FindLyrics)
	})

	log.Println("Server started on on port 8081")
	http.ListenAndServe(":8081", r)
}

func loadSpotifyClientHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if r.URL.Path == "/api/callback" {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if current.SpotifyClient(ctx) != nil || r.URL.Path == "/api/login" {
				ctx = current.WithSpotifyClient(ctx, current.SpotifyClient(ctx))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			clientID := os.Getenv("CLIENT_ID")
			clientSecret := os.Getenv("CLIENT_SECRET")
			redirectURI := os.Getenv("REDIRECT_URI")

			client, err := authenticate.ConnectAccount(redirectURI, clientID, clientSecret)
			if err != nil {
				log.Fatal("Could not connect account")
				return
			}

			ctx = current.WithSpotifyClient(ctx, client)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
