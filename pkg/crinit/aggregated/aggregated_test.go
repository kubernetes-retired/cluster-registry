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
	"bytes"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/util/cert/triple"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	"k8s.io/cluster-registry/pkg/crinit/options"
	fakeclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/fake"
)

var (
	apiSvcName = v1alpha1.SchemeGroupVersion.Version + "." + v1alpha1.GroupName
	apiGroup   = v1alpha1.GroupName
	apiVersion = v1alpha1.SchemeGroupVersion.Version
)

func TestCreateRBACObjects(t *testing.T) {
	testCases := []struct {
		desc        string
		opts        *aggregatedClusterRegistryOptions
		objExpected bool
	}{
		{
			desc: "should create RBAC object",
			opts: &aggregatedClusterRegistryOptions{
				SubcommandOptions: options.SubcommandOptions{
					Name: "test1",
					ClusterRegistryNamespace: "cr-namespace",
					DryRun: false,
				},
			},
			objExpected: true,
		},
		{
			desc: "dry run should not create RBAC objects",
			opts: &aggregatedClusterRegistryOptions{
				SubcommandOptions: options.SubcommandOptions{
					Name: "test2",
					ClusterRegistryNamespace: "dry-namespace",
					DryRun: true,
				},
			},
			objExpected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			client := fake.NewSimpleClientset()
			rbacObj, err := createRBACObjects(buffer, client, tc.opts)
			if err != nil {
				t.Error("Failed to create an RBAC Object")
			}
			createdRBACObj, _ := client.CoreV1().ServiceAccounts(tc.opts.ClusterRegistryNamespace).Get(rbacObj.Name, metav1.GetOptions{})
			if tc.objExpected {
				if createdRBACObj == nil {
					t.Errorf("Failed to create RBAC Object")
				}
			} else {
				if createdRBACObj != nil {
					t.Errorf("Expected no RBAC Object but got: %v", createdRBACObj)
				}
			}
		})
	}
}

func TestCreateAPIService(t *testing.T) {
	testCases := []struct {
		desc        string
		opts        *aggregatedClusterRegistryOptions
		svcExpected bool
	}{
		{
			desc: "should create an API Service",
			opts: &aggregatedClusterRegistryOptions{
				SubcommandOptions: options.SubcommandOptions{
					Name: "test1",
					ClusterRegistryNamespace: "cr-namespace",
					DryRun: false,
				},
			},
			svcExpected: true,
		},
		{
			desc: "dry-run should not create an API Service",
			opts: &aggregatedClusterRegistryOptions{
				SubcommandOptions: options.SubcommandOptions{
					Name: "test2",
					ClusterRegistryNamespace: "dry-namespace",
					DryRun: true,
				},
			},
			svcExpected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			ca, err := triple.NewCA(tc.opts.Name)
			if err != nil {
				t.Error("unable to create credential CA")
			}
			apisvcClient := fakeclientset.NewSimpleClientset()
			apisvcObj, err := createAPIService(buffer, apisvcClient, tc.opts, ca.Cert)
			if err != nil {
				t.Errorf("unexpected error creating API Service: got %v", err)
			}
			apiService, _ := apisvcClient.ApiregistrationV1beta1().APIServices().Get(apisvcObj.Name, metav1.GetOptions{})
			if tc.svcExpected {
				if apiService == nil {
					t.Errorf("Error in creating API Service")
				} else {
					if apiService.ObjectMeta.Name != apiSvcName {
						t.Errorf("Unexpected API Service Created. Expected %v, got %v", apiSvcName, apiService.ObjectMeta.Name)
					}
					if apiService.Spec.Group != apiGroup {
						t.Errorf("Unexpected API Service created. Expected %v, got %v", apiGroup, apiService.Spec.Group)
					}
					if apiService.Spec.Version != apiVersion {
						t.Errorf("Unexpected API Service created. Expected %v, got %v", apiVersion, apiService.Spec.Version)
					}
					if apiService.Spec.Service.Namespace != tc.opts.ClusterRegistryNamespace {
						t.Errorf("Unexpected API Service created. Expected %v, got %v", tc.opts.ClusterRegistryNamespace, apiService.Spec.Service.Namespace)
					}
					if apiService.Spec.Service.Name != tc.opts.Name {
						t.Errorf("Unexpected API Service created. Expected %v, got %v", tc.opts.Name, apiService.Spec.Service.Name)
					}
				}
			} else {
				if apiService != nil {
					t.Errorf("dry run should not create API Service but service is created")
				}
			}
		})
	}
}
