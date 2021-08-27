package models

type NodeStatusRequest struct {
	NodeName string `json:"node_name"`
	Containers []string `json:"containers"`
}

type NodeStatusResponse struct {
	IsOnline   bool        `json:"is_online"`
	Containers []ContainerStatus `json:"containers"`
}