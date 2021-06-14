package routes

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)
	if spotifyClient != nil {
		currentUser, err := spotifyClient.CurrentUser()
		sadpath.Check(err)
		render.JSON(w, r, currentUser)
		return
	}
	render.JSON(w, r, nil)
}
