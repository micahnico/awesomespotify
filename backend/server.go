package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/micahnico/awesomespotify/backend/authenticate"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/micahnico/awesomespotify/backend/routes"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURI  = os.Getenv("REDIRECT_URI")
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
		r.Get("/user/get", routes.GetUser)
		r.Get("/login", routes.Login)
		r.Get("/logout", routes.Logout)
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

			accessToken, _ := r.Cookie("AccessToken")
			refreshToken, _ := r.Cookie("RefreshToken")

			// same browsing session
			if accessToken != nil && refreshToken != nil && current.SpotifyClient(ctx) != nil {
				fmt.Println("client was set")
				token, _ := current.SpotifyClient(ctx).Token()
				authenticate.SetCookies(w, token)
			}

			// been less than an hour since last use
			if accessToken != nil && refreshToken != nil && current.SpotifyClient(ctx) == nil {
				currToken := &oauth2.Token{
					AccessToken:  accessToken.Value,
					RefreshToken: refreshToken.Value,
				}

				auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying)
				auth.SetAuthInfo(clientID, clientSecret)
				client := auth.NewClient(currToken)
				newToken, _ := client.Token()
				authenticate.SetCookies(w, newToken)

				current.SetSpotifyClient(&client)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
