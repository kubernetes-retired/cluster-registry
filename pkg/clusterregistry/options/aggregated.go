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

// Package options contains flags and options for initializing the cluster registry API server.
package options

import (
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
)

// AggregatedServerRunOptions contains runtime options for the cluster registry
// deployed using the Kubernetes aggregator.
type AggregatedServerRunOptions struct {
	*ServerRunOptions
	authentication *options.DelegatingAuthenticationOptions
	authorization  *options.DelegatingAuthorizationOptions
}

// NewAggregatedServerRunOptions creates a new AggregatedServerRunOptions
// object with default values.
func NewAggregatedServerRunOptions() *AggregatedServerRunOptions {
	return &AggregatedServerRunOptions{
		ServerRunOptions: NewServerRunOptions(),
		authentication:   options.NewDelegatingAuthenticationOptions(),
		authorization:    options.NewDelegatingAuthorizationOptions(),
	}
}

// AddFlags adds flags for specific AggregatedServerRunOptions fields specified
// via FlagSet before calling the embedded AddFlags method to add the rest of
// the common options.
func (s *AggregatedServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	s.authentication.AddFlags(fs)
	s.authorization.AddFlags(fs)
	s.ServerRunOptions.AddFlags(fs)
}

// GenericServerRunOptions gets the embedded genericServerRunOptions field.
func (s *AggregatedServerRunOptions) GenericServerRunOptions() *options.ServerRunOptions {
	return s.genericServerRunOptions
}

// SecureServing gets the embedded secureServing field.
func (s *AggregatedServerRunOptions) SecureServing() *options.SecureServingOptions {
	return s.secureServing
}

// Validate checks any specific AggregatedServerRunOptions before calling the
// embedded Validate method for the common options.
func (s *AggregatedServerRunOptions) Validate() []error {
	return s.ServerRunOptions.Validate()
}

// Audit gets the embedded audit field.
func (s *AggregatedServerRunOptions) Audit() *options.AuditOptions {
	return s.audit
}

// Features gets the embedded features field.
func (s *AggregatedServerRunOptions) Features() *options.FeatureOptions {
	return s.features
}

// ApplyAuthentication applies the delegated authentication to the config.
func (s *AggregatedServerRunOptions) ApplyAuthentication(c *server.Config) error {
	return s.authentication.ApplyTo(c)
}

// ApplyAuthorization applies the delegated authorization to the config.
func (s *AggregatedServerRunOptions) ApplyAuthorization(c *server.Config) error {
	return s.authorization.ApplyTo(c)
}

// Etcd gets the embedded etcd field.
func (s *AggregatedServerRunOptions) Etcd() *options.EtcdOptions {
	return s.etcd
}
