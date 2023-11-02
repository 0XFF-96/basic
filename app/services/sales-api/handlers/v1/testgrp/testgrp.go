package testgrp

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Create adds a new user to the system.
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	json.NewEncoder(w).Encode(status)
	statusCode := http.StatusOK

	h.Log.Infow("readiness",
		"statusCode", statusCode,
		"method", r.Method,
		"path", r.URL.Path,
		"remoteAddress", r.RemoteAddr,
	)

	return
}
