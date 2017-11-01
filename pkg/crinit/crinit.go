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

package crinit

import (
	"flag"
	"io"

	apiserverflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/client-go/tools/clientcmd"
	crinitinit "k8s.io/cluster-registry/pkg/crinit/init"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewClusterregistryCommand creates the `clusterregistry` command.
func NewClusterregistryCommand(out io.Writer, defaultServerImage, defaultEtcdImage string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "crinit",
		Short: "crinit runs a cluster registry in a Kubernetes cluster",
		Long:  "crinit bootstraps and runs a cluster registry as a Deployment in an existing Kubernetes cluster.",
	}

	// Add the command line flags from other dependencies (e.g., glog), but do not
	// warn if they contain underscores.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(apiserverflag.WordSepNormalizeFunc)
	rootCmd.PersistentFlags().AddFlagSet(pflag.CommandLine)

	// Warn for other flags that contain underscores.
	rootCmd.SetGlobalNormalizationFunc(apiserverflag.WarnWordSepNormalizeFunc)

	rootCmd.AddCommand(crinitinit.NewCmdInit(out, clientcmd.NewDefaultPathOptions(), defaultServerImage, defaultEtcdImage))

	return rootCmd
}
