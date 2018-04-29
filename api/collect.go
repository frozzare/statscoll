package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/frozzare/statscoll/stat"
)

var (
	errCollectStat = errors.New("Could not collect new stat")
)

func (h *Handler) handleCollect(r *http.Request) (interface{}, interface{}) {
	var stat *stat.Stat

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&stat); err != nil {
		return nil, err
	}

	if stat.Timestamp == 0 {
		stat.Timestamp = time.Now().Unix()
	}

	if err := h.db.Create(stat).Error; err != nil {
		return nil, errCollectStat
	}

	if err := h.cache.RemovePrefix(fmt.Sprintf("%s_%s", stat.Metric, stat.Project)); err != nil {
		log.Println(err)
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}
