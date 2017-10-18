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

package main

import (
	"fmt"
	"os"

	"k8s.io/apiserver/pkg/util/logs"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/cluster-registry/pkg/crinit"
)

const (
	clusterregistryImageName = "clusterregistry:dev"
	defaultEtcdImage         = "gcr.io/google_containers/etcd:3.0.17"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	err := crinit.NewClusterregistryCommand(os.Stdout, clusterregistryImageName, defaultEtcdImage).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
