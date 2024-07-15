package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joergjo/go-probe/internal/controllers"
	"github.com/joergjo/go-probe/internal/templates"
	"github.com/joergjo/go-probe/internal/views"
)

type config struct {
	addr    string
	csrfKey string
	devMode bool
}

func main() {
	cfg, err := newConfig()
	var opts slog.HandlerOptions
	if cfg.devMode {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}
	if err := run(cfg); err != nil {
		slog.Error("application exiting with error", "error", err)
		os.Exit(1)
	}
}

func newConfig() (config, error) {
	var cfg config
	cfg.addr = os.Getenv("GO_PROBE_SERVER_ADDR")
	if cfg.addr == "" {
		cfg.addr = ":8080"
	}
	cfg.csrfKey = os.Getenv("GO_PROBE_CSRF_KEY")
	if cfg.csrfKey == "" {
		return cfg, fmt.Errorf("GO_PROBE_CSRF_KEY not set")
	}
	cfg.devMode = os.Getenv("GO_PROBE_DEVELOPER_MODE") == "true"
	return cfg, nil
}

func newRouter(cfg config) *mux.Router {
	indexV := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "index.gohtml"))
	indexC := controllers.Page(indexV)

	dnsV := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "dns-form.gohtml", "dns.gohtml"))
	dnsC := controllers.DNS(dnsV)
	dnsTestC := controllers.TestDNS(dnsV)

	pgV := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "pg-form.gohtml", "pg.gohtml"))
	pgC := controllers.Postgres(pgV)
	pgTestC := controllers.TestPostgres(pgV)

	oaiV := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "openai-form.gohtml", "openai.gohtml"))
	oaiC := controllers.OpenAI(oaiV)
	oaiTestC := controllers.TestOpenAI(oaiV)

	mux := mux.NewRouter()
	mux.StrictSlash(true)
	mux.Use(csrf.Protect([]byte(cfg.csrfKey), csrf.Secure(false)))
	mux.HandleFunc("/", indexC).Methods(http.MethodGet)
	mux.HandleFunc("/dns", dnsC).Methods(http.MethodGet)
	mux.HandleFunc("/dns", dnsTestC).Methods(http.MethodPost)
	mux.HandleFunc("/postgres", pgC).Methods(http.MethodGet)
	mux.HandleFunc("/postgres", pgTestC).Methods(http.MethodPost)
	mux.HandleFunc("/openai", oaiC).Methods(http.MethodGet)
	mux.HandleFunc("/openai", oaiTestC).Methods(http.MethodPost)

	return mux
}

func run(cfg config) error {
	if cfg.devMode {
		slog.Warn("running in developer mode")
	}
	slog.Info("starting server", "address", cfg.addr)
	mux := newRouter(cfg)
	if err := http.ListenAndServe(cfg.addr, mux); err != http.ErrServerClosed {
		return err
	}
	return nil
}
