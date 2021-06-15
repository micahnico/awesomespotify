package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
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
	BgHex    string
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

	// searchArtist := currentlyPlayingInfo.Item.Artists[0].Name // searching for the first/main artist should give the right result
	currentArtists := getArtistNames(currentlyPlayingInfo.Item.Artists)
	currentSong := currentlyPlayingInfo.Item.Name
	albumImageURL := currentlyPlayingInfo.Item.Album.Images[1].URL // get the one that is 300x300

	lyrics, err := search(ctx, currentArtists, currentSong)
	sadpath.Check(err)

	if lyrics == "" {
		result := findLyricsResponse{Error: "No lyrics found"}
		render.JSON(w, r, result)
		return
	}

	bgColorHex, err := getMainColorFromAlbumnCover(ctx, albumImageURL)
	sadpath.Check(err)

	result := findLyricsResponse{Artists: currentArtists, Song: currentSong, Lyrics: lyrics, ImageURL: albumImageURL, BgHex: bgColorHex}
	render.JSON(w, r, result)
}

func search(ctx context.Context, artists []string, song string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.genius.com/search?q=%v", formatURL(artists, formatSongName(song))), nil)
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

	var url string
	for _, hit := range hits {
		result := hit.(map[string]interface{})["result"].(map[string]interface{})
		// fmt.Println(formatSongName(removeRemixString(result["title"].(string))), "|", formatSongName(removeRemixString(song)))
		if formatSongName(removeRemixString(result["title"].(string))) == formatSongName(removeRemixString(song)) {
			url = result["url"].(string)
			break
		}
	}

	if url == "" {
		return "", nil // no lyrics found
	}

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

func formatURL(artists []string, song string) string {
	var str string
	for i, artist := range artists {
		str = fmt.Sprintf("%v %v", str, artist)
		if i == 1 {
			break
		}
	}

	str = str + " " + song
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "&", "and", -1)
	str = strings.Replace(str, " ", "%20", -1)

	return str
}

func formatSongName(song string) string {
	regexList := []string{
		` - .+`,                     // Removes values after " - ...". from song name.
		`(?i)\(.*?feat.*?\)`,        // Removes all (...feat...)s from song name.
		`(?i)\[.*?feat.*?\]`,        // Removes all [...feat..]s from song name.
		`(?i)\(.*?remastered.*?\)`,  // Removes all (...remastered...)s from song name.
		`(?i)\[.*?remastered.*?\)]`, // Removes all [...remastered...]s from song name.
		`(?i)\(.*?cover.*?\)`,       // Removes all (...cover...)s from song name.
		`(?i)\[.*?cover.*?\]`,       // Removes all [...cover...]s from song name.
		`(?i)\(.*?with.*?\)`,        // Removes all (...with...)s from song name.
		`(?i)\[.*?with.*?\]`,        // Removes all [...with...]s from song name
		// `(?i)\(.*?Mix.*?\)`,
	}

	// Run regexs.
	for _, value := range regexList {
		re := regexp.MustCompile(value)
		song = re.ReplaceAllString(song, "")
	}

	// Trim spaces & other stuff
	song = strings.Replace(song, "'", "", -1)
	song = strings.Replace(song, "’", "", -1)
	song = strings.Replace(song, ".", "", -1)
	song = strings.Replace(song, "-", "", -1)

	// remove any left over parentheses
	song = strings.Replace(song, "(", "", -1)
	song = strings.Replace(song, ")", "", -1)
	song = strings.Replace(song, "[", "", -1)
	song = strings.Replace(song, "]", "", -1)
	song = strings.TrimSpace(song)
	song = strings.ToLower(song)

	return song
}

func removeRemixString(song string) string {
	regexList := []string{
		`(?i)\(.*?Remix.*?\)`,
		`(?i)\(.*?remix.*?\)`,
	}

	for _, value := range regexList {
		re := regexp.MustCompile(value)
		song = re.ReplaceAllString(song, "")
	}

	return song
}

func getArtistNames(artists []spotify.SimpleArtist) []string {
	var names []string
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return names
}

func getMainColorFromAlbumnCover(ctx context.Context, url string) (string, error) {
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

	img, _, err := image.Decode(res.Body)
	colors, err := prominentcolor.Kmeans(img)

	if err != nil {
		return "", err
	}

	return ("#" + colors[0].AsString()), nil
}
