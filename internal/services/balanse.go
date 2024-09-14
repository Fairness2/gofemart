package services

import (
	"context"
	"database/sql"
	"errors"
	database "gofemart/internal/databse"
	"gofemart/internal/models"
	"gofemart/internal/repositories"
	"time"
)

var ErrorNotEnoughItems = errors.New("there are not enough resources")

// BalanceService безопасный сервис для списания средств
type BalanceService struct {
	ctx context.Context
}

// NewBalanceService получение нового сервиса трат
func NewBalanceService(ctx context.Context) *BalanceService {
	return &BalanceService{ctx: ctx}
}

// Spend списываем средства со счёта
func (s *BalanceService) Spend(user *models.User, sum float64, order *models.Order) error {
	userMutex, exists := UserMutexInstance.GetMutex(user.ID)
	if !exists {
		userMutex = UserMutexInstance.SetMutex(user.ID)
	}
	userMutex.Lock()
	defer func() {
		userMutex.Unlock()
		UserMutexInstance.DeleteMutex(user.ID)
	}()

	rep := s.getAccountRepository()
	balanceSum, err := rep.GetSum(user.ID)
	if err != nil {
		return err
	}

	if balanceSum < sum {
		return ErrorNotEnoughItems
	}

	newAcc := models.Account{
		UserID:     user.ID,
		Difference: -sum,
		OrderNumber: sql.NullString{
			String: order.Number,
			Valid:  true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return rep.CreateAccount(&newAcc)
}

// getAccountRepository создаём репозиторий для начислений
func (s *BalanceService) getAccountRepository() *repositories.AccountRepository {
	return repositories.NewAccountRepository(s.ctx, database.DBx)
}
