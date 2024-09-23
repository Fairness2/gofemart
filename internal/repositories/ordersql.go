package repositories

const (
	createOrderSQL                                       = "INSERT INTO t_order (number, user_id, status_code) VALUES (:number, :user_id, :status_code)"
	updateOrderSQL                                       = "UPDATE t_order SET user_id = :user_id, status_code = :status_code, last_checked_at = :last_checked_at, updated_at = :updated_at WHERE number = :number"
	getOrdersExcludeOrdersWhereStatusInWithNumbersSQL    = "SELECT * FROM t_order WHERE status_code IN (?) AND number NOT IN (?) AND ((last_checked_at NOTNULL AND last_checked_at <= ?) OR (last_checked_at IS NULL AND created_at <= ?)) LIMIT ?"
	getOrdersExcludeOrdersWhereStatusInWithoutNumbersSQL = "SELECT * FROM t_order WHERE status_code IN (?) AND ((last_checked_at NOTNULL AND last_checked_at <= ?) OR (last_checked_at IS NULL AND created_at <= ?)) LIMIT ?"
	getOrderByNumberSQL                                  = "SELECT * FROM t_order WHERE number = $1"
	getOrdersByUserWithAccrualSQL                        = "SELECT t.*, CASE WHEN ta.difference NOTNULL THEN difference ELSE 0 END accrual FROM t_order t LEFT JOIN t_account ta ON t.number = ta.order_number AND ta.difference > 0 WHERE t.user_id = $1"
	getOrdersByUserWithdrawSQL                           = "SELECT ta.order_number number, abs(ta.difference) accrual, ta.created_at processed_at FROM  public.t_account ta WHERE ta.user_id = $1 AND ta.difference < 0 AND ta.order_number NOTNULL"
)
