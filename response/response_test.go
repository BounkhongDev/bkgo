package response_test

import (
	"encoding/json"
	"testing"

	"github.com/BounkhongDev/bkgo/response"
)

func TestSuccess(t *testing.T) {
	r := response.Success(map[string]string{"id": "1"})
	if !r.Success {
		t.Error("expected success=true")
	}
	if r.Message != "ok" {
		t.Errorf("message: want ok got %s", r.Message)
	}
	if r.Meta != nil {
		t.Error("expected meta=nil for non-paginated response")
	}
}

func TestPaginated(t *testing.T) {
	r := response.Paginated([]string{"a", "b"}, 2, 10, 100)
	if !r.Success {
		t.Error("expected success=true")
	}
	if r.Meta == nil {
		t.Fatal("expected meta to be set")
	}
	if r.Meta.Page != 2 || r.Meta.Limit != 10 || r.Meta.Total != 100 {
		t.Errorf("meta: got page=%d limit=%d total=%d", r.Meta.Page, r.Meta.Limit, r.Meta.Total)
	}
}

func TestError(t *testing.T) {
	r := response.Error("NOT_FOUND", "user not found")
	if r.Success {
		t.Error("expected success=false")
	}
	if r.ErrorCode != "NOT_FOUND" {
		t.Errorf("error code: want NOT_FOUND got %s", r.ErrorCode)
	}
	if r.Message != "user not found" {
		t.Errorf("message: want 'user not found' got %s", r.Message)
	}
}

func TestResponseSerializesCleanly(t *testing.T) {
	r := response.Success(map[string]int{"count": 3})
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if v, hasError := out["error"]; hasError && v != "" {
		t.Error("success response should not include 'error' field")
	}
	if _, hasMeta := out["meta"]; hasMeta {
		t.Error("non-paginated response should not include 'meta' field")
	}
}
