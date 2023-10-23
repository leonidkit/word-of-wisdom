// Code generated by MockGen. DO NOT EDIT.
// Source: challenger.go

// Package challenger_repo_mocks is a generated GoMock package.
package challenger_repo_mocks

import (
	context "context"
	reflect "reflect"
	hashcash "github.com/leonidkit/word-of-wisdom/pkg/hashcash"

	gomock "github.com/golang/mock/gomock"
)

// MockchallengesRepository is a mock of challengesRepository interface.
type MockchallengesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockchallengesRepositoryMockRecorder
}

// MockchallengesRepositoryMockRecorder is the mock recorder for MockchallengesRepository.
type MockchallengesRepositoryMockRecorder struct {
	mock *MockchallengesRepository
}

// NewMockchallengesRepository creates a new mock instance.
func NewMockchallengesRepository(ctrl *gomock.Controller) *MockchallengesRepository {
	mock := &MockchallengesRepository{ctrl: ctrl}
	mock.recorder = &MockchallengesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchallengesRepository) EXPECT() *MockchallengesRepositoryMockRecorder {
	return m.recorder
}

// AddChallenge mocks base method.
func (m *MockchallengesRepository) AddChallenge(ctx context.Context, key string, header *hashcash.HashcashHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddChallenge", ctx, key, header)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddChallenge indicates an expected call of AddChallenge.
func (mr *MockchallengesRepositoryMockRecorder) AddChallenge(ctx, key, header interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChallenge", reflect.TypeOf((*MockchallengesRepository)(nil).AddChallenge), ctx, key, header)
}

// DeleteChallenge mocks base method.
func (m *MockchallengesRepository) DeleteChallenge(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChallenge", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChallenge indicates an expected call of DeleteChallenge.
func (mr *MockchallengesRepositoryMockRecorder) DeleteChallenge(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChallenge", reflect.TypeOf((*MockchallengesRepository)(nil).DeleteChallenge), ctx, key)
}

// IsChallengeExists mocks base method.
func (m *MockchallengesRepository) IsChallengeExists(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChallengeExists", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsChallengeExists indicates an expected call of IsChallengeExists.
func (mr *MockchallengesRepositoryMockRecorder) IsChallengeExists(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChallengeExists", reflect.TypeOf((*MockchallengesRepository)(nil).IsChallengeExists), ctx, key)
}
