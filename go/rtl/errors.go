package rtl

import (
	"fmt"
	"net/http"
	"strings"

	json "github.com/json-iterator/go"
)

type ResponseError struct {
	StatusCode int
	Err        error
}

func (e ResponseError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	Message          string                   `json:"message"`           // Error details
	Errors           *[]ValidationErrorDetail `json:"errors,omitempty"`  // Error detail array
	DocumentationUrl string                   `json:"documentation_url"` // Documentation link
}

type ValidationErrorDetail struct {
	Field            *string `json:"field,omitempty"`   // Field with error
	Code             *string `json:"code,omitempty"`    // Error code
	Message          *string `json:"message,omitempty"` // Error info message
	DocumentationUrl string  `json:"documentation_url"` // Documentation link
}

type Error struct {
	Message          string `json:"message"`           // Error details
	DocumentationUrl string `json:"documentation_url"` // Documentation link
}

func (e Error) Error() string {
	return e.Message
}

func (e ValidationError) Error() string {
	if e.Errors == nil {
		return e.Message
	}

	var errSlice []string
	for _, m := range *e.Errors {
		// need to check field and message are not nil
		errSlice = append(errSlice, fmt.Sprintf("error on %s field. %s", *m.Field, *m.Message))
	}
	return strings.Join(errSlice, ",")
}

func DeserializeBody(status int, body []byte) error {
	if len(body) == 0 {
		return fmt.Errorf("response error. status=%d. error parsing error body", status)
	}

	switch status {
	case http.StatusUnprocessableEntity:
		// Status 422 returns a json payload of type ValidationError
		var e ValidationError
		if err := json.Unmarshal(body, &e); err != nil {
			// don't love this
			return fmt.Errorf("error unmarshalling body with status: %d, body:%s, error:%s", status, body, err.Error())
		}
		return e

	default:
		// All other status codes return a json payload of type Error
		var e Error
		if err := json.Unmarshal(body, &e); err != nil {
			return fmt.Errorf("error unmarshalling body with status: %d, body:%s, error:%s", status, body, err.Error())
		}
		return e
	}
}
