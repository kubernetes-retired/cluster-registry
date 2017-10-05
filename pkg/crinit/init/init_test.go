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

package init

import (
	"bytes"
	"errors"
	"reflect"
	"sort"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	clientgotesting "k8s.io/client-go/testing"
)

func TestValidateOptions(t *testing.T) {
	testNodePort := int32(1000)

	testCases := []struct {
		desc        string
		initialOpts *initClusterRegistryOptions
		finalOpts   *initClusterRegistryOptions
		errExpected bool
	}{
		{
			desc:        "LoadBalancer service type supported",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeLoadBalancer)},
			finalOpts:   &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeLoadBalancer), apiServerServiceType: v1.ServiceTypeLoadBalancer},
			errExpected: false,
		},
		{
			desc:        "NodePort service type supported",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort)},
			finalOpts:   &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerServiceType: v1.ServiceTypeNodePort},
			errExpected: false,
		},
		{
			desc:        "other service type not supported",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeExternalName)},
			errExpected: true,
		},
		{
			desc:        "adveritse address supported with NodePort service type",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerAdvertiseAddress: "10.0.0.1"},
			errExpected: false,
		},
		{
			desc:        "adveritse address not supported with non-NodePort service type",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeLoadBalancer), apiServerAdvertiseAddress: "10.0.0.1"},
			errExpected: true,
		},
		{
			desc:        "adveritse address validated",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerAdvertiseAddress: "notAValidIP"},
			errExpected: true,
		},
		{
			desc:        "advertise port supported with NodePort service type",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerNodePortPort: testNodePort},
			finalOpts:   &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerServiceType: v1.ServiceTypeNodePort, apiServerNodePortPort: testNodePort, apiServerNodePortPortPtr: &testNodePort},
			errExpected: false,
		},
		{
			desc:        "advertise port rejected with non-NodePort service type",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeLoadBalancer), apiServerNodePortPort: testNodePort},
			errExpected: true,
		},
		{
			desc:        "advertise port rejected if out of range",
			initialOpts: &initClusterRegistryOptions{apiServerServiceTypeString: string(v1.ServiceTypeNodePort), apiServerNodePortPort: 100000},
			errExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := validateOptions(tc.initialOpts)
			if tc.errExpected {
				if err == nil {
					t.Errorf("expected error, but no error was returned")
				}
			} else {
				if err != nil {
					t.Errorf("no error expected, but got %v", err)
				}
				if tc.finalOpts != nil {
					if !reflect.DeepEqual(tc.initialOpts, tc.finalOpts) {
						t.Errorf("unexpected output: got: %v, want: %v", tc.initialOpts, tc.finalOpts)
					}
				}
			}
		})
	}
}

func TestMarshalOptions(t *testing.T) {
	testCases := []struct {
		desc           string
		overrideParams string
		expectedMap    map[string]string
		expectedErr    string
	}{
		{
			desc:           "valid arguments",
			overrideParams: "valid-format-param1=override1,valid-format-param2=override2",
			expectedMap:    map[string]string{"valid-format-param1": "override1", "valid-format-param2": "override2"},
			expectedErr:    "",
		},
		{
			desc:           "empty arguments",
			overrideParams: "",
			expectedMap:    nil,
			expectedErr:    "",
		},
		{
			desc:           "zero value arugment",
			overrideParams: "zero-value-arg=",
			expectedMap:    map[string]string{"zero-value-arg": ""},
			expectedErr:    "",
		},
		{
			// TODO: Multiple arg values separated by , are not supported yet
			desc:           "multiple equals characters in an argument",
			overrideParams: "multiple-equalto-char=first-key=1",
			expectedMap:    map[string]string{"multiple-equalto-char": "first-key=1"},
			expectedErr:    "",
		},
		{
			desc:           "incorrectly formatted argument",
			overrideParams: "wrong-format-arg",
			expectedErr:    "wrong format for override arg: wrong-format-arg",
		},
		{
			desc:           "incorrectly formatted argument with only a value",
			overrideParams: "=wrong-format-only-value",
			expectedErr:    "wrong format for override arg: =wrong-format-only-value, arg name cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			options := initClusterRegistryOptions{apiServerOverridesString: tc.overrideParams}
			err := marshalOptions(&options)
			if tc.expectedErr == "" {
				got := options.apiServerOverrides
				want := tc.expectedMap

				if !reflect.DeepEqual(got, want) {
					t.Errorf("unexpected output: got: %v, want: %v", got, want)
				}
			} else {
				if err.Error() != tc.expectedErr {
					t.Errorf("unexpected error output: got: %s, want: %s", err.Error(), tc.expectedErr)
				}
			}
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	t.Run("simple namespace creation", func(t *testing.T) {
		name := "test"
		client := fake.NewSimpleClientset()
		ns, err := createNamespace(client, name, false)
		if ns == nil {
			t.Error("namespace not created")
		}
		if ns.Name != name {
			t.Error("namespace has wrong name: got '%v', want '%v')", ns.Name, name)
		}
		if serverNs, _ := client.CoreV1().Namespaces().Get(name, metav1.GetOptions{}); serverNs == nil {
			t.Error("namespace not created on server")
		}
		if err != nil {
			t.Error("unexpected error: got %v", err)
		}
	})

	t.Run("dry run should not create namespace on server", func(t *testing.T) {
		name := "test2"
		client := fake.NewSimpleClientset()
		ns, _ := createNamespace(client, name, true)
		if ns == nil {
			t.Error("namespace not returned")
		}
		if ns.Name != name {
			t.Error("namespace has wrong name: got '%v', want '%v')", ns.Name, name)
		}
		if serverNs, _ := client.CoreV1().Namespaces().Get(name, metav1.GetOptions{}); serverNs != nil {
			t.Error("dry run should not create namespace")
		}
	})

	t.Run("should return client error", func(t *testing.T) {
		name := "test3"
		client := &fake.Clientset{}
		client.AddReactor("create", "namespaces", func(action clientgotesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("error")
		})
		ns, err := createNamespace(client, name, false)
		if err == nil {
			t.Error("expected error, got none")
		}
		if ns != nil {
			t.Error("namespace should not be created")
		}
	})
}

func TestCreateService(t *testing.T) {
	t.Run("should create service on server", func(t *testing.T) {
		name := "test"
		client := fake.NewSimpleClientset()
		buffer := &bytes.Buffer{}
		createService(buffer, client, "ns", name, "", nil, v1.ServiceTypeClusterIP, false)
		if serverSvc, _ := client.CoreV1().Services("ns").Get(name, metav1.GetOptions{}); serverSvc == nil {
			t.Error("should create service")
		}
	})

	t.Run("dry run should not create service", func(t *testing.T) {
		name := "test"
		client := fake.NewSimpleClientset()
		buffer := &bytes.Buffer{}
		svc, ips, hostnames, err := createService(buffer, client, "ns", name, "", nil, v1.ServiceTypeClusterIP, true)
		if svc == nil {
			t.Error("service not returned")
		}
		if svc.Name != name {
			t.Error("service has wrong name: got '%v', want '%v')", svc.Name, name)
		}
		if serverSvc, _ := client.CoreV1().Services("ns").Get(name, metav1.GetOptions{}); serverSvc != nil {
			t.Error("dry run should not create service")
		}
		if len(ips) > 0 {
			t.Error("ips not expected in dry-run mode: got %v", ips)
		}
		if len(hostnames) > 0 {
			t.Error("hostnames not expected in dry-run mode: got %v", hostnames)
		}
		if err != nil {
			t.Error("unexpected error: got %v", err)
		}
	})

	node := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				v1.NodeAddress{
					Type:    v1.NodeExternalIP,
					Address: "200.0.0.1",
				},
			},
		},
	}
	loadBalancerAddress := "150.0.0.1"
	loadBalancerHostname := "foo"

	testCases := []struct {
		desc              string
		name              string
		namespace         string
		serviceType       v1.ServiceType
		advertiseAddress  string
		advertisePort     int32
		expectedIPs       []string
		expectedHostnames []string
	}{
		{
			desc:              "simple service creation",
			name:              "service",
			namespace:         "ns",
			expectedIPs:       []string{"200.0.0.1"},
			expectedHostnames: []string{},
		},
		{
			desc:              "service with advertise address",
			name:              "service",
			namespace:         "ns",
			advertiseAddress:  "100.0.0.1",
			expectedIPs:       []string{"100.0.0.1"},
			expectedHostnames: []string{},
		},
		{
			desc:              "service with load balancer",
			name:              "service",
			namespace:         "ns",
			serviceType:       v1.ServiceTypeLoadBalancer,
			expectedIPs:       []string{loadBalancerAddress},
			expectedHostnames: []string{loadBalancerHostname},
		},
		// TODO: NodePort service with provided API server advertise address
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := fake.NewSimpleClientset(&node)

			// Modify the service to have appropriate load balancing status if necessary.
			client.PrependReactor("create", "services", func(action clientgotesting.Action) (bool, runtime.Object, error) {
				createAction := action.(clientgotesting.CreateAction)
				svc := createAction.GetObject().(*v1.Service)
				if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
					svc.Status.LoadBalancer.Ingress = []v1.LoadBalancerIngress{v1.LoadBalancerIngress{IP: loadBalancerAddress, Hostname: loadBalancerHostname}}
					return false, svc, nil
				}
				return false, nil, nil
			})

			buffer := &bytes.Buffer{}
			svc, ips, hostnames, _ := createService(buffer, client, tc.namespace, tc.name, tc.advertiseAddress, &tc.advertisePort, tc.serviceType, false)

			if svc == nil {
				t.Error("service not returned")
			}
			if svc.Name != tc.name {
				t.Error("service has wrong name: got '%v', want '%v')", svc.Name, tc.name)
			}
			if !reflect.DeepEqual(ips, tc.expectedIPs) {
				t.Error("unexpected ips returned: got %v, want: %v", ips, tc.expectedIPs)
			}
			if !reflect.DeepEqual(hostnames, tc.expectedHostnames) {
				t.Error("unexpected hostnames returned: got %v, want: %v", hostnames, tc.expectedHostnames)
			}
		})
	}

	errorTestCases := []struct {
		desc            string
		serviceType     v1.ServiceType
		errorActionVerb string
		errorActionType string
	}{
		{
			desc:            "failure creating service",
			errorActionVerb: "create",
			errorActionType: "services",
			serviceType:     v1.ServiceTypeClusterIP,
		},
		{
			desc:            "failure getting nodes",
			errorActionVerb: "list",
			errorActionType: "nodes",
			serviceType:     v1.ServiceTypeNodePort,
		},
		//		{
		//			desc:            "failure getting the service",
		//			errorActionVerb: "get",
		//			errorActionType: "services",
		//			serviceType:     v1.ServiceTypeLoadBalancer,
		//		},
	}

	for _, tc := range errorTestCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := fake.NewSimpleClientset()
			client.PrependReactor(tc.errorActionVerb, tc.errorActionType, func(action clientgotesting.Action) (bool, runtime.Object, error) {
				return true, nil, errors.New("error")
			})
			buffer := &bytes.Buffer{}
			svc, ips, hostnames, err := createService(buffer, client, "ns", "test", "", nil, tc.serviceType, false)
			if err == nil {
				t.Error("Expected error, got none")
			}
			if svc != nil {
				t.Error("service not expected with error: got %v", svc)
			}
			if ips != nil {
				t.Error("ips not expected with error: got %v", ips)
			}
			if hostnames != nil {
				t.Error("hostnames not expected with error: got %v", hostnames)
			}
		})
	}
}

func TestGetClusterNodeIPs(t *testing.T) {
	nodeOne := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodeOne",
		},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				v1.NodeAddress{
					Type:    v1.NodeExternalIP,
					Address: "200.0.0.1",
				},
			},
		},
	}

	nodeTwo := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodeTwo",
		},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				v1.NodeAddress{
					Type:    v1.NodeInternalIP,
					Address: "10.0.0.1",
				},
			},
		},
	}

	internalExternalNode := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "internaExternalNode",
		},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				v1.NodeAddress{
					Type:    v1.NodeExternalIP,
					Address: "200.0.0.2",
				},
				v1.NodeAddress{
					Type:    v1.NodeInternalIP,
					Address: "10.0.0.2",
				},
			},
		},
	}

	testCases := []struct {
		desc      string
		nodes     []runtime.Object
		wantedIPs []string
	}{
		{
			desc:      "no nodes",
			wantedIPs: []string{},
		},
		{
			desc:      "one node with external address",
			nodes:     []runtime.Object{&nodeOne},
			wantedIPs: []string{"200.0.0.1"},
		},
		{
			desc:      "two nodes, one internal, one external",
			nodes:     []runtime.Object{&nodeOne, &nodeTwo},
			wantedIPs: []string{"200.0.0.1", "10.0.0.1"},
		},
		{
			desc:      "internal and external address",
			nodes:     []runtime.Object{&internalExternalNode},
			wantedIPs: []string{"200.0.0.2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := fake.NewSimpleClientset(tc.nodes...)
			got, _ := getClusterNodeIPs(client)
			want := tc.wantedIPs
			sort.Strings(got)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected ip addresses: got: %v, want: %v", got, want)
			}
		})
	}
}
