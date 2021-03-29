package main

import (
	"os"

	"github.com/briandowns/rancher-sys-drift/pkg/apis/sys-drift.io/v1"
	"github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/briandowns/rancher-sys-drift/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"sys-drift.io": {
				Types: []interface{}{
					v1.SysDrift{},
				},
				GenerateTypes: true,
			},
		},
	})
}
