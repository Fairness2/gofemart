package payloads

// Withdraw представляет собой запрос на вывод средств с номером заказа и суммой.
// OrderNumber — обязательное строковое поле, представляющее уникальный идентификатор заказа.
// Sum — обязательное поле float64, представляющее сумму для вывода.
type Withdraw struct {
	OrderNumber string  `json:"order" valid:"required,type(string)"`
	Sum         float64 `json:"sum" valid:"required,type(float64)"`
}
