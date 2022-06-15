package entity

const (
	StatusNew      = "new"
	StatusError    = "error"
	StatusSuccess  = "success"
	StatusFailure  = "failure"
	StatusCanceled = "canceled"
)

type Payment struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Amount    float64 `json:"amount" db:"amount"` // TODO: выделить 2 сущности под разбиение amount на целочисленную часть и дробную
	UserEmail string  `json:"user_email" db:"user_email"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
	Status    string  `json:"status" db:"status"`
}

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

type PaymentsData struct {
	Data []Payment `json:"data"`
}

type PaymentStatus struct {
	ID     int64  `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}
