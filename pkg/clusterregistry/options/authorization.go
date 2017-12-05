/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"github.com/spf13/pflag"

	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

// StandaloneAuthorizerConfig configures an authorizer for the standalone
// cluster registry.
type StandaloneAuthorizerConfig struct {
	AlwaysAllow bool
}

func (s StandaloneAuthorizerConfig) New() (authorizer.Authorizer, error) {
	if s.AlwaysAllow {
		return authorizerfactory.NewAlwaysAllowAuthorizer(), nil
	}
	return authorizerfactory.NewAlwaysDenyAuthorizer(), nil
}

// StandaloneAuthorizationOptions configure authorization in the cluster registry
// when it is run as a standalone API server.
type StandaloneAuthorizationOptions struct {
	AlwaysAllow bool
}

func NewStandaloneAuthorizationOptions() *StandaloneAuthorizationOptions {
	return &StandaloneAuthorizationOptions{AlwaysAllow: true}
}

// Validate checks that the configuration is valid.
func (s *StandaloneAuthorizationOptions) Validate() []error {
	return []error{}
}

func (s *StandaloneAuthorizationOptions) AddFlags(fs *pflag.FlagSet) {
	// TODO: Add a flag to set AlwaysAllow.
}

func (s *StandaloneAuthorizationOptions) ApplyTo(c *genericapiserver.Config) error {
	authorizer, err := s.toAuthorizationConfig().New()
	if err != nil {
		return err
	}
	c.Authorizer = authorizer
	return nil
}

func (s *StandaloneAuthorizationOptions) toAuthorizationConfig() StandaloneAuthorizerConfig {
	return StandaloneAuthorizerConfig{AlwaysAllow: s.AlwaysAllow}
}
