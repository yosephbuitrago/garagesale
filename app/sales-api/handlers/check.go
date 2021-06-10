package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// func (c check) readiness(w http.ResponseWriter, r *http.Request) {

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)
	log.Println(r, status)
	return nil
}
