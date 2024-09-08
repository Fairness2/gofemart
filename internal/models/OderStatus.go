package models

type OderStatus struct {
	Code        string `db:"code"`
	Description string `db:"description"`
}

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)
