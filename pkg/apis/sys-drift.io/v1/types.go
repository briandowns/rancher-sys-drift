package v1

import (
	"github.com/briandowns/rancher-sys-drift/pkg/files"
	"github.com/briandowns/rancher-sys-drift/pkg/kernel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SysDrift
type SysDrift struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SysDriftSpec `json:"spec"`
}

// CPU
type CPU struct {
	Count int `json:"count"`
}

// play with different ideas. Worst case is mirroring the type here
// as well as the kernel package.
//
// Module represents the state of a kernel module.
type Module kernel.Module

// BinaryStateStore contains the previously recorded
// state and the current binary state.
type BinaryStateStore files.BinaryStateStore

// BinaryState
type BinaryState files.BinaryState

//SysDriftSpec
type SysDriftSpec struct {
	Timestamp     metav1.Time             `json:"timestamp,omitempty"`
	UID           string                  `json:"uid,omitempty"`
	Hostname      string                  `json:"hostname,omitempty"`
	KernelVersion string                  `json:"kernel_version,omitempty"`
	Arch          string                  `json:"arch,omitempty"`
	TZ            string                  `json:"tz,omitempty"`
	NICs          []string                `json:"nics,omitempty"`
	Modules       []Module                `json:"modules,omitempty"`
	BinaryHashes  *files.BinaryStateStore `json:"binary_hashes,omitempty"`
	CPU           *CPU                    `json:"cpu,omitempty"`
}
