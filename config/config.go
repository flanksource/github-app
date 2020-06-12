/*
Copyright © 2020 Flanksource
This file is part of Flanksource github-app
*/
package config

import (
	"io/ioutil"

	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"gopkg.in/flanksource/yaml.v3"
)

type Config struct {
	Server baseapp.HTTPConfig `yaml:"server"`
	Github githubapp.Config   `yaml:"github"`

	AppConfig ApplicationConfig `yaml:"app_configuration"`
}

type ApplicationConfig struct {
	SampleKey string `yaml:"sampleKey"`
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