package payloads

type Withdraw struct {
	OrderNumber string  `json:"order" valid:"required,type(string)"`
	Sum         float64 `json:"sum" valid:"required,type(float64)"`
}
