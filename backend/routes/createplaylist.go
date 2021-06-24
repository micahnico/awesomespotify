package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/autoplaylist"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/zmb3/spotify"
)

type createPlaylistResponse struct {
	Name string `json:"name"`
}

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	sadpath.Check(err)

	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)
	p, err := autoplaylist.NewAutoPlaylist(spotifyClient, "DISCOVER", "by Awesome Spotify", 50, spotify.ID(body["id"].(string)))
	sadpath.Check(err)

	render.JSON(w, r, createPlaylistResponse{Name: p.Name})
}
