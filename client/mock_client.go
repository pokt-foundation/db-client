// Code generated by mockery v2.15.0. DO NOT EDIT.

package dbclient

import (
	context "context"

	types "github.com/pokt-foundation/portal-db/types"
	mock "github.com/stretchr/testify/mock"
)

// MockIDBClient is an autogenerated mock type for the IDBClient type
type MockIDBClient struct {
	mock.Mock
}

// ActivateBlockchain provides a mock function with given fields: ctx, blockchainID, active
func (_m *MockIDBClient) ActivateBlockchain(ctx context.Context, blockchainID string, active bool) (bool, error) {
	ret := _m.Called(ctx, blockchainID, active)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) bool); ok {
		r0 = rf(ctx, blockchainID, active)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, blockchainID, active)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateApplication provides a mock function with given fields: ctx, application
func (_m *MockIDBClient) CreateApplication(ctx context.Context, application types.Application) (*types.Application, error) {
	ret := _m.Called(ctx, application)

	var r0 *types.Application
	if rf, ok := ret.Get(0).(func(context.Context, types.Application) *types.Application); ok {
		r0 = rf(ctx, application)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.Application) error); ok {
		r1 = rf(ctx, application)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateBlockchain provides a mock function with given fields: ctx, blockchain
func (_m *MockIDBClient) CreateBlockchain(ctx context.Context, blockchain types.Blockchain) (*types.Blockchain, error) {
	ret := _m.Called(ctx, blockchain)

	var r0 *types.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context, types.Blockchain) *types.Blockchain); ok {
		r0 = rf(ctx, blockchain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.Blockchain) error); ok {
		r1 = rf(ctx, blockchain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateBlockchainRedirect provides a mock function with given fields: ctx, redirect
func (_m *MockIDBClient) CreateBlockchainRedirect(ctx context.Context, redirect types.Redirect) (*types.Redirect, error) {
	ret := _m.Called(ctx, redirect)

	var r0 *types.Redirect
	if rf, ok := ret.Get(0).(func(context.Context, types.Redirect) *types.Redirect); ok {
		r0 = rf(ctx, redirect)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Redirect)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.Redirect) error); ok {
		r1 = rf(ctx, redirect)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateLoadBalancer provides a mock function with given fields: ctx, loadBalancer
func (_m *MockIDBClient) CreateLoadBalancer(ctx context.Context, loadBalancer types.LoadBalancer) (*types.LoadBalancer, error) {
	ret := _m.Called(ctx, loadBalancer)

	var r0 *types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, types.LoadBalancer) *types.LoadBalancer); ok {
		r0 = rf(ctx, loadBalancer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.LoadBalancer) error); ok {
		r1 = rf(ctx, loadBalancer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplicationByID provides a mock function with given fields: ctx, applicationID
func (_m *MockIDBClient) GetApplicationByID(ctx context.Context, applicationID string) (*types.Application, error) {
	ret := _m.Called(ctx, applicationID)

	var r0 *types.Application
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.Application); ok {
		r0 = rf(ctx, applicationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, applicationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplications provides a mock function with given fields: ctx
func (_m *MockIDBClient) GetApplications(ctx context.Context) ([]*types.Application, error) {
	ret := _m.Called(ctx)

	var r0 []*types.Application
	if rf, ok := ret.Get(0).(func(context.Context) []*types.Application); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetApplicationsByUserID provides a mock function with given fields: ctx, userID
func (_m *MockIDBClient) GetApplicationsByUserID(ctx context.Context, userID string) ([]*types.Application, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*types.Application
	if rf, ok := ret.Get(0).(func(context.Context, string) []*types.Application); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockchainByID provides a mock function with given fields: ctx, blockchainID
func (_m *MockIDBClient) GetBlockchainByID(ctx context.Context, blockchainID string) (*types.Blockchain, error) {
	ret := _m.Called(ctx, blockchainID)

	var r0 *types.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.Blockchain); ok {
		r0 = rf(ctx, blockchainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blockchainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockchains provides a mock function with given fields: ctx
func (_m *MockIDBClient) GetBlockchains(ctx context.Context) ([]*types.Blockchain, error) {
	ret := _m.Called(ctx)

	var r0 []*types.Blockchain
	if rf, ok := ret.Get(0).(func(context.Context) []*types.Blockchain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Blockchain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoadBalancerByID provides a mock function with given fields: ctx, loadBalancerID
func (_m *MockIDBClient) GetLoadBalancerByID(ctx context.Context, loadBalancerID string) (*types.LoadBalancer, error) {
	ret := _m.Called(ctx, loadBalancerID)

	var r0 *types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.LoadBalancer); ok {
		r0 = rf(ctx, loadBalancerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, loadBalancerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoadBalancers provides a mock function with given fields: ctx
func (_m *MockIDBClient) GetLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error) {
	ret := _m.Called(ctx)

	var r0 []*types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context) []*types.LoadBalancer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoadBalancersByUserID provides a mock function with given fields: ctx, userID
func (_m *MockIDBClient) GetLoadBalancersByUserID(ctx context.Context, userID string) ([]*types.LoadBalancer, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, string) []*types.LoadBalancer); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPayPlanByType provides a mock function with given fields: ctx, payPlanType
func (_m *MockIDBClient) GetPayPlanByType(ctx context.Context, payPlanType types.PayPlanType) (*types.PayPlan, error) {
	ret := _m.Called(ctx, payPlanType)

	var r0 *types.PayPlan
	if rf, ok := ret.Get(0).(func(context.Context, types.PayPlanType) *types.PayPlan); ok {
		r0 = rf(ctx, payPlanType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.PayPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.PayPlanType) error); ok {
		r1 = rf(ctx, payPlanType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPayPlans provides a mock function with given fields: ctx
func (_m *MockIDBClient) GetPayPlans(ctx context.Context) ([]*types.PayPlan, error) {
	ret := _m.Called(ctx)

	var r0 []*types.PayPlan
	if rf, ok := ret.Get(0).(func(context.Context) []*types.PayPlan); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.PayPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveApplication provides a mock function with given fields: ctx, id
func (_m *MockIDBClient) RemoveApplication(ctx context.Context, id string) (*types.Application, error) {
	ret := _m.Called(ctx, id)

	var r0 *types.Application
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.Application); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveLoadBalancer provides a mock function with given fields: ctx, id
func (_m *MockIDBClient) RemoveLoadBalancer(ctx context.Context, id string) (*types.LoadBalancer, error) {
	ret := _m.Called(ctx, id)

	var r0 *types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.LoadBalancer); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAppFirstDateSurpassed provides a mock function with given fields: ctx, updateInput
func (_m *MockIDBClient) UpdateAppFirstDateSurpassed(ctx context.Context, updateInput types.UpdateFirstDateSurpassed) ([]*types.Application, error) {
	ret := _m.Called(ctx, updateInput)

	var r0 []*types.Application
	if rf, ok := ret.Get(0).(func(context.Context, types.UpdateFirstDateSurpassed) []*types.Application); ok {
		r0 = rf(ctx, updateInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.UpdateFirstDateSurpassed) error); ok {
		r1 = rf(ctx, updateInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateApplication provides a mock function with given fields: ctx, id, update
func (_m *MockIDBClient) UpdateApplication(ctx context.Context, id string, update types.UpdateApplication) (*types.Application, error) {
	ret := _m.Called(ctx, id, update)

	var r0 *types.Application
	if rf, ok := ret.Get(0).(func(context.Context, string, types.UpdateApplication) *types.Application); ok {
		r0 = rf(ctx, id, update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, types.UpdateApplication) error); ok {
		r1 = rf(ctx, id, update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLoadBalancer provides a mock function with given fields: ctx, id, lbUpdate
func (_m *MockIDBClient) UpdateLoadBalancer(ctx context.Context, id string, lbUpdate types.UpdateLoadBalancer) (*types.LoadBalancer, error) {
	ret := _m.Called(ctx, id, lbUpdate)

	var r0 *types.LoadBalancer
	if rf, ok := ret.Get(0).(func(context.Context, string, types.UpdateLoadBalancer) *types.LoadBalancer); ok {
		r0 = rf(ctx, id, lbUpdate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoadBalancer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, types.UpdateLoadBalancer) error); ok {
		r1 = rf(ctx, id, lbUpdate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockIDBClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIDBClient creates a new instance of MockIDBClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIDBClient(t mockConstructorTestingTNewMockIDBClient) *MockIDBClient {
	mock := &MockIDBClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
