package validator

import "github.com/go-playground/validator/v10"

var v = validator.New()

// FieldError holds a single field's validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate checks struct fields against their `validate` tags.
// Returns nil if valid, or a slice of FieldErrors describing what failed.
//
//	type CreateUserInput struct {
//	    Name  string `validate:"required,min=2"`
//	    Email string `validate:"required,email"`
//	    Age   int    `validate:"min=18"`
//	}
//
//	if errs := validator.Validate(input); errs != nil {
//	    return c.Status(422).JSON(response.Error("VALIDATION_FAILED", errs))
//	}
func Validate(s any) []FieldError {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var out []FieldError
	for _, fe := range err.(validator.ValidationErrors) {
		out = append(out, FieldError{
			Field:   fe.Field(),
			Message: msgFor(fe),
		})
	}
	return out
}

// RegisterTag registers a custom validation tag.
//
//	validator.RegisterTag("lao_phone", func(fl validator.FieldLevel) bool {
//	    return strings.HasPrefix(fl.Field().String(), "020")
//	})
func RegisterTag(tag string, fn validator.Func) error {
	return v.RegisterValidation(tag, fn)
}

func msgFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email"
	case "min":
		return fe.Field() + " is too short (min " + fe.Param() + ")"
	case "max":
		return fe.Field() + " is too long (max " + fe.Param() + ")"
	case "len":
		return fe.Field() + " must be exactly " + fe.Param() + " characters"
	case "numeric":
		return fe.Field() + " must be a number"
	case "url":
		return fe.Field() + " must be a valid URL"
	case "oneof":
		return fe.Field() + " must be one of: " + fe.Param()
	default:
		return fe.Field() + " failed validation (" + fe.Tag() + ")"
	}
}
