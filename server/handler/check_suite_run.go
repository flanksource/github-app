/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package handler

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v32/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type CheckSuiteHandler struct {
	githubapp.ClientCreator

	preamble string
}

func (h *CheckSuiteHandler) Handles() []string {
	return []string{"check_suite"}
}

func (h *CheckSuiteHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.CheckSuiteEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	zerolog.Ctx(ctx).Printf("%v, %v", event.GetAction(), event.GetCheckSuite().GetApp().GetName())

	return nil
}
