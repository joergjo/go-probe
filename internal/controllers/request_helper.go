package controllers

import (
	"net/http"

	"github.com/joergjo/go-probe/internal/views"
)

func isHtmx(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func requiredField(r *http.Request, field string) (string, views.ValidationError) {
	val := r.PostFormValue(field)
	if val == "" {
		return val, views.RequiredFieldError{Field: field}
	}
	return val, nil
}
