// Code generated by MockGen. DO NOT EDIT.
// Source: ./application.go
//
// Generated by this command:
//
//	mockgen -source ./application.go -destination ./mock/mock.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	dtos "account/internal/application/dtos"
	domain "account/internal/domain"
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
	isgomock struct{}
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserService) Create(c context.Context, username, phone, password string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", c, username, phone, password)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserServiceMockRecorder) Create(c, username, phone, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserService)(nil).Create), c, username, phone, password)
}

// FindOneByAccessToken mocks base method.
func (m *MockUserService) FindOneByAccessToken(c context.Context, accessToken string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByAccessToken", c, accessToken)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByAccessToken indicates an expected call of FindOneByAccessToken.
func (mr *MockUserServiceMockRecorder) FindOneByAccessToken(c, accessToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByAccessToken", reflect.TypeOf((*MockUserService)(nil).FindOneByAccessToken), c, accessToken)
}

// FindOneByID mocks base method.
func (m *MockUserService) FindOneByID(c context.Context, ID int) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByID", c, ID)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByID indicates an expected call of FindOneByID.
func (mr *MockUserServiceMockRecorder) FindOneByID(c, ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByID", reflect.TypeOf((*MockUserService)(nil).FindOneByID), c, ID)
}

// FindOneByPhone mocks base method.
func (m *MockUserService) FindOneByPhone(c context.Context, phone string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByPhone", c, phone)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByPhone indicates an expected call of FindOneByPhone.
func (mr *MockUserServiceMockRecorder) FindOneByPhone(c, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByPhone", reflect.TypeOf((*MockUserService)(nil).FindOneByPhone), c, phone)
}

// IsPhoneExists mocks base method.
func (m *MockUserService) IsPhoneExists(c context.Context, phone string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPhoneExists", c, phone)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPhoneExists indicates an expected call of IsPhoneExists.
func (mr *MockUserServiceMockRecorder) IsPhoneExists(c, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPhoneExists", reflect.TypeOf((*MockUserService)(nil).IsPhoneExists), c, phone)
}

// UpdatePassword mocks base method.
func (m *MockUserService) UpdatePassword(c context.Context, userId int, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", c, userId, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockUserServiceMockRecorder) UpdatePassword(c, userId, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserService)(nil).UpdatePassword), c, userId, password)
}

// MockBanService is a mock of BanService interface.
type MockBanService struct {
	ctrl     *gomock.Controller
	recorder *MockBanServiceMockRecorder
	isgomock struct{}
}

// MockBanServiceMockRecorder is the mock recorder for MockBanService.
type MockBanServiceMockRecorder struct {
	mock *MockBanService
}

// NewMockBanService creates a new mock instance.
func NewMockBanService(ctrl *gomock.Controller) *MockBanService {
	mock := &MockBanService{ctrl: ctrl}
	mock.recorder = &MockBanServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBanService) EXPECT() *MockBanServiceMockRecorder {
	return m.recorder
}

// Ban mocks base method.
func (m *MockBanService) Ban(c context.Context, callerId, userId int, reason string, expiry time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ban", c, callerId, userId, reason, expiry)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ban indicates an expected call of Ban.
func (mr *MockBanServiceMockRecorder) Ban(c, callerId, userId, reason, expiry any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ban", reflect.TypeOf((*MockBanService)(nil).Ban), c, callerId, userId, reason, expiry)
}

// FindOneByID mocks base method.
func (m *MockBanService) FindOneByID(c context.Context, id int) (*domain.Ban, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByID", c, id)
	ret0, _ := ret[0].(*domain.Ban)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByID indicates an expected call of FindOneByID.
func (mr *MockBanServiceMockRecorder) FindOneByID(c, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByID", reflect.TypeOf((*MockBanService)(nil).FindOneByID), c, id)
}

// IsUserBanned mocks base method.
func (m *MockBanService) IsUserBanned(c context.Context, userId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserBanned", c, userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserBanned indicates an expected call of IsUserBanned.
func (mr *MockBanServiceMockRecorder) IsUserBanned(c, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserBanned", reflect.TypeOf((*MockBanService)(nil).IsUserBanned), c, userId)
}

// UnBan mocks base method.
func (m *MockBanService) UnBan(c context.Context, banId, unBannedByID int, reason string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnBan", c, banId, unBannedByID, reason)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnBan indicates an expected call of UnBan.
func (mr *MockBanServiceMockRecorder) UnBan(c, banId, unBannedByID, reason any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnBan", reflect.TypeOf((*MockBanService)(nil).UnBan), c, banId, unBannedByID, reason)
}

// MockCodeService is a mock of CodeService interface.
type MockCodeService struct {
	ctrl     *gomock.Controller
	recorder *MockCodeServiceMockRecorder
	isgomock struct{}
}

// MockCodeServiceMockRecorder is the mock recorder for MockCodeService.
type MockCodeServiceMockRecorder struct {
	mock *MockCodeService
}

// NewMockCodeService creates a new mock instance.
func NewMockCodeService(ctrl *gomock.Controller) *MockCodeService {
	mock := &MockCodeService{ctrl: ctrl}
	mock.recorder = &MockCodeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodeService) EXPECT() *MockCodeServiceMockRecorder {
	return m.recorder
}

// RemoveAll mocks base method.
func (m *MockCodeService) RemoveAll(c context.Context, phone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", c, phone)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockCodeServiceMockRecorder) RemoveAll(c, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockCodeService)(nil).RemoveAll), c, phone)
}

// Send mocks base method.
func (m *MockCodeService) Send(c context.Context, phone, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", c, phone, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockCodeServiceMockRecorder) Send(c, phone, ip any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockCodeService)(nil).Send), c, phone, ip)
}

// Verify mocks base method.
func (m *MockCodeService) Verify(c context.Context, phone, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", c, phone, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockCodeServiceMockRecorder) Verify(c, phone, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockCodeService)(nil).Verify), c, phone, code)
}

// MockSessionService is a mock of SessionService interface.
type MockSessionService struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServiceMockRecorder
	isgomock struct{}
}

// MockSessionServiceMockRecorder is the mock recorder for MockSessionService.
type MockSessionServiceMockRecorder struct {
	mock *MockSessionService
}

// NewMockSessionService creates a new mock instance.
func NewMockSessionService(ctrl *gomock.Controller) *MockSessionService {
	mock := &MockSessionService{ctrl: ctrl}
	mock.recorder = &MockSessionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionService) EXPECT() *MockSessionServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionService) Create(c context.Context, userId int) (*dtos.AuthOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", c, userId)
	ret0, _ := ret[0].(*dtos.AuthOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionServiceMockRecorder) Create(c, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionService)(nil).Create), c, userId)
}
