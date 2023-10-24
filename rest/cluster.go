package rest

import (
	"context"
	"encoding/json"
	"github.com/sp-yduck/proxmox-go/api"
)

func (c *RESTClient) GetNextID(ctx context.Context) (int, error) {
	var res json.Number
	if err := c.Get(ctx, "/cluster/nextid", &res); err != nil {
		return 0, err
	}
	nextid, err := res.Int64()
	if err != nil {
		return 0, err
	}
	return int(nextid), nil
}

func (c *RESTClient) GetStatus(ctx context.Context) ([]*api.ClusterStatus, error) {
	var res []*api.ClusterStatus
	if err := c.Get(ctx, "/cluster/status", &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *RESTClient) GetLocalNode(ctx context.Context) (local *api.ClusterStatus, err error) {
	nodes, err := c.GetStatus(ctx)
	if err != nil {
		return
	}

	for _, node := range nodes {
		if node.Type == "node" && node.Local == 1 {
			local = node
		}
	}

	return
}

func (c *RESTClient) GetJoinConfig(ctx context.Context) (*api.ClusterJoinConfig, error) {
	var config *api.ClusterJoinConfig
	if err := c.Get(ctx, "/cluster/config/join", &config); err != nil {
		return nil, err
	}
	return config, nil
}
