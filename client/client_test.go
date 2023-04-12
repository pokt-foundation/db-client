package dbclient

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pokt-foundation/portal-db/types"
	"github.com/stretchr/testify/suite"
)

var (
	testCtx       = context.Background()
	mockTimestamp = time.Date(2022, 11, 11, 11, 11, 11, 0, time.UTC)
)

type DBClientTestSuite struct {
	suite.Suite
	client IDBClient
	mu     sync.Mutex
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
		BaseURL: "http://localhost:8080",
		APIKey:  "test_api_key_6789",
		Version: V1,
		Retries: 1,
		Timeout: 10 * time.Second,
	}

	client, err := NewDBClient(config)
	if err != nil {
		return err
	}
	ts.client = client

	return nil
}

func boolPointer(value bool) *bool {
	return &value
}

func intPointer(value int) *int {
	return &value
}

// Runs all the read-only endpoint tests first to compare to test DB seed data only
// ie. not yet including data written to the test DB by the test suite
func (ts *DBClientTestSuite) Test_ReadTests() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	// TODO: UPDATE payplans once they're set in the portal-db

	ts.Run("Test_GetBlockchains", func() {
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
						Altruist:          "https://test_pocket:auth123456@altruist-0001.com:1234", // pragma: allowlist secret
						Blockchain:        "mainnet",
						Description:       "Pocket Network Mainnet",
						EnforceResult:     "JSON",
						Path:              "/v1/query/height",
						Ticker:            "POKT",
						BlockchainAliases: []string{"mainnet"},
						Active:            true,
						Redirects: []types.Redirect{
							{
								Alias:          "altruist-0001",
								Domain:         "pokt-rpc.gateway.pokt.network",
								LoadBalancerID: "test_app_1",
							},
						},
						SyncCheckOptions: types.SyncCheckOptions{
							Body:      `{"id":1,"jsonrpc":"2.0","method":"query"}`,
							ResultKey: "result.sync_info",
							Allowance: 1,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "0021",
						Altruist:          "https://test_pocket:auth123456@altruist-0021.com:1234", // pragma: allowlist secret
						Blockchain:        "eth-mainnet",
						ChainID:           "1",
						ChainIDCheck:      `{"method":"eth_chainId","id":1,"jsonrpc":"2.0"}`,
						Description:       "Ethereum Mainnet",
						EnforceResult:     "JSON",
						Ticker:            "ETH",
						BlockchainAliases: []string{"eth-mainnet"},
						LogLimitBlocks:    100_000,
						Active:            true,
						Redirects: []types.Redirect{
							{
								Alias:          "altruist-0021",
								Domain:         "eth-rpc.gateway.pokt.network",
								LoadBalancerID: "test_app_3",
							},
						},
						SyncCheckOptions: types.SyncCheckOptions{
							Body:      `{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
							ResultKey: "result",
							Allowance: 5,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "0040",
						Altruist:          "https://test_pocket:auth123456@altruist-0040.com:1234", // pragma: allowlist secret
						Blockchain:        "harmony-0",
						ChainID:           "",
						Description:       "Harmony Shard 0",
						EnforceResult:     "JSON",
						Ticker:            "HMY",
						BlockchainAliases: []string{"harmony-0"},
						Active:            true,
						Redirects: []types.Redirect{
							{
								Alias:          "altruist-0040",
								Domain:         "hmy-rpc.gateway.pokt.network",
								LoadBalancerID: "test_app_3",
							},
						},
						SyncCheckOptions: types.SyncCheckOptions{
							Body:      `{"id":1,"jsonrpc":"2.0","method":"hmy_blockNumber","params":[]}`,
							ResultKey: "result",
							Allowance: 8,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "0053",
						Altruist:          "https://test_pocket:auth123456@altruist-0053.com:1234", // pragma: allowlist secret
						Blockchain:        "optimism-mainnet",
						Description:       "Optimism Mainnet",
						EnforceResult:     "JSON",
						Ticker:            "OP",
						BlockchainAliases: []string{"optimism-mainnet"},
						LogLimitBlocks:    100_000,
						Active:            true,
						Redirects: []types.Redirect{
							{
								Alias:          "altruist-0053",
								Domain:         "op-rpc.gateway.pokt.network",
								LoadBalancerID: "test_app_2",
							},
						},
						SyncCheckOptions: types.SyncCheckOptions{
							Body:      `{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
							ResultKey: "result",
							Allowance: 2,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "0064",
						Altruist:          "https://test_pocket:auth123456@altruist-0064.com:1234", // pragma: allowlist secret
						Blockchain:        "sui-testnet",
						Description:       "Sui Testnet",
						EnforceResult:     "JSON",
						Ticker:            "SUI-TESTNET",
						BlockchainAliases: []string{"sui-testnet"},
						LogLimitBlocks:    100_000,
						RequestTimeout:    60_000,
						Active:            false,
						Redirects:         nil,
						SyncCheckOptions: types.SyncCheckOptions{
							Body:      `{"id":1,"jsonrpc":"2.0","method":"sui_blockNumber","params":[]}`,
							ResultKey: "result",
							Allowance: 7,
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
			cmp.Equal(test.expectedBlockchains, blockchains)
		}
	})

	ts.Run("Test_GetBlockchain", func() {
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
					Altruist:          "https://test_pocket:auth123456@altruist-0021.com:1234", // pragma: allowlist secret
					Blockchain:        "eth-mainnet",
					ChainID:           "1",
					ChainIDCheck:      `{"method":"eth_chainId","id":1,"jsonrpc":"2.0"}`,
					Description:       "Ethereum Mainnet",
					EnforceResult:     "JSON",
					Ticker:            "ETH",
					BlockchainAliases: []string{"eth-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []types.Redirect{
						{
							Alias:          "altruist-0021",
							Domain:         "eth-rpc.gateway.pokt.network",
							LoadBalancerID: "test_app_3",
						},
					},
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      `{"id":1,"jsonrpc":"2.0","method":"eth_blockNumber","params":[]}`,
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
			blockchain, err := ts.client.GetBlockchainByID(testCtx, test.blockchainID)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedBlockchain, blockchain)
		}
	})

	ts.Run("Test_GetApplications", func() {
		tests := []struct {
			name                 string
			expectedApplications []*types.Application
			err                  error
		}{
			{
				name: "Should fetch all applications in the database",
				expectedApplications: []*types.Application{
					{
						ID:                 "test_protocol_app_1",
						UserID:             "auth0|james_holden",
						Name:               "pokt_app_123",
						FirstDateSurpassed: mockTimestamp,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_34715cae753e67c75fbb340442e7de8e",
							ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
							ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
							ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
							PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
							Version:              "0.0.1",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
							SecretKeyRequired:   true,
							WhitelistOrigins:    []string{"https://test.com"},
							WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
							WhitelistContracts: []types.WhitelistContracts{
								{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
							},
							WhitelistMethods: []types.WhitelistMethods{
								{BlockchainID: "0001", Methods: []string{"GET"}},
							},
							WhitelistBlockchains: []string{"0053"},
						},
						Limit: types.AppLimit{
							PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
						},

						NotificationSettings: types.NotificationSettings{
							SignedUp:      false,
							Quarter:       true,
							Half:          false,
							ThreeQuarters: true,
							Full:          true,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                 "test_protocol_app_2",
						UserID:             "auth0|ellen_ripley",
						Name:               "pokt_app_456",
						FirstDateSurpassed: mockTimestamp,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_8237c72345f12d1b1a8b64a1a7f66fa4",
							ApplicationPublicKey: "test_8237c72345f12d1b1a8b64a1a7f66fa4",
							ApplicationSignature: "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
							ClientPublicKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
							PrivateKey:           "test_2e83c836a29b423a47d8e18c779fd422",
							Version:              "0.0.1",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:           "test_9c9e3b193cfba5348f93bb2f3e3fb794",
							SecretKeyRequired:   false,
							WhitelistOrigins:    []string{"https://example.com"},
							WhitelistUserAgents: []string{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36"},
							WhitelistContracts: []types.WhitelistContracts{
								{BlockchainID: "0064", Contracts: []string{"0x0987654321abcdef"}},
							},
							WhitelistMethods: []types.WhitelistMethods{
								{BlockchainID: "0064", Methods: []string{"POST"}},
							},
							WhitelistBlockchains: []string{"0021"},
						},
						Limit: types.AppLimit{
							PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000},
						},
						NotificationSettings: types.NotificationSettings{
							SignedUp:      false,
							Quarter:       false,
							Half:          true,
							ThreeQuarters: false,
							Full:          true,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                 "test_protocol_app_3",
						UserID:             "auth0|chrisjen_avasarala",
						Name:               "pokt_app_789",
						FirstDateSurpassed: mockTimestamp,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_b5e07928fc80083c13ad0201b81bae9b",
							ApplicationPublicKey: "test_f608500e4fe3e09014fe2411b4a560b5",
							ApplicationSignature: "test_c3cd8be16ba32e24dd49fdb0247fc9b8",
							ClientPublicKey:      "test_328a9cf1b35085eeaa669aa858f6fba9",
							PrivateKey:           "test_8663e187c19f3c6e27317eab4ed6d7d5",
							Version:              "0.0.1",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:         "test_9f48b13e2bc5fd31ab367841f11495c1",
							SecretKeyRequired: false,
						},
						Limit: types.AppLimit{
							PayPlan: types.PayPlan{Type: types.PayPlanType("startup_plan"), Limit: 500},
						},
						NotificationSettings: types.NotificationSettings{
							SignedUp:      false,
							Quarter:       false,
							Half:          false,
							ThreeQuarters: false,
							Full:          false,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                 "test_protocol_app_4",
						UserID:             "auth0|chrisjen_avasarala",
						Name:               "pokt_app_789",
						FirstDateSurpassed: mockTimestamp,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_eb2e5bcba557cfe8fa76fd7fff54f9d1",
							ApplicationPublicKey: "test_f6a5d8690ecb669865bd752b7796a920",
							ApplicationSignature: "test_cf05cf9bb26111c548e88fb6157af708",
							ClientPublicKey:      "test_6ee5ea553408f0895923fd1569dc5072",
							PrivateKey:           "test_838d29d61a65401f7d56d084cb6e4783",
							Version:              "0.0.1",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:         "test_9f48b13e2bc5fd31ab367841f11495c1",
							SecretKeyRequired: false,
						},
						Limit: types.AppLimit{
							PayPlan: types.PayPlan{Type: types.PayPlanType("startup_plan"), Limit: 500},
						},
						NotificationSettings: types.NotificationSettings{
							SignedUp:      false,
							Quarter:       false,
							Half:          false,
							ThreeQuarters: false,
							Full:          false,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
		}

		for _, test := range tests {
			applications, err := ts.client.GetApplications(testCtx)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedApplications, applications)
		}
	})

	ts.Run("Test_GetApplicationByID", func() {
		tests := []struct {
			name                string
			applicationID       string
			expectedApplication *types.Application
			err                 error
		}{
			{
				name:          "Should fetch one application by ID",
				applicationID: "test_protocol_app_1",
				expectedApplication: &types.Application{
					ID:                 "test_protocol_app_1",
					UserID:             "auth0|james_holden",
					Name:               "pokt_app_123",
					FirstDateSurpassed: mockTimestamp,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_34715cae753e67c75fbb340442e7de8e",
						ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
						ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
						ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
						PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
						Version:              "0.0.1",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
						SecretKeyRequired:   true,
						WhitelistOrigins:    []string{"https://test.com"},
						WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
						WhitelistContracts: []types.WhitelistContracts{
							{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
						},
						WhitelistMethods: []types.WhitelistMethods{
							{BlockchainID: "0001", Methods: []string{"GET"}},
						},
						WhitelistBlockchains: []string{"0053"},
					},
					Limit: types.AppLimit{
						PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
					},

					NotificationSettings: types.NotificationSettings{
						SignedUp:      false,
						Quarter:       true,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
			{
				name:          "Should fail if the application does not exist in the DB",
				applicationID: "test_not_real_app",
				// TODO - fix this error string in PHD, should say `application`
				err: fmt.Errorf("Response not OK. 404 Not Found: portal app not found for app ID test_not_real_app"),
			},
		}

		for _, test := range tests {
			applicationByID, err := ts.client.GetApplicationByID(testCtx, test.applicationID)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedApplication, applicationByID)
		}
	})

	ts.Run("Test_GetApplicationsByUserID", func() {
		tests := []struct {
			name                 string
			userID               string
			expectedApplications []*types.Application
			err                  error
		}{
			{
				name:   "Should fetch all applications for a single user ID",
				userID: "auth0|james_holden",
				expectedApplications: []*types.Application{
					{
						ID:                 "test_protocol_app_1",
						UserID:             "auth0|james_holden",
						Name:               "pokt_app_123",
						FirstDateSurpassed: mockTimestamp,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_34715cae753e67c75fbb340442e7de8e",
							ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
							ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
							ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
							PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
							Version:              "0.0.1",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
							SecretKeyRequired:   true,
							WhitelistOrigins:    []string{"https://test.com"},
							WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
							WhitelistContracts: []types.WhitelistContracts{
								{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
							},
							WhitelistMethods: []types.WhitelistMethods{
								{BlockchainID: "0001", Methods: []string{"GET"}},
							},
							WhitelistBlockchains: []string{"0053"},
						},
						Limit: types.AppLimit{
							PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
						},

						NotificationSettings: types.NotificationSettings{
							SignedUp:      false,
							Quarter:       true,
							Half:          false,
							ThreeQuarters: true,
							Full:          true,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
			{
				name:   "Should fail if the user does not have any applications in the DB",
				userID: "auth0|bernard_marx",
				err:    fmt.Errorf(`Response not OK. 404 Not Found: portal app not found for user ID user_11`),
			},
		}

		for _, test := range tests {
			applicationsByUserID, err := ts.client.GetApplicationsByUserID(testCtx, test.userID)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedApplications, applicationsByUserID)
		}
	})

	ts.Run("Test_GetLoadBalancers", func() {
		tests := []struct {
			name                  string
			expectedLoadBalancers []*types.LoadBalancer
			err                   error
		}{
			{
				name: "Should fetch all load balancers in the database",
				expectedLoadBalancers: []*types.LoadBalancer{
					{
						ID:                "test_app_1",
						Name:              "pokt_app_123",
						UserID:            "auth0|james_holden",
						RequestTimeout:    5_000,
						Gigastake:         true,
						GigastakeRedirect: true,
						StickyOptions: types.StickyOptions{
							Duration:      "60",
							StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
							StickyMax:     300,
							Stickiness:    true,
						},
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_1",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_1",
								UserID:             "auth0|james_holden",
								Name:               "pokt_app_123",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
									ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
									PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
									SecretKeyRequired:   true,
									WhitelistOrigins:    []string{"https://test.com"},
									WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
									WhitelistContracts: []types.WhitelistContracts{
										{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
									},
									WhitelistMethods: []types.WhitelistMethods{
										{BlockchainID: "0001", Methods: []string{"GET"}},
									},
									WhitelistBlockchains: []string{"0053"},
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
								},

								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       true,
									Half:          false,
									ThreeQuarters: true,
									Full:          true,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_8", Email: "rick.deckard456@test.com", Accepted: false},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "test_app_2",
						Name:              "pokt_app_456",
						UserID:            "auth0|ellen_ripley",
						RequestTimeout:    10_000,
						Gigastake:         false,
						GigastakeRedirect: false,
						StickyOptions: types.StickyOptions{
							Duration:      "30",
							StickyOrigins: []string{"https://example.com", "https://test.com"},
							StickyMax:     600,
							Stickiness:    true,
						},
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_2",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_2",
								UserID:             "auth0|ellen_ripley",
								Name:               "pokt_app_456",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_8237c72345f12d1b1a8b64a1a7f66fa4",
									ApplicationPublicKey: "test_8237c72345f12d1b1a8b64a1a7f66fa4",
									ApplicationSignature: "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
									ClientPublicKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
									PrivateKey:           "test_2e83c836a29b423a47d8e18c779fd422",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:           "test_9c9e3b193cfba5348f93bb2f3e3fb794",
									SecretKeyRequired:   false,
									WhitelistOrigins:    []string{"https://example.com"},
									WhitelistUserAgents: []string{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36"},
									WhitelistContracts: []types.WhitelistContracts{
										{BlockchainID: "0064", Contracts: []string{"0x0987654321abcdef"}},
									},
									WhitelistMethods: []types.WhitelistMethods{
										{BlockchainID: "0064", Methods: []string{"POST"}},
									},
									WhitelistBlockchains: []string{"0021"},
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000},
								},
								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       false,
									Half:          true,
									ThreeQuarters: false,
									Full:          true,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_3", Email: "ellen.ripley789@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_4", Email: "ulfric.stormcloak123@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_9", Email: "tyrion.lannister789@test.com", Accepted: false},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
					{
						ID:                "test_app_3",
						Name:              "pokt_app_789",
						UserID:            "auth0|chrisjen_avasarala",
						RequestTimeout:    10_000,
						Gigastake:         false,
						GigastakeRedirect: false,
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_3",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_3",
								UserID:             "auth0|chrisjen_avasarala",
								Name:               "pokt_app_789",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_b5e07928fc80083c13ad0201b81bae9b",
									ApplicationPublicKey: "test_f608500e4fe3e09014fe2411b4a560b5",
									ApplicationSignature: "test_c3cd8be16ba32e24dd49fdb0247fc9b8",
									ClientPublicKey:      "test_328a9cf1b35085eeaa669aa858f6fba9",
									PrivateKey:           "test_8663e187c19f3c6e27317eab4ed6d7d5",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:         "test_9f48b13e2bc5fd31ab367841f11495c1",
									SecretKeyRequired: false,
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("startup_plan"), Limit: 500},
								},
								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       false,
									Half:          false,
									ThreeQuarters: false,
									Full:          false,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
							{
								ID:                 "test_protocol_app_4",
								UserID:             "auth0|chrisjen_avasarala",
								Name:               "pokt_app_789",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_eb2e5bcba557cfe8fa76fd7fff54f9d1",
									ApplicationPublicKey: "test_f6a5d8690ecb669865bd752b7796a920",
									ApplicationSignature: "test_cf05cf9bb26111c548e88fb6157af708",
									ClientPublicKey:      "test_6ee5ea553408f0895923fd1569dc5072",
									PrivateKey:           "test_838d29d61a65401f7d56d084cb6e4783",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:         "test_9f48b13e2bc5fd31ab367841f11495c1",
									SecretKeyRequired: false,
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("startup_plan"), Limit: 500},
								},
								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       false,
									Half:          false,
									ThreeQuarters: false,
									Full:          false,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_5", Email: "chrisjen.avasarala1@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_6", Email: "amos.burton789@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_10", Email: "daenerys.targaryen123@test.com", Accepted: false},
							{RoleName: types.RoleMember, UserID: "user_7", Email: "frodo.baggins123@test.com", Accepted: true},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
		}

		for _, test := range tests {
			loadBalancers, err := ts.client.GetLoadBalancers(testCtx)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedLoadBalancers, loadBalancers)
		}
	})

	ts.Run("Test_GetLoadBalancerByID", func() {
		tests := []struct {
			name                 string
			loadBalancerID       string
			expectedLoadBalancer *types.LoadBalancer
			err                  error
		}{
			{
				name:           "Should fetch one load balancer by ID",
				loadBalancerID: "test_app_1",
				expectedLoadBalancer: &types.LoadBalancer{
					ID:                "test_app_1",
					Name:              "pokt_app_123",
					UserID:            "auth0|james_holden",
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
					Integrations: types.AccountIntegrations{
						CovalentAPIKeyFree: "covalent_api_key_1",
					},
					Applications: []*types.Application{
						{
							ID:                 "test_protocol_app_1",
							UserID:             "auth0|james_holden",
							Name:               "pokt_app_123",
							FirstDateSurpassed: mockTimestamp,
							GatewayAAT: types.GatewayAAT{
								Address:              "test_34715cae753e67c75fbb340442e7de8e",
								ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
								ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
								ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
								PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
								Version:              "0.0.1",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
								SecretKeyRequired:   true,
								WhitelistOrigins:    []string{"https://test.com"},
								WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
								WhitelistContracts: []types.WhitelistContracts{
									{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
								},
								WhitelistMethods: []types.WhitelistMethods{
									{BlockchainID: "0001", Methods: []string{"GET"}},
								},
								WhitelistBlockchains: []string{"0053"},
							},
							Limit: types.AppLimit{
								PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
							},

							NotificationSettings: types.NotificationSettings{
								SignedUp:      false,
								Quarter:       true,
								Half:          false,
								ThreeQuarters: true,
								Full:          true,
							},
							CreatedAt: mockTimestamp,
							UpdatedAt: mockTimestamp,
						},
					},
					Users: []types.UserAccess{
						{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
						{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
						{RoleName: types.RoleAdmin, UserID: "user_8", Email: "rick.deckard456@test.com", Accepted: false},
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
			{
				name:           "Should fail if the load balancer does not exist in the DB",
				loadBalancerID: "test_not_real_load_balancer",
				err:            fmt.Errorf("Response not OK. 404 Not Found: portal app not found for load balancer ID test_not_real_load_balancer"),
			},
		}

		for _, test := range tests {
			loadBalancerByID, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedLoadBalancer, loadBalancerByID)
		}
	})

	ts.Run("Test_GetLoadBalancersByUserID", func() {
		tests := []struct {
			name                  string
			userID                string
			expectedLoadBalancers []*types.LoadBalancer
			roleNameFilter        types.RoleName
			err                   error
		}{
			{
				name:   "Should fetch all load balancers for a single user ID when no filter provided",
				userID: "auth0|james_holden",
				expectedLoadBalancers: []*types.LoadBalancer{
					{
						ID:                "test_app_1",
						Name:              "pokt_app_123",
						UserID:            "auth0|james_holden",
						RequestTimeout:    5_000,
						Gigastake:         true,
						GigastakeRedirect: true,
						StickyOptions: types.StickyOptions{
							Duration:      "60",
							StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
							StickyMax:     300,
							Stickiness:    true,
						},
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_1",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_1",
								UserID:             "auth0|james_holden",
								Name:               "pokt_app_123",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
									ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
									PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
									SecretKeyRequired:   true,
									WhitelistOrigins:    []string{"https://test.com"},
									WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
									WhitelistContracts: []types.WhitelistContracts{
										{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
									},
									WhitelistMethods: []types.WhitelistMethods{
										{BlockchainID: "0001", Methods: []string{"GET"}},
									},
									WhitelistBlockchains: []string{"0053"},
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
								},

								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       true,
									Half:          false,
									ThreeQuarters: true,
									Full:          true,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_8", Email: "rick.deckard456@test.com", Accepted: false},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
			{
				name:           "Should fetch all load balancers for a single user ID and role when a valid filter provided",
				userID:         "auth0|paul_atreides",
				roleNameFilter: types.RoleAdmin,
				expectedLoadBalancers: []*types.LoadBalancer{
					{
						ID:                "test_app_1",
						Name:              "pokt_app_123",
						UserID:            "auth0|james_holden",
						RequestTimeout:    5_000,
						Gigastake:         true,
						GigastakeRedirect: true,
						StickyOptions: types.StickyOptions{
							Duration:      "60",
							StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
							StickyMax:     300,
							Stickiness:    true,
						},
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_1",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_1",
								UserID:             "auth0|james_holden",
								Name:               "pokt_app_123",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationPublicKey: "test_34715cae753e67c75fbb340442e7de8e",
									ApplicationSignature: "test_1dc39a2e5a84a35bf030969a0b3231f7",
									ClientPublicKey:      "test_89a3af6a587aec02cfade6f5000424c2",
									PrivateKey:           "test_11b8d394ca331d7c7a71ca1896d630f6",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:           "test_40f482d91a5ef2300ebb4e2308c",
									SecretKeyRequired:   true,
									WhitelistOrigins:    []string{"https://test.com"},
									WhitelistUserAgents: []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
									WhitelistContracts: []types.WhitelistContracts{
										{BlockchainID: "0001", Contracts: []string{"0x1234567890abcdef"}},
									},
									WhitelistMethods: []types.WhitelistMethods{
										{BlockchainID: "0001", Methods: []string{"GET"}},
									},
									WhitelistBlockchains: []string{"0053"},
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("basic_plan"), Limit: 1000},
								},

								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       true,
									Half:          false,
									ThreeQuarters: true,
									Full:          true,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
							{RoleName: types.RoleAdmin, UserID: "user_8", Email: "rick.deckard456@test.com", Accepted: false},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
			{
				name:           "Should fail if an invalid role name provided as a filter",
				userID:         "auth0|james_holden",
				roleNameFilter: types.RoleName("fake"),
				err:            fmt.Errorf("invalid role name filter"),
			},
			{
				name:                  "Should fail if the user does not have any load balancers in the DB",
				userID:                "auth0|bernard_marx",
				expectedLoadBalancers: []*types.LoadBalancer{},
			},
		}

		for _, test := range tests {
			filter := &test.roleNameFilter
			if test.roleNameFilter == "" {
				filter = nil
			}

			loadBalancersByUserID, err := ts.client.GetLoadBalancersByUserID(testCtx, test.userID, filter)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedLoadBalancers, loadBalancersByUserID)
		}
	})

	ts.Run("Test_GetPendingLoadBalancersByUserID", func() {
		tests := []struct {
			name                  string
			userID                string
			expectedLoadBalancers []*types.LoadBalancer
			err                   error
		}{
			{
				name:   "Should fetch all pending load balancers for a single user ID",
				userID: "user_9",
				expectedLoadBalancers: []*types.LoadBalancer{
					{
						ID:                "test_app_2",
						Name:              "pokt_app_456",
						UserID:            "auth0|ellen_ripley",
						RequestTimeout:    10_000,
						Gigastake:         false,
						GigastakeRedirect: false,
						StickyOptions: types.StickyOptions{
							Duration:      "30",
							StickyOrigins: []string{"https://example.com", "https://test.com"},
							StickyMax:     600,
							Stickiness:    true,
						},
						Integrations: types.AccountIntegrations{
							CovalentAPIKeyFree: "covalent_api_key_2",
						},
						Applications: []*types.Application{
							{
								ID:                 "test_protocol_app_2",
								UserID:             "auth0|ellen_ripley",
								Name:               "pokt_app_456",
								FirstDateSurpassed: mockTimestamp,
								GatewayAAT: types.GatewayAAT{
									Address:              "test_8237c72345f12d1b1a8b64a1a7f66fa4",
									ApplicationPublicKey: "test_8237c72345f12d1b1a8b64a1a7f66fa4",
									ApplicationSignature: "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
									ClientPublicKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
									PrivateKey:           "test_2e83c836a29b423a47d8e18c779fd422",
									Version:              "0.0.1",
								},
								GatewaySettings: types.GatewaySettings{
									SecretKey:           "test_9c9e3b193cfba5348f93bb2f3e3fb794",
									SecretKeyRequired:   false,
									WhitelistOrigins:    []string{"https://example.com"},
									WhitelistUserAgents: []string{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36"},
									WhitelistContracts: []types.WhitelistContracts{
										{BlockchainID: "0064", Contracts: []string{"0x0987654321abcdef"}},
									},
									WhitelistMethods: []types.WhitelistMethods{
										{BlockchainID: "0064", Methods: []string{"POST"}},
									},
									WhitelistBlockchains: []string{"0021"},
								},
								Limit: types.AppLimit{
									PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000},
								},
								NotificationSettings: types.NotificationSettings{
									SignedUp:      false,
									Quarter:       false,
									Half:          true,
									ThreeQuarters: false,
									Full:          true,
								},
								CreatedAt: mockTimestamp,
								UpdatedAt: mockTimestamp,
							},
						},
						Users: []types.UserAccess{
							{RoleName: types.RoleOwner, UserID: "user_3", Email: "ellen.ripley789@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_4", Email: "ulfric.stormcloak123@test.com", Accepted: true},
							{RoleName: types.RoleMember, UserID: "user_9", Email: "tyrion.lannister789@test.com", Accepted: false},
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
			},
			{
				name:                  "Should fail if the email does not have any pending load balancers in the DB",
				userID:                "test_not_real",
				expectedLoadBalancers: []*types.LoadBalancer{},
			},
		}

		for _, test := range tests {
			pendingLoadBalancersByUserID, err := ts.client.GetPendingLoadBalancersByUserID(testCtx, test.userID)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedLoadBalancers, pendingLoadBalancersByUserID)
		}
	})

	ts.Run("Test_GetLoadBalancersCountByUserID", func() {
		tests := []struct {
			name          string
			userID        string
			expectedCount int
			err           error
		}{
			{
				name:          "Should return the number of loadBalancers owned by email",
				userID:        "user_1",
				expectedCount: 1,
			},
			{
				name:          "return 0 if there's no loadbalancer binded with the email",
				userID:        "random@test.com",
				expectedCount: 0,
			},
		}

		for _, test := range tests {
			loadBalancerCount, err := ts.client.GetLoadBalancersCountByUserID(testCtx, test.userID)
			ts.Equal(test.err, err)
			ts.Equal(test.expectedCount, loadBalancerCount)
		}
	})

	ts.Run("Test_GetPayPlans", func() {
		tests := []struct {
			name             string
			expectedPayPlans []*types.PayPlan
			err              error
		}{
			{
				name: "Should fetch all pay plans in the DB",
				expectedPayPlans: []*types.PayPlan{ // TODO: UPDATE payplans once they're set in the portal-db
					{Type: types.PayPlanType("basic_plan"), Limit: 1000},
					{Type: types.PayPlanType("pro_plan"), Limit: 5000},
					{Type: types.PayPlanType("enterprise_plan"), Limit: 10000},
					{Type: types.PayPlanType("developer_plan"), Limit: 100},
					{Type: types.PayPlanType("startup_plan"), Limit: 500},
				},
			},
		}

		for _, test := range tests {
			payPlans, err := ts.client.GetPayPlans(testCtx)
			ts.Equal(test.err, err)
			cmp.Equal(test.expectedPayPlans, payPlans)
		}
	})

	ts.Run("Test_GetPayPlanByType", func() {
		tests := []struct {
			name            string
			payPlanType     types.PayPlanType
			expectedPayPlan *types.PayPlan
			err             error
		}{
			{

				name:        "Should fetch a single pay plan by type",
				payPlanType: types.PayPlanType("basic_plan"),
				expectedPayPlan: &types.PayPlan{
					Type: types.PayPlanType("basic_plan"), Limit: 1000,
				},
			},
			{
				name:        "Should fail if passed a pay plan type that is not in the DB",
				payPlanType: types.PayPlanType("not_a_real_plan"),
				err:         fmt.Errorf("Response not OK. 404 Not Found: plan not found for type not_a_real_plan"),
			},
		}

		for _, test := range tests {
			payPlanByType, err := ts.client.GetPayPlanByType(testCtx, test.payPlanType)
			ts.Equal(test.err, err)
			ts.Equal(test.expectedPayPlan, payPlanByType)
		}
	})
	ts.Run("Test_GetUserPermissionsByUserID", func() {
		tests := []struct {
			name                string
			userID              types.UserID
			expectedPermissions *types.UserPermissions
			err                 error
		}{
			{

				name:   "Should fetch a single users load balancer permissions",
				userID: "auth0|james_holden",
				expectedPermissions: &types.UserPermissions{
					UserID: "auth0|james_holden",
					LoadBalancers: map[types.LoadBalancerID]types.LoadBalancerPermissions{
						"test_app_1": {
							RoleName:    types.RoleOwner,
							Permissions: []types.PermissionsEnum{types.ReadEndpoint, types.WriteEndpoint, types.DeleteEndpoint, types.TransferEndpoint},
						},
					},
				},
			},
			{

				name:   "Should fetch another single users load balancer permissions",
				userID: "auth0|ellen_ripley",
				expectedPermissions: &types.UserPermissions{
					UserID: "auth0|ellen_ripley",
					LoadBalancers: map[types.LoadBalancerID]types.LoadBalancerPermissions{
						"test_app_2": {
							RoleName:    types.RoleOwner,
							Permissions: []types.PermissionsEnum{types.ReadEndpoint, types.WriteEndpoint, types.DeleteEndpoint, types.TransferEndpoint},
						},
					},
				},
			},
			{
				name:   "Should return an empty list if the user exists but has not accepted their invite",
				userID: "auth0|rick_deckard",
				expectedPermissions: &types.UserPermissions{
					UserID:        "auth0|rick_deckard",
					LoadBalancers: map[types.LoadBalancerID]types.LoadBalancerPermissions{},
				},
			},
			{
				name:   "Should fail if the user does not have any permissions",
				userID: "test_user_hey_who_am_i_wow",
				err:    fmt.Errorf("Response not OK. 404 Not Found: user not found for provider user ID test_user_hey_who_am_i_wow"),
			},
		}

		for _, test := range tests {
			permissionsByUserID, err := ts.client.GetUserPermissionsByUserID(testCtx, test.userID)
			ts.Equal(test.err, err)
			ts.Equal(test.expectedPermissions, permissionsByUserID)
		}
	})
}

// Runs all the write endpoint tests after the read tests
// This ensures the write tests to not modify the seed data expected by the read tests
func (ts *DBClientTestSuite) Test_WriteTests() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.Run("Test_CreateBlockchain", func() {
		tests := []struct {
			name                        string
			blockchainInput, blockchain types.Blockchain
			err                         error
		}{
			{
				name: "Should create a single blockchain in the DB",
				blockchainInput: types.Blockchain{
					ID:                "003",
					Altruist:          "https://test:test_fg332f@shared-test3.nodes.pol.network:12345", // pragma: allowlist secret
					Blockchain:        "pol-mainnet",
					Description:       "Polygon Mainnet",
					EnforceResult:     "JSON",
					Ticker:            "POL",
					BlockchainAliases: []string{"pol-mainnet"},
					LogLimitBlocks:    100000,
					Active:            false,
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      "{}",
						ResultKey: "result",
						Allowance: 3,
					},
				},
				blockchain: types.Blockchain{
					ID:                "003",
					Altruist:          "https://test:test_fg332f@shared-test3.nodes.pol.network:12345", // pragma: allowlist secret
					Blockchain:        "pol-mainnet",
					Description:       "Polygon Mainnet",
					EnforceResult:     "JSON",
					Ticker:            "POL",
					BlockchainAliases: []string{"pol-mainnet"},
					LogLimitBlocks:    100000,
					Active:            false,
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      "{}",
						ResultKey: "result",
						Allowance: 3,
					},
				},
			},
			{
				name: "Should fail if attempting to create a duplicate record",
				blockchainInput: types.Blockchain{
					ID:                "003",
					Altruist:          "https://test:test_fg332f@shared-test3.nodes.pol.network:12345", // pragma: allowlist secret
					Blockchain:        "pol-mainnet",
					Description:       "Polygon Mainnet",
					EnforceResult:     "JSON",
					Ticker:            "POL",
					BlockchainAliases: []string{"pol-mainnet"},
					LogLimitBlocks:    100000,
					Active:            false,
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      "{}",
						ResultKey: "result",
						Allowance: 3,
					},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: CreateLegacyBlockchain failed: error chain already exists for chain ID '003'"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.CreateBlockchain(testCtx, test.blockchainInput)
			ts.Equal(test.err, err)
			if test.err == nil {
				blockchain, err := ts.client.GetBlockchainByID(testCtx, test.blockchainInput.ID)
				ts.Equal(test.err, err)
				ts.Equal(test.blockchain.ID, blockchain.ID)
				ts.Equal(test.blockchain.Altruist, blockchain.Altruist)
				ts.Equal(test.blockchain.Blockchain, blockchain.Blockchain)
				ts.Equal(test.blockchain.ChainID, blockchain.ChainID)
				ts.Equal(test.blockchain.ChainIDCheck, blockchain.ChainIDCheck)
				ts.Equal(test.blockchain.Description, blockchain.Description)
				ts.Equal(test.blockchain.EnforceResult, blockchain.EnforceResult)
				ts.Equal(test.blockchain.Path, blockchain.Path)
				ts.Equal(test.blockchain.Ticker, blockchain.Ticker)
				ts.Equal(test.blockchain.BlockchainAliases, blockchain.BlockchainAliases)
				ts.Equal(test.blockchain.LogLimitBlocks, blockchain.LogLimitBlocks)
				ts.Equal(test.blockchain.RequestTimeout, blockchain.RequestTimeout)
				ts.Equal(test.blockchain.Active, blockchain.Active)
				ts.Equal(test.blockchain.SyncCheckOptions, blockchain.SyncCheckOptions)
				ts.NotEmpty(blockchain.CreatedAt)
				ts.NotEmpty(blockchain.UpdatedAt)
			}
		}
	})

	ts.Run("Test_CreateRedirect", func() {
		tests := []struct {
			name          string
			redirectInput types.Redirect
			redirects     []types.Redirect
			err           error
		}{
			{
				name: "Should create a new redirect for an existing blockchain in the DB",
				redirectInput: types.Redirect{
					BlockchainID:   "0001",
					Alias:          "test-mainnet-2",
					Domain:         "test-rpc2.testnet.pokt.network",
					LoadBalancerID: "test_app_1",
					UpdatedAt:      mockTimestamp,
					CreatedAt:      mockTimestamp,
				},
				redirects: []types.Redirect{
					{
						Alias:          "altruist-0001",
						Domain:         "test-rpc1.testnet.pokt.network",
						LoadBalancerID: "test_app_1",
					},
					{
						Alias:          "test-mainnet-2",
						Domain:         "test-rpc2.testnet.pokt.network",
						LoadBalancerID: "test_app_1",
					},
				},
			},
		}

		for _, test := range tests {
			_, err := ts.client.CreateBlockchainRedirect(testCtx, test.redirectInput)
			ts.Equal(test.err, err)
			if test.err == nil {
				blockchain, err := ts.client.GetBlockchainByID(testCtx, test.redirectInput.BlockchainID)
				ts.Equal(test.err, err)
				ts.Len(blockchain.Redirects, len(test.redirects))
				for i, redirect := range blockchain.Redirects {
					cmp.Equal(test.redirects[i].Alias, redirect.Alias)
					cmp.Equal(test.redirects[i].Domain, redirect.Domain)
					cmp.Equal(test.redirects[i].LoadBalancerID, redirect.LoadBalancerID)
				}
			}
		}
	})

	ts.Run("Test_CreateLoadBalancer", func() {
		tests := []struct {
			name                            string
			loadBalancerInput, loadBalancer types.LoadBalancer
			err                             error
		}{
			{
				name: "Should create a single loadBalancer in the DB",
				loadBalancerInput: types.LoadBalancer{
					Name:              "pokt_app_7899",
					UserID:            "auth0|ellen_ripley",
					RequestTimeout:    5000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "70",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
					Users: []types.UserAccess{
						{
							UserID:   "auth0|ellen_ripley",
							RoleName: types.RoleOwner,
							Email:    "ellen.ripley789@test.com",
							Accepted: true,
						},
					},
					Applications: []*types.Application{
						{
							ID:                 "c58cdba6",
							UserID:             "auth0|ellen_ripley",
							Name:               "pokt_app_7899",
							FirstDateSurpassed: mockTimestamp,
							GatewayAAT: types.GatewayAAT{
								Address:              "test_8237c72345f12d1b1a8b64a1a7f66fa4",
								ApplicationPublicKey: "test_8237c72345f12d1b1a8b64a1a7f66fa4",
								ApplicationSignature: "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
								ClientPublicKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
								PrivateKey:           "test_2e83c836a29b423a47d8e18c779fd422",
								Version:              "0.0.1",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey:           "test_9c9e3b193cfba5348f93bb2f3e3fb794",
								SecretKeyRequired:   false,
								WhitelistOrigins:    []string{"https://example.com"},
								WhitelistUserAgents: []string{"Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36"},
								WhitelistContracts: []types.WhitelistContracts{
									{BlockchainID: "0064", Contracts: []string{"0x0987654321abcdef"}},
								},
								WhitelistMethods: []types.WhitelistMethods{
									{BlockchainID: "0064", Methods: []string{"POST"}},
								},
								WhitelistBlockchains: []string{"0021"},
							},
							Limit: types.AppLimit{
								PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000},
							},
							NotificationSettings: types.NotificationSettings{
								SignedUp:      false,
								Quarter:       false,
								Half:          true,
								ThreeQuarters: false,
								Full:          true,
							},
							CreatedAt: mockTimestamp,
							UpdatedAt: mockTimestamp,
						},
					},
				},
				loadBalancer: types.LoadBalancer{
					Name:              "pokt_app_7899",
					UserID:            "auth0|ellen_ripley",
					RequestTimeout:    5000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "70",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
					Users: []types.UserAccess{
						{
							UserID:   "user_3",
							RoleName: types.RoleOwner,
							Email:    "ellen.ripley789@test.com",
							Accepted: true,
						},
						{
							UserID:   "user_2",
							RoleName: types.RoleMember,
							Email:    "paul.atreides456@test.com",
							Accepted: true,
						},
						{
							UserID:   "user_4",
							RoleName: types.RoleMember,
							Email:    "ulfric.stormcloak123@test.com",
							Accepted: true,
						},
						{
							UserID:   "user_9",
							RoleName: types.RoleMember,
							Email:    "tyrion.lannister789@test.com",
							Accepted: false,
						},
					},
					Applications: []*types.Application{
						{
							UserID: "auth0|ellen_ripley",
							Name:   "pokt_app_7899",
							GatewayAAT: types.GatewayAAT{
								Address:              "test_8237c72345f12d1b1a8b64a1a7f66fa4",
								ApplicationPublicKey: "test_8237c72345f12d1b1a8b64a1a7f66fa4",
								ApplicationSignature: "test_f48d33b30ddaf60a1e5bb50d2ba8da5a",
								ClientPublicKey:      "test_04c71d90a92f40416b6f1d7d8af17e02",
								PrivateKey:           "test_2e83c836a29b423a47d8e18c779fd422",
								Version:              "0.0.1",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey: "test_9c9e3b193cfba5348f93bb2f3e3fb794",
							},
							Limit: types.AppLimit{
								PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000},
							},
							NotificationSettings: types.NotificationSettings{
								SignedUp:      true,
								Quarter:       false,
								Half:          false,
								ThreeQuarters: true,
								Full:          true,
							},
						},
					},
				},
			},
		}

		for _, test := range tests {
			createdLB, err := ts.client.CreateLoadBalancer(testCtx, test.loadBalancerInput)
			ts.Equal(test.err, err)
			if test.err == nil {
				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, createdLB.ID)
				test.loadBalancer.Applications[0].ID = loadBalancer.Applications[0].ID
				test.loadBalancer.Applications[0].CreatedAt = loadBalancer.Applications[0].CreatedAt
				test.loadBalancer.Applications[0].UpdatedAt = loadBalancer.Applications[0].UpdatedAt
				ts.Equal(test.err, err)
				ts.Equal(createdLB.ID, loadBalancer.ID)
				ts.Equal(test.loadBalancer.UserID, loadBalancer.UserID)
				ts.Equal(test.loadBalancer.Name, loadBalancer.Name)
				ts.Equal(test.loadBalancer.UserID, loadBalancer.UserID)
				ts.Equal(test.loadBalancer.RequestTimeout, loadBalancer.RequestTimeout)
				ts.Equal(test.loadBalancer.Gigastake, loadBalancer.Gigastake)
				ts.Equal(test.loadBalancer.GigastakeRedirect, loadBalancer.GigastakeRedirect)
				ts.Equal(test.loadBalancer.ApplicationIDs, loadBalancer.ApplicationIDs)
				ts.Equal(test.loadBalancer.Applications, loadBalancer.Applications)
				ts.Equal(test.loadBalancer.StickyOptions, loadBalancer.StickyOptions)
				ts.Equal(test.loadBalancer.Users, loadBalancer.Users)
				ts.NotEmpty(loadBalancer.CreatedAt)
				ts.NotEmpty(loadBalancer.UpdatedAt)
			}
		}
	})

	ts.Run("Test_CreateLoadBalancerUser", func() {
		tests := []struct {
			name              string
			loadBalancerID    string
			user              types.UserAccess
			loadBalancerUsers []types.UserAccess
			err               error
		}{
			{
				name:           "Should add a single user to an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				user: types.UserAccess{
					RoleName: types.RoleMember,
					Email:    "member_new@test.com",
				},
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "test_user_1dbffbdfeeb225", Email: "owner1@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "test_user_admin1234", Email: "admin1@test.com", Accepted: true},
					{RoleName: types.RoleMember, UserID: "test_user_member1234", Email: "member1@test.com", Accepted: true},
					{RoleName: types.RoleMember, UserID: "", Email: "member_new@test.com", Accepted: false},
				},
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "sir_not_appearing_in_this_film",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: portal app not found for load balancer ID sir_not_appearing_in_this_film"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.CreateLoadBalancerUser(testCtx, test.loadBalancerID, test.user)
			ts.Equal(test.err, err)
			if test.err == nil {
				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
				ts.Equal(test.err, err)
				cmp.Equal(test.loadBalancerUsers, loadBalancer.Users)
			}
		}
	})

	ts.Run("Test_ActivateBlockchain", func() {
		tests := []struct {
			name         string
			blockchainID string
			active       bool
			err          error
		}{
			{
				name:         "Should set a blockchain's active field to false",
				blockchainID: "0001",
				active:       false,
			},
			{
				name:         "Should set a blockchain's active field to true",
				blockchainID: "0001",
				active:       true,
			},
			{
				name:         "Should fail if blockchain cannot be found",
				blockchainID: "5440",
				err:          fmt.Errorf("Response not OK. 500 Internal Server Error: blockchain not found for chain ID 5440"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.ActivateBlockchain(testCtx, test.blockchainID, test.active)
			ts.Equal(test.err, err)
			if test.err == nil {
				blockchain, err := ts.client.GetBlockchainByID(testCtx, test.blockchainID)
				ts.Equal(test.err, err)
				ts.Equal(test.active, blockchain.Active)
			}
		}
	})

	ts.Run("Test_UpdateLoadBalancer", func() {
		tests := []struct {
			name                   string
			applicationID          string
			applicationUpdate      types.UpdateApplication
			applicationAfterUpdate types.Application
			err                    error
		}{
			{
				name:          "Should update a single application in the DB",
				applicationID: "test_app_2",
				applicationUpdate: types.UpdateApplication{
					Name: "pokt_app_updated_lb",
					GatewaySettings: &types.UpdateGatewaySettings{
						WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
						WhitelistUserAgents:  []string{"test-agent1"},
						WhitelistContracts:   []types.WhitelistContracts{{BlockchainID: "01", Contracts: []string{"test-contract1"}}},
						WhitelistMethods:     []types.WhitelistMethods{{BlockchainID: "01", Methods: []string{"test-method1"}}},
						WhitelistBlockchains: []string{"test-chain1"},
					},
					NotificationSettings: &types.UpdateNotificationSettings{SignedUp: boolPointer(false), Quarter: boolPointer(true), Half: boolPointer(true), ThreeQuarters: boolPointer(false), Full: boolPointer(false)},
					Limit:                &types.AppLimit{PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000}},
				},
				applicationAfterUpdate: types.Application{
					Name: "pokt_app_updated_lb",
					GatewaySettings: types.GatewaySettings{
						SecretKey:            "test_9c9e3b193cfba5348f93bb2f3e3fb794",
						SecretKeyRequired:    false,
						WhitelistOrigins:     []string{"test-origin1", "test-origin2"},
						WhitelistUserAgents:  []string{"test-agent1"},
						WhitelistContracts:   []types.WhitelistContracts{{BlockchainID: "01", Contracts: []string{"test-contract1"}}},
						WhitelistMethods:     []types.WhitelistMethods{{BlockchainID: "01", Methods: []string{"test-method1"}}},
						WhitelistBlockchains: []string{"test-chain1"},
					},
					NotificationSettings: types.NotificationSettings{SignedUp: false, Quarter: true, Half: true, ThreeQuarters: false, Full: false},
					Limit:                types.AppLimit{PayPlan: types.PayPlan{Type: types.PayPlanType("pro_plan"), Limit: 5000}},
				},
			},
			{
				name:          "Should fail if application cannot be found",
				applicationID: "test_app_fake",
				err:           fmt.Errorf("Response not OK. 404 Not Found: portal app not found for load balancer ID test_app_fake"),
			},
		}

		for _, test := range tests {
			createdApp, err := ts.client.UpdateLoadBalancer(testCtx, test.applicationID, test.applicationUpdate)
			ts.Equal(test.err, err)
			if err == nil {
				// Get the app inside the loadbalancer
				appID := createdApp.Applications[0].ID
				application, err := ts.client.GetApplicationByID(testCtx, appID)
				ts.NoError(err)
				ts.Equal(test.applicationAfterUpdate.Name, application.Name)
				ts.Equal(test.applicationAfterUpdate.GatewaySettings, application.GatewaySettings)
				ts.Equal(test.applicationAfterUpdate.NotificationSettings, application.NotificationSettings)
				ts.Equal(test.applicationAfterUpdate.Limit, application.Limit)
			}
		}
	})

	ts.Run("Test_UpdateAppFirstDateSurpassed", func() {
		tests := []struct {
			name         string
			update       types.UpdateFirstDateSurpassed
			expectedDate time.Time
			err          error
		}{
			{
				name: "Should update the app first date suprassed for the provided slice of app IDs",
				update: types.UpdateFirstDateSurpassed{
					ApplicationIDs:     []string{"test_app_2", "test_app_2"},
					FirstDateSurpassed: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC),
				},
				expectedDate: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC),
				err:          nil,
			},
			{
				name: "Should fail if update contains no application IDs cannot be found",
				update: types.UpdateFirstDateSurpassed{
					ApplicationIDs:     []string{},
					FirstDateSurpassed: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC),
				},
				err: fmt.Errorf("Response not OK. 400 Bad Request: no application IDs on input"),
			},
			{
				name: "Should fail if application cannot be found",
				update: types.UpdateFirstDateSurpassed{
					ApplicationIDs:     []string{"9000"},
					FirstDateSurpassed: time.Date(2022, time.December, 13, 5, 15, 0, 0, time.UTC),
				},
				err: fmt.Errorf("Response not OK. 400 Bad Request: UpdateFirstDateSurpassed failed: 9000 not found"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.UpdateAppFirstDateSurpassed(testCtx, test.update)
			ts.Equal(test.err, err)
			if test.err == nil {
				for _, appID := range test.update.ApplicationIDs {
					application, err := ts.client.GetLoadBalancerByID(testCtx, appID)
					ts.NoError(err)
					ts.Equal(test.expectedDate, application.Applications[0].FirstDateSurpassed)
				}
			}
		}
	})

	ts.Run("Test_RemoveApplication", func() {
		tests := []struct {
			name           string
			applicationID  string
			expectedStatus types.AppStatus
			err            error
		}{
			{
				name:           "should remove one application by setting its status to AWAITING_GRACE_PERIOD",
				applicationID:  "test_app_3",
				expectedStatus: types.AwaitingGracePeriod,
			},
			{
				name:          "Should fail if application cannot be found",
				applicationID: "2348",
				err:           fmt.Errorf("Response not OK. 404 Not Found: portal app not found for load balancer ID 2348"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.RemoveApplication(testCtx, test.applicationID)
			ts.Equal(test.err, err)
			if test.err == nil {
				application, err := ts.client.GetApplicationByID(testCtx, test.applicationID)
				ts.Equal(fmt.Errorf("Response not OK. 404 Not Found: portal app not found for app ID test_app_3"), err)
				ts.Nil(application)
			}
		}
	})

	ts.Run("Test_UpdateBlockchain", func() {
		tests := []struct {
			name                  string
			blockchainID          string
			blockchainUpdate      types.UpdateBlockchain
			blockchainAfterUpdate types.Blockchain
			err                   error
		}{
			{
				name:         "Should update a single loadBalancer in the DB",
				blockchainID: "0021",
				blockchainUpdate: types.UpdateBlockchain{
					Altruist:          "https://test:test_fg332f@test.nodes.pol.network:12345", // pragma: allowlist secret
					Blockchain:        "pokt-mainnet-updated",
					Description:       "POKT Network Mainnet Updated",
					RequestTimeout:    66_654,
					ResultKey:         "updated-key",
					LogLimitBlocks:    100_010,
					Ticker:            "SUCH-WOW",
					BlockchainAliases: []string{"pokt-mainnet", "another-one"},
					EnforceResult:     "JSON",
					Path:              "new-path",
					Body:              `{"new-body": "alliance"}`,
					Allowance:         intPointer(5),
				},
				blockchainAfterUpdate: types.Blockchain{
					ID:                "0021",
					Altruist:          "https://test:test_fg332f@test.nodes.pol.network:12345", // pragma: allowlist secret
					Blockchain:        "pokt-mainnet-updated",
					Description:       "POKT Network Mainnet Updated",
					EnforceResult:     "JSON",
					Ticker:            "SUCH-WOW",
					Path:              "new-path",
					BlockchainAliases: []string{"pokt-mainnet", "another-one"},
					LogLimitBlocks:    100_010,
					Active:            true,
					RequestTimeout:    66_654,
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
						Body:      `{"new-body": "alliance"}`,
						ResultKey: "updated-key",
						Allowance: 1,
					},
				},
			},
			{
				name:         "Should fail if blockchain cannot be found",
				blockchainID: "9000",
				err:          fmt.Errorf("Response not OK. 500 Internal Server Error: blockchain not found for chain ID 9000"),
			},
		}

		for _, test := range tests {
			updatedBlockchain, err := ts.client.UpdateBlockchain(testCtx, test.blockchainID, test.blockchainUpdate)
			ts.Equal(test.err, err)
			if err == nil {
				ts.Equal(test.blockchainAfterUpdate.Blockchain, updatedBlockchain.Blockchain)
				ts.Equal(test.blockchainAfterUpdate.Description, updatedBlockchain.Description)
				ts.Equal(test.blockchainAfterUpdate.RequestTimeout, updatedBlockchain.RequestTimeout)
				ts.Equal(test.blockchainAfterUpdate.LogLimitBlocks, updatedBlockchain.LogLimitBlocks)
				ts.Equal(test.blockchainAfterUpdate.Ticker, updatedBlockchain.Ticker)
				ts.Equal(test.blockchainAfterUpdate.BlockchainAliases, updatedBlockchain.BlockchainAliases)
				ts.Equal(test.blockchainAfterUpdate.EnforceResult, updatedBlockchain.EnforceResult)
				ts.Equal(test.blockchainAfterUpdate.Path, updatedBlockchain.Path)
				ts.Equal(test.blockchainAfterUpdate.SyncCheckOptions.ResultKey, updatedBlockchain.SyncCheckOptions.ResultKey)
				ts.Equal(test.blockchainAfterUpdate.SyncCheckOptions.Body, updatedBlockchain.SyncCheckOptions.Body)
			}
		}
	})

	ts.Run("Test_UpdateLoadBalancerUserRole", func() {
		tests := []struct {
			name              string
			loadBalancerID    string
			update            types.UpdateUserAccess
			loadBalancerUsers []types.UserAccess
			err               error
		}{
			{
				name:           "Should update a single user for an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_8",
					RoleName: types.RoleMember,
				},
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
					{RoleName: types.RoleMember, UserID: "user_8", Email: "", Accepted: false},
				},
			},
			{
				name:           "Should transfer ownership for an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_2",
					Email:    "james.holden123@test.com",
					RoleName: types.RoleOwner,
				},
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleAdmin, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
					{RoleName: types.RoleOwner, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
					{RoleName: types.RoleMember, UserID: "user_8", Email: "", Accepted: false},
				},
			},
			{
				name:           "Should transfer ownership back to original owner for an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_1",
					Email:    "paul.atreides456@test.com",
					RoleName: types.RoleOwner,
				},
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
					{RoleName: types.RoleMember, UserID: "user_8", Email: "", Accepted: false},
				},
			},
			{
				name:           "Should update a single unaccepted user to ADMIN for an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_8",
					RoleName: types.RoleAdmin,
				},
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_8", Email: "", Accepted: false},
				},
			},
			{
				name:           "Should fail if attempting to transfer ownership and the UpdaterEmail is not provided",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_2",
					RoleName: types.RoleOwner,
				},
				err: errOwnerRequiresUpdateEmail,
			},
			{
				name:           "Should fail if attempting to transfer ownership and the user has not accepted their invite",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_8",
					Email:    "paul.atreides456@test.com",
					RoleName: types.RoleOwner,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error cannot transfer ownership to user ID 'user_8' for account ID 'account_1' because the user has not accepted their invite"),
			},
			{
				name:           "Should fail if load balancer ID not provided",
				loadBalancerID: "",
				err:            errNoLoadBalancerID,
			},
			{
				name:           "Should fail if invalid role name provided",
				loadBalancerID: "test_app_1",
				update: types.UpdateUserAccess{
					UserID:   "user_8",
					RoleName: types.RoleName("wrong_one"),
				},
				err: errInvalidRoleName,
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "im_not_here",
				update: types.UpdateUserAccess{
					UserID:   "user_8",
					RoleName: types.RoleMember,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: portal app not found for load balancer ID im_not_here"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client.UpdateLoadBalancerUserRole(testCtx, test.loadBalancerID, test.update)
				ts.Equal(test.err, err)
				if test.err == nil {
					loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
					ts.Equal(test.err, err)
					cmp.Equal(test.loadBalancerUsers, loadBalancer.Users)
				}
			})
		}
	})

	ts.Run("Test_AcceptLoadBalancerUser", func() {
		tests := []struct {
			name                   string
			loadBalancerID, userID string
			loadBalancerUsers      []types.UserAccess
			err                    error
		}{
			{
				name:           "Should update a single user's ID and Accepted field for an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				userID:         "user_8",
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "user_1", Email: "", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_2", Email: "", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_8", Email: "", Accepted: true},
				},
			},
			{
				name:           "Should fail if load balancer ID not provided",
				loadBalancerID: "",
				userID:         "user_8",
				err:            errNoLoadBalancerID,
			},
			{
				name:           "Should fail if user ID not provided",
				loadBalancerID: "test_app_1",
				userID:         "",
				err:            errNoUserID,
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "im_not_here",
				userID:         "test_user_accept_member",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: portal app not found for load balancer ID im_not_here"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.AcceptLoadBalancerUser(testCtx, test.loadBalancerID, test.userID)
			ts.Equal(test.err, err)
			if test.err == nil {
				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
				ts.Equal(test.err, err)
				cmp.Equal(test.loadBalancerUsers, loadBalancer.Users)
			}
		}
	})

	ts.Run("Test_RemoveLoadBalancer", func() {
		tests := []struct {
			name           string
			loadBalancerID string
			expectedUserID string
			err            error
		}{
			{
				name:           "should remove one load balancer by setting its delete flag into true",
				loadBalancerID: "test_app_2",
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "9000",
				err:            fmt.Errorf("Response not OK. 404 Not Found: portal app not found for load balancer ID 9000"),
			},
		}

		for _, test := range tests {
			_, err := ts.client.RemoveLoadBalancer(testCtx, test.loadBalancerID)
			ts.Equal(test.err, err)
			if test.err == nil {
				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
				ts.Equal("Response not OK. 404 Not Found: portal app not found for load balancer ID test_app_2", err.Error())
				ts.Nil(loadBalancer)
			}
		}
	})

	ts.Run("Test_DeleteLoadBalancerUser", func() {
		tests := []struct {
			name                   string
			loadBalancerID, userID string
			loadBalancerUsers      []types.UserAccess
			err                    error
		}{
			{
				name:           "Should remove a single user from an existing load balancer in the DB",
				loadBalancerID: "test_app_1",
				userID:         "user_8",
				loadBalancerUsers: []types.UserAccess{
					{RoleName: types.RoleOwner, UserID: "user_1", Email: "james.holden123@test.com", Accepted: true},
					{RoleName: types.RoleAdmin, UserID: "user_2", Email: "paul.atreides456@test.com", Accepted: true},
				},
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "why_am_i_not_a_load_balancer",
				userID:         "user_9",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: portal app not found for load balancer ID why_am_i_not_a_load_balancer"),
			},
			{
				name:           "Should fail if load balancer ID not provided",
				loadBalancerID: "",
				userID:         "user_9",
				err:            errNoLoadBalancerID,
			},
			{
				name:           "Should fail if user ID not provided",
				loadBalancerID: "test_app_1",
				userID:         "",
				err:            errNoUserID,
			},
		}

		for _, test := range tests {
			_, err := ts.client.DeleteLoadBalancerUser(testCtx, test.loadBalancerID, test.userID)
			ts.Equal(test.err, err)
			if test.err == nil {
				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
				ts.Equal(test.err, err)
				cmp.Equal(test.loadBalancerUsers, loadBalancer.Users)
			}
		}
	})
	/* TODO: verify if this method is actually removed
	 ts.Run("Test_CreateLoadBalancerIntegration", func() {
		tests := []struct {
			name                 string
			loadBalancerID       string
			integrationsInput    types.AccountIntegrations
			expectedIntegrations types.AccountIntegrations
			err                  error
		}{
			{
				name:           "Should add new account integrations to a load balancer in the DB",
				loadBalancerID: "test_app_2",
				integrationsInput: types.AccountIntegrations{
					AccountID:          "account3",
					CovalentAPIKeyFree: "free_api_key_123",
					CovalentAPIKeyPaid: "paid_api_key_123",
				},
				expectedIntegrations: types.AccountIntegrations{
					AccountID:          "account3",
					CovalentAPIKeyFree: "free_api_key_123",
					CovalentAPIKeyPaid: "paid_api_key_123",
				},
			},
			{
				name:           "Should update existing account integrations for a load balancer in the DB",
				loadBalancerID: "test_app_1",
				integrationsInput: types.AccountIntegrations{
					AccountID:          "account1",
					CovalentAPIKeyFree: "free_api_key_456",
					CovalentAPIKeyPaid: "paid_api_key_456",
				},
				expectedIntegrations: types.AccountIntegrations{
					AccountID:          "account1",
					CovalentAPIKeyFree: "free_api_key_456",
					CovalentAPIKeyPaid: "paid_api_key_456",
				},
			},
			{
				name:           "Should fail if load balancer cannot be found",
				loadBalancerID: "sir_not_appearing_in_this_film",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: portal app not found for load balancer ID sir_not_appearing_in_this_film"),
			},
		}

		for _, test := range tests {
			updatedLB, err := ts.client.CreateLoadBalancerIntegration(testCtx, test.loadBalancerID, test.integrationsInput)
			ts.Equal(test.err, err)
			if test.err == nil {
				ts.Equal(test.expectedIntegrations, updatedLB.Integrations)

				loadBalancer, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
				ts.Equal(test.err, err)
				ts.Equal(test.expectedIntegrations, loadBalancer.Integrations)
			}
		}
	}) */
}
