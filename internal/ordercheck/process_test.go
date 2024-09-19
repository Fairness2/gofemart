package ordercheck

import (
	"errors"
	"github.com/golang/mock/gomock"
	"gofemart/internal/ordercheck/mock"
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
