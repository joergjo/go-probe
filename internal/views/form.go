package views

import (
	"fmt"
)

type SelectListItem struct {
	Text     string
	Value    string
	Disabled bool
	Selected bool
}

type ValidationError interface {
	Element() string
}

type RequiredFieldError struct {
	Field string
}

func (e RequiredFieldError) Error() string {
	return fmt.Sprintf("The %s field is required.", e.Field)
}

func (e RequiredFieldError) Element() string {
	return e.Field
}

type InvalidFieldError struct {
	Field string
	Value string
}

func (e InvalidFieldError) Error() string {
	return fmt.Sprintf("The value '%s' is not valid for %s.", e.Value, e.Field)
}

func (e InvalidFieldError) Element() string {
	return e.Field
}

type FormErrors map[string]ValidationError
