package routes

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/zmb3/spotify"
)

type simpleSimplePlaylist struct {
	Name   string     `json:"name"`
	ID     spotify.ID `json:"id"`
	ImgURL string     `json:"imgUrl"`
}

type getPlaylistsResponse struct {
	Playlists []simpleSimplePlaylist `json:"playlists"`
	Error     string                 `json:"error"`
}

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)
	if spotifyClient != nil {
		playlistPage, err := spotifyClient.CurrentUsersPlaylists()
		sadpath.Check(err)

		playlists := convertToSimpleSimplePlaylists(playlistPage.Playlists)
		response := getPlaylistsResponse{Playlists: playlists, Error: ""}
		render.JSON(w, r, response) // need to transform into the response type
		return
	}
	render.JSON(w, r, nil)
}

func convertToSimpleSimplePlaylists(in []spotify.SimplePlaylist) []simpleSimplePlaylist {
	var out []simpleSimplePlaylist
	for _, playlist := range in {
		var url string
		if len(playlist.Images) > 0 {
			url = playlist.Images[0].URL
		}
		out = append(out, simpleSimplePlaylist{Name: playlist.Name, ID: playlist.ID, ImgURL: url})
	}
	return out
}
