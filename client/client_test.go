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

func (ts *DBClientTestSuite) Test_GetApplications() {
	tests := []struct {
		name                 string
		expectedApplications []*types.Application
		err                  error
	}{
		{
			name: "Should fetch all applications in the database",
			expectedApplications: []*types.Application{
				{
					ID:     "test_app_47hfnths73j2se",
					UserID: "test_user_1dbffbdfeeb225",
					Name:   "pokt_app_123",
					URL:    "https://test.app123.io",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_34715cae753e67c75fbb340442e7de8e",
						ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
						ApplicationSignature: "test_89a3af6a587aec02cfade6f5000424c2",
						ClientPublicKey:      "test_1dc39a2e5a84a35bf030969a0b3231f7",
						PrivateKey:           "test_d2ce53f115f4ecb2208e9188800a85cf",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
						SecretKeyRequired: true,
					},
					Limit: types.AppLimit{
						PayPlan: types.PayPlan{Type: types.FreetierV0, Limit: 250_000},
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
				{
					ID:     "test_app_5hdf7sh23jd828",
					UserID: "test_user_04228205bd261a",
					Name:   "pokt_app_456",
					URL:    "https://test.app456.io",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
						ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
						ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
						ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
						PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
						SecretKeyRequired: false,
					},
					Limit: types.AppLimit{
						PayPlan:     types.PayPlan{Type: types.Enterprise},
						CustomLimit: 2_000_000,
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
						Half:          false,
						ThreeQuarters: true,
						Full:          true,
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
		},
	}

	for _, test := range tests {
		applications, err := ts.client.GetApplications(testCtx)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedApplications, applications)
	}
}

func (ts *DBClientTestSuite) Test_GetApplicationByID() {
	tests := []struct {
		name                string
		applicationID       string
		expectedApplication *types.Application
		err                 error
	}{
		{
			name:          "Should fetch one application by ID",
			applicationID: "test_app_5hdf7sh23jd828",
			expectedApplication: &types.Application{
				ID:     "test_app_5hdf7sh23jd828",
				UserID: "test_user_04228205bd261a",
				Name:   "pokt_app_456",
				URL:    "https://test.app456.io",
				Dummy:  true,
				Status: types.InService,
				GatewayAAT: types.GatewayAAT{
					Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
					ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
					ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
					ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
					PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
				},
				GatewaySettings: types.GatewaySettings{
					SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
					SecretKeyRequired: false,
				},
				Limit: types.AppLimit{
					PayPlan:     types.PayPlan{Type: types.Enterprise},
					CustomLimit: 2_000_000,
				},
				NotificationSettings: types.NotificationSettings{
					SignedUp:      true,
					Quarter:       false,
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
			err: fmt.Errorf("Response not OK. 404 Not Found: applications not found"),
		},
	}

	for _, test := range tests {
		applicationByID, err := ts.client.GetApplicationByID(testCtx, test.applicationID)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedApplication, applicationByID)
	}
}

func (ts *DBClientTestSuite) Test_GetApplicationsByUserID() {
	tests := []struct {
		name                 string
		userID               string
		expectedApplications []*types.Application
		err                  error
	}{
		{
			name:   "Should fetch all applications for a single user ID",
			userID: "test_user_04228205bd261a",
			expectedApplications: []*types.Application{
				{
					ID:     "test_app_5hdf7sh23jd828",
					UserID: "test_user_04228205bd261a",
					Name:   "pokt_app_456",
					URL:    "https://test.app456.io",
					Dummy:  true,
					Status: types.InService,
					GatewayAAT: types.GatewayAAT{
						Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
						ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
						ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
						ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
						PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
					},
					GatewaySettings: types.GatewaySettings{
						SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
						SecretKeyRequired: false,
					},
					Limit: types.AppLimit{
						PayPlan:     types.PayPlan{Type: types.Enterprise},
						CustomLimit: 2_000_000,
					},
					NotificationSettings: types.NotificationSettings{
						SignedUp:      true,
						Quarter:       false,
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
			userID: "test_not_real_user",
			err:    fmt.Errorf("Response not OK. 404 Not Found: applications not found"),
		},
	}

	for _, test := range tests {
		applicationsByUserID, err := ts.client.GetApplicationsByUserID(testCtx, test.userID)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedApplications, applicationsByUserID)
	}
}

func (ts *DBClientTestSuite) Test_GetLoadBalancers() {
	tests := []struct {
		name                  string
		expectedLoadBalancers []*types.LoadBalancer
		err                   error
	}{
		{
			name: "Should fetch all load balancers in the database",
			expectedLoadBalancers: []*types.LoadBalancer{
				{
					ID:                "test_lb_34987u329rfn23f",
					Name:              "pokt_app_123",
					UserID:            "test_user_1dbffbdfeeb225",
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
					Applications: []*types.Application{
						{
							ID:     "test_app_47hfnths73j2se",
							UserID: "test_user_1dbffbdfeeb225",
							Name:   "pokt_app_123",
							URL:    "https://test.app123.io",
							Dummy:  true,
							Status: types.InService,
							GatewayAAT: types.GatewayAAT{
								Address:              "test_34715cae753e67c75fbb340442e7de8e",
								ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
								ApplicationSignature: "test_89a3af6a587aec02cfade6f5000424c2",
								ClientPublicKey:      "test_1dc39a2e5a84a35bf030969a0b3231f7",
								PrivateKey:           "test_d2ce53f115f4ecb2208e9188800a85cf",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
								SecretKeyRequired: true,
							},
							Limit: types.AppLimit{
								PayPlan: types.PayPlan{Type: types.FreetierV0, Limit: 250_000},
							},
							NotificationSettings: types.NotificationSettings{
								SignedUp:      true,
								Quarter:       false,
								Half:          false,
								ThreeQuarters: true,
								Full:          true,
							},
							CreatedAt: mockTimestamp,
							UpdatedAt: mockTimestamp,
						},
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_1dbffbdfeeb225", Email: "owner1@test.com", Accepted: true},
						{RoleName: "ADMIN", UserID: "test_user_admin1234", Email: "admin1@test.com", Accepted: true},
						{RoleName: "MEMBER", UserID: "test_user_member1234", Email: "member1@test.com", Accepted: true},
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
				{
					ID:                "test_lb_34gg4g43g34g5hh",
					Name:              "test_lb_redirect",
					UserID:            "test_user_redirect233344",
					RequestTimeout:    5_000,
					Gigastake:         false,
					GigastakeRedirect: false,
					StickyOptions: types.StickyOptions{
						Duration:      "20",
						StickyOrigins: []string{"test-extension://", "test-extension2://"},
						StickyMax:     600,
						Stickiness:    false,
					},
					Applications: []*types.Application{nil},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_redirect233344", Email: "owner3@test.com", Accepted: true},
						{RoleName: "MEMBER", UserID: "test_user_member5678", Email: "member2@test.com", Accepted: true},
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
				{
					ID:                "test_lb_3890ru23jfi32fj",
					Name:              "pokt_app_456",
					UserID:            "test_user_04228205bd261a",
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "40",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
					Applications: []*types.Application{
						{
							ID:     "test_app_5hdf7sh23jd828",
							UserID: "test_user_04228205bd261a",
							Name:   "pokt_app_456",
							URL:    "https://test.app456.io",
							Dummy:  true,
							Status: types.InService,
							GatewayAAT: types.GatewayAAT{
								Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
								ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
								ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
								ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
								PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
								SecretKeyRequired: false,
							},
							Limit: types.AppLimit{
								PayPlan:     types.PayPlan{Type: types.Enterprise},
								CustomLimit: 2_000_000,
							},
							NotificationSettings: types.NotificationSettings{
								SignedUp:      true,
								Quarter:       false,
								Half:          false,
								ThreeQuarters: true,
								Full:          true,
							},
							CreatedAt: mockTimestamp,
							UpdatedAt: mockTimestamp,
						},
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_04228205bd261a", Email: "owner2@test.com", Accepted: true},
						{RoleName: "ADMIN", UserID: "test_user_admin5678", Email: "admin2@test.com", Accepted: true},
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
		},
	}

	for _, test := range tests {
		loadBalancers, err := ts.client.GetLoadBalancers(testCtx)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedLoadBalancers, loadBalancers)
	}
}

func (ts *DBClientTestSuite) Test_GetLoadBalancerByID() {
	tests := []struct {
		name                 string
		loadBalancerID       string
		expectedLoadBalancer *types.LoadBalancer
		err                  error
	}{
		{
			name:           "Should fetch one load balancer by ID",
			loadBalancerID: "test_lb_3890ru23jfi32fj",
			expectedLoadBalancer: &types.LoadBalancer{
				ID:                "test_lb_3890ru23jfi32fj",
				Name:              "pokt_app_456",
				UserID:            "test_user_04228205bd261a",
				RequestTimeout:    5_000,
				Gigastake:         true,
				GigastakeRedirect: true,
				StickyOptions: types.StickyOptions{
					Duration:      "40",
					StickyOrigins: []string{"chrome-extension://"},
					StickyMax:     400,
					Stickiness:    true,
				},
				Applications: []*types.Application{
					{
						ID:     "test_app_5hdf7sh23jd828",
						UserID: "test_user_04228205bd261a",
						Name:   "pokt_app_456",
						URL:    "https://test.app456.io",
						Dummy:  true,
						Status: types.InService,
						GatewayAAT: types.GatewayAAT{
							Address:              "test_558c0225c7019e14ccf2e7379ad3eb50",
							ApplicationPublicKey: "test_96c981db344ab6920b7e87853838e285",
							ApplicationSignature: "test_1272a8ab4cbbf636f09bf4fa5395b885",
							ClientPublicKey:      "test_d709871777b89ed3051190f229ea3f01",
							PrivateKey:           "test_53e50765d8bc1fb41b3b0065dd8094de",
						},
						GatewaySettings: types.GatewaySettings{
							SecretKey:         "test_90210ac4bdd3423e24877d1ff92",
							SecretKeyRequired: false,
						},
						Limit: types.AppLimit{
							PayPlan:     types.PayPlan{Type: types.Enterprise},
							CustomLimit: 2_000_000,
						},
						NotificationSettings: types.NotificationSettings{
							SignedUp:      true,
							Quarter:       false,
							Half:          false,
							ThreeQuarters: true,
							Full:          true,
						},
						CreatedAt: mockTimestamp,
						UpdatedAt: mockTimestamp,
					},
				},
				Users: []types.UserAccess{
					{RoleName: "OWNER", UserID: "test_user_04228205bd261a", Email: "owner2@test.com", Accepted: true},
					{RoleName: "ADMIN", UserID: "test_user_admin5678", Email: "admin2@test.com", Accepted: true},
				},
				CreatedAt: mockTimestamp,
				UpdatedAt: mockTimestamp,
			},
		},
		{
			name:           "Should fail if the load balancer does not exist in the DB",
			loadBalancerID: "test_not_real_load_balancer",
			err:            fmt.Errorf("Response not OK. 404 Not Found: load balancer not found"),
		},
	}

	for _, test := range tests {
		loadBalancerByID, err := ts.client.GetLoadBalancerByID(testCtx, test.loadBalancerID)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedLoadBalancer, loadBalancerByID)
	}
}

func (ts *DBClientTestSuite) Test_GetLoadBalancersByUserID() {
	tests := []struct {
		name                  string
		userID                string
		expectedLoadBalancers []*types.LoadBalancer
		err                   error
	}{
		{
			name:   "Should fetch all load balancers for a single user ID",
			userID: "test_user_1dbffbdfeeb225",
			expectedLoadBalancers: []*types.LoadBalancer{
				{
					ID:                "test_lb_34987u329rfn23f",
					Name:              "pokt_app_123",
					UserID:            "test_user_1dbffbdfeeb225",
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
					Applications: []*types.Application{
						{
							ID:     "test_app_47hfnths73j2se",
							UserID: "test_user_1dbffbdfeeb225",
							Name:   "pokt_app_123",
							URL:    "https://test.app123.io",
							Dummy:  true,
							Status: types.InService,
							GatewayAAT: types.GatewayAAT{
								Address:              "test_34715cae753e67c75fbb340442e7de8e",
								ApplicationPublicKey: "test_11b8d394ca331d7c7a71ca1896d630f6",
								ApplicationSignature: "test_89a3af6a587aec02cfade6f5000424c2",
								ClientPublicKey:      "test_1dc39a2e5a84a35bf030969a0b3231f7",
								PrivateKey:           "test_d2ce53f115f4ecb2208e9188800a85cf",
							},
							GatewaySettings: types.GatewaySettings{
								SecretKey:         "test_40f482d91a5ef2300ebb4e2308c",
								SecretKeyRequired: true,
							},
							Limit: types.AppLimit{
								PayPlan: types.PayPlan{Type: types.FreetierV0, Limit: 250_000},
							},
							NotificationSettings: types.NotificationSettings{
								SignedUp:      true,
								Quarter:       false,
								Half:          false,
								ThreeQuarters: true,
								Full:          true,
							},
							CreatedAt: mockTimestamp,
							UpdatedAt: mockTimestamp,
						},
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_1dbffbdfeeb225", Email: "owner1@test.com", Accepted: true},
						{RoleName: "ADMIN", UserID: "test_user_admin1234", Email: "admin1@test.com", Accepted: true},
						{RoleName: "MEMBER", UserID: "test_user_member1234", Email: "member1@test.com", Accepted: true},
					},
					CreatedAt: mockTimestamp,
					UpdatedAt: mockTimestamp,
				},
			},
		},
		{
			name:   "Should fail if the user does not have any load balancers in the DB",
			userID: "test_not_real_user",
			// TODO - fix this error string in PHD, should say `load balancers`
			err: fmt.Errorf("Response not OK. 404 Not Found: load balancer not found"),
		},
	}

	for _, test := range tests {
		loadBalancersByUserID, err := ts.client.GetLoadBalancersByUserID(testCtx, test.userID)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedLoadBalancers, loadBalancersByUserID)
	}
}

func (ts *DBClientTestSuite) Test_GetPayPlans() {
	tests := []struct {
		name             string
		expectedPayPlans []*types.PayPlan
		err              error
	}{
		{
			name: "Should fetch all pay plans in the DB",
			expectedPayPlans: []*types.PayPlan{
				{Type: types.Enterprise, Limit: 0},
				{Type: types.FreetierV0, Limit: 250000},
				{Type: types.PayAsYouGoV0, Limit: 0},
				{Type: types.TestPlan10K, Limit: 10000},
				{Type: types.TestPlan90k, Limit: 90000},
				{Type: types.TestPlanV0, Limit: 100},
			},
		},
	}

	for _, test := range tests {
		payPlans, err := ts.client.GetPayPlans(testCtx)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedPayPlans, payPlans)
	}
}

func (ts *DBClientTestSuite) Test_GetPayPlanByType() {
	tests := []struct {
		name            string
		payPlanType     types.PayPlanType
		expectedPayPlan *types.PayPlan
		err             error
	}{
		{

			name:        "Should fetch a single pay plans by type",
			payPlanType: types.FreetierV0,
			expectedPayPlan: &types.PayPlan{
				Type: types.FreetierV0, Limit: 250000,
			},
		},
	}

	for _, test := range tests {
		payPlanByType, err := ts.client.GetPayPlanByType(testCtx, test.payPlanType)
		ts.Equal(err, test.err)
		ts.Equal(test.expectedPayPlan, payPlanByType)
	}
}

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
