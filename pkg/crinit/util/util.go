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

package util

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/spf13/pflag"
)

const (
	// DefaultClusterRegistryNamespace is the default namespace in which
	// cluster registry cmponents are hosted.
	DefaultClusterRegistryNamespace = "clusterregistry"
)

// SubcommandOptions holds the configuration required by the subcommands of
// `clusterregistry`.
type SubcommandOptions struct {
	Name                     string
	Host                     string
	ClusterRegistryNamespace string
	Kubeconfig               string
}

func (o *SubcommandOptions) Bind(flags *pflag.FlagSet) {
	flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	flags.StringVar(&o.Host, "host-cluster-context", "", "Context of the cluster in which to host the cluster registry.")
	flags.StringVar(&o.ClusterRegistryNamespace, "cluster-registry-namespace", DefaultClusterRegistryNamespace, "Namespace in the host cluster where the cluster registry components are installed")
}

func (o *SubcommandOptions) SetName(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("NAME is required")
	}
	o.Name = args[0]
	return nil
}

func GetClientConfig(pathOptions *clientcmd.PathOptions, context, kubeconfigPath string) clientcmd.ClientConfig {
	loadingRules := *pathOptions.LoadingRules
	loadingRules.Precedence = pathOptions.GetLoadingPrecedence()
	loadingRules.ExplicitPath = kubeconfigPath
	overrides := &clientcmd.ConfigOverrides{
		CurrentContext: context,
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&loadingRules, overrides)
}
