package httputil

import (
	"net/http"
	"strconv"
)

type OffsetPagination struct {
	Limit  int  `json:"limit"`
	Offset int  `json:"offset"`
	Count  bool `json:"count"`
}

func parsePagination(r *http.Request) OffsetPagination {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	count := query.Get("count") == "true"

	return OffsetPagination{
		Limit:  limit,
		Offset: offset,
		Count:  count,
	}
}
