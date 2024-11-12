// Code generated by MockGen. DO NOT EDIT.
// Source: src/v1/handler/init.go
//
// Generated by this command:
//
//	mockgen -source=src/v1/handler/init.go -destination=src/v1/handler/mock/init.go -package=mock
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	entity "github.com/usagifm/dating-app/src/entity"
	contract "github.com/usagifm/dating-app/src/v1/contract"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// GetProfile mocks base method.
func (m *MockAuthService) GetProfile(ctx context.Context) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", ctx)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockAuthServiceMockRecorder) GetProfile(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockAuthService)(nil).GetProfile), ctx)
}

// Login mocks base method.
func (m *MockAuthService) Login(ctx context.Context, param contract.LoginRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, param)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceMockRecorder) Login(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), ctx, param)
}

// SignUp mocks base method.
func (m *MockAuthService) SignUp(ctx context.Context, param contract.SignUpRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceMockRecorder) SignUp(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthService)(nil).SignUp), ctx, param)
}

// MockDatingService is a mock of DatingService interface.
type MockDatingService struct {
	ctrl     *gomock.Controller
	recorder *MockDatingServiceMockRecorder
}

// MockDatingServiceMockRecorder is the mock recorder for MockDatingService.
type MockDatingServiceMockRecorder struct {
	mock *MockDatingService
}

// NewMockDatingService creates a new mock instance.
func NewMockDatingService(ctrl *gomock.Controller) *MockDatingService {
	mock := &MockDatingService{ctrl: ctrl}
	mock.recorder = &MockDatingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatingService) EXPECT() *MockDatingServiceMockRecorder {
	return m.recorder
}

// BuyPackage mocks base method.
func (m *MockDatingService) BuyPackage(ctx context.Context, param contract.BuyPackageRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyPackage", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyPackage indicates an expected call of BuyPackage.
func (mr *MockDatingServiceMockRecorder) BuyPackage(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyPackage", reflect.TypeOf((*MockDatingService)(nil).BuyPackage), ctx, param)
}

// GetPackages mocks base method.
func (m *MockDatingService) GetPackages(ctx context.Context) ([]*entity.Package, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPackages", ctx)
	ret0, _ := ret[0].([]*entity.Package)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPackages indicates an expected call of GetPackages.
func (mr *MockDatingServiceMockRecorder) GetPackages(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPackages", reflect.TypeOf((*MockDatingService)(nil).GetPackages), ctx)
}

// GetProfilesByPreference mocks base method.
func (m *MockDatingService) GetProfilesByPreference(ctx context.Context) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfilesByPreference", ctx)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfilesByPreference indicates an expected call of GetProfilesByPreference.
func (mr *MockDatingServiceMockRecorder) GetProfilesByPreference(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfilesByPreference", reflect.TypeOf((*MockDatingService)(nil).GetProfilesByPreference), ctx)
}

// GetUserMatches mocks base method.
func (m *MockDatingService) GetUserMatches(ctx context.Context) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserMatches", ctx)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserMatches indicates an expected call of GetUserMatches.
func (mr *MockDatingServiceMockRecorder) GetUserMatches(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserMatches", reflect.TypeOf((*MockDatingService)(nil).GetUserMatches), ctx)
}

// GetUserPreference mocks base method.
func (m *MockDatingService) GetUserPreference(ctx context.Context) (*entity.UserPreference, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPreference", ctx)
	ret0, _ := ret[0].(*entity.UserPreference)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserPreference indicates an expected call of GetUserPreference.
func (mr *MockDatingServiceMockRecorder) GetUserPreference(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPreference", reflect.TypeOf((*MockDatingService)(nil).GetUserPreference), ctx)
}

// Swipe mocks base method.
func (m *MockDatingService) Swipe(ctx context.Context, param contract.SwipeRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Swipe", ctx, param)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Swipe indicates an expected call of Swipe.
func (mr *MockDatingServiceMockRecorder) Swipe(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Swipe", reflect.TypeOf((*MockDatingService)(nil).Swipe), ctx, param)
}

// UpdateUserPreference mocks base method.
func (m *MockDatingService) UpdateUserPreference(ctx context.Context, param contract.UpdateUserPreferenceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPreference", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPreference indicates an expected call of UpdateUserPreference.
func (mr *MockDatingServiceMockRecorder) UpdateUserPreference(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPreference", reflect.TypeOf((*MockDatingService)(nil).UpdateUserPreference), ctx, param)
}
