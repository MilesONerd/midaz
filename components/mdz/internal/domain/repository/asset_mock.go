// Code generated by MockGen. DO NOT EDIT.
// Source: /home/max/Workspace/midaz/components/mdz/internal/domain/repository/asset.go
//
// Generated by this command:
//
//	mockgen -source=/home/max/Workspace/midaz/components/mdz/internal/domain/repository/asset.go -destination=/home/max/Workspace/midaz/components/mdz/internal/domain/repository/asset_mock.go -package repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	mmodel "github.com/LerianStudio/midaz/common/mmodel"
	gomock "go.uber.org/mock/gomock"
)

// MockAsset is a mock of Asset interface.
type MockAsset struct {
	ctrl     *gomock.Controller
	recorder *MockAssetMockRecorder
	isgomock struct{}
}

// MockAssetMockRecorder is the mock recorder for MockAsset.
type MockAssetMockRecorder struct {
	mock *MockAsset
}

// NewMockAsset creates a new mock instance.
func NewMockAsset(ctrl *gomock.Controller) *MockAsset {
	mock := &MockAsset{ctrl: ctrl}
	mock.recorder = &MockAssetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAsset) EXPECT() *MockAssetMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAsset) Create(organizationID, ledgerID string, inp mmodel.CreateAssetInput) (*mmodel.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", organizationID, ledgerID, inp)
	ret0, _ := ret[0].(*mmodel.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAssetMockRecorder) Create(organizationID, ledgerID, inp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAsset)(nil).Create), organizationID, ledgerID, inp)
}

// Get mocks base method.
func (m *MockAsset) Get(organizationID, ledgerID string, limit, page int) (*mmodel.Assets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", organizationID, ledgerID, limit, page)
	ret0, _ := ret[0].(*mmodel.Assets)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAssetMockRecorder) Get(organizationID, ledgerID, limit, page any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAsset)(nil).Get), organizationID, ledgerID, limit, page)
}

// GetByID mocks base method.
func (m *MockAsset) GetByID(organizationID, ledgerID, assetID string) (*mmodel.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", organizationID, ledgerID, assetID)
	ret0, _ := ret[0].(*mmodel.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAssetMockRecorder) GetByID(organizationID, ledgerID, assetID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAsset)(nil).GetByID), organizationID, ledgerID, assetID)
}

// Update mocks base method.
func (m *MockAsset) Update(organizationID, ledgerID, assetID string, inp mmodel.UpdateAssetInput) (*mmodel.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", organizationID, ledgerID, assetID, inp)
	ret0, _ := ret[0].(*mmodel.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAssetMockRecorder) Update(organizationID, ledgerID, assetID, inp any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAsset)(nil).Update), organizationID, ledgerID, assetID, inp)
}
