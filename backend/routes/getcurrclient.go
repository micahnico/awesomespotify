package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/micahnico/awesomespotify/backend/current"
)

// type Employee struct {
// 	Name string `json:"name"`
// 	ID   int    `json:"id"`
// }

// func TestRoute(w http.ResponseWriter, r *http.Request) {
// 	employee := Employee{Name: "Micah Nicodemus", ID: 7}
// 	render.JSON(w, r, employee)
// }

func GetCurrClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	render.JSON(w, r, current.Client(ctx))
}
