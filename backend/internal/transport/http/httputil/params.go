package httputil

import "net/http"

type RequestParams struct {
	PaginateBy OffsetPagination `json:"paginate"`
	FilterBy   []Filter         `json:"filters"`
	SortBy     []Sort           `json:"sorters"`
}

// ParseQueryParams parses query parameters for pagination, filtering, and sorting
func ParseQueryParams(r *http.Request) (*RequestParams, error) {
	if len(r.URL.Query()) == 0 {
		return nil, nil
	}

	params := &RequestParams{}
	params.PaginateBy = parsePagination(r)

	filters, err := parseFilters(r)
	if err != nil {
		return nil, err
	}

	params.FilterBy = filters

	sorters, err := parseSorting(r)
	if err != nil {
		return nil, err
	}

	params.SortBy = sorters
	return params, nil
}
