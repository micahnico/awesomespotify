package authenticate

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jackc/sadpath"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var state = randStringRunes(42)
var auth spotify.Authenticator

func StartConnectAccount(redirectURI string, clientID string, secretKey string) (string, error) {
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopePlaylistReadCollaborative, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
	auth.SetAuthInfo(clientID, secretKey)
	url := auth.AuthURL(state)
	// url := auth.AuthURLWithDialog(state)

	return url, nil
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)

	token, err := client.Token()
	sadpath.Check(err)
	SetCookies(w, token)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SetCookies(w http.ResponseWriter, token *oauth2.Token) {
	accessTokenCookie := &http.Cookie{Name: "AccessToken", Value: token.AccessToken, Path: "/", HttpOnly: false, Expires: time.Now().Add(time.Hour)}
	http.SetCookie(w, accessTokenCookie)
	refreshTokenCookie := &http.Cookie{Name: "RefreshToken", Value: token.RefreshToken, Path: "/", HttpOnly: false, Expires: time.Now().Add(time.Hour)}
	http.SetCookie(w, refreshTokenCookie)
}
