// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/interfaces.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=internal/services/interfaces.go -destination=internal/test/mocks/interfaces.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	repository "github.com/jcserv/rivalslfg/internal/repository"
	gomock "go.uber.org/mock/gomock"
)

// MockIGroup is a mock of IGroup interface.
type MockIGroup struct {
	ctrl     *gomock.Controller
	recorder *MockIGroupMockRecorder
	isgomock struct{}
}

// MockIGroupMockRecorder is the mock recorder for MockIGroup.
type MockIGroupMockRecorder struct {
	mock *MockIGroup
}

// NewMockIGroup creates a new mock instance.
func NewMockIGroup(ctrl *gomock.Controller) *MockIGroup {
	mock := &MockIGroup{ctrl: ctrl}
	mock.recorder = &MockIGroupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGroup) EXPECT() *MockIGroupMockRecorder {
	return m.recorder
}

// CreateGroup mocks base method.
func (m *MockIGroup) CreateGroup(ctx context.Context, arg repository.CreateGroupParams) (repository.CreateGroupRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, arg)
	ret0, _ := ret[0].(repository.CreateGroupRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockIGroupMockRecorder) CreateGroup(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockIGroup)(nil).CreateGroup), ctx, arg)
}

// GetGroupByID mocks base method.
func (m *MockIGroup) GetGroupByID(ctx context.Context, id string, isGroupOwner bool) (*repository.GroupWithPlayers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByID", ctx, id, isGroupOwner)
	ret0, _ := ret[0].(*repository.GroupWithPlayers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByID indicates an expected call of GetGroupByID.
func (mr *MockIGroupMockRecorder) GetGroupByID(ctx, id, isGroupOwner any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByID", reflect.TypeOf((*MockIGroup)(nil).GetGroupByID), ctx, id, isGroupOwner)
}

// GetGroups mocks base method.
func (m *MockIGroup) GetGroups(ctx context.Context, arg repository.GetGroupsParams) ([]repository.GroupWithPlayers, int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", ctx, arg)
	ret0, _ := ret[0].([]repository.GroupWithPlayers)
	ret1, _ := ret[1].(int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockIGroupMockRecorder) GetGroups(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockIGroup)(nil).GetGroups), ctx, arg)
}

// MockIPlayer is a mock of IPlayer interface.
type MockIPlayer struct {
	ctrl     *gomock.Controller
	recorder *MockIPlayerMockRecorder
	isgomock struct{}
}

// MockIPlayerMockRecorder is the mock recorder for MockIPlayer.
type MockIPlayerMockRecorder struct {
	mock *MockIPlayer
}

// NewMockIPlayer creates a new mock instance.
func NewMockIPlayer(ctrl *gomock.Controller) *MockIPlayer {
	mock := &MockIPlayer{ctrl: ctrl}
	mock.recorder = &MockIPlayerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPlayer) EXPECT() *MockIPlayerMockRecorder {
	return m.recorder
}

// JoinGroup mocks base method.
func (m *MockIPlayer) JoinGroup(ctx context.Context, arg repository.JoinGroupParams) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinGroup", ctx, arg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JoinGroup indicates an expected call of JoinGroup.
func (mr *MockIPlayerMockRecorder) JoinGroup(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinGroup", reflect.TypeOf((*MockIPlayer)(nil).JoinGroup), ctx, arg)
}
