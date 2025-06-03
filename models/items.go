package models

type Items struct {
	ID          int     `json:"id"`
	Cod         string  `json:"cod"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	CreatedBy   int     `json:"created_by"`
	UpdatedInfo string  `json:"updated_info"`
}
