// Code generated by mockery v2.33.0. DO NOT EDIT.

package dbclient

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	types "github.com/pokt-foundation/portal-db/v2/types"
)

// MockIDBWriter is an autogenerated mock type for the IDBWriter type
type MockIDBWriter struct {
	mock.Mock
}

// ActivateChain provides a mock function with given fields: ctx, chainID, active
func (_m *MockIDBWriter) ActivateChain(ctx context.Context, chainID types.RelayChainID, active bool) (bool, error) {
	ret := _m.Called(ctx, chainID, active)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID, bool) (bool, error)); ok {
		return rf(ctx, chainID, active)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.RelayChainID, bool) bool); ok {
		r0 = rf(ctx, chainID, active)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.RelayChainID, bool) error); ok {
		r1 = rf(ctx, chainID, active)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAccount provides a mock function with given fields: ctx, userID, account, timestamp
func (_m *MockIDBWriter) CreateAccount(ctx context.Context, userID types.UserID, account types.Account, timestamp time.Time) (*types.Account, error) {
	ret := _m.Called(ctx, userID, account, timestamp)

	var r0 *types.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, types.Account, time.Time) (*types.Account, error)); ok {
		return rf(ctx, userID, account, timestamp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID, types.Account, time.Time) *types.Account); ok {
		r0 = rf(ctx, userID, account, timestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UserID, types.Account, time.Time) error); ok {
		r1 = rf(ctx, userID, account, timestamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAccountIntegration provides a mock function with given fields: ctx, accountID, integration
func (_m *MockIDBWriter) CreateAccountIntegration(ctx context.Context, accountID types.AccountID, integration types.AccountIntegrations) (*types.AccountIntegrations, error) {
	ret := _m.Called(ctx, accountID, integration)

	var r0 *types.AccountIntegrations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.AccountIntegrations) (*types.AccountIntegrations, error)); ok {
		return rf(ctx, accountID, integration)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.AccountIntegrations) *types.AccountIntegrations); ok {
		r0 = rf(ctx, accountID, integration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.AccountIntegrations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.AccountID, types.AccountIntegrations) error); ok {
		r1 = rf(ctx, accountID, integration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateChainAndGigastakeApps provides a mock function with given fields: ctx, newChainInput
func (_m *MockIDBWriter) CreateChainAndGigastakeApps(ctx context.Context, newChainInput types.NewChainInput) (*types.NewChainInput, error) {
	ret := _m.Called(ctx, newChainInput)

	var r0 *types.NewChainInput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.NewChainInput) (*types.NewChainInput, error)); ok {
		return rf(ctx, newChainInput)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.NewChainInput) *types.NewChainInput); ok {
		r0 = rf(ctx, newChainInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.NewChainInput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.NewChainInput) error); ok {
		r1 = rf(ctx, newChainInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateGigastakeApp provides a mock function with given fields: ctx, gigastakeAppInput
func (_m *MockIDBWriter) CreateGigastakeApp(ctx context.Context, gigastakeAppInput types.GigastakeApp) (*types.GigastakeApp, error) {
	ret := _m.Called(ctx, gigastakeAppInput)

	var r0 *types.GigastakeApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeApp) (*types.GigastakeApp, error)); ok {
		return rf(ctx, gigastakeAppInput)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeApp) *types.GigastakeApp); ok {
		r0 = rf(ctx, gigastakeAppInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GigastakeApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.GigastakeApp) error); ok {
		r1 = rf(ctx, gigastakeAppInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePortalApp provides a mock function with given fields: ctx, portalAppInput
func (_m *MockIDBWriter) CreatePortalApp(ctx context.Context, portalAppInput types.PortalApp) (*types.PortalApp, error) {
	ret := _m.Called(ctx, portalAppInput)

	var r0 *types.PortalApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalApp) (*types.PortalApp, error)); ok {
		return rf(ctx, portalAppInput)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalApp) *types.PortalApp); ok {
		r0 = rf(ctx, portalAppInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.PortalApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.PortalApp) error); ok {
		r1 = rf(ctx, portalAppInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *MockIDBWriter) CreateUser(ctx context.Context, user types.CreateUser) (*types.CreateUserResponse, error) {
	ret := _m.Called(ctx, user)

	var r0 *types.CreateUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.CreateUser) (*types.CreateUserResponse, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.CreateUser) *types.CreateUserResponse); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.CreateUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.CreateUser) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccount provides a mock function with given fields: ctx, accountID
func (_m *MockIDBWriter) DeleteAccount(ctx context.Context, accountID types.AccountID) (map[string]string, error) {
	ret := _m.Called(ctx, accountID)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID) (map[string]string, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID) map[string]string); ok {
		r0 = rf(ctx, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.AccountID) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePortalApp provides a mock function with given fields: ctx, portalAppID
func (_m *MockIDBWriter) DeletePortalApp(ctx context.Context, portalAppID types.PortalAppID) (map[string]string, error) {
	ret := _m.Called(ctx, portalAppID)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalAppID) (map[string]string, error)); ok {
		return rf(ctx, portalAppID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.PortalAppID) map[string]string); ok {
		r0 = rf(ctx, portalAppID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.PortalAppID) error); ok {
		r1 = rf(ctx, portalAppID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *MockIDBWriter) DeleteUser(ctx context.Context, userID types.UserID) (map[string]string, error) {
	ret := _m.Called(ctx, userID)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID) (map[string]string, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UserID) map[string]string); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UserID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveAccountUser provides a mock function with given fields: ctx, removeUser
func (_m *MockIDBWriter) RemoveAccountUser(ctx context.Context, removeUser types.UpdateRemoveAccountUser) (map[string]string, error) {
	ret := _m.Called(ctx, removeUser)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateRemoveAccountUser) (map[string]string, error)); ok {
		return rf(ctx, removeUser)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateRemoveAccountUser) map[string]string); ok {
		r0 = rf(ctx, removeUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateRemoveAccountUser) error); ok {
		r1 = rf(ctx, removeUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveBlockedContract provides a mock function with given fields: ctx, address
func (_m *MockIDBWriter) RemoveBlockedContract(ctx context.Context, address types.BlockedAddress) (map[string]string, error) {
	ret := _m.Called(ctx, address)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedAddress) (map[string]string, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedAddress) map[string]string); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.BlockedAddress) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAccountUserRole provides a mock function with given fields: ctx, updateUser, _a2
func (_m *MockIDBWriter) SetAccountUserRole(ctx context.Context, updateUser types.UpdateAccountUserRole, _a2 time.Time) (map[string]string, error) {
	ret := _m.Called(ctx, updateUser, _a2)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAccountUserRole, time.Time) (map[string]string, error)); ok {
		return rf(ctx, updateUser, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAccountUserRole, time.Time) map[string]string); ok {
		r0 = rf(ctx, updateUser, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateAccountUserRole, time.Time) error); ok {
		r1 = rf(ctx, updateUser, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAcceptAccountUser provides a mock function with given fields: ctx, acceptUser, _a2
func (_m *MockIDBWriter) UpdateAcceptAccountUser(ctx context.Context, acceptUser types.UpdateAcceptAccountUser, _a2 time.Time) (map[string]string, error) {
	ret := _m.Called(ctx, acceptUser, _a2)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAcceptAccountUser, time.Time) (map[string]string, error)); ok {
		return rf(ctx, acceptUser, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAcceptAccountUser, time.Time) map[string]string); ok {
		r0 = rf(ctx, acceptUser, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateAcceptAccountUser, time.Time) error); ok {
		r1 = rf(ctx, acceptUser, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAccount provides a mock function with given fields: ctx, account
func (_m *MockIDBWriter) UpdateAccount(ctx context.Context, account types.UpdateAccount) (*types.Account, error) {
	ret := _m.Called(ctx, account)

	var r0 *types.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAccount) (*types.Account, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateAccount) *types.Account); ok {
		r0 = rf(ctx, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateAccount) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAccountIntegration provides a mock function with given fields: ctx, accountID, integration
func (_m *MockIDBWriter) UpdateAccountIntegration(ctx context.Context, accountID types.AccountID, integration types.AccountIntegrations) (*types.AccountIntegrations, error) {
	ret := _m.Called(ctx, accountID, integration)

	var r0 *types.AccountIntegrations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.AccountIntegrations) (*types.AccountIntegrations, error)); ok {
		return rf(ctx, accountID, integration)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.AccountID, types.AccountIntegrations) *types.AccountIntegrations); ok {
		r0 = rf(ctx, accountID, integration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.AccountIntegrations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.AccountID, types.AccountIntegrations) error); ok {
		r1 = rf(ctx, accountID, integration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBlockedContractActive provides a mock function with given fields: ctx, address, isActive
func (_m *MockIDBWriter) UpdateBlockedContractActive(ctx context.Context, address types.BlockedAddress, isActive bool) (map[string]bool, error) {
	ret := _m.Called(ctx, address, isActive)

	var r0 map[string]bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedAddress, bool) (map[string]bool, error)); ok {
		return rf(ctx, address, isActive)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedAddress, bool) map[string]bool); ok {
		r0 = rf(ctx, address, isActive)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]bool)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.BlockedAddress, bool) error); ok {
		r1 = rf(ctx, address, isActive)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateChain provides a mock function with given fields: ctx, chainUpdate
func (_m *MockIDBWriter) UpdateChain(ctx context.Context, chainUpdate types.UpdateChain) (*types.Chain, error) {
	ret := _m.Called(ctx, chainUpdate)

	var r0 *types.Chain
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateChain) (*types.Chain, error)); ok {
		return rf(ctx, chainUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateChain) *types.Chain); ok {
		r0 = rf(ctx, chainUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Chain)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateChain) error); ok {
		r1 = rf(ctx, chainUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateGigastakeApp provides a mock function with given fields: ctx, id, updateGigastakeApp
func (_m *MockIDBWriter) UpdateGigastakeApp(ctx context.Context, id types.GigastakeAppID, updateGigastakeApp types.UpdateGigastakeApp) (*types.UpdateGigastakeApp, error) {
	ret := _m.Called(ctx, id, updateGigastakeApp)

	var r0 *types.UpdateGigastakeApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeAppID, types.UpdateGigastakeApp) (*types.UpdateGigastakeApp, error)); ok {
		return rf(ctx, id, updateGigastakeApp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.GigastakeAppID, types.UpdateGigastakeApp) *types.UpdateGigastakeApp); ok {
		r0 = rf(ctx, id, updateGigastakeApp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.UpdateGigastakeApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.GigastakeAppID, types.UpdateGigastakeApp) error); ok {
		r1 = rf(ctx, id, updateGigastakeApp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePortalApp provides a mock function with given fields: ctx, portalAppUpdate
func (_m *MockIDBWriter) UpdatePortalApp(ctx context.Context, portalAppUpdate types.UpdatePortalApp) (*types.UpdatePortalApp, error) {
	ret := _m.Called(ctx, portalAppUpdate)

	var r0 *types.UpdatePortalApp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdatePortalApp) (*types.UpdatePortalApp, error)); ok {
		return rf(ctx, portalAppUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdatePortalApp) *types.UpdatePortalApp); ok {
		r0 = rf(ctx, portalAppUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.UpdatePortalApp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdatePortalApp) error); ok {
		r1 = rf(ctx, portalAppUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePortalAppsFirstDateSurpassed provides a mock function with given fields: ctx, firstDateSurpassedUpdate
func (_m *MockIDBWriter) UpdatePortalAppsFirstDateSurpassed(ctx context.Context, firstDateSurpassedUpdate types.UpdateFirstDateSurpassed) (map[string]string, error) {
	ret := _m.Called(ctx, firstDateSurpassedUpdate)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateFirstDateSurpassed) (map[string]string, error)); ok {
		return rf(ctx, firstDateSurpassedUpdate)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateFirstDateSurpassed) map[string]string); ok {
		r0 = rf(ctx, firstDateSurpassedUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateFirstDateSurpassed) error); ok {
		r1 = rf(ctx, firstDateSurpassedUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, user
func (_m *MockIDBWriter) UpdateUser(ctx context.Context, user types.UpdateUser) (*types.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateUser) (*types.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateUser) *types.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateUser) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteAccountUser provides a mock function with given fields: ctx, createUser, _a2
func (_m *MockIDBWriter) WriteAccountUser(ctx context.Context, createUser types.CreateAccountUserAccess, _a2 time.Time) (map[string]types.UserID, error) {
	ret := _m.Called(ctx, createUser, _a2)

	var r0 map[string]types.UserID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.CreateAccountUserAccess, time.Time) (map[string]types.UserID, error)); ok {
		return rf(ctx, createUser, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.CreateAccountUserAccess, time.Time) map[string]types.UserID); ok {
		r0 = rf(ctx, createUser, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]types.UserID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.CreateAccountUserAccess, time.Time) error); ok {
		r1 = rf(ctx, createUser, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteBlockedContract provides a mock function with given fields: ctx, blockedContract
func (_m *MockIDBWriter) WriteBlockedContract(ctx context.Context, blockedContract types.BlockedContract) (map[string]string, error) {
	ret := _m.Called(ctx, blockedContract)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedContract) (map[string]string, error)); ok {
		return rf(ctx, blockedContract)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.BlockedContract) map[string]string); ok {
		r0 = rf(ctx, blockedContract)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.BlockedContract) error); ok {
		r1 = rf(ctx, blockedContract)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockIDBWriter creates a new instance of MockIDBWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIDBWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIDBWriter {
	mock := &MockIDBWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
