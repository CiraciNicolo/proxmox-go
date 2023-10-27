package rest

import (
	"context"
	"fmt"

	"github.com/sp-yduck/proxmox-go/api"
)

func (c *RESTClient) GetVirtualMachines(ctx context.Context, node string) ([]*api.VirtualMachine, error) {
	path := fmt.Sprintf("/nodes/%s/qemu", node)
	var vms []*api.VirtualMachine
	if err := c.Get(ctx, path, &vms); err != nil {
		return nil, err
	}
	return vms, nil
}

func (c *RESTClient) GetVirtualMachine(ctx context.Context, node string, vmid int) (*api.VirtualMachine, error) {
	vms, err := c.GetVirtualMachines(ctx, node)
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if vm.VMID == vmid {
			return vm, nil
		}
	}
	return nil, NotFoundErr
}

func (c *RESTClient) CreateVirtualMachine(ctx context.Context, node string, vmid int, options api.VirtualMachineCreateOptions) (*string, error) {
	options.VMID = &vmid
	path := fmt.Sprintf("/nodes/%s/qemu", node)
	var upid *string
	if err := c.Post(ctx, path, options, nil, &upid); err != nil {
		return nil, err
	}
	return upid, nil
}

func (c *RESTClient) DeleteVirtualMachine(ctx context.Context, node string, vmid int) (*string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d", node, vmid)
	var upid *string
	if err := c.Delete(ctx, path, nil, upid); err != nil {
		return nil, err
	}
	return upid, nil
}

func (c *RESTClient) GetVirtualMachineConfig(ctx context.Context, node string, vmid int) (*api.VirtualMachineConfig, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid)
	var config *api.VirtualMachineConfig
	if err := c.Get(ctx, path, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *RESTClient) SetVirtualMachineConfig(ctx context.Context, node string, vmid int, config api.VirtualMachineConfig) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid)
	config.VMGenID = "" // Value should not be set
	if err := c.Put(ctx, path, &config, nil); err != nil {
		return err
	}
	return nil
}

func (c *RESTClient) GetVirtualMachineStatus(ctx context.Context, node string, vmid int) (*api.VirtualMachineStatus, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/current", node, vmid)
	var status *api.VirtualMachineStatus
	if err := c.Get(ctx, path, &status); err != nil {
		return nil, err
	}
	return status, nil
}
