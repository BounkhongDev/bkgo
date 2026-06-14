package errs_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bounkhongdev/kbgo/errs"
)

func TestNew(t *testing.T) {
	e := errs.New(http.StatusBadRequest, "BAD_REQUEST", "something went wrong")
	if e.Status != http.StatusBadRequest {
		t.Errorf("status: want %d got %d", http.StatusBadRequest, e.Status)
	}
	if e.Code != "BAD_REQUEST" {
		t.Errorf("code: want BAD_REQUEST got %s", e.Code)
	}
	if e.Error() != "something went wrong" {
		t.Errorf("message: want 'something went wrong' got %s", e.Error())
	}
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *errs.AppError
		wantStatus int
		wantCode   string
	}{
		{"NotFound", errs.NotFound("not found"), http.StatusNotFound, "NOT_FOUND"},
		{"BadRequest", errs.BadRequest("bad"), http.StatusBadRequest, "BAD_REQUEST"},
		{"Unauthorized", errs.Unauthorized("unauth"), http.StatusUnauthorized, "UNAUTHORIZED"},
		{"Forbidden", errs.Forbidden("forbidden"), http.StatusForbidden, "FORBIDDEN"},
		{"Conflict", errs.Conflict("conflict"), http.StatusConflict, "CONFLICT"},
		{"Internal", errs.Internal("oops"), http.StatusInternalServerError, "INTERNAL_ERROR"},
		{"Unprocessable", errs.Unprocessable("invalid"), http.StatusUnprocessableEntity, "UNPROCESSABLE"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Status != tc.wantStatus {
				t.Errorf("status: want %d got %d", tc.wantStatus, tc.err.Status)
			}
			if tc.err.Code != tc.wantCode {
				t.Errorf("code: want %s got %s", tc.wantCode, tc.err.Code)
			}
		})
	}
}

func TestIsAppError(t *testing.T) {
	t.Run("app error returns true", func(t *testing.T) {
		err := errs.NotFound("missing")
		ae, ok := errs.IsAppError(err)
		if !ok {
			t.Fatal("expected ok=true")
		}
		if ae.Code != "NOT_FOUND" {
			t.Errorf("code: want NOT_FOUND got %s", ae.Code)
		}
	})

	t.Run("non-app error returns false", func(t *testing.T) {
		_, ok := errs.IsAppError(fmt.Errorf("plain error"))
		if ok {
			t.Fatal("expected ok=false for plain error")
		}
	})
}
