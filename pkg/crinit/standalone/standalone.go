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

// Package standalone contains the implementation of the `standalone` subcommand
// of `crinit`, which deploys a cluster registry as a standalone API server in
// a host Kubernetes cluster.
package standalone

import (
	"io"

	"k8s.io/api/core/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/cluster-registry/pkg/crinit/options"
	"k8s.io/cluster-registry/pkg/crinit/util"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	longInitCommandDescription = `
	Init initializes a standalone cluster registry.

	The standalone cluster registry is hosted inside a Kubernetes
	cluster but handles its own authentication and authorization.
	The host cluster must be specified using the
        --host-cluster-context flag.`
	initCommandExample = `
	# Initialize a standalone cluster registry named foo
	# in the host cluster whose local kubeconfig
	# context is bar.
	crinit standalone init foo --host-cluster-context=bar`
)

type standaloneClusterRegistryOptions struct {
	options.SubcommandOptions
	apiServerServiceTypeString   string
	apiServerEnableHTTPBasicAuth bool
	apiServerEnableTokenAuth     bool
}

func (o *standaloneClusterRegistryOptions) Bind(flags *pflag.FlagSet) {
	flags.StringVar(&o.apiServerServiceTypeString, options.APIServerServiceTypeFlag,
		string(v1.ServiceTypeLoadBalancer),
		"The type of service to create for the cluster registry. Options: 'LoadBalancer', 'NodePort'.")
	flags.BoolVar(&o.apiServerEnableHTTPBasicAuth, "apiserver-enable-basic-auth", false,
		"Enables HTTP Basic authentication for the cluster registry API server. Defaults to false.")
	flags.BoolVar(&o.apiServerEnableTokenAuth, "apiserver-enable-token-auth", false,
		"Enables token authentication for the cluster registry API server. Defaults to false.")
}

// NewCmdStandalone defines the `standalone` command that bootstraps a cluster registry
// inside a host Kubernetes cluster.
func NewCmdStandalone(cmdOut io.Writer, pathOptions *clientcmd.PathOptions, defaultServerImage, defaultEtcdImage string) *cobra.Command {
	opts := &standaloneClusterRegistryOptions{}

	cmd := &cobra.Command{
		Use:   "standalone",
		Short: "Subcommands to manage a standalone cluster registry",
		Long:  "Commands used to manage a standalone cluster registry. That is, a cluster registry that is not aggregated with another Kubernetes API server.",
	}

	initCmd := &cobra.Command{
		Use:     "init CLUSTER_REGISTRY_NAME --host-cluster-context=HOST_CONTEXT",
		Short:   "Initialize a standalone cluster registry",
		Long:    longInitCommandDescription,
		Example: initCommandExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := opts.SetName(args)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}

			err = validateOptions(opts)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}

			err = opts.MarshalOptions()
			if err != nil {
				glog.Fatalf("error: %v", err)
			}

			hostConfig, err := util.GetClientConfig(pathOptions, opts.Host, opts.Kubeconfig).ClientConfig()
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
			hostClientset, err := client.NewForConfig(hostConfig)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
			err = Run(opts, cmdOut, hostClientset, pathOptions)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
		},
	}

	flags := initCmd.Flags()
	opts.BindCommon(flags, defaultServerImage, defaultEtcdImage)
	opts.Bind(flags)

	cmd.AddCommand(initCmd)
	return cmd
}

// validateOptions ensures that options are valid.
func validateOptions(opts *standaloneClusterRegistryOptions) error {
	opts.APIServerServiceType = v1.ServiceType(opts.apiServerServiceTypeString)
	return opts.ValidateCommonOptions()
}

// Run initializes a cluster registry.
func Run(opts *standaloneClusterRegistryOptions, cmdOut io.Writer,
	hostClientset client.Interface, pathOptions *clientcmd.PathOptions) error {

	err := opts.CreateNamespace(cmdOut, hostClientset)
	if err != nil {
		return err
	}

	svc, ips, hostnames, err := opts.CreateService(cmdOut, hostClientset)
	if err != nil {
		return err
	}

	credentials, err := opts.GenerateCredentials(cmdOut, svc.Name, ips, hostnames,
		opts.apiServerEnableHTTPBasicAuth, opts.apiServerEnableTokenAuth)
	if err != nil {
		return err
	}

	err = opts.CreateAPIServerCredentialsSecret(hostClientset, credentials)
	if err != nil {
		return err
	}

	pvc, err := opts.CreatePVC(cmdOut, hostClientset, svc.Name)
	if err != nil {
		return err
	}

	err = opts.CreateAPIServer(cmdOut, hostClientset, opts.apiServerEnableHTTPBasicAuth,
		opts.apiServerEnableTokenAuth, false, ips, pvc, "default")
	if err != nil {
		return err
	}

	err = opts.UpdateKubeconfig(cmdOut, pathOptions, svc, ips, hostnames,
		credentials)
	if err != nil {
		return err
	}

	return opts.WaitForAPIServer(cmdOut, hostClientset, pathOptions, ips,
		hostnames, svc)
}
