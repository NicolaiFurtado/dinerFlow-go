package models

type Tables struct {
	ID          int    `json:"id"`
	TableName   string `json:"table_name"`
	Seats       int    `json:"seats"`
	Status      string `json:"status"`
	CreatedBy   int    `json:"created_by"`
	UpdatedInfo string `json:"updated_info"`
}
