package httputil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryParams(t *testing.T) {
	defaultPagination := &OffsetPagination{Limit: 250, Offset: 0, Count: false}

	tests := []struct {
		name          string
		query         string
		expected      *QueryParams
		expectedError bool
	}{
		{
			name:          "No Query Params",
			query:         "",
			expected:      nil,
			expectedError: false,
		},
		{
			name:  "Pagination Only",
			query: "limit=50&offset=10&count=true",
			expected: &QueryParams{
				PaginateBy: &OffsetPagination{Limit: 50, Offset: 10, Count: true},
			},
			expectedError: false,
		},
		{
			name:  "Filters Only",
			query: `filter=region eq "na" and active eq true`,
			expected: &QueryParams{
				FilterBy: []Filter{
					{Field: "region", Op: Eq, Value: "na"},
					{Field: "active", Op: Eq, Value: true},
				},
				PaginateBy: defaultPagination,
			},
			expectedError: false,
		},
		{
			name:  "Sorting Only",
			query: "sort=-region,gamemode",
			expected: &QueryParams{
				SortBy: []Sort{
					{Field: "region", Ascending: false},
					{Field: "gamemode", Ascending: true},
				},
				PaginateBy: defaultPagination,
			},
			expectedError: false,
		},
		{
			name:  "Combined Query",
			query: `filter=region eq "na" and active eq true&limit=100&offset=0&count=true&sort=-region,gamemode`,
			expected: &QueryParams{
				PaginateBy: &OffsetPagination{Limit: 100, Offset: 0, Count: true},
				FilterBy: []Filter{
					{Field: "region", Op: Eq, Value: "na"},
					{Field: "active", Op: Eq, Value: true},
				},
				SortBy: []Sort{
					{Field: "region", Ascending: false},
					{Field: "gamemode", Ascending: true},
				},
			},
			expectedError: false,
		},
		{
			name:          "Invalid Filter Format",
			query:         `filter=region eq`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Empty Filter Field",
			query:         `filter= eq "value"`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Empty Filter Operator",
			query:         `filter=field "" "value"`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Empty Filter Value",
			query:         `filter=field eq`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Malformed Filter String",
			query:         `filter=field eq value and`,
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "Special Characters in Filter",
			query: `filter=name eq "John Doe @work"`,
			expected: &QueryParams{
				FilterBy: []Filter{
					{Field: "name", Op: Eq, Value: "John Doe @work"},
				},
				PaginateBy: defaultPagination,
			},
			expectedError: false,
		},
		{
			name:          "Empty Sort Field",
			query:         `sort=,-field`,
			expected:      &QueryParams{SortBy: []Sort{{Field: "field", Ascending: false}}, PaginateBy: defaultPagination},
			expectedError: false,
		},
		{
			name:          "Negative Limit",
			query:         "limit=-10&offset=-5",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "Non-Numeric Pagination Values",
			query:         "limit=abc&offset=xyz",
			expected:      &QueryParams{PaginateBy: &OffsetPagination{Limit: 0, Offset: 0, Count: false}},
			expectedError: false,
		},
		{
			name:          "Case Insensitivity",
			query:         `filter=Region EQ "na"`,
			expected:      &QueryParams{FilterBy: []Filter{{Field: "Region", Op: Eq, Value: "na"}}, PaginateBy: defaultPagination},
			expectedError: false,
		},
		{
			name:  "Whitespace Handling",
			query: `filter=  field   eq   "value"  `,
			expected: &QueryParams{
				FilterBy:   []Filter{{Field: "field", Op: Eq, Value: "value"}},
				PaginateBy: defaultPagination,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			result, err := ParseQueryParams(req)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
