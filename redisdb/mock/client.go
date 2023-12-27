// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/slamchillz/xchange/redisdb (interfaces: RedisClient)
//
// Generated by this command:
//
//	mockgen -package mockredisdb -destination redisdb/mock/client.go github.com/slamchillz/xchange/redisdb RedisClient
//
// Package mockredisdb is a generated GoMock package.
package mockredisdb

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockRedisClient is a mock of RedisClient interface.
type MockRedisClient struct {
	ctrl     *gomock.Controller
	recorder *MockRedisClientMockRecorder
}

// MockRedisClientMockRecorder is the mock recorder for MockRedisClient.
type MockRedisClientMockRecorder struct {
	mock *MockRedisClient
}

// NewMockRedisClient creates a new mock instance.
func NewMockRedisClient(ctrl *gomock.Controller) *MockRedisClient {
	mock := &MockRedisClient{ctrl: ctrl}
	mock.recorder = &MockRedisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisClient) EXPECT() *MockRedisClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRedisClient) Get(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisClient)(nil).Get), arg0, arg1)
}

// Set mocks base method.
func (m *MockRedisClient) Set(arg0 context.Context, arg1 string, arg2 any, arg3 time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockRedisClientMockRecorder) Set(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisClient)(nil).Set), arg0, arg1, arg2, arg3)
}