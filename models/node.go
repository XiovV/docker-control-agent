package models

type NodeStatusRequest struct {
	Container string `json:"container"`
}

type NodeStatusResponse struct {
	ContainerStatus
}