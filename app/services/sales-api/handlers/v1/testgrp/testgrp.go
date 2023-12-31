package testgrp

import (
	"context"
	"github.com/yourusername/basic-a/foundation/web"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Create adds a new user to the system.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	randn := rand.Intn(100)
	if randn%2 == 0 {
		panic("panic error")
	}

	if randn%5 == 0 {
		return web.NewShutdownError("ShutdownError error")
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	statusCode := http.StatusOK

	h.Log.Infow("readiness",
		"statusCode", statusCode,
		"method", r.Method,
		"path", r.URL.Path,
		"remoteAddress", r.RemoteAddr,
	)

	return web.Respond(ctx, w, status, http.StatusOK)
}
