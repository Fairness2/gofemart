package ordercheck

import (
	"database/sql"
	"errors"
	"gofemart/internal/accrual"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/payloads"
	"time"
)

// processOder обрабатываем заказ, запрашивая информацию у внешней системы
func (p *Pool) processOder(number string) {
	logger.Log.Infow("Process order", "number", number)
	order, ok := p.poolInWork(number)
	// Если заказ не найден, то пропускаем его
	if order == nil || !ok {
		return
	}
	defer p.removeFromWork(number)
	accrualResponse, err := p.accrualProxy.Accrual(order.model)
	if err != nil {
		logger.Log.Error(err)
		var tmrErr accrual.TooManyRequestError
		if ok := errors.Is(err, &tmrErr); ok {
			p.accrualProxy.Pause(tmrErr.PauseDuration)
		}
		return
	}
	if err := p.processOrderAccrual(accrualResponse, order.model); err != nil {
		logger.Log.Error(err)
	}
}

// processOrderAccrual обрабатываем ответ системы начислений, обновляем заказ и создаём запись в счёте пользователя
func (p *Pool) processOrderAccrual(accrual *payloads.Accrual, order *models.Order) error {
	logger.Log.Infow("Process order accrual", "order", order.Number, "status", accrual.Status)
	order.LastCheckedAt = sql.NullTime{Time: time.Now(), Valid: true}
	switch accrual.Status {
	case payloads.StatusAccrualProcessing, payloads.StatusAccrualRegistered:
		order.StatusCode = models.StatusProcessing
	case payloads.StatusAccrualInvalid:
		order.StatusCode = models.StatusInvalid
	case payloads.StatusAccrualProcessed:
		if _, err := p.createNewAccount(order.Number, order.UserID, accrual.Accrual); err != nil {
			return err
		}
		order.StatusCode = models.StatusProcessed
	}
	orderRep := p.orderRepo
	return orderRep.UpdateOrder(order)
}

// createNewAccount создаём новую запись о начислении
func (p *Pool) createNewAccount(orderNumber string, userID int64, diff float64) (*models.Account, error) {
	logger.Log.Infow("Create new account", "orderNumber", orderNumber, "userID", userID, "diff", diff)
	repository := p.accountRepo
	account := models.NewAccount(sql.NullString{String: orderNumber, Valid: true}, userID, diff)
	if err := repository.CreateAccount(account); err != nil {
		return nil, err
	}
	return account, nil
}
