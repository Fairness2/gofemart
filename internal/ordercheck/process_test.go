package ordercheck

import (
	"errors"
	"github.com/golang/mock/gomock"
	"gofemart/internal/models"
	"gofemart/internal/ordercheck/mock"
	"gofemart/internal/payloads"
	"testing"
)

func TestCreateNewAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	testCases := []struct {
		name             string
		inputOrderNumber string
		inputUserID      int64
		inputDiff        float64
		wantErr          bool
		setup            func() aRepo
	}{
		{
			name:             "success_create_account",
			inputOrderNumber: "ORD123",
			inputUserID:      123,
			inputDiff:        45.6,
			wantErr:          false,
			setup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				repo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
		},
		{
			name:             "wrong_input_order_number",
			inputOrderNumber: "",
			inputUserID:      0,
			inputDiff:        0,
			wantErr:          true,
			setup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				repo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(errors.New("wrong input"))
				return repo
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Pool{
				accountRepo: tc.setup(),
			}
			_, err := p.createNewAccount(tc.inputOrderNumber, tc.inputUserID, tc.inputDiff)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestProcessOrderAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	testCases := []struct {
		name       string
		accrual    *payloads.Accrual
		order      *models.Order
		setup      func() oRepo
		aSetup     func() aRepo
		wantErr    bool
		wantStatus string
	}{
		{
			name: "status_processing",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualProcessing,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				return repo
			},
			wantErr:    false,
			wantStatus: models.StatusProcessing,
		},
		{
			name: "status_registered",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualRegistered,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				return repo
			},
			wantErr:    false,
			wantStatus: models.StatusProcessing,
		},
		{
			name: "status_invalid",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualInvalid,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				return repo
			},
			wantErr:    false,
			wantStatus: models.StatusInvalid,
		},
		{
			name: "status_processed",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualProcessed,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				repo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			wantErr:    false,
			wantStatus: models.StatusProcessed,
		},
		{
			name: "status_processed_and_account_error",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualProcessed,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(nil)
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				repo.EXPECT().
					CreateAccount(gomock.Any()).
					AnyTimes().
					Return(errors.New("account error"))
				return repo
			},
			wantErr:    true,
			wantStatus: models.StatusProcessing,
		},
		{
			name: "status_invalid_and_order_error",
			accrual: &payloads.Accrual{
				Order:   "1",
				Status:  payloads.StatusAccrualInvalid,
				Accrual: 11,
			},
			order: &models.Order{
				Number:     "1",
				StatusCode: models.StatusProcessing,
			},
			setup: func() oRepo {
				repo := mock.NewMockoRepo(ctrl)
				repo.EXPECT().
					UpdateOrder(gomock.Any()).
					AnyTimes().
					Return(errors.New("order error"))
				return repo
			},
			aSetup: func() aRepo {
				repo := mock.NewMockaRepo(ctrl)
				return repo
			},
			wantErr:    true,
			wantStatus: models.StatusInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Pool{
				orderRepo:   tc.setup(),
				accountRepo: tc.aSetup(),
			}
			err := p.processOrderAccrual(tc.accrual, tc.order)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.wantStatus != tc.order.StatusCode {
				t.Errorf("expected status %s, got %s", tc.wantStatus, tc.order.StatusCode)
			}
		})
	}
}
