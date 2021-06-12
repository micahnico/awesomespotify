package routes

import (
	"net/http"

	"github.com/go-chi/render"
)

type Employee struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestRoute(w http.ResponseWriter, r *http.Request) {
	employee := Employee{Name: "Micah Nicodemus", ID: 7}
	render.JSON(w, r, employee)
}
