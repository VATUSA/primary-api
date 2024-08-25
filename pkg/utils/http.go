package utils

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func Response(r *http.Request, code int) {
	render.Status(r, code)
}

func Render(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		log.Printf("Error rendering response: %v", err)
	}
}

func JSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	render.Status(r, code)
	render.JSON(w, r, data)
}

func TempRedirect(w http.ResponseWriter, r *http.Request, location string) {
	http.Redirect(w, r, location, http.StatusTemporaryRedirect)
}
