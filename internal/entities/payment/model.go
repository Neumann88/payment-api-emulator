package payment

type Payment struct {
	ID                int64   `json:"id"`
	UserID            int64   `json:"user_id"`
	UserEmail         string  `json:"user_email"`
	Amount            float64 `json:"amout"`
	Currency          string  `json:"currency"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	TransactionStatus string  `json:"transaction_status"`
}
