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
	"github.com/palantir/go-githubapp/githubapp"
	"io/ioutil"

	"github.com/palantir/go-baseapp/baseapp"

	"github.com/pkg/errors"
	"gopkg.in/flanksource/yaml.v3"
)

type Config struct {
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
}

type SecretsConfig struct {
	GhPat string `yaml:"gh-pat" json:"gh-pat"`
}

type RunnersConfig struct {
	Owner string `yaml:"owner" json:"owner"`
	Repo  string `yaml:"repo" json:"repo"`
}

type LoggingConfig struct {
	Level string `yaml:"level" json:"level"`
	Text  bool   `yaml:"text" json:"text"`
}

type WorkerConfig struct {
	Workers   int `yaml:"workers"`
	QueueSize int `yaml:"queue_size"`
}

type SessionsConfig struct {
	Key      string `yaml:"key"`
	Lifetime string `yaml:"lifetime"`
}

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
