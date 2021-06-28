// Package handlers contains the full set of handler functions and routes
// supported by the web api.package handlers
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/yosephbuitrago/garagesale/business/auth"
	"github.com/yosephbuitrago/garagesale/business/mid"
	"github.com/yosephbuitrago/garagesale/foundation/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, a *auth.Auth) *web.App {

	tm := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	check := checkGroup{
		build: build,
		db:    db,
	}

	tm.Handle(http.MethodGet, "/readiness", check.readiness)
	tm.Handle(http.MethodGet, "/liveness", check.liveness)
	return tm
}
