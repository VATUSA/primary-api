package utils

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		log.Printf("Error rendering response: %v", err)
	}
}
