package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/zmb3/spotify"
)

type findLyricsResponse struct {
	Artists  []string
	Song     string
	Lyrics   string
	ImageURL string
	Error    string
}

func FindLyrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)

	currentlyPlayingInfo, err := spotifyClient.PlayerCurrentlyPlaying()
	sadpath.Check(err)

	if !currentlyPlayingInfo.Playing {
		result := findLyricsResponse{Error: "No currently playing song"}
		render.JSON(w, r, result)
		return
	}

	searchArtist := currentlyPlayingInfo.Item.Artists[0].Name // searching for the first/main artist should give the right result
	currentArtists := getArtistNames(currentlyPlayingInfo.Item.Artists)
	currentSong := currentlyPlayingInfo.Item.Name
	albumImageURL := currentlyPlayingInfo.Item.Album.Images[1].URL // get the one that is 300x300

	lyrics, err := search(ctx, searchArtist, currentSong)
	sadpath.Check(err)

	result := findLyricsResponse{Artists: currentArtists, Song: currentSong, Lyrics: lyrics, ImageURL: albumImageURL}
	render.JSON(w, r, result)
}

func search(ctx context.Context, artist string, song string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.genius.com/search?q=%v", formatURL(artist, song)), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("GENIUS_API_ACCESS_TOKEN")))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	hits := body["response"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return "", errors.New("No lyrics found")
	}
	url := hits[0].(map[string]interface{})["result"].(map[string]interface{})["url"].(string)
	html, err := scrape(ctx, url)
	if err != nil {
		return "", err
	}

	return html, nil
}

func scrape(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:48.0) Gecko/20100101 Firefox/48.0")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var html string
	document.Find(".dVtOne").Each(func(i int, s *goquery.Selection) {
		section, _ := s.Html()
		html += section
	})

	return html, nil
}

func formatURL(x string, y string) string {
	str := x + " " + y
	str = strings.TrimSpace(str)
	str = strings.Replace(str, " ", "%20", -1)
	str = strings.Replace(str, "'", "", -1)
	return str
}

func getArtistNames(artists []spotify.SimpleArtist) []string {
	var names []string
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return names
}
