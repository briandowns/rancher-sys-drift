package kernel

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type state uint8

const (
	Live state = iota
	Loading
	Unloading
)

// Module represents the state of a kernel module.
type Module struct {
	Name            string   `json:"name"`
	MemSize         int64    `json:"mem_size"`
	InstancesLoaded int64    `json:"instances_loaded"`
	Dependencies    []string `json:"dependencies"`
	State           string   `json:"state"`
	MemOffset       string   `json:"mem_offset"`
}

// Hash hashes the given module struct.
func (m *Module) Hash() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", m)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// String
func (m *Module) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

// Modules collects all kernel modules, their dependencies,
// and their state.
func Modules() ([]Module, error) {
	f, err := os.Open("/proc/modules")
	if err != nil {
		return nil, err
	}

	var modules []Module

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sep := strings.Split(scanner.Text(), " ")
		memSize, err := strconv.ParseInt(sep[1], 10, 64)
		if err != nil {
			return nil, err
		}

		instancesLoaded, err := strconv.ParseInt(sep[2], 10, 64)
		if err != nil {
			return nil, err
		}

		deps := strings.Split(sep[3], ",")
		for i := range deps {
			deps[i] = strings.TrimSpace(deps[i])
		}

		modules = append(modules, Module{
			Name:            sep[0],
			MemSize:         memSize,
			InstancesLoaded: instancesLoaded,
			Dependencies:    deps[:len(deps)-1],
			State:           sep[4],
			MemOffset:       sep[5],
		})
	}

	return modules, nil
}
