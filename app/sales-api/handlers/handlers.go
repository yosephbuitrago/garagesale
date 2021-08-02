// Package handlers contains the full set of handler functions and routes
// supported by the web api.package handlers
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/yosephbuitrago/garagesale/business/auth"
	"github.com/yosephbuitrago/garagesale/business/data/product"
	"github.com/yosephbuitrago/garagesale/business/data/user"
	"github.com/yosephbuitrago/garagesale/business/mid"
	"github.com/yosephbuitrago/garagesale/foundation/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, a *auth.Auth) *web.App {

	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	cg := checkGroup{
		build: build,
		db:    db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveness", cg.liveness)

	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}

	app.Handle(http.MethodGet, "/v1/users", ug.query, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPost, "/v1/users", ug.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/users/:id", ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPut, "/v1/users/:id", ug.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/v1/users/:id", ug.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/users/token/:kid", ug.token)

	// Register products endpoints
	pg := productGroup{
		product: product.New(log, db),
	}

	app.Handle(http.MethodGet, "/v1/products", pg.create, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/v1/products", pg.create, mid.Authenticate(a))

	return app
}
