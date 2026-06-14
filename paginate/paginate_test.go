package paginate_test

import (
	"testing"

	"github.com/bounkhongdev/kbgo/paginate"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name      string
		input     paginate.Params
		wantPage  int
		wantLimit int
	}{
		{"valid params unchanged", paginate.Params{Page: 2, Limit: 10}, 2, 10},
		{"page 0 becomes 1", paginate.Params{Page: 0, Limit: 10}, 1, 10},
		{"negative page becomes 1", paginate.Params{Page: -5, Limit: 10}, 1, 10},
		{"limit 0 becomes default", paginate.Params{Page: 1, Limit: 0}, 1, paginate.DefaultLimit},
		{"limit over max becomes default", paginate.Params{Page: 1, Limit: 999}, 1, paginate.DefaultLimit},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Normalize()
			if tc.input.Page != tc.wantPage {
				t.Errorf("page: want %d got %d", tc.wantPage, tc.input.Page)
			}
			if tc.input.Limit != tc.wantLimit {
				t.Errorf("limit: want %d got %d", tc.wantLimit, tc.input.Limit)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		page       int
		limit      int
		wantOffset int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 10, 20},
		{5, 25, 100},
	}

	for _, tc := range tests {
		p := paginate.Params{Page: tc.page, Limit: tc.limit}
		got := p.Offset()
		if got != tc.wantOffset {
			t.Errorf("page=%d limit=%d: want offset %d got %d", tc.page, tc.limit, tc.wantOffset, got)
		}
	}
}
