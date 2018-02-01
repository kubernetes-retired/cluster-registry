/*
Copyright 2018 The Kubernetes Authors.

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

package aggregated

import (
	"io"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/cluster-registry/pkg/crinit/util"
	apiregclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

var (
	longDeleteCommandDescription = `
	Delete deletes an aggregated cluster registry.

	The aggregated cluster registry is hosted inside a Kubernetes
	cluster and has its API registered with the Kubernetes API aggregator.
	The host cluster must be specified using the --host-cluster-context flag.`
	deleteCommandExample = `
	# Delete an aggregated cluster registry named foo
	# in the host cluster whose local kubeconfig
	# context is bar.
	crinit aggregated delete foo --host-cluster-context=bar`
)

// newSubCmdDelete defines the `delete` subcommand to bootstrap a cluster registry
// inside a host Kubernetes cluster.
func newSubCmdDelete(cmdOut io.Writer, pathOptions *clientcmd.PathOptions) *cobra.Command {
	opts := &aggregatedClusterRegistryOptions{}

	delCmd := &cobra.Command{
		Use:     "delete CLUSTER_REGISTRY_NAME --host-cluster-context=HOST_CONTEXT",
		Short:   "Delete an aggregated cluster registry.",
		Long:    longDeleteCommandDescription,
		Example: deleteCommandExample,
		Run: func(cmd *cobra.Command, args []string) {
			err := opts.SetName(args)
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
			apiServiceClientset, err := apiregclient.NewForConfig(hostConfig)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
			err = RunDelete(opts, cmdOut, hostClientset, apiServiceClientset, pathOptions)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
		},
	}

	flags := delCmd.Flags()
	opts.BindBase(flags)
	return delCmd
}

// RunDelete deletes a cluster registry.
func RunDelete(opts *aggregatedClusterRegistryOptions, cmdOut io.Writer,
	hostClientset client.Interface, apiSvcClientset apiregclient.Interface,
	pathOptions *clientcmd.PathOptions) error {
	return nil
}
