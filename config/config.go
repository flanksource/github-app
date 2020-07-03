// Copyright 2018 Palantir Technologies, Inc.
// Modifications copyright (C) 2020 Flanksource

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

package config

import (
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"gopkg.in/flanksource/yaml.v3"
	"io/ioutil"
)

func ReadConfig(path string) (*Config, error) {
	var c Config

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return nil, errors.Wrap(err, "failed parsing configuration file")
	}

	return &c, nil
}


type (
	Config struct {
	// Server holds basic server configs
	Server baseapp.HTTPConfig `yaml:"server"`
	// Github holds github app related configs
	Github githubapp.Config `yaml:"github"`
	// Secrets configures required secrets not related to github apps
	Secrets  SecretsConfig  `yaml:"secrets"`
	Runners  RunnersConfig  `yaml:"runners"`
	Logging  LoggingConfig  `yaml:"logging"`
	Sessions SessionsConfig `yaml:"sessions"`
	Workers  WorkerConfig   `yaml:"workers"`
	// Auth allows the internal OAuth 2.0 auth server to be configured
	Auth AuthConfig         `yaml:"auth"`
}
)

type (
	SecretsConfig struct {
		GhPat string `yaml:"gh-pat" json:"gh-pat"`
	}


	RunnersConfig struct {
		Owner string `yaml:"owner" json:"owner"`
		Repo  string `yaml:"repo" json:"repo"`
	}

	LoggingConfig struct {
		Level string `yaml:"level" json:"level"`
		Text  bool   `yaml:"text" json:"text"`
	}

	WorkerConfig struct {
		Workers   int `yaml:"workers"`
		QueueSize int `yaml:"queue_size"`
	}

	SessionsConfig struct {
		Key      string `yaml:"key"`
		Lifetime string `yaml:"lifetime"`
	}

	AuthConfig struct {
		Clients []ClientSpec  `yaml:"clients"`
		// A shared secret as symmetric key
		SymmetricKey    string `yaml:"key"`
	}
)

// ClientSpec allows oauth2 clients
// (think system users) to be specified
type ClientSpec struct {
	// ID Identifies the client and serves as user for basic auth
	ID string		`yaml:"id"`
	// Secret authenticates the client and serves as password for basic auth
	Secret string	`yaml:"secret"`
	// Domain is a label for a "range" over which access is granted
	Domain string	`yaml:"domain"`
	// UserID serves as a friendly name for the client (extra and optional)
	UserID string	`yaml:"user"`
}

// GetClient is a conventience method to retrieve a
// github.com/go-oauth2/oauth2/v4/models Client
func (cs *ClientSpec) GetClient() *models.Client {
	return &models.Client{
		ID:     cs.ID,
		Secret: cs.Secret,
		Domain: cs.Domain,
		UserID: cs.UserID,
	}
}



