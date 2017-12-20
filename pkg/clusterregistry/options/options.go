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
	"time"

	"k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/storage/storagebackend"

	"github.com/spf13/pflag"
)

type OptionsGetter interface {
	GenericServerRunOptions() *genericoptions.ServerRunOptions
	SecureServing() *genericoptions.SecureServingOptions
	Validate() []error
	Audit() *genericoptions.AuditOptions
	Features() *genericoptions.FeatureOptions
	ApplyAuthentication(*server.Config) error
	ApplyAuthorization(c *server.Config) error
	Etcd() *genericoptions.EtcdOptions
}

// ServerRunOptions contains runtime options for the cluster registry.
type ServerRunOptions struct {
	genericServerRunOptions *genericoptions.ServerRunOptions
	etcd                    *genericoptions.EtcdOptions
	secureServing           *genericoptions.SecureServingOptions
	audit                   *genericoptions.AuditOptions
	features                *genericoptions.FeatureOptions
	//StandaloneAuthentication *StandaloneAuthenticationOptions
	//StandaloneAuthorization  *StandaloneAuthorizationOptions

	EventTTL time.Duration
}

// NewServerRunOptions creates a new ServerRunOptions object with default values.
func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{
		genericServerRunOptions: genericoptions.NewServerRunOptions(),
		etcd:          genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig("/registry/clusterregistry.kubernetes.io", nil)),
		secureServing: genericoptions.NewSecureServingOptions(),
		audit:         genericoptions.NewAuditOptions(),
		features:      genericoptions.NewFeatureOptions(),
		//StandaloneAuthentication: NewStandaloneAuthenticationOptions().WithAll(),
		//StandaloneAuthorization:  NewStandaloneAuthorizationOptions(),

		EventTTL: 1 * time.Hour,
	}
	return o
}

// AddFlags adds flags for ServerRunOptions fields to be specified via FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	s.genericServerRunOptions.AddUniversalFlags(fs)
	s.etcd.AddFlags(fs)
	s.secureServing.AddFlags(fs)
	s.audit.AddFlags(fs)
	s.features.AddFlags(fs)
	//s.StandaloneAuthentication.AddFlags(fs)
	//s.StandaloneAuthorization.AddFlags(fs)

	fs.DurationVar(&s.EventTTL, "event-ttl", s.EventTTL, "Amount of time to retain events.")
}
