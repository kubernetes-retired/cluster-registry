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
	"errors"

	"github.com/spf13/pflag"

	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
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

// BuiltInAuthorizationOptions configures authorization in the cluster registry.
type BuiltInAuthorizationOptions struct {
	Standalone *StandaloneAuthorizationOptions
	Delegating *genericoptions.DelegatingAuthorizationOptions
}

// StandaloneAuthorizationOptions configure authorization in the cluster registry
// when it is run as a standalone API server.
type StandaloneAuthorizationOptions struct {
	AlwaysAllow bool
}

func NewBuiltInAuthorizationOptions() *BuiltInAuthorizationOptions {
	return &BuiltInAuthorizationOptions{
		Standalone: &StandaloneAuthorizationOptions{AlwaysAllow: true},
		Delegating: genericoptions.NewDelegatingAuthorizationOptions(),
	}
}

// Validate checks that the configuration is valid.
func (b *BuiltInAuthorizationOptions) Validate() []error {
	if b.Delegating != nil {
		return b.Delegating.Validate()
	}
	return []error{}
}

func (b *BuiltInAuthorizationOptions) AddFlags(fs *pflag.FlagSet) {
	if b.Delegating != nil {
		b.Delegating.AddFlags(fs)
	}
}

func (b *BuiltInAuthorizationOptions) ApplyTo(c *genericapiserver.Config) error {
	if b.Delegating != nil {
		return b.Delegating.ApplyTo(c)
	}
	return nil
}

func (b *BuiltInAuthorizationOptions) ToStandaloneAuthorizationConfig() StandaloneAuthorizerConfig {
	config := StandaloneAuthorizerConfig{}
	if b.Standalone != nil {
		config.AlwaysAllow = b.Standalone.AlwaysAllow
	}
	return config
}

func (b *BuiltInAuthorizationOptions) ToDelegatedAuthorizationConfig() (authorizerfactory.DelegatingAuthorizerConfig, error) {
	if b.Delegating != nil {
		return b.Delegating.ToAuthorizationConfig()
	}
	return authorizerfactory.DelegatingAuthorizerConfig{}, errors.New("delegating authorization config requested, but no delegating authorization options present")
}
