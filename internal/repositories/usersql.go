package repositories

const (
	getUserByIDSQL    = "SELECT id, login, password_hash FROM t_user WHERE id = $1"
	getUserByLoginSQL = "SELECT id, login, password_hash FROM t_user WHERE login = $1"
	createUserSQL     = "INSERT INTO t_user (login, password_hash) VALUES (:login, :password_hash) RETURNING id"
	userExistsSQL     = "SELECT true FROM t_user WHERE login = $1"
)
