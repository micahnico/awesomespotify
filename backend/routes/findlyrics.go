package routes

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
)

type findLyricsResponse struct {
	Artist string
	Song   string
	Lyrics string
}

func FindLyrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	spotifyClient := current.SpotifyClient(ctx)

	currenlyPlayingInfo, err := spotifyClient.PlayerCurrentlyPlaying()
	sadpath.Check(err)

	var currentArtist string
	var currentSong string
	if currenlyPlayingInfo != nil {
		currentArtist = currenlyPlayingInfo.Item.Artists[0].Name
		currentSong = currenlyPlayingInfo.Item.Name
	}

	lyrics, err := search(ctx, currentArtist, currentSong)
	sadpath.Check(err)

	result := findLyricsResponse{Artist: currentArtist, Song: currentSong, Lyrics: lyrics}
	render.JSON(w, r, result)
}

func search(ctx context.Context, artist string, song string) (string, error) {
	url := fmt.Sprintf("https://genius.com/%v-%v-lyrics", formatURL(artist), formatURL(song))

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

func formatURL(x string) string {
	reg := regexp.MustCompile(`(feat.\s[a-zA-Z0-9_,&\s]*)`)
	featMatch := reg.FindStringSubmatch(x)

	str := x
	if featMatch != nil {
		str = strings.Replace(x, featMatch[1], "", -1)
	}
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, ")", "", -1)
	str = strings.TrimSpace(str)
	str = strings.Replace(str, " - ", " ", -1)
	str = strings.Replace(str, " ", "-", -1)
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "&", "and", -1)
	return str
}
