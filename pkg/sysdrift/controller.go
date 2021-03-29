package sysdrift

import (
	"context"

	v1 "github.com/briandowns/rancher-sys-drift/pkg/apis/sys-drift.io/v1"
	sdcontroller "github.com/briandowns/rancher-sys-drift/pkg/generated/controllers/sys-drift.io/v1"
)

type Controller struct {
	sdc sdcontroller.SysDriftController
}

func Register(ctx context.Context, sdc sdcontroller.SysDriftController) {
	controller := &Controller{
		sdc: sdc,
	}

	sdc.OnChange(ctx, "sys-drift-handler", controller.OnSysDriftChange)
	sdc.OnRemove(ctx, "sys-drift-handler", controller.OnSysDriftRemove)
}

func (c *Controller) OnSysDriftChange(key string, sd *v1.SysDrift) (*v1.SysDrift, error) {
	//change logic, return original sysdrift if no changes

	sysdriftCopy := sd.DeepCopy()
	//make changes to sysdriftCopy

	d, err := Load()
	if err != nil {
		return nil, err
	}
	sd.Spec = *d

	return c.sdc.Update(sysdriftCopy)
}

func (c *Controller) OnSysDriftRemove(key string, sd *v1.SysDrift) (*v1.SysDrift, error) {
	//remove logic, return original sysdrift if no changes

	sysdriftCopy := sd.DeepCopy()
	//make changes to sysdriftCopy
	return c.sdc.Update(sysdriftCopy)
}
