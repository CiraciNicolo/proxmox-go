package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sp-yduck/proxmox-go/api"
	"io"
	"mime/multipart"
	"os"
)

func (c *RESTClient) GetStorages(ctx context.Context) ([]*api.Storage, error) {
	var storages []*api.Storage
	if err := c.Get(ctx, "/storage", &storages); err != nil {
		return nil, err
	}
	return storages, nil
}

func (c *RESTClient) GetStorage(ctx context.Context, name string) (*api.Storage, error) {
	storages, err := c.GetStorages(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range storages {
		if s.Storage == name {
			return s, nil
		}
	}
	return nil, NotFoundErr
}

func (c *RESTClient) CreateStorage(ctx context.Context, name, storageType string, options api.StorageCreateOptions) (*api.Storage, error) {
	options.Storage = name
	options.StorageType = storageType
	var storage *api.Storage
	if err := c.Post(ctx, "/storage", options, nil, &storage); err != nil {
		return nil, err
	}
	return storage, nil
}

func (c *RESTClient) DeleteStorage(ctx context.Context, name string) error {
	path := fmt.Sprintf("/storage/%s", name)
	if err := c.Delete(ctx, path, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *RESTClient) GetContents(ctx context.Context, node string, storage string) ([]*api.StorageContent, error) {
	var contents []*api.StorageContent
	path := fmt.Sprintf("/nodes/%s/storage/%s/content", node, storage)
	if err := c.Get(ctx, path, &contents); err != nil {
		return nil, err
	}
	return contents, nil
}

func (c *RESTClient) GetContent(ctx context.Context, node string, storage string, volumeID string) (*api.StorageContent, error) {
	contents, err := c.GetContents(ctx, node, storage)
	if err != nil {
		return nil, err
	}
	for _, content := range contents {
		if content.VolID == volumeID {
			return content, nil
		}
	}
	return nil, NotFoundErr
}

func (c *RESTClient) GetVolume(ctx context.Context, node string, storage string, volumeID string) (*api.StorageVolume, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/content/%s", node, storage, volumeID)
	var volume *api.StorageVolume
	if err := c.Get(ctx, path, &volume); err != nil {
		return nil, err
	}
	return volume, nil
}

func (c *RESTClient) DeleteVolume(ctx context.Context, node string, storage string, volumeID string) (*string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/content/%s", node, storage, volumeID)
	var upid *string
	if err := c.Delete(ctx, path, nil, &upid); err != nil {
		return nil, err
	}
	return upid, nil
}

// UploadToStorage TODO: Add other parameters such as checksum
func (c *RESTClient) UploadToStorage(ctx context.Context, options api.StorageUpload, file io.Reader) error {
	var buf bytes.Buffer
	var fileSize int64
	body := Body{}

	if f, s := file.(*os.File); s {
		fs, err := f.Stat()
		if err != nil {
			return errors.Wrap(err, "unable to get file info")
		}
		fileSize = fs.Size()
		body.StreamingUpload = true
	}

	writer := multipart.NewWriter(&buf)
	err := writer.WriteField("content", options.Content)
	if err != nil {
		return errors.Wrap(err, "unable to set content type")
	}

	fw, err := writer.CreateFormFile("filename", options.Filename)
	if err != nil {
		return errors.Wrap(err, "unable to set filename")
	}

	if body.StreamingUpload {
		headerSize := buf.Len()

		err = writer.Close()
		if err != nil {
			return errors.Wrap(err, "unable to close writer")
		}

		body.ContentLength = int64(buf.Len()) + fileSize
		body.Reader = io.MultiReader(
			bytes.NewReader(buf.Bytes()[:headerSize]),
			file,
			bytes.NewReader(buf.Bytes()[headerSize:]),
		)

	} else {
		_, err = io.Copy(fw, file)
		if err != nil {
			return err
		}

		err = writer.Close()
		if err != nil {
			return errors.Wrap(err, "unable to close writer")
		}

		body.Reader = &buf
	}

	body.ContentType = writer.FormDataContentType()
	path := fmt.Sprintf("/nodes/%s/storage/%s/upload", options.Node, options.Storage)
	if err := c.Post(ctx, path, options, &body, nil); err != nil {
		return err
	}
	return nil
}

func (c *RESTClient) DownloadToStorage(ctx context.Context, options api.StorageDownload) (*string, error) {
	path := fmt.Sprintf("/nodes/%s/storage/%s/download-url", options.Node, options.Storage)
	var upid *string
	if err := c.Post(ctx, path, options, nil, &upid); err != nil {
		return nil, err
	}
	return upid, nil
}
