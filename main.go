//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"flag"
	"os"

	sysdrift "github.com/briandowns/rancher-sys-drift/pkg/generated/controllers/sys-drift.io"
	sd "github.com/briandowns/rancher-sys-drift/pkg/sysdrift"
	"github.com/briandowns/rancher-sys-drift/pkg/version"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	kubeapiext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
)

var KubeConfig string

func main() {
	app := cli.NewApp()
	app.Name = "sys-drift"
	app.Version = version.FriendlyVersion()
	app.Usage = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			EnvVar:      "KUBECONFIG",
			Destination: &KubeConfig,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

// registerCRD
func registerCRD(ctx context.Context, cfg *rest.Config) error {
	xcs, err := kubeapiext.NewForConfig(cfg)
	if err != nil {
		return err
	}

	factory := crd.NewFactoryFromClientGetter(xcs)

	var crds []crd.CRD
	for _, crdFn := range []func() (*crd.CRD, error){
		sd.CRD,
	} {
		crdef, err := crdFn()
		if err != nil {
			return err
		}
		crds = append(crds, *crdef)
	}

	return factory.BatchCreateCRDs(ctx, crds...).BatchWait()
}

func run(c *cli.Context) {
	flag.Parse()

	logrus.Info("Starting sys-drift controller")
	ctx := signals.SetupSignalHandler(context.Background())

	kubeConfig, err := kubeconfig.GetNonInteractiveClientConfig(KubeConfig).ClientConfig()
	if err != nil {
		logrus.Fatalf("failed to find kubeconfig: %v", err)
	}

	if err := registerCRD(ctx, kubeConfig); err != nil {
		logrus.Fatal(err)
	}

	sdf, err := sysdrift.NewFactoryFromConfig(kubeConfig)
	if err != nil {
		logrus.Fatalf("Error building controllers: %s", err.Error())
	}

	sd.Register(ctx, sdf.Sysdrift().V1().SysDrift())

	if err := start.All(ctx, 2, sdf); err != nil {
		logrus.Fatalf("Error starting: %s", err.Error())
	}

	<-ctx.Done()
}
