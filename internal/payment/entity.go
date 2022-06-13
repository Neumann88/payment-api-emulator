package payment

type payment struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Amount    float64 `json:"amount" db:"amount"`
	UserEmail string  `json:"user_email" db:"user_email"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
	Status    string  `json:"status" db:"status"`
}

type paymentInput struct {
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	UserEmail string  `json:"user_email"`
	Currency  string  `json:"currency"`
}

type paymentUser struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type paymentsData struct {
	Data []payment `json:"data"`
}

type paymentStatus struct {
	ID     int64  `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}
