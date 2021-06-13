package routes

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
)

func TestRoute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)
	currentUser, err := spotifyClient.CurrentUser()
	sadpath.Check(err)
	render.JSON(w, r, currentUser)
}
