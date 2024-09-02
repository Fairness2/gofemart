package payloads

type Register struct {
	Login    string `json:"login" valid:"required,type(string),stringlength(3)"`
	Password string `json:"password" valid:"required,type(string),stringlength(6)"`
}

type ErrorResponseBody struct {
	Status  int    `json:"status"` // Успешный или не успешный результат
	Message string `json:"message,omitempty"`
}

type Authorization struct {
	Token string `json:"token"`
}
