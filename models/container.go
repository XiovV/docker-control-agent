package models

type ContainerStatus struct {
	Name string `json:"name"`
	ID string `json:"id"`
	Image string `json:"image"`
	Status string `json:"status"`
}
