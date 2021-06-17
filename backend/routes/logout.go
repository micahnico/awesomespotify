package routes

import (
	"net/http"

	"github.com/go-chi/render"
)

type logoutResponse struct {
	Status int32
}

func Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie := &http.Cookie{Name: "AccessToken", MaxAge: -1, Path: "/"}
	http.SetCookie(w, accessTokenCookie)
	refreshTokenCookie := &http.Cookie{Name: "RefreshToken", MaxAge: -1, Path: "/"}
	http.SetCookie(w, refreshTokenCookie)

	render.JSON(w, r, logoutResponse{Status: 200})
}
