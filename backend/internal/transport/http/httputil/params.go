package httputil

import "net/http"

type QueryParams struct {
	PaginateBy *OffsetPagination `json:"paginate"`
	FilterBy   []Filter          `json:"filters"`
	SortBy     []Sort            `json:"sorters"`
}

// ParseQueryParams parses query parameters for pagination, filtering, and sorting
// TODO: Provide set of valid filters, sorters
func ParseQueryParams(r *http.Request) (*QueryParams, error) {
	if len(r.URL.Query()) == 0 {
		return nil, nil
	}

	paginateBy, err := parsePagination(r)
	if err != nil {
		return nil, err
	}

	filters, err := parseFilters(r)
	if err != nil {
		return nil, err
	}

	sorters, err := parseSorting(r)
	if err != nil {
		return nil, err
	}

	return &QueryParams{
		PaginateBy: paginateBy,
		FilterBy:   filters,
		SortBy:     sorters,
	}, nil
}
