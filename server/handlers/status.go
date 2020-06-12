/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v30/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

type StatusHandler struct {
	githubapp.ClientCreator

	preamble string
}

func (h *StatusHandler) Handles() []string {
	return []string{"status"}
}

func (h *StatusHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.StatusEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	return nil
}


