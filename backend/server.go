package main

import (
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
		r.Post("/login", routes.Login)
		r.Post("/logout", routes.Logout)
		r.Get("/callback", authenticate.CompleteAuth)
		r.Get("/find", routes.FindLyrics)
		r.Get("/playlists/get", routes.GetPlaylists)
		r.Post("/playlists/create", routes.CreatePlaylist)
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

			if accessToken != nil && refreshToken != nil {
				currToken := &oauth2.Token{
					AccessToken:  accessToken.Value,
					RefreshToken: refreshToken.Value,
				}

				auth := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopePlaylistReadCollaborative, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
				auth.SetAuthInfo(clientID, clientSecret)
				client := auth.NewClient(currToken)
				newToken, _ := client.Token()
				authenticate.SetCookies(w, newToken)

				ctx = current.WithSpotifyClient(ctx, &client)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
