package httputil

import (
	"fmt"
	"net/http"
	"strings"
)

type Sort struct {
	Field     string `json:"field"`
	Ascending bool   `json:"ascending"` // ?sort=-type
}

func parseSorting(r *http.Request) ([]Sort, error) {
	query := r.URL.Query()
	sortString := query.Get("sort")
	if sortString == "" {
		return nil, nil
	}

	sortFields := strings.Split(sortString, ",")
	var sorters []Sort

	for _, field := range sortFields {
		ascending := true
		if field == "" {
			continue
		}
		if strings.HasPrefix(field, "-") {
			ascending = false
			field = strings.TrimPrefix(field, "-")
		}
		if field == "" {
			return nil, fmt.Errorf("empty sort field")
		}
		sorters = append(sorters, Sort{
			Field:     field,
			Ascending: ascending,
		})
	}

	return sorters, nil
}
