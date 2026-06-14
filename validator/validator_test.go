package validator_test

import (
	"testing"

	"github.com/BounkhongDev/bkgo/validator"
)

type createInput struct {
	Name  string `validate:"required,min=2"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18"`
}

func TestValidate_Valid(t *testing.T) {
	input := createInput{Name: "Bounkhong", Email: "bk@example.com", Age: 25}
	if errs := validator.Validate(input); errs != nil {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	input := createInput{Email: "bk@example.com", Age: 25} // Name missing
	errs := validator.Validate(input)
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	found := false
	for _, e := range errs {
		if e.Field == "Name" {
			found = true
		}
	}
	if !found {
		t.Error("expected error on Name field")
	}
}

func TestValidate_InvalidEmail(t *testing.T) {
	input := createInput{Name: "Boun", Email: "not-an-email", Age: 25}
	errs := validator.Validate(input)
	if errs == nil {
		t.Fatal("expected validation errors")
	}
	found := false
	for _, e := range errs {
		if e.Field == "Email" {
			found = true
		}
	}
	if !found {
		t.Error("expected error on Email field")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	input := createInput{} // all fields invalid
	errs := validator.Validate(input)
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(errs))
	}
}
