package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

//	func (f *Form) Has(field string, r *http.Request) bool {
//		fv := r.Form.Get(field)
//		if fv == "" {
//			f.Errors.Add(field, "This field is Required")
//			return false
//		}
//		return true
//	}
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This Field is required")
		}
	}
}
func (f *Form) MinLength(field string, length int) bool {
	value := f.Get(field)
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("Must be at least %d characters long", length))
		return false
	}
	return true

}
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Enter a Valid Email address")
	}
}
