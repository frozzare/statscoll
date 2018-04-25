package api

import (
	"errors"
	"net/http"

	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/stat"
)

var (
	errNoTotalStatsFound = errors.New("No total stats found")
)

func (h *Handler) handleTotal(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
	query, err := h.statsQuery(r, ps)
	if err != nil {
		return nil, errNoTotalStatsFound
	}

	var result struct {
		Total int64 `json:"total"`
	}

	if err := query.Model(&stat.Stat{}).Select("sum(count) as total").Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
