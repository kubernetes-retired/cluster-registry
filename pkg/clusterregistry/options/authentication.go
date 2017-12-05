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

package options

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/cluster-registry/pkg/clusterregistry/authenticator"
)

type StandaloneAuthenticationOptions struct {
	Anonymous      *AnonymousAuthenticationOptions
	BootstrapToken *BootstrapTokenAuthenticationOptions
	ClientCert     *genericoptions.ClientCertAuthenticationOptions
	PasswordFile   *PasswordFileAuthenticationOptions
	TokenFile      *TokenFileAuthenticationOptions
	WebHook        *WebHookAuthenticationOptions

	TokenSuccessCacheTTL time.Duration
	TokenFailureCacheTTL time.Duration
}

type AnonymousAuthenticationOptions struct {
	Allow bool
}

type BootstrapTokenAuthenticationOptions struct {
	Enable bool
}

type PasswordFileAuthenticationOptions struct {
	BasicAuthFile string
}

type TokenFileAuthenticationOptions struct {
	TokenFile string
}

type WebHookAuthenticationOptions struct {
	ConfigFile string
	CacheTTL   time.Duration
}

func NewStandaloneAuthenticationOptions() *StandaloneAuthenticationOptions {
	return &StandaloneAuthenticationOptions{
		TokenSuccessCacheTTL: 10 * time.Second,
		TokenFailureCacheTTL: 0 * time.Second,
	}
}

func (s *StandaloneAuthenticationOptions) WithAll() *StandaloneAuthenticationOptions {
	return s.
		WithAnonymous().
		WithBootstrapToken().
		WithClientCert().
		WithPasswordFile().
		WithTokenFile().
		WithWebHook()
}

func (s *StandaloneAuthenticationOptions) WithAnonymous() *StandaloneAuthenticationOptions {
	s.Anonymous = &AnonymousAuthenticationOptions{Allow: true}
	return s
}

func (s *StandaloneAuthenticationOptions) WithBootstrapToken() *StandaloneAuthenticationOptions {
	s.BootstrapToken = &BootstrapTokenAuthenticationOptions{}
	return s
}

func (s *StandaloneAuthenticationOptions) WithClientCert() *StandaloneAuthenticationOptions {
	s.ClientCert = &genericoptions.ClientCertAuthenticationOptions{}
	return s
}

func (s *StandaloneAuthenticationOptions) WithPasswordFile() *StandaloneAuthenticationOptions {
	s.PasswordFile = &PasswordFileAuthenticationOptions{}
	return s
}

func (s *StandaloneAuthenticationOptions) WithTokenFile() *StandaloneAuthenticationOptions {
	s.TokenFile = &TokenFileAuthenticationOptions{}
	return s
}

func (s *StandaloneAuthenticationOptions) WithWebHook() *StandaloneAuthenticationOptions {
	s.WebHook = &WebHookAuthenticationOptions{
		CacheTTL: 2 * time.Minute,
	}
	return s
}

// Validate checks invalid config combination
func (s *StandaloneAuthenticationOptions) Validate() []error {
	return []error{}
}

func (s *StandaloneAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	if s.Anonymous != nil {
		fs.BoolVar(&s.Anonymous.Allow, "anonymous-auth", s.Anonymous.Allow, ""+
			"Enables anonymous requests to the secure port of the API server. "+
			"Requests that are not rejected by another authentication method are treated as anonymous requests. "+
			"Anonymous requests have a username of system:anonymous, and a group name of system:unauthenticated.")
	}

	if s.BootstrapToken != nil {
		fs.BoolVar(&s.BootstrapToken.Enable, "experimental-bootstrap-token-auth", s.BootstrapToken.Enable, ""+
			"Deprecated (use --enable-bootstrap-token-auth).")
		fs.MarkDeprecated("experimental-bootstrap-token-auth", "use --enable-bootstrap-token-auth instead.")

		fs.BoolVar(&s.BootstrapToken.Enable, "enable-bootstrap-token-auth", s.BootstrapToken.Enable, ""+
			"Enable to allow secrets of type 'bootstrap.kubernetes.io/token' in the 'kube-system' "+
			"namespace to be used for TLS bootstrapping authentication.")
	}

	// genericoptions.DelegatingAuthenticationOptions adds flags for its own copy
	// of the ClientCert and WebHook authentication options.
	//
	// TODO: Determine if there is a better way to split the personalities of the
	// delegated and standalone modes. The wart here is that these flags must be
	// added before the flag that tells the cluster registry to use delegated
	// auth is parsed.

	if s.PasswordFile != nil {
		fs.StringVar(&s.PasswordFile.BasicAuthFile, "basic-auth-file", s.PasswordFile.BasicAuthFile, ""+
			"If set, the file that will be used to admit requests to the secure port of the API server "+
			"via http basic authentication.")
	}

	if s.TokenFile != nil {
		fs.StringVar(&s.TokenFile.TokenFile, "token-auth-file", s.TokenFile.TokenFile, ""+
			"If set, the file that will be used to secure the secure port of the API server "+
			"via token authentication.")
	}

}

func (o *StandaloneAuthenticationOptions) ApplyTo(c *genericapiserver.Config) error {
	var err error
	if o.ClientCert != nil {
		c, err = c.ApplyClientCert(o.ClientCert.ClientCA)
		if err != nil {
			return fmt.Errorf("unable to load client CA file: %v", err)
		}
	}

	c.SupportsBasicAuth = o.PasswordFile != nil && len(o.PasswordFile.BasicAuthFile) > 0

	authenticator, securityDefinitions, err := o.toAuthenticationConfig().New()
	if err != nil {
		return err
	}
	c.Authenticator = authenticator
	if c.OpenAPIConfig != nil {
		c.OpenAPIConfig.SecurityDefinitions = securityDefinitions
	}
	return nil
}

func (s *StandaloneAuthenticationOptions) toAuthenticationConfig() authenticator.AuthenticatorConfig {
	ret := authenticator.AuthenticatorConfig{
		TokenSuccessCacheTTL: s.TokenSuccessCacheTTL,
		TokenFailureCacheTTL: s.TokenFailureCacheTTL,
	}

	if s.Anonymous != nil {
		ret.Anonymous = s.Anonymous.Allow
	}

	if s.BootstrapToken != nil {
		ret.BootstrapToken = s.BootstrapToken.Enable
	}

	if s.ClientCert != nil {
		ret.ClientCAFile = s.ClientCert.ClientCA
	}

	if s.PasswordFile != nil {
		ret.BasicAuthFile = s.PasswordFile.BasicAuthFile
	}

	if s.TokenFile != nil {
		ret.TokenAuthFile = s.TokenFile.TokenFile
	}

	if s.WebHook != nil {
		ret.WebhookTokenAuthnConfigFile = s.WebHook.ConfigFile
		ret.WebhookTokenAuthnCacheTTL = s.WebHook.CacheTTL

		if len(s.WebHook.ConfigFile) > 0 && s.WebHook.CacheTTL > 0 {
			if s.TokenSuccessCacheTTL > 0 && s.WebHook.CacheTTL < s.TokenSuccessCacheTTL {
				glog.Warningf("the webhook cache ttl of %s is shorter than the overall cache ttl of %s for successful token authentication attempts.", s.WebHook.CacheTTL, s.TokenSuccessCacheTTL)
			}
			if s.TokenFailureCacheTTL > 0 && s.WebHook.CacheTTL < s.TokenFailureCacheTTL {
				glog.Warningf("the webhook cache ttl of %s is shorter than the overall cache ttl of %s for failed token authentication attempts.", s.WebHook.CacheTTL, s.TokenFailureCacheTTL)
			}
		}
	}

	return ret
}
