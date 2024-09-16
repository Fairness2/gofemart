package payloads

// Accrual представляет собой структуру ответа системы по начислению.
// Order — идентификатор заказа в системе.
// Status указывает на текущий статус заказа в системе начисления.
// Accrual представляет собой начисленное значение для заказа.
type Accrual struct {
	Order   string  `json:"order" valid:"required,type(string)"`
	Status  string  `json:"status" valid:"required,type(string)"`
	Accrual float64 `json:"accrual" valid:"required,type(float64)"`
}

const (
	StatusAccrualRegistered = "REGISTERED"
	StatusAccrualInvalid    = "INVALID"
	StatusAccrualProcessing = "PROCESSING"
	StatusAccrualProcessed  = "PROCESSED"
)
