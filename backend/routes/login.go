package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/jackc/sadpath"
	"github.com/micahnico/awesomespotify/backend/authenticate"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURI  = os.Getenv("REDIRECT_URI")
)

type loginResponse struct {
	URL string `json:"url"`
	Err error  `json:"err"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	url, err := authenticate.StartConnectAccount(redirectURI, clientID, clientSecret)
	sadpath.Check(err)

	render.JSON(w, r, loginResponse{URL: url, Err: err})
}
