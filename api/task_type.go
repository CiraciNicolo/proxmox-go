package api

type Task struct {
	Id         string     `json:"id,omitempty"`
	Node       string     `json:"node,omitempty"`
	PID        int        `json:"pid"`
	StartTime  int        `json:"starttime"`
	Status     TaskStatus `json:"status"`
	Type       string     `json:"type"`
	UPID       string     `json:"upid"`
	User       string     `json:"user"`
	Exitstatus string     `json:"exitstatus"`
}

type Tasks struct {
	Endtime   int    `json:"endtime"`
	PID       int    `json:"pid"`
	PStart    int    `json:"pstart"`
	StartTime int    `json:"starttime"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	UPID      string `json:"upid"`
	User      string `json:"user"`
}

type TaskStatus string

const (
	TaskStatusRunning TaskStatus = "running"
	TaskStatusStopped TaskStatus = "stopped"
)
