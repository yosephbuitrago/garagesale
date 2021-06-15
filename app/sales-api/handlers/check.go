package handlers

import (
	"context"
	"log"
	"math/rand"
	"net/http"

	"github.com/yosephbuitrago/garagesale/foundation/web"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		// return errors.New("untrusted error")
		// return web.NewRequestError(errors.New("bad request"), http.StatusBadRequest)
		panic("panic")
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
