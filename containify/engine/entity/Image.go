package entity

type Image struct {
	ID          string            `json:"id"`
	ParentID    string            `json:"parentID"`
	RepoTags    []string          `json:"repoTags"`
	RepoDigests []string          `json:"repoDigests"`
	Created     int64             `json:"created"`
	Size        int64             `json:"size"`
	SharedSize  int               `json:"sharedSize"`
	VirtualSize int64             `json:"virtualSize"`
	Labels      map[string]string `json:"labels"`
	Containers  int               `json:"containers"`
}

type InspectImage struct {
	ExposedPorts map[string]struct{} `json:"exposedPorts"`
	Env          []string            `json:"env"`
	Entrypoint   []string            `json:"entrypoint"`
	Cmd          []string            `json:"cmd"`
	Volumes      map[string]struct{} `json:"volumes"`
	WorkingDir   string              `json:"workingDir"`
	Labels       map[string]string   `json:"labels"`
}
