package probes

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

func Postgres(fqdn, login, password string, useTLS bool) (string, error) {
	sslMode := "require"
	if !useTLS {
		sslMode = "disable"
	}
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres sslmode=%s",
		fqdn,
		login,
		password,
		sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("initializing database driver", "error", err)
		return "", err
	}
	slog.Debug("connected to database", "fqdn", fqdn, "login", login, "use_tls", useTLS)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return "", err
	}
	return "Successfully connected", nil
}
