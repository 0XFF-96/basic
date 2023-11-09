package testgrp

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Create adds a new user to the system.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	return json.NewEncoder(w).Encode(status)
}
