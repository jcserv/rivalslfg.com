package httputil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// source: https://bookstack.soffid.com/books/scim/page/scim-query-syntax

type FilterOp string

const (
	Eq FilterOp = "eq"
)

type Filter struct {
	Field string   `json:"field"`
	Op    FilterOp `json:"op"`
	Value any      `json:"value"`
}

func parseFilters(r *http.Request) ([]Filter, error) {
	query := r.URL.Query()
	filterString := query.Get("filter")
	if filterString == "" {
		return nil, nil
	}

	var filters []Filter
	filterParts := strings.Split(filterString, "and")

	for _, part := range filterParts {
		fields := strings.Fields(part)
		if len(fields) < 3 {
			return nil, fmt.Errorf("invalid filter format: %s", part)
		}

		field := fields[0]
		op := strings.ToLower(fields[1])
		if op == "" {
			return nil, fmt.Errorf("missing operator in filter: %s", part)
		}
		if op != string(Eq) {
			return nil, fmt.Errorf("unsupported operator: %s", op)
		}
		value := parseFilterValue(strings.Join(fields[2:], " "))

		filters = append(filters, Filter{
			Field: field,
			Op:    FilterOp(op),
			Value: value,
		})
	}

	return filters, nil
}

func parseFilterValue(value string) any {
	if value == "true" {
		return true
	}

	if value == "false" {
		return false
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}

	return value
}
