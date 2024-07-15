package controllers

import (
	"html/template"

	"github.com/joergjo/go-probe/internal/views"
)

type dnsModel struct {
	Type        string
	Value       string
	Result      string
	ModelErrors views.FormErrors
	CsrfField   template.HTML
}

type postgresModel struct {
	FQDN        string
	Login       string
	Password    string
	UseTLS      bool
	Result      string
	ModelErrors views.FormErrors
	CsrfField   template.HTML
}

type openAIModel struct {
	Prompt      string
	Endpoint    string
	Key         string
	Deployment  string
	Result      string
	ModelErrors views.FormErrors
	CsrfField   template.HTML
}
