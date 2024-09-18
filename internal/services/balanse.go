package services

import (
	"context"
	"database/sql"
	"errors"
	"gofemart/internal/logger"
	"gofemart/internal/models"
	"gofemart/internal/repositories"
	"sync"
	"time"
)

// ErrorNotEnoughItems Ошибка, что не счету пользователя недостаточно ресурсов
var ErrorNotEnoughItems = errors.New("there are not enough resources")

// BalanceRepository интерфейс для репозитория для работы с балансом пользователя
type BalanceRepository interface {
	GetSum(userID int64) (float64, error)
	CreateAccount(account *models.Account) error
}

// MutexService интерфейс сервиса для работы с мьютексами пользователя
type MutexService interface {
	SetMutex(userID int64) *sync.Mutex
	GetMutex(userID int64) (*sync.Mutex, bool)
	DeleteMutex(userID int64) error
}

// BalanceService безопасный сервис для списания средств
type BalanceService struct {
	ctx        context.Context
	repository BalanceRepository
	userMutex  MutexService
}

// NewBalanceService получение нового сервиса трат
func NewBalanceService(ctx context.Context, dbPool repositories.SQLExecutor) *BalanceService {
	logger.Log.Debug("NewBalanceService")
	return &BalanceService{
		ctx:        ctx,
		repository: getAccountRepository(ctx, dbPool),
		userMutex:  GetUserMutexInstance(),
	}
}

// Spend списываем средства со счёта
func (s *BalanceService) Spend(user *models.User, sum float64, order *models.Order) error {
	logger.Log.Debugw("Spend", "user", user.ID, "sum", sum, "order", order.Number)
	userMutex, exists := s.userMutex.GetMutex(user.ID)
	if !exists {
		userMutex = s.userMutex.SetMutex(user.ID)
	}
	userMutex.Lock()
	defer func() {
		userMutex.Unlock()
		if errUM := s.userMutex.DeleteMutex(user.ID); errUM != nil {
			logger.Log.Info(errUM)
		}
	}()

	balanceSum, err := s.repository.GetSum(user.ID)
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

	return s.repository.CreateAccount(&newAcc)
}

// getAccountRepository создаём репозиторий для начислений
func getAccountRepository(ctx context.Context, dbPool repositories.SQLExecutor) *repositories.AccountRepository {
	return repositories.NewAccountRepository(ctx, dbPool)
}
