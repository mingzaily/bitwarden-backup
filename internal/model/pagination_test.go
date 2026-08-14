package model

import "testing"

func TestPaginationClampsLimitBeforeOffset(t *testing.T) {
	params := PaginationParams{Page: 1000000, PageSize: 1 << 30}
	if got := params.GetOffset(); got != 99999900 {
		t.Fatalf("GetOffset() = %d, want 99999900", got)
	}
	if got := params.GetLimit(); got != 100 {
		t.Fatalf("GetLimit() = %d, want 100", got)
	}
}
