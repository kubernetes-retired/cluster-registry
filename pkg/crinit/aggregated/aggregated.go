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

// Package aggregated contains the implementation of the `aggregated` subcommand
// of crinit, which deploys a cluster registry as an aggregated API server in a
// Kubernetes cluster.
package aggregated

import (
	"crypto/x509"
	"fmt"
	"io"
	"strings"

	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/cert"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	"k8s.io/cluster-registry/pkg/crinit/util"
	apiregv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	apiregclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	longInitCommandDescription = `
	Init initializes an aggregated cluster registry.

	The aggregated cluster registry is hosted inside a Kubernetes
	cluster and registers its API with the Kubernetes API aggregator.
	The host cluster must be specified using the --host-cluster-context flag.`
	initCommandExample = `
	# Initialize an aggregated cluster registry named foo
	# in the host cluster whose local kubeconfig
	# context is bar.
	crinit aggregated init foo --host-cluster-context=bar`

	// Set priorities for our API in the APIService object.
	apiServiceGroupPriorityMinimum int32 = 10000
	apiServiceVersionPriority      int32 = 20

	// Name used for our cluster registry APIService object to register with
	// the K8s API aggregator.
	apiServiceName = v1alpha1.SchemeGroupVersion.Version + "." + v1alpha1.GroupName

	// Name used for our cluster registry service account object to be used
	// with our cluster role objects.
	serviceAccountName = strings.Replace(v1alpha1.GroupName, ".", "-", -1) + "-apiserver"

	// Name used for our cluster role object to subsequently specify what
	// operations we want to allow on our API resources via cluster role
	// bindings.
	clusterRoleName = v1alpha1.GroupName + ":apiserver"

	// The name of our cluster registry API group used in the creation of the
	// cluster role binding.
	clusterRoleAPIGroup = []string{v1alpha1.GroupName}

	// The list of our cluster registry API resources to which the rule applies
	// in the cluster role object.
	clusterRoleResources = []string{"cluster"}

	// The list of verbs that apply to our cluster registry API resources.
	clusterRoleVerbs = []string{"get", "list", "watch", "create", "update", "patch", "delete"}

	// Name used for our cluster registry cluster role binding (CRB) object to
	// specify the operations we want to allow on our API resources.
	apiServerCRBName = v1alpha1.GroupName + ":apiserver"

	// Name used for our cluster registry cluster role binding (CRB) object that
	// allows delegated authentication and authorization checks.
	authDelegatorCRBName = v1alpha1.GroupName + ":apiserver-auth-delegator"

	// Name used for the cluster registry role binding that allows the cluster
	// registry service account to access the extension-apiserver-authentication
	// ConfigMap.
	extensionAPIServerRBName = v1alpha1.GroupName + ":extension-apiserver-authentication-reader"
)

type aggregatedClusterRegistryOptions struct {
	util.SubcommandOptions
	apiServerServiceTypeString string
}

func (o *aggregatedClusterRegistryOptions) Bind(flags *pflag.FlagSet) {
	flags.StringVar(&o.apiServerServiceTypeString, util.APIServerServiceTypeFlag,
		string(v1.ServiceTypeNodePort),
		"The type of service to create for the cluster registry. Options: 'LoadBalancer', 'NodePort'.")
}

// NewCmdAggregated defines the `aggregated` command that bootstraps a cluster registry
// inside a host Kubernetes cluster.
func NewCmdAggregated(cmdOut io.Writer, pathOptions *clientcmd.PathOptions, defaultServerImage, defaultEtcdImage string) *cobra.Command {
	opts := &aggregatedClusterRegistryOptions{}

	cmd := &cobra.Command{
		Use:   "aggregated",
		Short: "Subcommands to manage an aggregated cluster registry",
		Long:  "Commands used to manage an aggregated cluster registry. That is, a cluster registry that is aggregated with another Kubernetes API server.",
	}

	initCmd := &cobra.Command{
		Use:     "init CLUSTER_REGISTRY_NAME --host-cluster-context=HOST_CONTEXT",
		Short:   "Initialize an aggregated cluster registry.",
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
			apiServiceClientset, err := apiregclient.NewForConfig(hostConfig)
			if err != nil {
				glog.Fatalf("error: %v", err)
			}
			err = Run(opts, cmdOut, hostClientset, apiServiceClientset, pathOptions)
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
func validateOptions(opts *aggregatedClusterRegistryOptions) error {
	opts.APIServerServiceType = v1.ServiceType(opts.apiServerServiceTypeString)
	return opts.ValidateCommonOptions()
}

// Run initializes a cluster registry.
func Run(opts *aggregatedClusterRegistryOptions, cmdOut io.Writer,
	hostClientset client.Interface, apiSvcClientset apiregclient.Interface, pathOptions *clientcmd.PathOptions) error {

	err := opts.CreateNamespace(cmdOut, hostClientset)
	if err != nil {
		return err
	}

	svc, ips, hostnames, err := opts.CreateService(cmdOut, hostClientset)
	if err != nil {
		return err
	}

	credentials, err := opts.GenerateCredentials(cmdOut, svc.Name, ips, hostnames,
		false, false)
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

	sa, err := createRBACObjects(cmdOut, hostClientset, opts)
	if err != nil {
		return err
	}

	err = opts.CreateAPIServer(cmdOut, hostClientset, false, false, true, ips, pvc, sa.Name)
	if err != nil {
		return err
	}

	_, err = createAPIService(cmdOut, apiSvcClientset, opts,
		util.GetCAKeyPair(credentials).Cert)
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

// createRBACObjects handles the creation of all the RBAC objects necessary
// to deploy the cluster registry in aggregated mode.
func createRBACObjects(cmdOut io.Writer, clientset client.Interface,
	opts *aggregatedClusterRegistryOptions) (*v1.ServiceAccount, error) {

	fmt.Fprintf(cmdOut, "Creating RBAC objects...")

	// Create a Kubernetes service account in our namespace.
	glog.V(4).Infof(
		"Creating service account %v for cluster registry apiserver in the host cluster",
		serviceAccountName)

	sa, err := createServiceAccount(clientset, opts.ClusterRegistryNamespace, opts.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create service account %v: %v", sa, err)
		return nil, err
	}

	glog.V(4).Info("Successfully created service account")

	// Create a Kubernetes cluster role to allow REST operations on our
	// cluster registry API resources e.g. Cluster.
	glog.V(4).Infof("Creating cluster role %v", clusterRoleName)

	cr, err := createClusterRole(clientset, opts.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create cluster role %v: %v", cr, err)
		return nil, err
	}

	glog.V(4).Info("Successfully created cluster role")

	// Create a Kubernetes cluster role binding from the default service account
	// in our namespace to the cluster role we just created.
	glog.V(4).Infof("Creating cluster role bindings %v and %v", apiServerCRBName, authDelegatorCRBName)

	err = createClusterRoleBindings(clientset, opts.ClusterRegistryNamespace, opts.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create cluster role bindings")
		return nil, err
	}

	// Create a role binding to allow the cluster registry service account to
	// access the extension-apiserver-authentication configmap.
	glog.V(4).Infof("Creating role %v for accessing extension-apiserver-authentication ConfigMap", extensionAPIServerRBName)

	_, err = createExtensionAPIServerAuthenticationRoleBinding(clientset, extensionAPIServerRBName, opts.ClusterRegistryNamespace, opts.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create extension-apiserver-authentication ConfigMap reader role binding")
		return nil, err
	}

	glog.V(4).Info("Successfully created cluster role bindings")
	fmt.Fprintln(cmdOut, " done")
	return sa, nil
}

// createServiceAccount handles the creation of the service account for
// the cluster registry to be used with RBAC.
func createServiceAccount(clientset client.Interface,
	namespace string, dryRun bool) (*v1.ServiceAccount, error) {

	sa := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
			Labels:    util.ComponentLabel,
		},
	}

	if dryRun {
		return sa, nil
	}

	return clientset.CoreV1().ServiceAccounts(namespace).Create(sa)
}

// createClusterRole creates the cluster role for the operations we will allow
// on our cluster registry API resources e.g. Cluster.
func createClusterRole(clientset client.Interface,
	dryRun bool) (*rbacv1.ClusterRole, error) {

	rule := rbacv1.PolicyRule{
		Verbs:     clusterRoleVerbs,
		APIGroups: clusterRoleAPIGroup,
		Resources: clusterRoleResources,
	}

	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   clusterRoleName,
			Labels: util.ComponentLabel,
		},
		Rules: []rbacv1.PolicyRule{rule},
	}

	if dryRun {
		return cr, nil
	}

	return clientset.RbacV1().ClusterRoles().Create(cr)
}

// createClusterRoleBindings creates the cluster role bindings for the
// operations we will allow on our cluster registry API resources.
func createClusterRoleBindings(clientset client.Interface,
	namespace string, dryRun bool) error {

	// Create cluster role binding for the clusterregistry.k8s.io:apiserver
	// cluster role.
	crb, err := createClusterRoleBindingObject(clientset, apiServerCRBName,
		rbacv1.ServiceAccountKind, serviceAccountName, namespace, rbacv1.GroupName,
		"ClusterRole", clusterRoleName, util.ComponentLabel, dryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create cluster role binding %v: %v", crb, err)
		return err
	}

	// Create cluster role binding for the system:auth-delegator cluster role.
	crb, err = createClusterRoleBindingObject(clientset, authDelegatorCRBName,
		rbacv1.ServiceAccountKind, serviceAccountName, namespace, rbacv1.GroupName,
		"ClusterRole", "system:auth-delegator", util.ComponentLabel, dryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create cluster role binding %v: %v", crb, err)
		return err
	}

	return nil
}

// createClusterRoleBindingObject creates and returns the cluster role binding
// object.
func createClusterRoleBindingObject(clientset client.Interface, name, subjectKind,
	subjectName, subjectNamespace, roleAPIGroup, roleKind, roleName string,
	labels map[string]string, dryRun bool) (*rbacv1.ClusterRoleBinding, error) {

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      subjectKind,
				Name:      subjectName,
				Namespace: subjectNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: roleAPIGroup,
			Kind:     roleKind,
			Name:     roleName,
		},
	}

	if dryRun {
		return crb, nil
	}

	return clientset.RbacV1().ClusterRoleBindings().Create(crb)
}

// createExtensionApiserverAuthenticationRoleBinding creates and returns a rolebinding
// object to allow the cluster registry to access the extension-apiserver-authentication
// ConfigMap.
func createExtensionAPIServerAuthenticationRoleBinding(clientset client.Interface, name, namespace string, dryRun bool) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: util.ComponentLabel,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     "extension-apiserver-authentication-reader",
		},
	}

	if dryRun {
		return rb, nil
	}

	return clientset.RbacV1().RoleBindings("kube-system").Create(rb)
}

// createAPIService creates the Kubernetes API Service to handle the cluster
// registry objects.
func createAPIService(cmdOut io.Writer, clientset apiregclient.Interface,
	opts *aggregatedClusterRegistryOptions,
	ca *x509.Certificate) (*apiregv1beta1.APIService, error) {

	fmt.Fprint(cmdOut, "Creating cluster registry Kubernetes API Service...")
	glog.V(4).Infof("Creating cluster registry Kubernetes API Service %v", apiServiceName)

	caBundle := cert.EncodeCertPEM(ca)

	apiSvc, err := createAPIServiceObject(clientset, opts.Name,
		opts.ClusterRegistryNamespace, opts.DryRun, caBundle)

	if err != nil {
		glog.V(4).Infof("Failed to create cluster registry Kubernetes API Service %v: %v",
			apiSvc, err)
		return nil, err
	}

	fmt.Fprintln(cmdOut, " done")
	glog.V(4).Info("Successfully created cluster registry Kubernetes API Service")

	return apiSvc, nil
}

// createAPIServiceObject creates and returns the cluster registry API Service
// object.
func createAPIServiceObject(clientset apiregclient.Interface,
	clusterRegistryName, namespace string, dryRun bool,
	caBundle []byte) (*apiregv1beta1.APIService, error) {

	apiSvc := &apiregv1beta1.APIService{
		ObjectMeta: metav1.ObjectMeta{
			Name:   apiServiceName,
			Labels: util.ComponentLabel,
		},
		Spec: apiregv1beta1.APIServiceSpec{
			Service: &apiregv1beta1.ServiceReference{
				Namespace: namespace,
				Name:      clusterRegistryName,
			},
			Group:                v1alpha1.GroupName,
			Version:              v1alpha1.SchemeGroupVersion.Version,
			CABundle:             caBundle,
			GroupPriorityMinimum: apiServiceGroupPriorityMinimum,
			VersionPriority:      apiServiceVersionPriority,
		},
	}

	if dryRun {
		return apiSvc, nil
	}

	return clientset.ApiregistrationV1beta1().APIServices().Create(apiSvc)
}
