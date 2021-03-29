package sysdrift

import (
	v1 "github.com/briandowns/rancher-sys-drift/pkg/apis/sys-drift.io/v1"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/schemas/openapi"
)

func CRD() (*crd.CRD, error) {
	prototype := v1.NewSysDrift("", "", v1.SysDrift{})
	schema, err := openapi.ToOpenAPIFromStruct(*prototype)
	if err != nil {
		return nil, err
	}
	return &crd.CRD{
		GVK:        prototype.GroupVersionKind(),
		PluralName: v1.SysDriftResourceName,
		Status:     true,
		Schema:     schema,
		Categories: []string{"sys-state"},
	}, nil
}
