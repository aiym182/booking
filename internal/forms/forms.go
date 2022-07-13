package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

// It creates a custom form struct and it embeds a url.Values object

type Form struct {
	url.Values
	Errors   errors
	Validate *validator.Validate
}

// New initializes a form struct
func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
		validator.New(),
	}
}

// Has Checks if required field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	return x != ""
}

// Valid returns true if there are no errors , otherwise return false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Rquired checks for required fields
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Checks for minimum length
func (f *Form) MinLen(field string, length int) bool {

	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be longer than %d characters", length))
		return false
	}
	return true
}

func (f *Form) ValidatEmail(field string) {
	x := f.Get(field)

	err := f.Validate.Var(x, "email")
	if err != nil {
		f.Errors.Add(field, "This field requires a valid email adress")
	}

}
