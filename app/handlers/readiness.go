package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/mistermoe/go-starter/framework"
)

type readiness struct {
	log *log.Logger
}

func (_ readiness) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return framework.Respond(ctx, w, status, http.StatusOK)
}
