package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Page(w http.ResponseWriter, r *http.Request, model any) {
	t.Partial(w, r, "", model)
}

func (t Template) Partial(w http.ResponseWriter, r *http.Request, name string, model any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	var err error
	if name == "" {
		err = t.htmlTpl.Execute(&buf, model)
	} else {
		err = t.htmlTpl.ExecuteTemplate(&buf, name, model)
	}
	if err != nil {
		slog.Error("executing template", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (t Template) ParseClientTemplate(fs fs.FS, name, left, right string) (Template, error) {
	_, err := t.htmlTpl.Delims(left, right).ParseFS(fs, name)
	if err != nil {
		return Template{}, fmt.Errorf("parsing client template: %w", err)
	}
	return t, nil
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(filepath.Base(patterns[0])).Funcs(template.FuncMap{
		"ticks": func(i int) int64 {
			return time.Now().UnixMilli() + int64(i)
		},
		"timeString": func() string {
			return time.Now().Format("03:04:05")
		},
		"isSelectedCSS": func(s, t, css string) string {
			if s != t {
				return ""
			}
			return css
		},
	})
	_, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{htmlTpl: tpl}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
