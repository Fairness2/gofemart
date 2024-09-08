package payloads

type Accrual struct {
	Order   string `json:"order" valid:"required,type(string)"`
	Status  string `json:"status" valid:"required,type(string)"`
	Accrual int    `json:"accrual" valid:"required,type(int)"`
}

const (
	StatusAccrualRegistered = "REGISTERED"
	StatusAccrualInvalid    = "INVALID"
	StatusAccrualProcessing = "PROCESSING"
	StatusAccrualProcessed  = "PROCESSED"
)
