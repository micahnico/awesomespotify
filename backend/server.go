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
		r.Get("/callback", authenticate.CompleteAuth)
		r.Get("/lyrics/find", routes.FindLyrics)
	})

	log.Println("Server started on on port 8081")
	http.ListenAndServe(":8081", r)
}

// func loadSpotifyClientHandler() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {
// 			ctx := r.Context()
// 			ctx = current.WithSpotifyClient(ctx)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		}
// 		return http.HandlerFunc(fn)
// 	}
// }

func loadSpotifyClientHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			expiry, err := r.Cookie("ExpiryToken")
			if err == nil {
				accessToken, _ := r.Cookie("AccessToken")
				refreshToken, _ := r.Cookie("RefreshToken")

				token := &oauth2.Token{
					Expiry:       expiry.Expires,
					TokenType:    "Bearer",
					AccessToken:  accessToken.Value,
					RefreshToken: refreshToken.Value,
				}

				clientID := os.Getenv("CLIENT_ID")
				clientSecret := os.Getenv("CLIENT_SECRET")
				redirectURI := os.Getenv("REDIRECT_URI")

				auth := spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserTopRead, spotify.ScopeUserLibraryRead, spotify.ScopeUserFollowRead, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadCurrentlyPlaying)
				auth.SetAuthInfo(clientID, clientSecret)
				client := auth.NewClient(token)

				ctx = current.WithSpotifyClient(ctx, &client)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
