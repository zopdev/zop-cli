package service

type DeploymentSpaceOptions struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type DeploymentOption struct {
	Option []map[string]any `json:"options"`
	Next   *NextPage        `json:"nextPage"`
}

type NextPage struct {
	Name   string            `json:"name"`
	Path   string            `json:"path"`
	Params map[string]string `json:"params"`
}

type apiResponse struct {
	Data *DeploymentOption `json:"data"`
}
