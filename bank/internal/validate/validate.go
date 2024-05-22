package validate

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ErrorResponse struct {
	FailedField string      `json:"failed_field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

func Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func ParseAndValidate(body []byte, out any) []ErrorResponse {
	err := json.Unmarshal(body, out)
	if err != nil {
		return []ErrorResponse{
			{
				FailedField: "body",
				Tag:         "json",
				Value:       err.Error(),
			},
		}
	}

	return Validate(out)
}
