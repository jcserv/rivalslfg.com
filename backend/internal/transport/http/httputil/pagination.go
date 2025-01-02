package httputil

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	MAX_LIMIT = 250
)

type OffsetPagination struct {
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	Count  bool `json:"count"`
}

func parsePagination(r *http.Request) (*OffsetPagination, error) {
	query := r.URL.Query()
	limitVal := query.Get("limit")
	limit, _ := strconv.Atoi(limitVal)
	if limit < 0 {
		return nil, fmt.Errorf("limit must be greater than or equal to 0")
	}
	if limit > MAX_LIMIT {
		return nil, fmt.Errorf("limit must be less than or equal to %d", MAX_LIMIT)
	}

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		return nil, fmt.Errorf("offset must be greater than or equal to 0")
	}

	// If no limit is provided, use the max limit
	if limitVal == "" {
		limit = MAX_LIMIT
	}

	count := query.Get("count") == "true"
	return &OffsetPagination{
		Limit:  limit,
		Offset: offset,
		Count:  count,
	}, nil
}
