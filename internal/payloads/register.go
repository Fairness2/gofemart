package payloads

// Register пэйлоад для регистрации пользователя, в данный момент и для авторизации)
type Register struct {
	Login    string `json:"login" valid:"required,type(string),minstringlength(3)"`
	Password string `json:"password" valid:"required,type(string),minstringlength(6)"`
}

// Authorization ответ с токеном авторизации
type Authorization struct {
	Token string `json:"token"`
}
