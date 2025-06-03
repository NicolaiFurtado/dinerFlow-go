package models

type Payment struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	ClientName  string  `json:"client_name"`
	TotalPrice  float64 `json:"total_price"`
	ReceiptData string  `json:"receipt_data"`
	Status      string  `json:"status"`
	TypePayment string  `json:"type_payment"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   string  `json:"created_by"`
	ClosedAt    string  `json:"closed_at"`
	ClosedBy    string  `json:"closed_by"`
}
