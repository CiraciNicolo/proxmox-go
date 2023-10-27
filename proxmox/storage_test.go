package proxmox

import (
	"context"

	"github.com/sp-yduck/proxmox-go/api"
)

func (s *TestSuite) getTestStorage() (*api.Node, *Storage) {
	node := s.getTestNode()
	storage, err := s.service.Storage(context.TODO(), "local")
	if err != nil {
		s.T().Fatalf("failed to get tasks: %v", err)
	}
	return node, storage
}

func (s *TestSuite) TestUpload() {
	testNode, testStorage := s.getTestStorage()
	testStorage.Node = testNode.Node

	options := api.StorageDownload{
		Content:  "iso",
		Filename: "ubuntu.iso",
		Url:      "https://releases.ubuntu.com/jammy/ubuntu-22.04.3-live-server-amd64.iso",
	}
	err := testStorage.Download(context.TODO(), options)
	if err != nil {
		s.T().Fatalf("failed to get task: %v", err)
	}
}
