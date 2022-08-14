// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudquery/cloudquery/plugins/source/azure/client/services (interfaces: ADGroupsClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	graphrbac "github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	gomock "github.com/golang/mock/gomock"
)

// MockADGroupsClient is a mock of ADGroupsClient interface.
type MockADGroupsClient struct {
	ctrl     *gomock.Controller
	recorder *MockADGroupsClientMockRecorder
}

// MockADGroupsClientMockRecorder is the mock recorder for MockADGroupsClient.
type MockADGroupsClientMockRecorder struct {
	mock *MockADGroupsClient
}

// NewMockADGroupsClient creates a new mock instance.
func NewMockADGroupsClient(ctrl *gomock.Controller) *MockADGroupsClient {
	mock := &MockADGroupsClient{ctrl: ctrl}
	mock.recorder = &MockADGroupsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockADGroupsClient) EXPECT() *MockADGroupsClientMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockADGroupsClient) List(arg0 context.Context, arg1 string) (graphrbac.GroupListResultPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(graphrbac.GroupListResultPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockADGroupsClientMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockADGroupsClient)(nil).List), arg0, arg1)
}
