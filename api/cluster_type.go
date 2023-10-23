package api

type ClusterStatus struct {
	Id    string `json:"id"`
	Local int    `json:"local"`
	Type  string `json:"type"`
	Name  string `json:"name"`
}
