package dbclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pokt-foundation/portal-db/types"
	"github.com/stretchr/testify/suite"
)

var (
	testCtx       = context.Background()
	mockTimestamp = time.Date(2022, 11, 11, 11, 11, 11, 0, time.UTC)
)

type DBClientTestSuite struct {
	suite.Suite
	client *DBClient
}

func Test_RunDBClientTestSuite(t *testing.T) {
	suite.Run(t, new(DBClientTestSuite))
}

// SetupSuite runs before each test suite run
func (ts *DBClientTestSuite) SetupSuite() {
	err := ts.initDBClient()
	ts.NoError(err)
}

// Initializes
func (ts *DBClientTestSuite) initDBClient() error {
	config := Config{
		baseURL: "http://localhost:8080",
		apiKey:  "test_api_key_6789",
		version: V1,
		retries: 1,
		timeout: 10 * time.Second,
	}

	client, err := NewDBClient(config)
	if err != nil {
		return err
	}
	ts.client = client

	return nil
}

func (ts *DBClientTestSuite) Test_GetBlockchains() {
	tests := []struct {
		name                string
		expectedBlockchains []*types.Blockchain
		err                 error
	}{
		{
			name: "Should fetch all blockchains in the database",
			expectedBlockchains: []*types.Blockchain{
				{
					ID:                "0001",
					Altruist:          "https://test:test_93uhfniu23f8@shared-test2.nodes.pokt.network:12345",
					Blockchain:        "pokt-mainnet",
					Description:       "POKT Network Mainnet",
					EnforceResult:     "JSON",
					Network:           "POKT-mainnet",
					Ticker:            "POKT",
					BlockchainAliases: []string{"pokt-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []types.Redirect{
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc1.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc2.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      `{}`,
						Path:      "/v1/query/height",
						ResultKey: "height",
						Allowance: 1,
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
				{
					ID:                "0021",
					Altruist:          "https://test:test_u32fh239hf@shared-test2.nodes.eth.network:12345",
					Blockchain:        "eth-mainnet",
					ChainID:           "1",
					ChainIDCheck:      `{\"method\":\"eth_chainId\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
					Description:       "Ethereum Mainnet",
					EnforceResult:     "JSON",
					Network:           "ETH-1",
					Ticker:            "ETH",
					BlockchainAliases: []string{"eth-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []types.Redirect{
						{
							Alias:          "eth-mainnet",
							Domain:         "test-rpc.testnet.eth.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      `{\"method\":\"eth_blockNumber\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
						ResultKey: "result",
						Allowance: 5,
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
		},
	}

	for _, test := range tests {
		blockchains, err := ts.client.GetBlockchains(testCtx)
		ts.ErrorIs(test.err, err)
		ts.Equal(test.expectedBlockchains, blockchains)
	}
}

func (ts *DBClientTestSuite) Test_GetBlockchain() {
	tests := []struct {
		name               string
		blockchainID       string
		expectedBlockchain *types.Blockchain
		err                error
	}{
		{
			name:         "Should fetch one blockchain by ID",
			blockchainID: "0021",
			expectedBlockchain: &types.Blockchain{
				ID:                "0021",
				Altruist:          "https://test:test_u32fh239hf@shared-test2.nodes.eth.network:12345",
				Blockchain:        "eth-mainnet",
				ChainID:           "1",
				ChainIDCheck:      `{\"method\":\"eth_chainId\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
				Description:       "Ethereum Mainnet",
				EnforceResult:     "JSON",
				Network:           "ETH-1",
				Ticker:            "ETH",
				BlockchainAliases: []string{"eth-mainnet"},
				LogLimitBlocks:    100_000,
				Active:            true,
				Redirects: []types.Redirect{
					{
						Alias:          "eth-mainnet",
						Domain:         "test-rpc.testnet.eth.network",
						LoadBalancerID: "test_lb_34gg4g43g34g5hh",
					},
				},
				SyncCheckOptions: types.SyncCheckOptions{
					Body:      `{\"method\":\"eth_blockNumber\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
					ResultKey: "result",
					Allowance: 5,
				},
				CreatedAt: mockTimestamp,
				UpdatedAt: mockTimestamp,
			},
		},
		{
			name:         "Should fail if the blockchain does not exist in the DB",
			blockchainID: "666",
			err:          fmt.Errorf("Response not OK. 404 Not Found: blockchain not found"),
		},
	}

	for _, test := range tests {
		blockchain, err := ts.client.GetBlockchain(testCtx, test.blockchainID)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedBlockchain, blockchain)
	}
}

// func (ts *DBClientTestSuite) Test_GetApplications() {
// 	tests := []struct {
// 		name                 string
// 		expectedApplications []*types.Application
// 		err                  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		applications, err := ts.client.GetApplications(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetApplicationByID() {
// 	tests := []struct {
// 		name                string
// 		expectedApplication *types.Application
// 		err                 error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		applicationByID, err := ts.client.GetApplicationByID(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetApplicationsByUserID() {
// 	tests := []struct {
// 		name                 string
// 		expectedApplications []*types.Application
// 		err                  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		applicationsByUserID, err := ts.client.GetApplicationsByUserID(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetLoadBalancers() {
// 	tests := []struct {
// 		name                  string
// 		expectedLoadBalancers []*types.LoadBalancer
// 		err                   error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancers, err := ts.client.GetLoadBalancers(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetLoadBalancerByID() {
// 	tests := []struct {
// 		name                 string
// 		expectedLoadBalancer *types.LoadBalancer
// 		err                  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancerByID, err := ts.client.GetLoadBalancerByID(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetLoadBalancersByUserID() {
// 	tests := []struct {
// 		name                  string
// 		expectedLoadBalancers []*types.LoadBalancer
// 		err                   error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancersByUserID, err := ts.client.GetLoadBalancersByUserID(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetPayPlans() {
// 	tests := []struct {
// 		name             string
// 		expectedPayPlans []*types.PayPlan
// 		err              error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		payPlans, err := ts.client.GetPayPlans(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_GetPayPlanByType() {
// 	tests := []struct {
// 		name            string
// 		expectedPayPlan *types.PayPlan
// 		err             error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		payPlanByType, err := ts.client.GetPayPlanByType(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_CreateBlockchain() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		blockchain, err := ts.client.CreateBlockchain(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_CreateRedirect() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		redirect, err := ts.client.CreateRedirect(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_CreateApplication() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		application, err := ts.client.CreateApplication(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_CreateLoadBalancer() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancer, err := ts.client.CreateLoadBalancer(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_ActivateBlockchain() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		blockchain, err := ts.client.ActivateBlockchain(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_UpdateApplication() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		application, err := ts.client.UpdateApplication(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_UpdateAppFirstDateSurpassed() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		appFirstDateSurpassed, err := ts.client.UpdateAppFirstDateSurpassed(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_RemoveApplication() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		application, err := ts.client.RemoveApplication(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_UpdateLoadBalancer() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancer, err := ts.client.UpdateLoadBalancer(testCtx)
// 	}
// }

// func (ts *DBClientTestSuite) Test_RemoveLoadBalancer() {
// 	tests := []struct {
// 		name string
// 		err  error
// 	}{
// 		{},
// 	}

// 	for _, test := range tests {
// 		loadBalancers, err := ts.client.RemoveLoadBalancer(testCtx)
// 	}
// }
