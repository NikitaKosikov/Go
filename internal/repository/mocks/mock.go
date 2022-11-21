// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	domain "test/internal/domain"
	api "test/pkg/api"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, user)
}

// Delete mocks base method.
func (m *MockUserRepository) Delete(ctx context.Context, oid primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, oid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserRepositoryMockRecorder) Delete(ctx, oid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), ctx, oid)
}

// FindAll mocks base method.
func (m *MockUserRepository) FindAll(ctx context.Context, pagination api.Pagination, filters []api.Filters, sortOptions []api.Options) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, pagination, filters, sortOptions)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockUserRepositoryMockRecorder) FindAll(ctx, pagination, filters, sortOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockUserRepository)(nil).FindAll), ctx, pagination, filters, sortOptions)
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// FindOne mocks base method.
func (m *MockUserRepository) FindOne(ctx context.Context, oid primitive.ObjectID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, oid)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockUserRepositoryMockRecorder) FindOne(ctx, oid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockUserRepository)(nil).FindOne), ctx, oid)
}

// GetUserByRefreshToken mocks base method.
func (m *MockUserRepository) GetUserByRefreshToken(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByRefreshToken", ctx, id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByRefreshToken indicates an expected call of GetUserByRefreshToken.
func (mr *MockUserRepositoryMockRecorder) GetUserByRefreshToken(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByRefreshToken", reflect.TypeOf((*MockUserRepository)(nil).GetUserByRefreshToken), ctx, id)
}

// SetSession mocks base method.
func (m *MockUserRepository) SetSession(ctx context.Context, oid primitive.ObjectID, session domain.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSession", ctx, oid, session)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSession indicates an expected call of SetSession.
func (mr *MockUserRepositoryMockRecorder) SetSession(ctx, oid, session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSession", reflect.TypeOf((*MockUserRepository)(nil).SetSession), ctx, oid, session)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, user)
}