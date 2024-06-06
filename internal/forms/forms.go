package forms

import (
	"net/http"
	"net/url"
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
func (f *Form) Has(field string, r *http.Request) bool {
	fv := r.Form.Get(field)
	if fv == "" {
		f.Errors.Add(field, "This field is Required")
		return false
	}
	return true
}
