package tableconvertor

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
)

var swaggerMetadataDescriptions = metav1.ObjectMeta{}.SwaggerDoc()

func New() (rest.TableConvertor, error) {
	c := &convertor{}
	return c, nil
}

type convertor struct {
}

func (c *convertor) ConvertToTable(ctx genericapirequest.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1beta1.Table, error) {
	table := &metav1beta1.Table{
		ColumnDefinitions: []metav1beta1.TableColumnDefinition{
			{Name: "Name", Type: "string", Format: "name", Description: swaggerMetadataDescriptions["name"]},
			//			{Name: "Ready", Type: "integer", Description: swaggerMetadataDescriptions["numberReady"]},
			{Name: "External-IP", Type: "string", Description: swaggerMetadataDescriptions["externalIPs"]},
			{Name: "Age", Type: "string", Description: swaggerMetadataDescriptions["creationTimestamp"]},
		},
	}
	if m, err := meta.ListAccessor(obj); err == nil {
		table.ResourceVersion = m.GetResourceVersion()
		table.SelfLink = m.GetSelfLink()
		table.Continue = m.GetContinue()
	} else {
		if m, err := meta.CommonAccessor(obj); err == nil {
			table.ResourceVersion = m.GetResourceVersion()
			table.SelfLink = m.GetSelfLink()
		}
	}

	var err error
	table.Rows, err = metatable.MetaToTableRow(obj, func(obj runtime.Object, m metav1.Object, name, age string) ([]interface{}, error) {
		clusterRegistry := obj.(*v1alpha1.Cluster)
		cells := []interface{}{
			name, clusterRegistry.Spec.KubernetesAPIEndpoints.ServerEndpoints.ServerAddress, age,
		}
		return cells, nil
	})
	return table, err
}
