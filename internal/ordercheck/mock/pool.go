// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ordercheck/pool.go

// Package mock is a generated GoMock package.
package mock

import (
	models "gofemart/internal/models"
	payloads "gofemart/internal/payloads"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockoRepo is a mock of oRepo interface.
type MockoRepo struct {
	ctrl     *gomock.Controller
	recorder *MockoRepoMockRecorder
}

// MockoRepoMockRecorder is the mock recorder for MockoRepo.
type MockoRepoMockRecorder struct {
	mock *MockoRepo
}

// NewMockoRepo creates a new mock instance.
func NewMockoRepo(ctrl *gomock.Controller) *MockoRepo {
	mock := &MockoRepo{ctrl: ctrl}
	mock.recorder = &MockoRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoRepo) EXPECT() *MockoRepoMockRecorder {
	return m.recorder
}

// GetOrdersExcludeOrdersWhereStatusIn mocks base method.
func (m *MockoRepo) GetOrdersExcludeOrdersWhereStatusIn(limit int, excludedNumbers []string, olderThen time.Time, statuses ...string) ([]models.Order, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{limit, excludedNumbers, olderThen}
	for _, a := range statuses {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOrdersExcludeOrdersWhereStatusIn", varargs...)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersExcludeOrdersWhereStatusIn indicates an expected call of GetOrdersExcludeOrdersWhereStatusIn.
func (mr *MockoRepoMockRecorder) GetOrdersExcludeOrdersWhereStatusIn(limit, excludedNumbers, olderThen interface{}, statuses ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{limit, excludedNumbers, olderThen}, statuses...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersExcludeOrdersWhereStatusIn", reflect.TypeOf((*MockoRepo)(nil).GetOrdersExcludeOrdersWhereStatusIn), varargs...)
}

// UpdateOrder mocks base method.
func (m *MockoRepo) UpdateOrder(order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockoRepoMockRecorder) UpdateOrder(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockoRepo)(nil).UpdateOrder), order)
}

// MockaRepo is a mock of aRepo interface.
type MockaRepo struct {
	ctrl     *gomock.Controller
	recorder *MockaRepoMockRecorder
}

// MockaRepoMockRecorder is the mock recorder for MockaRepo.
type MockaRepoMockRecorder struct {
	mock *MockaRepo
}

// NewMockaRepo creates a new mock instance.
func NewMockaRepo(ctrl *gomock.Controller) *MockaRepo {
	mock := &MockaRepo{ctrl: ctrl}
	mock.recorder = &MockaRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockaRepo) EXPECT() *MockaRepoMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockaRepo) CreateAccount(account *models.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", account)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockaRepoMockRecorder) CreateAccount(account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockaRepo)(nil).CreateAccount), account)
}

// MockAccrual is a mock of Accrual interface.
type MockAccrual struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualMockRecorder
}

// MockAccrualMockRecorder is the mock recorder for MockAccrual.
type MockAccrualMockRecorder struct {
	mock *MockAccrual
}

// NewMockAccrual creates a new mock instance.
func NewMockAccrual(ctrl *gomock.Controller) *MockAccrual {
	mock := &MockAccrual{ctrl: ctrl}
	mock.recorder = &MockAccrualMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrual) EXPECT() *MockAccrualMockRecorder {
	return m.recorder
}

// Accrual mocks base method.
func (m *MockAccrual) Accrual(order *models.Order) (*payloads.Accrual, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accrual", order)
	ret0, _ := ret[0].(*payloads.Accrual)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accrual indicates an expected call of Accrual.
func (mr *MockAccrualMockRecorder) Accrual(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accrual", reflect.TypeOf((*MockAccrual)(nil).Accrual), order)
}

// Pause mocks base method.
func (m *MockAccrual) Pause(duration time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Pause", duration)
}

// Pause indicates an expected call of Pause.
func (mr *MockAccrualMockRecorder) Pause(duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pause", reflect.TypeOf((*MockAccrual)(nil).Pause), duration)
}
