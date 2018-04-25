package api

import (
	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/db"
)

// Handler represents a the api handler.
type Handler struct {
	*httpapi.Router
	db *db.DB
}

// NewHandler creates a new handler.
func NewHandler(db *db.DB) (*Handler, error) {
	h := &Handler{
		Router: httpapi.NewRouter(),
		db:     db,
	}

	h.Post("/stats", h.handleCreate)
	h.Get("/stats/:metric", h.handleList)

	return h, nil
}
