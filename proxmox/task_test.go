package proxmox

import (
	"context"

	"github.com/sp-yduck/proxmox-go/api"
)

func (s *TestSuite) getTestTask() (*api.Node, *api.Tasks) {
	node := s.getTestNode()
	tasks, err := s.service.restclient.GetTasks(context.TODO(), node.Node)
	if err != nil {
		s.T().Fatalf("failed to get tasks: %v", err)
	}
	return node, tasks[0]
}

func (s *TestSuite) TestEnsureTaskDone() {
	testNode, testTask := s.getTestTask()
	err := s.service.EnsureTaskDone(context.TODO(), testNode.Node, testTask.UPID)
	if err != nil {
		s.T().Fatalf("failed to get task: %v", err)
	}
}
