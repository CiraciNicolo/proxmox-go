package rest

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sp-yduck/proxmox-go/api"
)

func (s *TestSuite) TestGetStorages() {
	storages, err := s.restclient.GetStorages(context.TODO())
	if err != nil {
		s.T().Fatalf("failed to get storages: %v", err)
	}
	s.T().Logf("get storages: %v", storages)
}

func (s *TestSuite) GetTestStorage() *api.Storage {
	storages, err := s.restclient.GetStorages(context.TODO())
	if err != nil {
		s.T().Fatalf("failed to get storages: %v", err)
	}
	return storages[0]
}

func (s *TestSuite) TestGetStorage() {
	testStorageName := s.GetTestStorage().Storage

	storage, err := s.restclient.GetStorage(context.TODO(), testStorageName)
	if err != nil {
		s.T().Fatalf("failed to get storage(name=%s): %v", testStorageName, err)
	}
	s.T().Logf("get storage: %v", *storage)
}

func (s *TestSuite) EnsureNoStorage(name string) {
	storage, err := s.restclient.GetStorage(context.TODO(), name)
	if err == nil {
		s.T().Logf("error: %v", err)
		if err := s.restclient.DeleteStorage(context.TODO(), storage.Storage); err != nil {
			s.T().Fatalf("failed to ensure no storage (name=%s): %v", storage.Storage, err)
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil && !IsNotFound(err) {
		s.T().Logf("failed to get storage(name=%s): %v", name, err)
	}
}

func (s *TestSuite) TestCreateDeleteStorage() {
	testStorageName := "test-proxmox-go"
	s.EnsureNoStorage(testStorageName)

	// create
	mkdir := true
	testOptions := api.StorageCreateOptions{
		Content: "images",
		Mkdir:   &mkdir,
		Path:    "/var/lib/vz/test",
	}
	storage, err := s.restclient.CreateStorage(context.TODO(), testStorageName, "dir", testOptions)
	if err != nil {
		s.T().Fatalf("failed to create storage(name=%s): %v", testStorageName, err)
	}
	s.T().Logf("create storage: %v", *storage)
	time.Sleep(2 * time.Second)

	// delete
	err = s.restclient.DeleteStorage(context.TODO(), testStorageName)
	if err != nil {
		s.T().Fatalf("failed to delete storage(name=%s): %v", testStorageName, err)
	}
}

func (s *TestSuite) TestCreateUploadDeleteStorage() {
	testStorageName := "test-upload-proxmox-go"
	s.EnsureNoStorage(testStorageName)

	// create
	mkdir := true
	testOptions := api.StorageCreateOptions{
		Content: "images,iso",
		Mkdir:   &mkdir,
		Path:    "/var/lib/vz/test",
	}
	storage, err := s.restclient.CreateStorage(context.TODO(), testStorageName, "dir", testOptions)
	if err != nil {
		s.T().Fatalf("failed to create storage(name=%s): %v", testStorageName, err)
	}
	s.T().Logf("create storage: %v", *storage)
	time.Sleep(2 * time.Second)

	if err != nil {
		s.T().Fatalf("failed to get nodes: %v", err)
	}

	node, err := s.restclient.GetLocalNode(context.TODO())
	uploadOptions := api.StorageUpload{
		Content:  "iso",
		Filename: "tlc.iso",
		Node:     node.Name,
		Storage:  testStorageName,
	}

	volumeID := fmt.Sprintf("%s:iso/%s", testStorageName, uploadOptions.Filename)
	f, err := os.Open("../tlc.iso")
	if err != nil {
		s.T().Fatalf("failed to open file(name=%s): %v", testStorageName, err)
	}

	err = s.restclient.UploadToStorage(context.TODO(), uploadOptions, f)
	if err != nil {
		s.T().Fatalf("failed to upload storage(name=%s): %v", testStorageName, err)
	}
	time.Sleep(2 * time.Second)

	_, err = s.restclient.DeleteVolume(context.TODO(), node.Name, testStorageName, volumeID)
	time.Sleep(2 * time.Second)

	err = s.restclient.DeleteStorage(context.TODO(), testStorageName)
	if err != nil {
		s.T().Fatalf("failed to delete storage(name=%s): %v", testStorageName, err)
	}
}

func (s *TestSuite) TestCreateDownloadDeleteStorage() {
	testStorageName := "test-upload-proxmox-go"
	s.EnsureNoStorage(testStorageName)

	// create
	mkdir := true
	testOptions := api.StorageCreateOptions{
		Content: "images,iso",
		Mkdir:   &mkdir,
		Path:    "/var/lib/vz/test",
	}
	storage, err := s.restclient.CreateStorage(context.TODO(), testStorageName, "dir", testOptions)
	if err != nil {
		s.T().Fatalf("failed to create storage(name=%s): %v", testStorageName, err)
	}
	s.T().Logf("create storage: %v", *storage)
	time.Sleep(2 * time.Second)

	if err != nil {
		s.T().Fatalf("failed to get nodes: %v", err)
	}

	node, err := s.restclient.GetLocalNode(context.TODO())
	uploadOptions := api.StorageDownload{
		Content:  "iso",
		Filename: "tlc.iso",
		Node:     node.Name,
		Storage:  testStorageName,
		Url:      "http://tinycorelinux.net/14.x/x86/release/Core-current.iso",
	}

	volumeID := fmt.Sprintf("%s:iso/%s", testStorageName, uploadOptions.Filename)
	_, err = s.restclient.DownloadToStorage(context.TODO(), uploadOptions)
	if err != nil {
		s.T().Fatalf("failed to download to storage(name=%s): %v", testStorageName, err)
	}
	time.Sleep(30 * time.Second)

	_, err = s.restclient.DeleteVolume(context.TODO(), node.Name, testStorageName, volumeID)
	time.Sleep(2 * time.Second)

	err = s.restclient.DeleteStorage(context.TODO(), testStorageName)
	if err != nil {
		s.T().Fatalf("failed to delete storage(name=%s): %v", testStorageName, err)
	}
}
