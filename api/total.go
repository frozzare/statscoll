package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/stat"
)

var (
	errNoTotalStatsFound = errors.New("No total stats found")
)

func (h *Handler) handleTotal(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
	key := fmt.Sprintf("%s_%s_%s", r.URL.Query().Get("project"), ps.ByName("metric"), r.URL.String())
	if v, err := h.cache.Get(key); err == nil {
		return v, nil
	}

	query, err := h.statsQuery(r, ps)
	if err != nil {
		return nil, errNoTotalStatsFound
	}

	var result struct {
		Total float64 `json:"total"`
	}

	if err := query.Model(&stat.Stat{}).Select("sum(value) as total").Scan(&result).Error; err != nil {
		return nil, err
	}

	if err := h.cache.Set(key, result); err != nil {
		return nil, err
	}

	return result, nil
}
