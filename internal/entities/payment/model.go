package payment

type PaymentInput struct {
	UserEmail string  `json:"user_email"`
	Currency  string  `json:"currency"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type PaymentResonse struct {
	Status string `json:"status,omitempty"`
	ID     int64  `json:"id,omitempty"`
}

type PaymentUserID int64
type PaymentUserEmail string
