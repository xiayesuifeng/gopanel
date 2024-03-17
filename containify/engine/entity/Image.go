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
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`
	Env          []string            `json:"Env,omitempty"`
	Entrypoint   []string            `json:"Entrypoint,omitempty"`
	Cmd          []string            `json:"Cmd,omitempty"`
	Volumes      map[string]struct{} `json:"Volumes,omitempty"`
	WorkingDir   string              `json:"WorkingDir,omitempty"`
	Labels       map[string]string   `json:"Labels,omitempty"`
}
