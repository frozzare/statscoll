package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/stat"
)

var (
	errCreateStat = errors.New("Could not create new stat")
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

	// Execute query and find any errors.
	if err := query.Find(&stats).Error; err != nil || len(stats) == 0 {
		return nil, fmt.Errorf("No stats found for: %s", metric)
	}

	// Sort stats so the last one is listed first.
	sort.Slice(stats, func(i, j int) bool {
		return time.Unix(stats[i].Timestamp, 0).After(time.Unix(stats[j].Timestamp, 0))
	})

	return stats, nil
}
