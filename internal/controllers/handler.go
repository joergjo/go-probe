package controllers

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/joergjo/go-probe/internal/probes"
	"github.com/joergjo/go-probe/internal/views"
)

type View interface {
	Page(w http.ResponseWriter, r *http.Request, data any)
	Partial(w http.ResponseWriter, r *http.Request, name string, data any)
}

func Page(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v.Page(w, r, nil)
	}
}

func File(fs fs.FS, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fs, path)
	}
}

func DNS(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		model := dnsModel{
			ModelErrors: views.FormErrors{},
			CsrfField:   csrfField,
		}
		v.Page(w, r, model)
	}
}

func TestDNS(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		modelErrs := views.FormErrors{}
		typ, err := requiredField(r, "type")
		if err != nil {
			modelErrs["type"] = err
		}

		val, err := requiredField(r, "value")
		if err != nil {
			modelErrs["value"] = err
		}
		model := dnsModel{
			Type:        typ,
			Value:       val,
			ModelErrors: modelErrs,
			CsrfField:   csrfField,
		}

		if len(modelErrs) == 0 {
			res, err := probes.DNS(typ, val)
			if err != nil {
				model.Result = err.Error()
			} else {
				model.Result = res
			}
		}

		if isHtmx(r) {
			v.Partial(w, r, "dns-form", model)
			return
		}
		v.Page(w, r, model)
	}
}

func Postgres(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		model := postgresModel{
			UseTLS:      true,
			ModelErrors: views.FormErrors{},
			CsrfField:   csrfField,
		}
		v.Page(w, r, model)
	}
}

func TestPostgres(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		modelErrs := views.FormErrors{}
		fqdn, err := requiredField(r, "fqdn")
		if err != nil {
			modelErrs["fqdn"] = err
		}

		login, err := requiredField(r, "login")
		if err != nil {
			modelErrs["login"] = err
		}

		password, err := requiredField(r, "password")
		if err != nil {
			modelErrs["password"] = err
		}

		useTLS := r.FormValue("useTLS") == "on"

		model := postgresModel{
			FQDN:        fqdn,
			Login:       login,
			Password:    password,
			UseTLS:      useTLS,
			ModelErrors: modelErrs,
			CsrfField:   csrfField,
		}

		if len(modelErrs) == 0 {
			res, err := probes.Postgres(fqdn, login, password, useTLS)
			if err != nil {
				model.Result = err.Error()
			} else {
				model.Result = res
			}
		}

		if isHtmx(r) {
			v.Partial(w, r, "pg-form", model)
			return
		}
		v.Page(w, r, model)
	}
}

func OpenAI(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		model := openAIModel{
			ModelErrors: views.FormErrors{},
			CsrfField:   csrfField,
		}
		v.Page(w, r, model)
	}
}

func TestOpenAI(v View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfField := csrf.TemplateField(r)
		modelErrs := views.FormErrors{}

		endpoint, err := requiredField(r, "endpoint")
		if err != nil {
			modelErrs["endpoint"] = err
		}

		key, err := requiredField(r, "key")
		if err != nil {
			modelErrs["key"] = err
		}

		deployment, err := requiredField(r, "deployment")
		if err != nil {
			modelErrs["deployment"] = err
		}

		prompt, err := requiredField(r, "prompt")
		if err != nil {
			modelErrs["prompt"] = err
		}

		model := openAIModel{
			Endpoint:    endpoint,
			Key:         key,
			Deployment:  deployment,
			Prompt:      prompt,
			ModelErrors: modelErrs,
			CsrfField:   csrfField,
		}

		if len(modelErrs) == 0 {
			res, err := probes.OpenAI(prompt, endpoint, key, deployment)
			if err != nil {
				model.Result = err.Error()
			} else {
				model.Result = res
			}
		}

		if isHtmx(r) {
			v.Partial(w, r, "openai-form", model)
			return
		}
		v.Page(w, r, model)
	}
}
