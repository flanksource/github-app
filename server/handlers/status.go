// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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


