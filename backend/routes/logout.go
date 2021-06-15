package routes

import (
	"net/http"

	"github.com/go-chi/render"
)

type logoutResponse struct {
	Status int32
}

func Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie := &http.Cookie{Name: "AccessToken", MaxAge: -1}
	http.SetCookie(w, accessTokenCookie)
	refreshTokenCookie := &http.Cookie{Name: "RefreshToken", MaxAge: -1}
	http.SetCookie(w, refreshTokenCookie)
	expiryToken := &http.Cookie{Name: "ExpiryToken", MaxAge: -1}
	http.SetCookie(w, expiryToken)

	render.JSON(w, r, logoutResponse{Status: 200})
}
