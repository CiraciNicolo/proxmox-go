package proxmox

import (
	"context"
	"errors"
	"github.com/sp-yduck/proxmox-go/api"
	"io"
)

type Storage struct {
	Service
	Storage *api.Storage
	Node    string
}

func (s *Service) Storage(ctx context.Context, name string) (*Storage, error) {
	storage, err := s.restclient.GetStorage(ctx, name)
	if err != nil {
		return nil, err
	}
	return &Storage{Service: *s, Storage: storage}, nil
}

func (s *Service) CreateStorage(ctx context.Context, name, storageType string, options api.StorageCreateOptions) (*Storage, error) {
	storage, err := s.restclient.CreateStorage(ctx, name, storageType, options)
	if err != nil {
		return nil, err
	}
	return &Storage{Service: *s, Storage: storage}, nil
}

func (s *Storage) Delete(ctx context.Context) error {
	return s.restclient.DeleteStorage(ctx, s.Storage.Storage)
}

func (s *Storage) GetContents(ctx context.Context) ([]*api.StorageContent, error) {
	err := ensureStorage(s)
	if err != nil {
		return nil, err
	}

	return s.restclient.GetContents(ctx, s.Node, s.Storage.Storage)
}

func (s *Storage) GetContent(ctx context.Context, volumeID string) (*api.StorageContent, error) {
	err := ensureStorage(s)
	if err != nil {
		return nil, err
	}

	return s.restclient.GetContent(ctx, s.Node, s.Storage.Storage, volumeID)
}

func (s *Storage) GetVolume(ctx context.Context, volumeID string) (*api.StorageVolume, error) {
	err := ensureStorage(s)
	if err != nil {
		return nil, err
	}

	return s.restclient.GetVolume(ctx, s.Node, s.Storage.Storage, volumeID)
}

func (s *Storage) DeleteVolume(ctx context.Context, volumeID string) error {
	err := ensureStorage(s)
	if err != nil {
		return err
	}
	var taskid *string

	if taskid, err = s.restclient.DeleteVolume(ctx, s.Node, s.Storage.Storage, volumeID); err != nil {
		return err
	}
	if err = s.EnsureTaskDone(ctx, s.Node, *taskid); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Upload(ctx context.Context, option api.StorageUpload, file io.Reader) error {
	option.Node = s.Node
	option.Storage = s.Storage.Storage

	if err := s.restclient.UploadToStorage(ctx, option, file); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Download(ctx context.Context, option api.StorageDownload) error {
	option.Node = s.Node
	option.Storage = s.Storage.Storage
	var taskid *string
	var err error

	if taskid, err = s.restclient.DownloadToStorage(ctx, option); err != nil {
		return err
	}
	if err = s.EnsureTaskDone(ctx, s.Node, *taskid); err != nil {
		return err
	}

	return nil
}

func ensureStorage(s *Storage) error {
	if s.Node == "" {
		return errors.New("Node must not be empty")
	}
	if s.Storage.Storage == "" {
		return errors.New("Storage not specified")
	}
	return nil
}
