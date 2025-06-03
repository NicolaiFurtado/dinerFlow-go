package models

type RemoveOrderRequest struct {
	ID    int `json:"id"`
	Order struct {
		Remove []int `json:"remove"`
	} `json:"order"`
}

type Tab struct {
	ID          int       `json:"id"`
	TableId     int       `json:"table_id"`
	ClientName  string    `json:"client_name"`
	Order       OrderData `json:"order"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedInfo string    `json:"updated_info"`
}

type OrderItem struct {
	ItemCod string `json:"item_cod"`
	Qty     int    `json:"qtd"`
	Notes   string `json:"notes"`
}

type OrderData struct {
	Items []OrderItem `json:"items"`
}
