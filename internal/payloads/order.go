package payloads

type Withdraw struct {
	OrderNumber string `json:"order" valid:"required,type(string)"`
	Sum         int    `json:"sum" valid:"required,type(int)"`
}
