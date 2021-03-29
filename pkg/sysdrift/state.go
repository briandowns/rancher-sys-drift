package sysdrift

import (
	"net"
	"os"
	"runtime"

	v1 "github.com/briandowns/rancher-sys-drift/pkg/apis/sys-drift.io/v1"
	"github.com/briandowns/rancher-sys-drift/pkg/files"
	"github.com/briandowns/rancher-sys-drift/pkg/kernel"
	"github.com/elastic/go-sysinfo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// nics collects the names of the interfaces
// found on the system.
func nics() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var nics []string
	for _, i := range ifaces {
		nics = append(nics, i.Name)
	}
	return nics, nil
}

// Load collects all of the information necessary to
// determine if there is any drift.
func Load() (*v1.SysDriftSpec, error) {
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	h, err := sysinfo.Host()
	if err != nil {
		return nil, err
	}

	nics, err := nics()
	if err != nil {
		return nil, err
	}

	ss := v1.SysDriftSpec{
		Timestamp:     metav1.Now(),
		UID:           h.Info().UniqueID,
		Hostname:      hn,
		KernelVersion: h.Info().KernelVersion,
		Arch:          h.Info().Architecture,
		TZ:            h.Info().Timezone,
		NICs:          nics,
		CPU: &v1.CPU{
			Count: runtime.NumCPU(),
		},
	}

	modules, err := kernel.Modules()
	if err != nil {
		return nil, err
	}
	var nm []v1.Module
	for _, m := range modules {
		nm = append(nm, v1.Module(m))
	}
	ss.Modules = nm

	cbs := files.NewBinaryState()
	if err := cbs.GenerateBinaryHashes(); err != nil {
		return nil, err
	}
	ss.BinaryHashes.Current = cbs

	return &ss, nil
}
