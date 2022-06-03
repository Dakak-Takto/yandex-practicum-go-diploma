// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/storage (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	entity "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetOrderByNumber mocks base method.
func (m *MockStorage) GetOrderByNumber(arg0 string) (*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", arg0)
	ret0, _ := ret[0].(*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockStorageMockRecorder) GetOrderByNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockStorage)(nil).GetOrderByNumber), arg0)
}

// GetUserByID mocks base method.
func (m *MockStorage) GetUserByID(arg0 int) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStorageMockRecorder) GetUserByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStorage)(nil).GetUserByID), arg0)
}

// GetUserByLogin mocks base method.
func (m *MockStorage) GetUserByLogin(arg0 string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStorageMockRecorder) GetUserByLogin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorage)(nil).GetUserByLogin), arg0)
}

// SaveUser mocks base method.
func (m *MockStorage) SaveUser(arg0 *entity.User) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", arg0)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUser indicates an expected call of SaveUser.
func (mr *MockStorageMockRecorder) SaveUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUser", reflect.TypeOf((*MockStorage)(nil).SaveUser), arg0)
}

// SaveUserOrder mocks base method.
func (m *MockStorage) SaveUserOrder(arg0 string, arg1 int) (*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUserOrder", arg0, arg1)
	ret0, _ := ret[0].(*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUserOrder indicates an expected call of SaveUserOrder.
func (mr *MockStorageMockRecorder) SaveUserOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUserOrder", reflect.TypeOf((*MockStorage)(nil).SaveUserOrder), arg0, arg1)
}

// SaveWithdraw mocks base method.
func (m *MockStorage) SaveWithdraw(arg0 *entity.Withdraw) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWithdraw", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWithdraw indicates an expected call of SaveWithdraw.
func (mr *MockStorageMockRecorder) SaveWithdraw(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWithdraw", reflect.TypeOf((*MockStorage)(nil).SaveWithdraw), arg0)
}

// SelectNewOrders mocks base method.
func (m *MockStorage) SelectNewOrders() ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectNewOrders")
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectNewOrders indicates an expected call of SelectNewOrders.
func (mr *MockStorageMockRecorder) SelectNewOrders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectNewOrders", reflect.TypeOf((*MockStorage)(nil).SelectNewOrders))
}

// SelectOrdersByUserID mocks base method.
func (m *MockStorage) SelectOrdersByUserID(arg0 int) ([]*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectOrdersByUserID", arg0)
	ret0, _ := ret[0].([]*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectOrdersByUserID indicates an expected call of SelectOrdersByUserID.
func (mr *MockStorageMockRecorder) SelectOrdersByUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectOrdersByUserID", reflect.TypeOf((*MockStorage)(nil).SelectOrdersByUserID), arg0)
}

// SelectWithdrawals mocks base method.
func (m *MockStorage) SelectWithdrawals(arg0 int) ([]*entity.Withdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectWithdrawals", arg0)
	ret0, _ := ret[0].([]*entity.Withdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectWithdrawals indicates an expected call of SelectWithdrawals.
func (mr *MockStorageMockRecorder) SelectWithdrawals(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectWithdrawals", reflect.TypeOf((*MockStorage)(nil).SelectWithdrawals), arg0)
}

// UpdateOrder mocks base method.
func (m *MockStorage) UpdateOrder(arg0 *entity.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStorageMockRecorder) UpdateOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStorage)(nil).UpdateOrder), arg0)
}

// UpdateUser mocks base method.
func (m *MockStorage) UpdateUser(arg0 *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStorageMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStorage)(nil).UpdateUser), arg0)
}

// UserBalanceChange mocks base method.
func (m *MockStorage) UserBalanceChange(arg0 int, arg1 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserBalanceChange", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UserBalanceChange indicates an expected call of UserBalanceChange.
func (mr *MockStorageMockRecorder) UserBalanceChange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserBalanceChange", reflect.TypeOf((*MockStorage)(nil).UserBalanceChange), arg0, arg1)
}
