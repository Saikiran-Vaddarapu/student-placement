package models

type Student struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Branch  string  `json:"branch"`
	DOB     string  `json:"dob"`
	Status  string  `json:"status"`
	Company Company `json:"company"`
}
