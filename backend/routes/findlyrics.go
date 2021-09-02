package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/go-chi/render"
	"github.com/gocolly/colly/v2"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/current"
	"github.com/sorucoder/colorhelper"
	"github.com/zmb3/spotify"
)

type findLyricsResponse struct {
	Artists     []string
	Song        string
	URLSafeSong string
	Lyrics      string
	ImageURL    string
	BgHex       string
	TxtHex      string
	Error       string
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

	currentSong := r.URL.Query().Get("currentSong")
	if formatSongName(removeRemixString(currentlyPlayingInfo.Item.Name)) == currentSong {
		result := findLyricsResponse{Error: "Already fetched lyrics"}
		render.JSON(w, r, result)
		return
	}

	currentSong = currentlyPlayingInfo.Item.Name
	currentArtists := getArtistNames(currentlyPlayingInfo.Item.Artists)
	albumImageURL := currentlyPlayingInfo.Item.Album.Images[1].URL // get the one that is 300x300
	mainColors, bgColorHex, err := getMainColorFromAlbumnCover(ctx, albumImageURL)
	textColor, err := getBestTextColor(mainColors)
	sadpath.Check(err)

	lyrics, err := search(ctx, currentArtists, currentSong)
	sadpath.Check(err)

	if lyrics == "" {
		result := findLyricsResponse{Artists: currentArtists, Song: currentSong, URLSafeSong: formatSongName(removeRemixString(currentSong)), ImageURL: albumImageURL, BgHex: bgColorHex, TxtHex: textColor, Error: "No lyrics found"}
		render.JSON(w, r, result)
		return
	}

	result := findLyricsResponse{Artists: currentArtists, Song: currentSong, URLSafeSong: formatSongName(removeRemixString(currentSong)), Lyrics: lyrics, ImageURL: albumImageURL, BgHex: bgColorHex, TxtHex: textColor}
	render.JSON(w, r, result)
}

func search(ctx context.Context, artists []string, song string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.genius.com/search?q=%v", formatURL(artists, formatSongName(song))), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("GENIUS_API_ACCESS_TOKEN")))
	client := &http.Client{}
	res, err := client.Do(req)
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
		if formatSongName(removeRemixString(result["title"].(string))) == formatSongName(removeRemixString(song)) {
			url = result["url"].(string)
			break
		}
	}

	if url == "" {
		return "", nil // no lyrics found
	}

	html := scrape(url)
	if err != nil {
		return "", err
	}

	return html, nil
}

func scrape(url string) string {
	var html string

	c := colly.NewCollector()
	c.OnHTML("div[class^='Lyrics__Container']", func(e *colly.HTMLElement) {
		section, _ := e.DOM.Html()
		html += (section + "<br>")
	})
	c.Visit(url)

	// remove all anchor tags, but keep the lyrics inside
	re := regexp.MustCompile(`<a[^>]+href=\"(.*?)\"[^>]*>(.*?)<\/a>`)
	html = re.ReplaceAllString(html, "$2") // $2 is the second capture group

	return html
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

func getMainColorFromAlbumnCover(ctx context.Context, url string) ([]color.Color, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:48.0) Gecko/20100101 Firefox/48.0")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	colors, err := prominentcolor.Kmeans(img)

	if err != nil {
		return nil, "", err
	}

	var colorsParsed []color.Color
	for _, co := range colors {
		colorsParsed = append(colorsParsed, color.RGBA{uint8(co.Color.R), uint8(co.Color.G), uint8(co.Color.B), 0xFF})
	}

	return colorsParsed, fmt.Sprintf("#%v", colors[0].AsString()), nil
}

func getBestTextColor(colors []color.Color) (string, error) {
	var bestTextColor color.Color
	if len(colors) > 0 {
		bestTextColor = colorhelper.PickBestTextColor(colors[0], colors[1:]...)
	} else {
		bestTextColor = colorhelper.PickBestTextColor(colors[0])
	}
	return colorhelper.MakeColorRepresentation(bestTextColor, colorhelper.HashedHexadecimalTripletRepresentation), nil
}
