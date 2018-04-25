package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/stat"
)

var (
	errCreateStat   = errors.New("Could not create new stat")
	errNoStatsFound = errors.New("No stats found")
)

func (h *Handler) handleCreate(r *http.Request) (interface{}, interface{}) {
	var stat *stat.Stat

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
		return nil, errCreateStat
	}

	if stat.Timestamp == 0 {
		stat.Timestamp = time.Now().Unix()
	}

	if err := h.db.Create(stat).Error; err != nil {
		return nil, errCreateStat
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (h *Handler) handleList(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
	var stats []*stat.Stat

	query, err := h.statsQuery(r, ps)
	if err != nil {
		return nil, errNoStatsFound
	}

	// Execute query and find any errors.
	if err := query.Find(&stats).Error; err != nil || len(stats) == 0 {
		return nil, errNoStatsFound
	}

	// Sort stats so the last one is listed first.
	sort.Slice(stats, func(i, j int) bool {
		return time.Unix(stats[i].Timestamp, 0).After(time.Unix(stats[j].Timestamp, 0))
	})

	return stats, nil
}
