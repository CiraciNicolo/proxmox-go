package proxmox

import (
	"context"
	"errors"
	"github.com/sp-yduck/proxmox-go/api"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

const (
	TaskStatusOK = "OK"
)

func (s *Service) EnsureTaskDone(ctx context.Context, node, upid string) error {
	checkVMCompleted := func() (bool, error) {
		task, err := s.restclient.GetTask(ctx, node, upid)
		if err != nil {
			return false, err
		}
		if task.Status == api.TaskStatusRunning {
			return false, nil
		}
		return task.Exitstatus == TaskStatusOK, nil
	}

	backoff := wait.Backoff{
		Duration: time.Second,
		Factor:   1.2,
		Jitter:   0,
		Steps:    32,
		Cap:      16 * time.Minute,
	}

	err := wait.ExponentialBackoff(backoff, checkVMCompleted)
	if err != nil {
		return errors.Join(err, errors.New("task wait deadline exceeded"))
	}

	return nil
}
