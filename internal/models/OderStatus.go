package models

// OderStatus статус заказа
type OderStatus struct {
	Code        string `db:"code"`
	Description string `db:"description"`
}

const (
	StatusNew        = "NEW"        // Новый заказ
	StatusProcessing = "PROCESSING" // Заказ обрабатывается
	StatusInvalid    = "INVALID"    // Заказу отказано в начислении
	StatusProcessed  = "PROCESSED"  // Заказ обработан, и ему начислены баллы
)
