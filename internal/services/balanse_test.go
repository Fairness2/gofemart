package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"gofemart/internal/models"
	"gofemart/internal/services/mock"
	"sync"
	"testing"
)

func TestSpend(t *testing.T) {
	ctrl := gomock.NewController(t)
	userMutex := mock.NewMockMutexService(ctrl)
	userMutex.EXPECT().
		GetMutex(gomock.Any()).
		AnyTimes().
		Return(&sync.Mutex{}, true)
	userMutex.EXPECT().
		SetMutex(gomock.Any()).
		AnyTimes().
		Return(&sync.Mutex{})
	userMutex.EXPECT().
		DeleteMutex(gomock.Any()).
		AnyTimes().
		Return(nil)

	tests := []struct {
		name           string
		sum            float64
		expectationErr error
		wantErr        bool
		setup          func() BalanceRepository
	}{
		{
			"enough balance",
			1000,
			nil,
			false,
			func() BalanceRepository {
				balanceRepo := mock.NewMockBalanceRepository(ctrl)
				balanceRepo.EXPECT().
					GetSum(gomock.Any()).
					AnyTimes().
					Return(float64(2000), nil)
				balanceRepo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return balanceRepo
			},
		},
		{
			"balance exact",
			1000,
			nil,
			false,
			func() BalanceRepository {
				balanceRepo := mock.NewMockBalanceRepository(ctrl)
				balanceRepo.EXPECT().
					GetSum(gomock.Any()).
					AnyTimes().
					Return(float64(1000), nil)
				balanceRepo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return balanceRepo
			},
		},
		{
			"not enough balance",
			2000,
			ErrorNotEnoughItems,
			true,
			func() BalanceRepository {
				balanceRepo := mock.NewMockBalanceRepository(ctrl)
				balanceRepo.EXPECT().
					GetSum(gomock.Any()).
					AnyTimes().
					Return(float64(1000), nil)
				balanceRepo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return balanceRepo
			},
		},
		{
			"balance retrieval error",
			2000,
			sql.ErrNoRows,
			true,
			func() BalanceRepository {
				balanceRepo := mock.NewMockBalanceRepository(ctrl)
				balanceRepo.EXPECT().
					GetSum(gomock.Any()).
					AnyTimes().
					Return(float64(0), sql.ErrNoRows)
				balanceRepo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return balanceRepo
			},
		},
		{
			"account creation error",
			1000,
			sql.ErrConnDone,
			true,
			func() BalanceRepository {
				balanceRepo := mock.NewMockBalanceRepository(ctrl)
				balanceRepo.EXPECT().
					GetSum(gomock.Any()).
					AnyTimes().
					Return(float64(2000), nil)
				balanceRepo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(sql.ErrConnDone)
				return balanceRepo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup()
			user := &models.User{ID: 1}
			order := &models.Order{Number: "2377225624"}

			service := &BalanceService{
				ctx:        context.Background(),
				repository: repo,
				userMutex:  userMutex,
			}

			err := service.Spend(user, tt.sum, order)
			if tt.wantErr && !errors.Is(err, tt.expectationErr) {
				t.Errorf("BalanceService.Spend() error = %v, wantErr %v", err, tt.expectationErr)
			} else if !tt.wantErr && err != nil {
				t.Errorf("BalanceService.Spend() error = %v, expect no errors", err)
			}
		})
	}
}
