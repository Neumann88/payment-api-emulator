package payment

type PaymentInput struct {
	UserEmail string  `json:"user_email"`
	Currency  string  `json:"currency"`
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type PaymentIDResponse struct {
	ID int64 `json:"id"`
}

type PaymentStatusResponse struct {
	Status string `json:"status"`
}

type PaymentStatus struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}
