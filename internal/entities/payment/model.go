package payment

type PaymentInput struct {
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	UserEmail string  `json:"user_email"`
	Currency  string  `json:"currency"`
}

type PaymentUser struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type PaymentData struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Amount    float64 `json:"amount" db:"amount"`
	UserEmail string  `json:"user_email" db:"user_email"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UreatedAt string  `json:"updated_at" db:"updated_at"`
	Status    string  `json:"status" db:"status"`
}

type PaymentsResonse struct {
	Data []PaymentData `json:"data"`
}

type PaymentIDResponse struct {
	ID int64 `json:"id"`
}

type PaymentStatusResponse struct {
	Status string `json:"status"`
}

type PaymentStatus struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}
