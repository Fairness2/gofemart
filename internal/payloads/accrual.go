package payloads

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
