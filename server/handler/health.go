/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"fmt"
	"github.com/flanksource/github-app/version"
	"net/http"
)

// HealthCheck is a handler for the API endpoint for healthchecks
type HealthCheck struct {
}

func (h *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, `{"status": "ok", "version": "%s"}`, version.GetVersion())
	return nil
}
