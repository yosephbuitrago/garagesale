// Package handlers contains the full set of handler functions and routes
// supported by the web api.package handlers
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/yosephbuitrago/garagesale/business/auth"
	"github.com/yosephbuitrago/garagesale/business/mid"
	"github.com/yosephbuitrago/garagesale/foundation/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth) *web.App {

	tm := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	check := check{
		log: log,
	}

	tm.Handle(http.MethodGet, "/readiness", check.readiness, mid.Authenticate(a), mid.Authorize(auth.RoleUser))
	return tm
}
