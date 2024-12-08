package forms

import (
	"net/url"
)

type Form struct {
	url.Values
	ValidForm bool
}

func NewForm(data url.Values) *Form {
	return &Form{
		data,
		true,
	}
}

func (f *Form) Has(field string) bool {
	if f.Get(field) == "" {
		f.ValidForm = false
		return false
	} else {
		return true
	}
}

func (f *Form) IsValid() bool {
	return f.ValidForm
}
