/*
Copyright 2016 The Kubernetes Authors.

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

	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/storage/storagebackend"

	"github.com/spf13/pflag"
)

// Runtime options for the clusterregistry-apiserver.
type ServerRunOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
	Etcd                    *genericoptions.EtcdOptions
	SecureServing           *genericoptions.SecureServingOptions
	Audit                   *genericoptions.AuditOptions
	Features                *genericoptions.FeatureOptions
	Admission               *genericoptions.AdmissionOptions
	Authentication          *BuiltInAuthenticationOptions

	EventTTL time.Duration
}

// NewServerRunOptions creates a new ServerRunOptions object with default values.
func NewServerRunOptions() *ServerRunOptions {
	o := &ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		Etcd:           genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig("/registry/clusterregistry.kubernetes.io", nil)),
		SecureServing:  genericoptions.NewSecureServingOptions(),
		Audit:          genericoptions.NewAuditOptions(),
		Features:       genericoptions.NewFeatureOptions(),
		Admission:      genericoptions.NewAdmissionOptions(),
		Authentication: NewBuiltInAuthenticationOptions().WithAll(),

		EventTTL: 1 * time.Hour,
	}
	o.Authentication.Anonymous.Allow = false
	return o
}

// AddFlags adds flags for ServerRunOptions fields to be specified via FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	s.GenericServerRunOptions.AddUniversalFlags(fs)
	s.Etcd.AddFlags(fs)
	s.SecureServing.AddFlags(fs)
	s.Audit.AddFlags(fs)
	s.Features.AddFlags(fs)
	s.Authentication.AddFlags(fs)
	s.Admission.AddFlags(fs)

	fs.DurationVar(&s.EventTTL, "event-ttl", s.EventTTL, "Amount of time to retain events.")
}
