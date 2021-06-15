package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/authenticate"
)

type loginResponse struct {
	Err error
}

func Login(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	client, err := authenticate.ConnectAccount(redirectURI, clientID, clientSecret)
	sadpath.Check(err)

	token, err := client.Token()
	sadpath.Check(err)

	// set cookies
	accessTokenCookie := &http.Cookie{Name: "AccessToken", Value: token.AccessToken, HttpOnly: false, Expires: token.Expiry}
	http.SetCookie(w, accessTokenCookie)
	refreshTokenCookie := &http.Cookie{Name: "RefreshToken", Value: token.RefreshToken, HttpOnly: false, Expires: token.Expiry}
	http.SetCookie(w, refreshTokenCookie)
	expiryToken := &http.Cookie{Name: "ExpiryToken", Value: token.Expiry.String(), HttpOnly: false, Expires: token.Expiry}
	http.SetCookie(w, expiryToken)

	render.JSON(w, r, loginResponse{Err: err})
}
