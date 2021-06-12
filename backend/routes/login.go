package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/micahnico/awesomespotify/backend/authenticate"
	"github.com/micahnico/awesomespotify/backend/current"
)

type response struct {
	OK  bool
	Err error
}

func Login(w http.ResponseWriter, r *http.Request) {
	var resp response

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")

	client, err := authenticate.ConnectAccount(redirectURI, clientID, clientSecret)
	if err != nil {
		resp = response{OK: false, Err: fmt.Errorf("could not connect account")}
		return
	}
	current.SetClient(client)

	resp = response{OK: true, Err: nil}
	render.JSON(w, r, resp)
}
