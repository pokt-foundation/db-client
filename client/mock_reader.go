// Code generated by mockery v2.33.0. DO NOT EDIT.

package dbclient

import (
	context "context"

	types "github.com/pokt-foundation/portal-db/v2/types"
	mock "github.com/stretchr/testify/mock"
)

// MockIDBReader is an autogenerated mock type for the IDBReader type
type MockIDBReader struct {
	mock.Mock
}

// GetAllAccounts provides a mock function with given fields: ctx, options
func (_m *MockIDBReader) GetAllAccounts(ctx context.Context, options ...AccountOptions) ([]*types.Account, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...AccountOptions) ([]*types.Account, error)); ok {
		return rf(ctx, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...AccountOptions) []*types.Account); ok {
		r0 = rf(ctx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...AccountOptions) error); ok {
		r1 = rf(ctx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllChains provides a mock function with given fields: ctx, options
func (_m *MockIDBReader) GetAllChains(ctx context.Context, options ...ChainOptions) ([]*types.Chain, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Chain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...ChainOptions) ([]*types.Chain, error)); ok {
		return rf(ctx, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...ChainOptions) []*types.Chain); ok {
		r0 = rf(ctx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Chain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...ChainOptions) error); ok {
		r1 = rf(ctx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllGigastakeApps provides a mock function with given fields: ctx, optionParams
func (_m *MockIDBReader) GetAllGigastakeApps(ctx context.Context, optionParams ...GigastakeAppOptions) ([]*types.GigastakeApp, error) {
	_va := make([]interface{}, len(optionParams))
	for _i := range optionParams {
		_va[_i] = optionParams[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.GigastakeApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...GigastakeAppOptions) ([]*types.GigastakeApp, error)); ok {
		return rf(ctx, optionParams...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...GigastakeAppOptions) []*types.GigastakeApp); ok {
		r0 = rf(ctx, optionParams...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.GigastakeApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...GigastakeAppOptions) error); ok {
		r1 = rf(ctx, optionParams...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllGigastakeAppsByChain provides a mock function with given fields: ctx, chainID
func (_m *MockIDBReader) GetAllGigastakeAppsByChain(ctx context.Context, chainID types.RelayChainID) ([]*types.GigastakeApp, error) {
	ret := _m.Called(ctx, chainID)

	var r0 []*types.GigastakeApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID) ([]*types.GigastakeApp, error)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID) []*types.GigastakeApp); ok {
		r0 = rf(ctx, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.GigastakeApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.RelayChainID) error); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPlans provides a mock function with given fields: ctx
func (_m *MockIDBReader) GetAllPlans(ctx context.Context) ([]types.Plan, error) {
	ret := _m.Called(ctx)

	var r0 []types.Plan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]types.Plan, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []types.Plan); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Plan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPortalApps provides a mock function with given fields: ctx, options
func (_m *MockIDBReader) GetAllPortalApps(ctx context.Context, options ...PortalAppOptions) ([]*types.PortalApp, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.PortalApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...PortalAppOptions) ([]*types.PortalApp, error)); ok {
		return rf(ctx, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...PortalAppOptions) []*types.PortalApp); ok {
		r0 = rf(ctx, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.PortalApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...PortalAppOptions) error); ok {
		r1 = rf(ctx, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockedContracts provides a mock function with given fields: ctx
func (_m *MockIDBReader) GetBlockedContracts(ctx context.Context) (types.GlobalBlockedContracts, error) {
	ret := _m.Called(ctx)

	var r0 types.GlobalBlockedContracts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (types.GlobalBlockedContracts, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) types.GlobalBlockedContracts); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(types.GlobalBlockedContracts)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChainByID provides a mock function with given fields: ctx, chainID
func (_m *MockIDBReader) GetChainByID(ctx context.Context, chainID types.RelayChainID) (*types.Chain, error) {
	ret := _m.Called(ctx, chainID)

	var r0 *types.Chain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID) (*types.Chain, error)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID) *types.Chain); ok {
		r0 = rf(ctx, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Chain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.RelayChainID) error); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGigastakeAppByID provides a mock function with given fields: ctx, gigastakeAppID
func (_m *MockIDBReader) GetGigastakeAppByID(ctx context.Context, gigastakeAppID types.GigastakeAppID) (*types.GigastakeApp, error) {
	ret := _m.Called(ctx, gigastakeAppID)

	var r0 *types.GigastakeApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeAppID) (*types.GigastakeApp, error)); ok {
		return rf(ctx, gigastakeAppID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeAppID) *types.GigastakeApp); ok {
		r0 = rf(ctx, gigastakeAppID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GigastakeApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.GigastakeAppID) error); ok {
		r1 = rf(ctx, gigastakeAppID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortalAppByID provides a mock function with given fields: ctx, portalAppID
func (_m *MockIDBReader) GetPortalAppByID(ctx context.Context, portalAppID types.PortalAppID) (*types.PortalApp, error) {
	ret := _m.Called(ctx, portalAppID)

	var r0 *types.PortalApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalAppID) (*types.PortalApp, error)); ok {
		return rf(ctx, portalAppID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalAppID) *types.PortalApp); ok {
		r0 = rf(ctx, portalAppID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.PortalApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.PortalAppID) error); ok {
		r1 = rf(ctx, portalAppID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortalAppsByUser provides a mock function with given fields: ctx, userID, options
func (_m *MockIDBReader) GetPortalAppsByUser(ctx context.Context, userID types.UserID, options ...PortalAppOptions) ([]*types.PortalApp, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.PortalApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, ...PortalAppOptions) ([]*types.PortalApp, error)); ok {
		return rf(ctx, userID, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, ...PortalAppOptions) []*types.PortalApp); ok {
		r0 = rf(ctx, userID, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.PortalApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UserID, ...PortalAppOptions) error); ok {
		r1 = rf(ctx, userID, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortalAppsForMiddleware provides a mock function with given fields: ctx
func (_m *MockIDBReader) GetPortalAppsForMiddleware(ctx context.Context) ([]*types.PortalAppLite, error) {
	ret := _m.Called(ctx)

	var r0 []*types.PortalAppLite
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*types.PortalAppLite, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*types.PortalAppLite); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.PortalAppLite)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortalUser provides a mock function with given fields: ctx, userID
func (_m *MockIDBReader) GetPortalUser(ctx context.Context, userID string) (*types.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*types.User, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPortalUserID provides a mock function with given fields: ctx, userID
func (_m *MockIDBReader) GetPortalUserID(ctx context.Context, userID string) (types.UserID, error) {
	ret := _m.Called(ctx, userID)

	var r0 types.UserID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.UserID, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.UserID); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(types.UserID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserAccount provides a mock function with given fields: ctx, accountID, userID, options
func (_m *MockIDBReader) GetUserAccount(ctx context.Context, accountID types.AccountID, userID types.UserID, options ...AccountOptions) (*types.Account, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, accountID, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.UserID, ...AccountOptions) (*types.Account, error)); ok {
		return rf(ctx, accountID, userID, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.UserID, ...AccountOptions) *types.Account); ok {
		r0 = rf(ctx, accountID, userID, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.AccountID, types.UserID, ...AccountOptions) error); ok {
		r1 = rf(ctx, accountID, userID, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserAccounts provides a mock function with given fields: ctx, userID, options
func (_m *MockIDBReader) GetUserAccounts(ctx context.Context, userID types.UserID, options ...AccountOptions) ([]*types.Account, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*types.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, ...AccountOptions) ([]*types.Account, error)); ok {
		return rf(ctx, userID, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, ...AccountOptions) []*types.Account); ok {
		r0 = rf(ctx, userID, options...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UserID, ...AccountOptions) error); ok {
		r1 = rf(ctx, userID, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockIDBReader creates a new instance of MockIDBReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIDBReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIDBReader {
	mock := &MockIDBReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
