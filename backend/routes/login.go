package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/micahnico/awesomespotify/backend/current"
)

type loginResponse struct {
	Err error
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)

	var resp loginResponse
	if spotifyClient != nil {
		resp = loginResponse{Err: nil}
	} else {
		resp = loginResponse{Err: fmt.Errorf("Could not connect account")}
	}
	render.JSON(w, r, resp)
}
