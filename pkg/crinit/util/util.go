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

package util

import (
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/cert/triple"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

const (
	// DefaultClusterRegistryNamespace is the default namespace in which
	// cluster registry components are hosted.
	DefaultClusterRegistryNamespace = "clusterregistry"
	APIServerCN                     = "clusterregistry"
	AdminCN                         = "admin"
	HostClusterLocalDNSZoneName     = "cluster.local."
	APIServerNameSuffix             = "apiserver"
	CredentialSuffix                = "credentials"

	lbAddrRetryInterval = 5 * time.Second
	podWaitInterval     = 2 * time.Second

	apiserverAdvertiseAddressFlag = "api-server-advertise-address"
	APIServerServiceTypeFlag      = "api-server-service-type"
	apiserverPortFlag             = "api-server-port"

	apiServerSecurePortName = "https"
	// Set the secure port to 8443 to avoid requiring root privileges
	// to bind to port < 1000.  The apiserver's service will still
	// expose on port 443.
	apiServerSecurePort = 8443
)

var (
	serverName     string
	serverCredName string

	ComponentLabel = map[string]string{
		"app": "clusterregistry",
	}

	apiserverSvcSelector = map[string]string{
		"app":    "clusterregistry",
		"module": "clusterregistry-apiserver",
	}

	apiserverPodLabels = map[string]string{
		"app":    "clusterregistry",
		"module": "clusterregistry-apiserver",
	}
)

// SubcommandOptions holds the configuration required by the subcommands of
// `clusterregistry`.
type SubcommandOptions struct {
	Name                      string
	Host                      string
	ClusterRegistryNamespace  string
	Kubeconfig                string
	ServerImage               string
	EtcdImage                 string
	EtcdPVCapacity            string
	EtcdPVStorageClass        string
	EtcdPersistentStorage     bool
	DryRun                    bool
	ApiServerOverridesString  string
	ApiServerOverrides        map[string]string
	ApiServerServiceType      v1.ServiceType
	ApiServerAdvertiseAddress string
	ApiServerNodePortPort     int32
	ApiServerNodePortPortPtr  *int32
}

type entityKeyPairs struct {
	ca     *triple.KeyPair
	server *triple.KeyPair
	admin  *triple.KeyPair
}

type credentials struct {
	username        string
	password        string
	token           string
	certEntKeyPairs *entityKeyPairs
}

// BindCommon adds the common options that are shared by different
// sub-commands to the list of flags.
func (o *SubcommandOptions) BindCommon(flags *pflag.FlagSet, defaultServerImage, defaultEtcdImage string) {
	flags.StringVar(&o.Kubeconfig, "kubeconfig", "",
		"Path to the kubeconfig file to use for CLI requests.")
	flags.StringVar(&o.Host, "host-cluster-context", "",
		"Context of the cluster in which to host the cluster registry.")
	flags.StringVar(&o.ClusterRegistryNamespace, "cluster-registry-namespace",
		DefaultClusterRegistryNamespace,
		"Namespace in the host cluster where the cluster registry components are installed")
	flags.StringVar(&o.ServerImage, "image", defaultServerImage,
		"Image to use for the cluster registry API server binary.")
	flags.StringVar(&o.EtcdImage, "etcd-image", defaultEtcdImage,
		"Image to use for the etcd server binary.")
	flags.StringVar(&o.EtcdPVCapacity, "etcd-pv-capacity", "10Gi",
		"Size of the persistent volume claim to be used for etcd.")
	flags.StringVar(&o.EtcdPVStorageClass, "etcd-pv-storage-class", "",
		"The storage class of the persistent volume claim used for etcd. Must be provided if a default storage class is not enabled for the host cluster.")
	flags.BoolVar(&o.EtcdPersistentStorage, "etcd-persistent-storage", true,
		"Use a persistent volume for etcd. Defaults to 'true'.")
	flags.BoolVar(&o.DryRun, "dry-run", false,
		"Run the command in dry-run mode, without making any server requests.")
	flags.StringVar(&o.ApiServerOverridesString, "apiserver-arg-overrides", "",
		"Comma-separated list of cluster registry API server arguments to override, e.g., \"--arg1=value1,--arg2=value2...\"")
	flags.StringVar(&o.ApiServerAdvertiseAddress, apiserverAdvertiseAddressFlag, "",
		"Preferred address at which to advertise the cluster registry API server NodePort service. Valid only if '"+APIServerServiceTypeFlag+"=NodePort'.")
	flags.Int32Var(&o.ApiServerNodePortPort, apiserverPortFlag, 0,
		"Preferred port to use for the cluster registry API server NodePort service. Set to 0 to randomly assign a port. Valid only if '"+APIServerServiceTypeFlag+"=NodePort'.")
}

// SetName sets the name of the cluster registry.
func (o *SubcommandOptions) SetName(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("NAME is required")
	}
	o.Name = args[0]
	return nil
}

// ValidateCommonOptions validates the options that are shared across
// different sub-commands.
func (o *SubcommandOptions) ValidateCommonOptions() error {
	serverName = fmt.Sprintf("%s-%s", o.Name, APIServerNameSuffix)
	serverCredName = fmt.Sprintf("%s-%s", serverName, CredentialSuffix)

	if o.ApiServerServiceType != v1.ServiceTypeLoadBalancer &&
		o.ApiServerServiceType != v1.ServiceTypeNodePort {
		return fmt.Errorf("invalid %s: %s, should be either %s or %s",
			APIServerServiceTypeFlag, o.ApiServerServiceType,
			v1.ServiceTypeLoadBalancer, v1.ServiceTypeNodePort)
	}

	if o.ApiServerAdvertiseAddress != "" {
		ip := net.ParseIP(o.ApiServerAdvertiseAddress)
		if ip == nil {
			return fmt.Errorf("invalid %s: %s, should be a valid ip address",
				apiserverAdvertiseAddressFlag, o.ApiServerAdvertiseAddress)
		}
		if o.ApiServerServiceType != v1.ServiceTypeNodePort {
			return fmt.Errorf("%s should be passed only with '%s=NodePort'",
				apiserverAdvertiseAddressFlag, APIServerServiceTypeFlag)
		}
	}

	if o.ApiServerNodePortPort != 0 {
		if o.ApiServerServiceType != v1.ServiceTypeNodePort {
			return fmt.Errorf("%s should be passed only with '%s=NodePort'",
				apiserverPortFlag, APIServerServiceTypeFlag)
		}
		o.ApiServerNodePortPortPtr = &o.ApiServerNodePortPort
	} else {
		o.ApiServerNodePortPortPtr = nil
	}

	if o.ApiServerNodePortPort < 0 || o.ApiServerNodePortPort > 65535 {
		return fmt.Errorf("Please provide a valid port number for %s", apiserverPortFlag)
	}

	return nil
}

// marshalOptions marshals options if necessary.
func (o *SubcommandOptions) MarshalOptions() error {
	if o.ApiServerOverridesString == "" {
		return nil
	}

	argsMap := make(map[string]string)
	overrideArgs := strings.Split(o.ApiServerOverridesString, ",")
	for _, overrideArg := range overrideArgs {
		splitArg := strings.SplitN(overrideArg, "=", 2)
		if len(splitArg) != 2 {
			return fmt.Errorf("wrong format for override arg: %s", overrideArg)
		}
		key := strings.TrimSpace(splitArg[0])
		val := strings.TrimSpace(splitArg[1])
		if len(key) == 0 {
			return fmt.Errorf("wrong format for override arg: %s, arg name cannot be empty", overrideArg)
		}
		argsMap[key] = val
	}

	o.ApiServerOverrides = argsMap

	return nil
}

// CreateNamespace creates the cluster registry namespace.
func (o *SubcommandOptions) CreateNamespace(cmdOut io.Writer,
	clientset client.Interface) error {

	fmt.Fprintf(cmdOut, "Creating a namespace %s for the cluster registry...",
		o.ClusterRegistryNamespace)
	glog.V(4).Infof("Creating a namespace %s for the cluster registry",
		o.ClusterRegistryNamespace)

	_, err := CreateNamespace(clientset, o.ClusterRegistryNamespace, o.DryRun)

	if err != nil {
		return err
	}

	fmt.Fprintln(cmdOut, " done")
	return err
}

// CreateNamespace helper to create the cluster registry namespace object and return
// the object.
func CreateNamespace(clientset client.Interface, namespace string,
	dryRun bool) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	if dryRun {
		return ns, nil
	}

	return clientset.CoreV1().Namespaces().Create(ns)
}

// CreateService creates the cluster registry apiserver service.
func (o *SubcommandOptions) CreateService(cmdOut io.Writer,
	clientset client.Interface) (*v1.Service, []string, []string, error) {

	fmt.Fprint(cmdOut, "Creating cluster registry API server service...")
	glog.V(4).Info("Creating cluster registry API server service")

	svc, ips, hostnames, err := CreateService(cmdOut, clientset,
		o.ClusterRegistryNamespace, o.Name, o.ApiServerAdvertiseAddress,
		o.ApiServerNodePortPortPtr, o.ApiServerServiceType, o.DryRun)

	if err != nil {
		return nil, nil, nil, err
	}

	fmt.Fprintln(cmdOut, " done")
	glog.V(4).Infof("Created service named %s with IP addresses %v, hostnames %v",
		svc.Name, ips, hostnames)

	return svc, ips, hostnames, err
}

// CreateService helper to create the cluster registry apiserver service object
// and return the object. If service type is load balancer, will wait for load
// balancer IP and return it and hostnames.
func CreateService(cmdOut io.Writer, clientset client.Interface, namespace,
	svcName, apiserverAdvertiseAddress string, apiserverPort *int32,
	apiserverServiceType v1.ServiceType,
	dryRun bool) (*v1.Service, []string, []string, error) {

	port := v1.ServicePort{
		Name:       "https",
		Protocol:   "TCP",
		Port:       443,
		TargetPort: intstr.FromString(apiServerSecurePortName),
	}

	if apiserverServiceType == v1.ServiceTypeNodePort && apiserverPort != nil {
		port.NodePort = *apiserverPort
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: namespace,
			Labels:    ComponentLabel,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceType(apiserverServiceType),
			Selector: apiserverSvcSelector,
			Ports:    []v1.ServicePort{port},
		},
	}

	if dryRun {
		return svc, nil, nil, nil
	}

	var err error
	svc, err = clientset.CoreV1().Services(namespace).Create(svc)
	if err != nil {
		return nil, nil, nil, err
	}

	ips := []string{}
	hostnames := []string{}
	if apiserverServiceType == v1.ServiceTypeLoadBalancer {
		ips, hostnames, err = waitForLoadBalancerAddress(cmdOut, clientset, svc, dryRun)
	} else {
		if apiserverAdvertiseAddress != "" {
			ips = append(ips, apiserverAdvertiseAddress)
		} else {
			ips, err = GetClusterNodeIPs(clientset)
		}
	}
	if err != nil {
		return nil, nil, nil, err
	}

	return svc, ips, hostnames, err
}

// GetClusterNodeIPs returns a list of the IP addresses of nodes in the cluster,
// with a preference for external IP addresses.
func GetClusterNodeIPs(clientset client.Interface) ([]string, error) {
	preferredAddressTypes := []v1.NodeAddressType{
		v1.NodeExternalIP,
		v1.NodeInternalIP,
	}
	nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	nodeAddresses := []string{}
	for _, node := range nodeList.Items {
	OuterLoop:
		for _, addressType := range preferredAddressTypes {
			for _, address := range node.Status.Addresses {
				if address.Type == addressType {
					nodeAddresses = append(nodeAddresses, address.Address)
					break OuterLoop
				}
			}
		}
	}

	return nodeAddresses, nil
}

// waitForLoadBalancerAddress polls the apiserver for load balancer status
// to retrieve IPs and hostnames for it.
func waitForLoadBalancerAddress(cmdOut io.Writer, clientset client.Interface, svc *v1.Service, dryRun bool) ([]string, []string, error) {
	ips := []string{}
	hostnames := []string{}

	if dryRun {
		return ips, hostnames, nil
	}

	err := wait.PollImmediateInfinite(lbAddrRetryInterval, func() (bool, error) {
		fmt.Fprint(cmdOut, ".")
		pollSvc, err := clientset.CoreV1().Services(svc.Namespace).Get(svc.Name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		if ings := pollSvc.Status.LoadBalancer.Ingress; len(ings) > 0 {
			for _, ing := range ings {
				if len(ing.IP) > 0 {
					ips = append(ips, ing.IP)
				}
				if len(ing.Hostname) > 0 {
					hostnames = append(hostnames, ing.Hostname)
				}
			}
			if len(ips) > 0 || len(hostnames) > 0 {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return ips, hostnames, nil
}

// GenerateCredentials creates the credentials for apiserver secret.
func (o *SubcommandOptions) GenerateCredentials(cmdOut io.Writer, svcName string,
	ips, hostnames []string, apiServerEnableHTTPBasicAuth,
	apiServerEnableTokenAuth bool) (*credentials, error) {

	fmt.Fprint(cmdOut,
		"Creating cluster registry objects (credentials, persistent volume claim)...")
	glog.V(4).Info("Generating TLS certificates and credentials for communicating with the cluster registry API server")

	credentials, err := generateCredentials(o.ClusterRegistryNamespace, o.Name,
		svcName, HostClusterLocalDNSZoneName, ips, hostnames,
		apiServerEnableHTTPBasicAuth, apiServerEnableTokenAuth)

	if err != nil {
		return nil, err
	}

	return credentials, nil
}

// generateCredentials helper to create the certs for the apiserver.
func generateCredentials(svcNamespace, name, svcName, localDNSZoneName string,
	ips, hostnames []string, enableHTTPBasicAuth, enableTokenAuth bool) (*credentials, error) {

	credentials := credentials{
		username: AdminCN,
	}
	if enableHTTPBasicAuth {
		credentials.password = string(uuid.NewUUID())
	}
	if enableTokenAuth {
		credentials.token = string(uuid.NewUUID())
	}

	entKeyPairs, err := genCerts(svcNamespace, name, svcName, localDNSZoneName, ips, hostnames)
	if err != nil {
		return nil, err
	}
	credentials.certEntKeyPairs = entKeyPairs
	return &credentials, nil
}

func genCerts(svcNamespace, name, svcName, localDNSZoneName string,
	ips, hostnames []string) (*entityKeyPairs, error) {
	ca, err := triple.NewCA(name)

	if err != nil {
		return nil, fmt.Errorf("failed to create CA key and certificate: %v", err)
	}
	server, err := triple.NewServerKeyPair(ca, APIServerCN, svcName, svcNamespace, localDNSZoneName, ips, hostnames)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster registry API server key and certificate: %v", err)
	}
	admin, err := triple.NewClientKeyPair(ca, AdminCN, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client key and certificate for an admin: %v", err)
	}
	return &entityKeyPairs{
		ca:     ca,
		server: server,
		admin:  admin,
	}, nil
}

// CreateAPIServerCredentialsSecret creates the secret containing the
// apiserver credentials passed in.
func (o *SubcommandOptions) CreateAPIServerCredentialsSecret(clientset client.Interface,
	credentials *credentials) error {

	_, err := createAPIServerCredentialsSecret(clientset,
		o.ClusterRegistryNamespace, serverCredName, credentials, o.DryRun)

	if err != nil {
		return err
	}

	glog.V(4).Info("Certificates and credentials generated")
	return nil
}

// createAPIServerCredentialsSecret helper to create secret object and return
// the object.
func createAPIServerCredentialsSecret(clientset client.Interface, namespace,
	credentialsName string, credentials *credentials, dryRun bool) (*v1.Secret, error) {
	// Build the secret object with API server credentials.
	data := map[string][]byte{
		"ca.crt":     certutil.EncodeCertPEM(credentials.certEntKeyPairs.ca.Cert),
		"server.crt": certutil.EncodeCertPEM(credentials.certEntKeyPairs.server.Cert),
		"server.key": certutil.EncodePrivateKeyPEM(credentials.certEntKeyPairs.server.Key),
	}
	if credentials.password != "" {
		data["basicauth.csv"] = authFileContents(credentials.username, credentials.password)
	}
	if credentials.token != "" {
		data["token.csv"] = authFileContents(credentials.username, credentials.token)
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credentialsName,
			Namespace: namespace,
		},
		Data: data,
	}

	if dryRun {
		return secret, nil
	}

	return clientset.CoreV1().Secrets(namespace).Create(secret)
}

func (o *SubcommandOptions) CreatePVC(cmdOut io.Writer,
	clientset client.Interface, svcName string) (*v1.PersistentVolumeClaim, error) {

	if !o.EtcdPersistentStorage {
		return nil, nil
	}

	glog.V(4).Info("Creating a persistent volume and a claim to store the cluster registry API server's state, including etcd data")

	pvc, err := createPVC(clientset, o.ClusterRegistryNamespace, svcName,
		o.EtcdPVCapacity, o.EtcdPVStorageClass, o.DryRun)

	if err != nil {
		return nil, err
	}

	glog.V(4).Info("Persistent volume and claim created")
	fmt.Fprintln(cmdOut, " done")

	return pvc, nil
}

// createPVC helper to create the persistent volume claim object and
// return the object.
func createPVC(clientset client.Interface, namespace, svcName, etcdPVCapacity,
	etcdPVStorageClass string, dryRun bool) (*v1.PersistentVolumeClaim, error) {
	capacity, err := resource.ParseQuantity(etcdPVCapacity)
	if err != nil {
		return nil, err
	}

	var storageClassName *string
	if len(etcdPVStorageClass) > 0 {
		storageClassName = &etcdPVStorageClass
	}

	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-etcd-claim", svcName),
			Namespace: namespace,
			Labels:    ComponentLabel,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: capacity,
				},
			},
			StorageClassName: storageClassName,
		},
	}

	if dryRun {
		return pvc, nil
	}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Create(pvc)
}

func (o *SubcommandOptions) CreateAPIServer(cmdOut io.Writer, clientset client.Interface,
	apiServerEnableHTTPBasicAuth, apiServerEnableTokenAuth bool, ips []string,
	pvc *v1.PersistentVolumeClaim, serviceAccountName string, aggregated bool) error {
	// Since only one IP address can be specified as advertise address,
	// we arbitrarily pick the first available IP address.
	// Pick user provided apiserverAdvertiseAddress over other available IP addresses.
	advertiseAddress := o.ApiServerAdvertiseAddress
	if advertiseAddress == "" && len(ips) > 0 {
		advertiseAddress = ips[0]
	}

	fmt.Fprint(cmdOut, "Creating cluster registry deployment...")
	glog.V(4).Info("Creating cluster registry deployment")

	_, err := createAPIServer(clientset, o.ClusterRegistryNamespace,
		serverName, o.ServerImage, o.EtcdImage, advertiseAddress, serverCredName,
		apiServerEnableHTTPBasicAuth, apiServerEnableTokenAuth,
		o.ApiServerOverrides, pvc, serviceAccountName, aggregated, o.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to create API server: %v", err)
		return err
	}

	fmt.Fprintln(cmdOut, " done")
	glog.V(4).Info("Successfully created cluster registry deployment")

	return nil
}

// createAPIServer helper to create the apiserver deployment object and
// return the object.
func createAPIServer(clientset client.Interface, namespace, name, serverImage,
	etcdImage, advertiseAddress, credentialsName string, hasHTTPBasicAuthFile,
	hasTokenAuthFile bool, argOverrides map[string]string,
	pvc *v1.PersistentVolumeClaim, serviceAccountName string, aggregated, dryRun bool) (*appsv1beta1.Deployment, error) {

	command := []string{"./clusterregistry"}
	argsMap := map[string]string{
		"--bind-address":         "0.0.0.0",
		"--etcd-servers":         "http://localhost:2379",
		"--secure-port":          fmt.Sprintf("%d", apiServerSecurePort),
		"--client-ca-file":       "/etc/clusterregistry/apiserver/ca.crt",
		"--tls-cert-file":        "/etc/clusterregistry/apiserver/server.crt",
		"--tls-private-key-file": "/etc/clusterregistry/apiserver/server.key",
	}

	if advertiseAddress != "" {
		argsMap["--advertise-address"] = advertiseAddress
	}
	if hasHTTPBasicAuthFile {
		argsMap["--basic-auth-file"] = "/etc/clusterregistry/apiserver/basicauth.csv"
	}
	if hasTokenAuthFile {
		argsMap["--token-auth-file"] = "/etc/clusterregistry/apiserver/token.csv"
	}

	args := argMapsToArgStrings(argsMap, argOverrides)
	command = append(command, args...)
	if aggregated {
		command = append(command, "--use-delegated-auth")
	}

	replicas := int32(1)
	dep := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    ComponentLabel,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: apiserverPodLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "clusterregistry",
							Image:           serverImage,
							ImagePullPolicy: v1.PullAlways,
							Command:         command,
							Ports: []v1.ContainerPort{
								{
									Name:          apiServerSecurePortName,
									ContainerPort: apiServerSecurePort,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      credentialsName,
									MountPath: "/etc/clusterregistry/apiserver",
									ReadOnly:  true,
								},
							},
						},
						{
							Name:  "etcd",
							Image: etcdImage,
							Command: []string{
								"/usr/local/bin/etcd",
								"--data-dir",
								"/var/etcd/data",
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: credentialsName,
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: credentialsName,
								},
							},
						},
					},
				},
			},
		},
	}

	if pvc != nil {
		dataVolumeName := "etcddata"
		etcdVolume := v1.Volume{
			Name: dataVolumeName,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvc.Name,
				},
			},
		}
		etcdVolumeMount := v1.VolumeMount{
			Name:      dataVolumeName,
			MountPath: "/var/etcd",
		}

		dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, etcdVolume)
		for i, container := range dep.Spec.Template.Spec.Containers {
			if container.Name == "etcd" {
				dep.Spec.Template.Spec.Containers[i].VolumeMounts = append(dep.Spec.Template.Spec.Containers[i].VolumeMounts, etcdVolumeMount)
			}
		}
	}

	if len(serviceAccountName) > 0 {
		dep.Spec.Template.Spec.ServiceAccountName = serviceAccountName
	}

	if dryRun {
		return dep, nil
	}

	return clientset.AppsV1beta1().Deployments(namespace).Create(dep)
}

func argMapsToArgStrings(argsMap, overrides map[string]string) []string {
	for key, val := range overrides {
		argsMap[key] = val
	}
	args := []string{}
	for key, value := range argsMap {
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}
	// This is needed for the unit test deep copy to get an exact match.
	sort.Strings(args)
	return args
}

// UpdateKubeconfig handles updating the kubeconfig by building up the endpoint
// while printing and logging progress.
func (o *SubcommandOptions) UpdateKubeconfig(cmdOut io.Writer,
	pathOptions *clientcmd.PathOptions, svc *v1.Service, ips, hostnames []string,
	credentials *credentials) error {

	fmt.Fprint(cmdOut, "Updating kubeconfig...")
	glog.V(4).Info("Updating kubeconfig")

	// Pick the first ip/hostname to update the api server endpoint in kubeconfig
	// and also to give information to user.
	// In case of NodePort Service for api server, ips are node external ips.
	endpoint := ""
	if len(ips) > 0 {
		endpoint = ips[0]
	} else if len(hostnames) > 0 {
		endpoint = hostnames[0]
	}

	// If the service is nodeport, need to append the port to endpoint as it is
	// non-standard port.
	if o.ApiServerServiceType == v1.ServiceTypeNodePort {
		endpoint = endpoint + ":" + strconv.Itoa(int(svc.Spec.Ports[0].NodePort))
	}

	err := updateKubeconfig(pathOptions, o.Name, endpoint, o.Kubeconfig,
		credentials, o.DryRun)

	if err != nil {
		glog.V(4).Infof("Failed to update kubeconfig: %v", err)
		return err
	}

	fmt.Fprintln(cmdOut, " done")
	glog.V(4).Info("Successfully updated kubeconfig")
	return nil
}

// updateKubeconfig helper to update the kubeconfig file based on input
// parameters.
func updateKubeconfig(pathOptions *clientcmd.PathOptions, name, endpoint,
	kubeConfigPath string, credentials *credentials, dryRun bool) error {

	pathOptions.LoadingRules.ExplicitPath = kubeConfigPath
	kubeconfig, err := pathOptions.GetStartingConfig()
	if err != nil {
		return err
	}

	// Populate API server endpoint info.
	cluster := clientcmdapi.NewCluster()

	// Prefix "https" as the URL scheme to endpoint.
	if !strings.HasPrefix(endpoint, "https://") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}

	cluster.Server = endpoint
	cluster.CertificateAuthorityData = certutil.EncodeCertPEM(credentials.certEntKeyPairs.ca.Cert)

	// Populate credentials.
	authInfo := clientcmdapi.NewAuthInfo()
	authInfo.ClientCertificateData = certutil.EncodeCertPEM(credentials.certEntKeyPairs.admin.Cert)
	authInfo.ClientKeyData = certutil.EncodePrivateKeyPEM(credentials.certEntKeyPairs.admin.Key)
	authInfo.Token = credentials.token

	var httpBasicAuthInfo *clientcmdapi.AuthInfo

	if credentials.password != "" {
		httpBasicAuthInfo = clientcmdapi.NewAuthInfo()
		httpBasicAuthInfo.Password = credentials.password
		httpBasicAuthInfo.Username = credentials.username
	}

	// Populate context.
	context := clientcmdapi.NewContext()
	context.Cluster = name
	context.AuthInfo = name

	// Update the config struct with API server endpoint info,
	// credentials and context.
	kubeconfig.Clusters[name] = cluster
	kubeconfig.AuthInfos[name] = authInfo

	if httpBasicAuthInfo != nil {
		kubeconfig.AuthInfos[fmt.Sprintf("%s-basic-auth", name)] = httpBasicAuthInfo
	}

	kubeconfig.Contexts[name] = context

	if !dryRun {
		// Write the update kubeconfig.
		if err := clientcmd.ModifyConfig(pathOptions, *kubeconfig, true); err != nil {
			return err
		}
	}

	return nil
}

func (o *SubcommandOptions) WaitForAPIServer(cmdOut io.Writer,
	clientset client.Interface, pathOptions *clientcmd.PathOptions,
	ips, hostnames []string, svc *v1.Service) error {

	if o.DryRun {
		_, err := fmt.Fprintln(cmdOut, "Cluster registry can be run (dry run)")
		glog.V(4).Info("Cluster registry can be run (dry run)")
		return err
	}

	fmt.Fprint(cmdOut, "Waiting for the cluster registry API server to come up...")
	glog.V(4).Info("Waiting for the cluster registry API server to come up")

	err := waitForPods(cmdOut, clientset, []string{serverName},
		o.ClusterRegistryNamespace)

	if err != nil {
		return err
	}

	crConfig, err := GetClientConfig(pathOptions, o.Name, o.Kubeconfig).ClientConfig()

	if err != nil {
		return err
	}
	crClientset, err := client.NewForConfig(crConfig)

	if err != nil {
		return err
	}

	err = waitSrvHealthy(cmdOut, crClientset)

	if err != nil {
		return err
	}

	glog.V(4).Info("Cluster registry running")
	fmt.Fprintln(cmdOut, " done")
	return printSuccess(cmdOut, ips, hostnames, svc)
}

func waitForPods(cmdOut io.Writer, clientset client.Interface, pods []string, namespace string) error {
	err := wait.PollInfinite(podWaitInterval, func() (bool, error) {
		fmt.Fprint(cmdOut, ".")
		podCheck := len(pods)
		podList, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		for _, pod := range podList.Items {
			for _, fedPod := range pods {
				if strings.HasPrefix(pod.Name, fedPod) && pod.Status.Phase == "Running" {
					podCheck -= 1
				}
			}
			// ensure that all pods are in running state or keep waiting
			if podCheck == 0 {
				return true, nil
			}
		}
		return false, nil
	})
	return err
}

func waitSrvHealthy(cmdOut io.Writer, crClientset client.Interface) error {
	discoveryClient := crClientset.Discovery()
	return wait.PollInfinite(podWaitInterval, func() (bool, error) {
		fmt.Fprint(cmdOut, ".")
		body, err := discoveryClient.RESTClient().Get().AbsPath("/healthz").Do().Raw()
		if err != nil {
			return false, nil
		}
		if strings.EqualFold(string(body), "ok") {
			return true, nil
		}
		return false, nil
	})
}

func printSuccess(cmdOut io.Writer, ips, hostnames []string, svc *v1.Service) error {
	svcEndpoints := append(ips, hostnames...)
	endpoints := strings.Join(svcEndpoints, ", ")
	if svc.Spec.Type == v1.ServiceTypeNodePort {
		endpoints = ips[0] + ":" + strconv.Itoa(int(svc.Spec.Ports[0].NodePort))
		if len(ips) > 1 {
			endpoints = endpoints + ", ..."
		}
	}

	_, err := fmt.Fprintf(cmdOut, "Cluster registry API server is running at: %s\n", endpoints)
	return err
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

// authFileContents returns a CSV string containing the contents of an
// authentication file in the format required by the cluster registry.
func authFileContents(username, authSecret string) []byte {
	return []byte(fmt.Sprintf("%s,%s,%s\n", authSecret, username, uuid.NewUUID()))
}

// GetCAKeyPair retrieves the CA key pair stored in the internal credentials
// structure.
func GetCAKeyPair(credentials *credentials) *triple.KeyPair {
	if credentials == nil {
		glog.V(4).Info("credentials argument is nil!")
		return nil
	}

	return credentials.certEntKeyPairs.ca
}
