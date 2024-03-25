package models

type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}
