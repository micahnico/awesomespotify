package main

import (
	"log"
	"net/http"
	"time"

	"github.com/micahnico/awesomespotify/backend/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	// r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api", func(r chi.Router) {
		r.Get("/test", routes.TestRoute)
	})

	log.Println("Server started on on port 8081")
	http.ListenAndServe(":8081", r)
}
