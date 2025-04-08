package trainingrecords

import (
	"context"
	"github.com/VATUSA/primary-api/pkg/database/models"
	middleware "github.com/VATUSA/primary-api/pkg/go-chi/middleware/auth"
	"github.com/VATUSA/primary-api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func Router(r chi.Router) {
	r.Use(middleware.NotGuest)

	r.Route("/ots", func(r chi.Router) {
		r.Use(middleware.CanCROTSRecord)
		r.Get("/", listOTSRecords)
		r.Post("/", createOTSRecord)
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(middleware.CanCreateReadTrainingNote)
		r.Get("/", listNotes)
		r.Post("/", createNote)

		r.Route("/{NoteID}", func(r chi.Router) {
			r.Use(noteCtx)
			r.Get("/", getNote)
			r.With(middleware.CanDeleteTrainingNote).Delete("/", deleteNote)
		})
	})
}

func noteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "NoteID")
		if id == "" {
			http.Error(w, "Invalid training note ID", http.StatusBadRequest)
			return
		}

		noteID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid training note ID", http.StatusBadRequest)
			return
		}

		note := &models.TrainingNotes{ID: uint(noteID)}
		if err = note.Get(); err != nil {
			http.Error(w, "Invalid training note ID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), utils.TrainingNoteKey{}, note)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
