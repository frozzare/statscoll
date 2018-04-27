package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/frozzare/go-httpapi"
	"github.com/jinzhu/gorm"
)

// Handler represents a the api handler.
type Handler struct {
	*httpapi.Router
	db *gorm.DB
}

// NewHandler creates a new handler.
func NewHandler(db *gorm.DB) (*Handler, error) {
	h := &Handler{
		Router: httpapi.NewRouter(),
		db:     db,
	}

	h.Post("/collect", h.handleCollect)
	h.Get("/stats/:metric", h.handleList)
	h.Get("/total/:metric", h.handleTotal)

	return h, nil
}

func (h *Handler) statsQuery(r *http.Request, ps httpapi.Params) (*gorm.DB, error) {
	qs := r.URL.Query()
	metric := ps.ByName("metric")

	query := h.db.Where("`metric` = ?", metric)

	// Parse start query string if any.
	if start := qs.Get("start"); len(start) > 0 {
		i, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			return nil, err
		}

		query = query.Where("timestamp >= ?", i)
	}

	// Parse end query string if any.
	if end := qs.Get("end"); len(end) > 0 {
		i, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			return nil, err
		}

		query = query.Where("timestamp <= ?", i)
	}

	// Add project query string if any.
	if project := qs.Get("project"); len(project) > 0 {
		query = query.Where("project = ?", project)
	}

	// Add last query string to filter last days.
	if last := qs.Get("last"); len(last) > 0 {
		i, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			return nil, err
		}
		query = query.Where("timestamp >= ?", time.Now().AddDate(0, 0, -int(i)).Unix())
	}

	return query, nil
}
