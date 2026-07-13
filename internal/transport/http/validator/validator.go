package validator

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateBody[T any](bodyReader io.Reader) (*T, map[string]string) {
	var body T

	decoder := json.NewDecoder(bodyReader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return nil, handleBindingError(err)
	}

	if err := validate.Struct(&body); err != nil {
		if valErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return nil, formatValidationErrors(valErr)
		}
		return nil, map[string]string{"errors": err.Error()}
	}

	return &body, nil
}

func handleBindingError(err error) map[string]string {
	if typeErr, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		return map[string]string{
			typeErr.Field: "invalid data type, expected " + typeErr.Type.String(),
		}
	}

	if syntaxErr, ok := errors.AsType[*json.SyntaxError](err); ok {
		return map[string]string{
			"error": "JSON syntax error at byte " + strconv.FormatInt(syntaxErr.Offset, 10),
		}
	}

	if errors.Is(err, io.EOF) {
		return map[string]string{"error": "request body is empty"}
	}

	if strings.HasPrefix(err.Error(), "json: unknown field") {
		return map[string]string{"error": err.Error()}
	}

	return map[string]string{"error": "failed to parse payload: " + err.Error()}
}

func formatValidationErrors(ve validator.ValidationErrors) map[string]string {
	errs := make(map[string]string, len(ve))
	for _, fe := range ve {
		errs[fe.Field()] = fieldErrorMessage(fe)
	}
	return errs
}

func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "len":
		return "must be exactly " + fe.Param() + " characters"
	case "oneof":
		return "must be one of: " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	case "lt":
		return "must be less than " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	default:
		return "failed validation: " + fe.Tag()
	}
}
