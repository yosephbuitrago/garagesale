package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/yosephbuitrago/garagesale/foundation/web"
	"go.opentelemetry.io/otel"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			tr := otel.GetTracerProvider().Tracer("business.mid.panic")
			_, span := tr.Start(ctx, "panic")
			defer span.End()
			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s  :    PANIC   :\n%s", v.TraceID, debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
