package repositories

const (
	createAccountSQL      = "INSERT INTO t_account (user_id, difference, order_number, created_at, updated_at) VALUES (:user_id, :difference, :order_number, :created_at, :updated_at) RETURNING id"
	getSumSQL             = "SELECT COALESCE(SUM(difference), 0) FROM t_account WHERE user_id = $1"
	getBalanceSQL         = "SELECT COALESCE(sum(difference), 0) current, COALESCE(sum(CASE WHEN difference < 0 THEN abs(difference) ELSE 0 END), 0) withdrawn FROM t_account WHERE user_id = $1"
	getWithdrawByOrderSQL = "SELECT * FROM t_account WHERE order_number = $1 AND difference < 0"
)
