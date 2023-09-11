package dbclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/pokt-foundation/portal-db/v2/testdata"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func Test_DBClientImplementsInterfaces(t *testing.T) {
	tests := []struct {
		name   string
		client interface{}
	}{
		{
			name:   "Should verify that DBClient implements the IDBClient interface",
			client: &DBClient{},
		},
		{
			name:   "Should verify that DBClient implements the IDBReader interface",
			client: &DBClient{},
		},
		{
			name:   "Should verify that DBClient implements the IDBWriter interface",
			client: &DBClient{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if dbClient, ok := test.client.(*DBClient); ok {
				dbClient.httpClient = &http.Client{}
				dbClient.config = Config{
					BaseURL: "http://localhost",
					APIKey:  "test-key",
					Retries: 3,
					Timeout: time.Duration(3 * time.Second),
				}
			}

			switch test.name {
			case "Should verify that DBClient implements the IDBClient interface":
				assert.Implements(t, (*IDBClient)(nil), test.client)
			case "Should verify that DBClient implements the IDBReader interface":
				assert.Implements(t, (*IDBReader)(nil), test.client)
			case "Should verify that DBClient implements the IDBWriter interface":
				assert.Implements(t, (*IDBWriter)(nil), test.client)
			}
		})
	}
}

func Test_V1_E2E_PocketHTTPDBTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end to end test")
	}

	readSuite := new(phdE2EReadTestSuite)
	writeSuite := new(phdE2EWriteTestSuite)

	err := initDBClient(readSuite)
	if err != nil {
		t.Fatal("Failed to initialize the DB client for read tests:", err)
	}

	err = initDBClient(writeSuite)
	if err != nil {
		t.Fatal("Failed to initialize the DB client for write tests:", err)
	}

	suite.Run(t, readSuite)
	suite.Run(t, writeSuite)
}

// Runs all the read-only endpoint tests first to compare to test DB seed data only
// ie. not yet including data written to the test DB by the test suite
func (ts *phdE2EReadTestSuite) Test_ReadTests() {

	/* ------ Health Check Endpoint ------ */

	ts.Run("Test_HealthCheckEndpoint", func() {
		tests := []struct {
			name         string
			url          string
			expectedBody string
			expectedCode int
		}{
			{
				name:         "Should return status 200 and correct body on port 1",
				url:          fmt.Sprintf("http://localhost:%s/healthz", phdPortOne),
				expectedBody: "DB Check Done. Pocket HTTP DB is up and running!\nImage Tag: development",
				expectedCode: http.StatusOK,
			},
			{
				name:         "Should return status 200 and correct body on port 2",
				url:          fmt.Sprintf("http://localhost:%s/healthz", phdPortTwo),
				expectedBody: "DB Check Done. Pocket HTTP DB is up and running!\nImage Tag: development",
				expectedCode: http.StatusOK,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				resp, err := http.Get(test.url)
				ts.NoError(err)

				ts.Equal(test.expectedCode, resp.StatusCode)

				bodyBytes, err := io.ReadAll(resp.Body)
				ts.NoError(err)
				ts.Equal(test.expectedBody, string(bodyBytes))

				err = resp.Body.Close()
				ts.NoError(err)
			})
		}
	})

	/* ------ V2 Chain Read Tests ------ */

	ts.Run("Test_GetChainByID", func() {
		tests := []struct {
			name          string
			chainID       types.RelayChainID
			err           error
			expectedChain *types.Chain
			gigastakeApp  *types.GigastakeApp
		}{
			{
				name:          "Should get chain by ID",
				chainID:       "0001",
				expectedChain: testdata.Chains["0001"],
				gigastakeApp:  testdata.GigastakeApps["test_gigastake_app_1"],
			},
			{
				name:    "Should return error if chain ID is empty",
				chainID: "",
				err:     fmt.Errorf("no chain ID"),
			},
			{
				name:    "Should return error if chain does not exist",
				chainID: "9999",
				err:     fmt.Errorf("Response not OK. 404 Not Found: error in getChainByID: chain not found"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				chain, err := ts.client1.GetChainByID(context.Background(), test.chainID)
				ts.Equal(test.err, err)

				if err == nil {
					test.expectedChain.GigastakeApps = make(map[types.GigastakeAppID]*types.GigastakeApp)
					test.expectedChain.GigastakeApps[test.gigastakeApp.ID] = test.gigastakeApp
					ts.Equal(test.expectedChain, chain)

					chain, err = ts.client2.GetChainByID(context.Background(), test.chainID)
					ts.Equal(test.err, err)
					test.expectedChain.GigastakeApps = make(map[types.GigastakeAppID]*types.GigastakeApp)
					test.expectedChain.GigastakeApps[test.gigastakeApp.ID] = test.gigastakeApp
					ts.Equal(test.expectedChain, chain)
				}
			})
		}
	})

	ts.Run("Test_GetGigastakeAppByID", func() {
		tests := []struct {
			name           string
			gigastakeAppID types.GigastakeAppID
			err            error
			expectedApp    *types.GigastakeApp
		}{
			{
				name:           "Should get GigastakeApp by ID",
				gigastakeAppID: "test_gigastake_app_1",
				expectedApp:    testdata.GigastakeApps["test_gigastake_app_1"],
			},
			{
				name:           "Should return error if GigastakeApp ID is empty",
				gigastakeAppID: "",
				err:            fmt.Errorf("no gigastake app ID"),
			},
			{
				name:           "Should return error if GigastakeApp does not exist",
				gigastakeAppID: "9999",
				err:            fmt.Errorf("Response not OK. 404 Not Found: error in getGigastakeAppByID: gigastake app not found"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				app, err := ts.client1.GetGigastakeAppByID(context.Background(), test.gigastakeAppID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApp, app)
					app, err = ts.client2.GetGigastakeAppByID(context.Background(), test.gigastakeAppID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApp, app)
				}
			})
		}
	})

	ts.Run("Test_GetAllChains", func() {
		tests := []struct {
			name           string
			expectedChains map[types.RelayChainID]*types.Chain
			gigastakeApps  map[types.GigastakeAppID]*types.GigastakeApp
			options        ChainOptions
			err            error
		}{
			{
				name:           "Should get all active chains",
				expectedChains: filterActiveChains(testdata.Chains),
				gigastakeApps:  testdata.GigastakeApps,
			},
			{
				name:           "Should get all chains including inactive",
				expectedChains: testdata.Chains,
				gigastakeApps:  testdata.GigastakeApps,
				options: ChainOptions{
					IncludeInactive: BoolPtr(true),
				},
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				// Set pointers to Chain GigastakeApps
				for aatID, gigastakeApp := range test.gigastakeApps {
					for chainID := range gigastakeApp.ChainIDs {
						if chain, ok := test.expectedChains[chainID]; ok {
							if chain.GigastakeApps == nil {
								chain.GigastakeApps = make(map[types.GigastakeAppID]*types.GigastakeApp)
							}
							chain.GigastakeApps[aatID] = gigastakeApp
						}
					}
				}

				var chains []*types.Chain
				var err error

				if test.options.IncludeInactive != nil && *test.options.IncludeInactive {
					chains, err = ts.client1.GetAllChains(context.Background(), test.options)
				} else {
					chains, err = ts.client1.GetAllChains(context.Background())
				}
				ts.Equal(test.err, err)

				if test.err == nil {
					ts.Equal(test.expectedChains, chainsToMap(chains))

					if test.options.IncludeInactive != nil && *test.options.IncludeInactive {
						chains, err = ts.client2.GetAllChains(context.Background(), test.options)
					} else {
						chains, err = ts.client2.GetAllChains(context.Background())
					}
					ts.Equal(test.err, err)
					ts.Equal(test.expectedChains, chainsToMap(chains))
				}
			})
		}
	})

	ts.Run("Test_GetAllGigastakeApps", func() {
		tests := []struct {
			name         string
			expectedApps map[types.GigastakeAppID]*types.GigastakeApp
			err          error
		}{
			{
				name:         "Should get all GigastakeApps",
				expectedApps: testdata.GigastakeApps,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				gigastakeApps, err := ts.client1.GetAllGigastakeApps(context.Background())
				ts.Equal(test.err, err)

				if test.err == nil {
					ts.Equal(test.expectedApps, gigastakeAppsToMap(gigastakeApps))

					gigastakeApps, err = ts.client2.GetAllGigastakeApps(context.Background())
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApps, gigastakeAppsToMap(gigastakeApps))
				}
			})
		}
	})

	ts.Run("Test_GetAllGigastakeAppsByChain", func() {
		tests := []struct {
			name         string
			chainID      types.RelayChainID
			err          error
			expectedApps []*types.GigastakeApp
		}{
			{
				name:    "Should get all GigastakeApps by chain ID",
				chainID: "0001",
				expectedApps: []*types.GigastakeApp{
					testdata.GigastakeApps["test_gigastake_app_1"],
				},
			},
			{
				name:    "Should return error if chain ID is empty",
				chainID: "",
				err:     fmt.Errorf("no chain ID"),
			},
			{
				name:    "Should return error if chain does not exist",
				chainID: "9999",
				err:     fmt.Errorf("Response not OK. 400 Bad Request: error in getAllGigastakeAppsByChain: chain not found"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				apps, err := ts.client1.GetAllGigastakeAppsByChain(context.Background(), test.chainID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApps, apps)

					apps, err = ts.client2.GetAllGigastakeAppsByChain(context.Background(), test.chainID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApps, apps)
				}
			})
		}
	})

	/* ------ V2 Portal App Read Tests ------ */

	ts.Run("Test_GetPortalAppByID", func() {
		tests := []struct {
			name        string
			portalAppID types.PortalAppID
			expectedApp *types.PortalApp
			err         error
		}{
			{
				name:        "Should get portal app by ID",
				portalAppID: "test_app_1",
				expectedApp: testdata.PortalApps["test_app_1"],
			},
			{
				name:        "Should return error if app ID is empty",
				portalAppID: "",
				err:         fmt.Errorf("no portal app ID"),
			},
			{
				name:        "Should return error if app does not exist",
				portalAppID: "9999",
				err:         fmt.Errorf("Response not OK. 404 Not Found: error in getPortalAppByID: portal app not found"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalApp, err := ts.client1.GetPortalAppByID(context.Background(), test.portalAppID)
				ts.Equal(test.err, err)

				if test.err == nil {
					ts.Equal(test.expectedApp, portalApp)

					portalApp, err = ts.client2.GetPortalAppByID(context.Background(), test.portalAppID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApp, portalApp)
				}
			})
		}
	})

	ts.Run("Test_GetAllPortalApps", func() {
		tests := []struct {
			name         string
			expectedApps map[types.PortalAppID]*types.PortalApp
			err          error
		}{
			{
				name:         "Should get all portal apps",
				expectedApps: testdata.PortalApps,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalApps, err := ts.client1.GetAllPortalApps(context.Background())
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApps, portalAppsToMap(portalApps))

					portalApps, err = ts.client2.GetAllPortalApps(context.Background())
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApps, portalAppsToMap(portalApps))
				}
			})
		}
	})

	ts.Run("Test_GetPortalAppsByUser", func() {
		tests := []struct {
			name         string
			userID       types.UserID
			options      PortalAppOptions
			expectedApps map[types.PortalAppID]*types.PortalApp
			err          error
		}{
			{
				name:   "Should get all portal apps for user_4 with no role filter",
				userID: "user_4",
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_2": testdata.PortalApps["test_app_2"],
				},
			},
			{
				name:   "Should get portal apps where user_1 is OWNER",
				userID: "user_1",
				options: PortalAppOptions{
					RoleNameFilters: []types.RoleName{
						types.RoleOwner,
					},
				},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_1": testdata.PortalApps["test_app_1"],
				},
			},
			{
				name:   "Should get portal apps where user_6 is ADMIN",
				userID: "user_6",
				options: PortalAppOptions{
					RoleNameFilters: []types.RoleName{
						types.RoleAdmin,
					},
				},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_3": testdata.PortalApps["test_app_3"],
				},
			},
			{
				name:   "Should get portal apps where user_7 is MEMBER",
				userID: "user_7",
				options: PortalAppOptions{
					RoleNameFilters: []types.RoleName{
						types.RoleMember,
					},
				},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_3": testdata.PortalApps["test_app_3"],
				},
			},
			{
				name:   "Should get portal apps where user_2 is ADMIN or MEMBER",
				userID: "user_2",
				options: PortalAppOptions{
					RoleNameFilters: []types.RoleName{
						types.RoleAdmin,
						types.RoleMember,
					},
				},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_1": testdata.PortalApps["test_app_1"],
					"test_app_2": testdata.PortalApps["test_app_2"],
				},
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				var portalApps []*types.PortalApp
				var err error

				if len(test.options.RoleNameFilters) > 0 {
					portalApps, err = ts.client1.GetPortalAppsByUser(context.Background(), test.userID, test.options)
				} else {
					portalApps, err = ts.client1.GetPortalAppsByUser(context.Background(), test.userID)
				}
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApps, portalAppsToMap(portalApps))

					if len(test.options.RoleNameFilters) > 0 {
						portalApps, err = ts.client2.GetPortalAppsByUser(context.Background(), test.userID, test.options)
					} else {
						portalApps, err = ts.client2.GetPortalAppsByUser(context.Background(), test.userID)
					}
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApps, portalAppsToMap(portalApps))
				}
			})
		}
	})

	ts.Run("Test_GetPortalAppsForMiddleware", func() {
		tests := []struct {
			name         string
			expectedApps map[types.PortalAppID]*types.PortalAppLite
			err          error
		}{
			{
				name:         "Should get all portal app lites for middleware",
				expectedApps: testdata.PortalAppLites,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalAppLites, err := ts.client1.GetPortalAppsForMiddleware(context.Background())
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApps, portalAppLitesToMap(portalAppLites))

					portalAppLites, err = ts.client2.GetPortalAppsForMiddleware(context.Background())
					ts.Equal(test.err, err)
					ts.Equal(test.expectedApps, portalAppLitesToMap(portalAppLites))
				}
			})
		}
	})

	/* ------ V2 Account Read Tests ------ */

	ts.Run("Test_GetAccounts", func() {
		tests := []struct {
			name           string
			expectedAccNum int
			err            error
		}{
			{
				name:           "Should get all accounts",
				expectedAccNum: 5,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				accounts, err := ts.client1.GetAllAccounts(context.Background())
				ts.Equal(test.err, err)

				if err == nil {
					ts.Len(accounts, test.expectedAccNum)

					accounts, err = ts.client2.GetAllAccounts(context.Background())
					ts.Equal(test.err, err)
					ts.Len(accounts, test.expectedAccNum)
				}
			})
		}
	})

	ts.Run("Test_GetAccountByID", func() {
		tests := []struct {
			name        string
			accountID   types.AccountID
			expectedAcc *types.Account
			assignPlan  *types.Plan
			assignApp   *types.PortalApp
			err         error
		}{
			{
				name:        "Should get an account by its account ID",
				accountID:   "account_1",
				assignPlan:  testdata.PayPlans["basic_plan"],
				assignApp:   testdata.PortalApps["test_app_1"],
				expectedAcc: testdata.Accounts["account_1"],
			},
			{
				name:        "Should get another account by its account ID",
				accountID:   "account_2",
				assignPlan:  testdata.PayPlans["pro_plan"],
				assignApp:   testdata.PortalApps["test_app_2"],
				expectedAcc: testdata.Accounts["account_2"],
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				// Assign plan and portal app to the account
				test.expectedAcc.Plan = test.assignPlan
				if test.expectedAcc.PortalApps == nil {
					test.expectedAcc.PortalApps = make(map[types.PortalAppID]*types.PortalApp)
				}
				test.expectedAcc.PortalApps[test.assignApp.ID] = test.assignApp

				account, err := ts.client1.GetAccountByID(context.Background(), test.accountID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedAcc, account)

					account, err = ts.client2.GetAccountByID(context.Background(), test.accountID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedAcc, account)
				}
			})
		}
	})

	ts.Run("Test_GetAccountsByUser", func() {
		tests := []struct {
			name         string
			userID       types.UserID
			expectedAccs map[types.AccountID]*types.Account
			plans        map[types.AccountID]*types.Plan
			portalApps   map[types.AccountID]map[types.PortalAppID]*types.PortalApp
			err          error
		}{
			{
				name:   "Should get accounts for user_2",
				userID: "user_2",
				expectedAccs: map[types.AccountID]*types.Account{
					"account_1": testdata.Accounts["account_1"],
					"account_2": testdata.Accounts["account_2"],
				},
				plans: map[types.AccountID]*types.Plan{
					"account_1": testdata.PayPlans["basic_plan"],
					"account_2": testdata.PayPlans["pro_plan"],
				},
				portalApps: map[types.AccountID]map[types.PortalAppID]*types.PortalApp{
					"account_1": {"test_app_1": testdata.PortalApps["test_app_1"]},
					"account_2": {"test_app_2": testdata.PortalApps["test_app_2"]},
				},
			},
			{
				name:   "Should get accounts for user_1",
				userID: "user_1",
				expectedAccs: map[types.AccountID]*types.Account{
					"account_1": testdata.Accounts["account_1"],
				},
				plans: map[types.AccountID]*types.Plan{
					"account_1": testdata.PayPlans["basic_plan"],
				},
				portalApps: map[types.AccountID]map[types.PortalAppID]*types.PortalApp{
					"account_1": {"test_app_1": testdata.PortalApps["test_app_1"]},
				},
			},
			{
				name:   "Should get accounts for user_4",
				userID: "user_4",
				expectedAccs: map[types.AccountID]*types.Account{
					"account_2": testdata.Accounts["account_2"],
					"account_4": testdata.Accounts["account_4"],
					"account_5": testdata.Accounts["account_5"],
				},
				plans: map[types.AccountID]*types.Plan{
					"account_2": testdata.PayPlans["pro_plan"],
					"account_4": testdata.PayPlans["enterprise_plan"],
					"account_5": testdata.PayPlans["basic_plan"],
				},
				portalApps: map[types.AccountID]map[types.PortalAppID]*types.PortalApp{
					"account_2": {"test_app_2": testdata.PortalApps["test_app_2"]},
				},
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				accounts, err := ts.client1.GetAccountsByUser(context.Background(), test.userID)
				ts.Equal(test.err, err)

				if err == nil {
					accountMap := convertAccountsToMap(accounts)
					// Assign plans and portal apps to expected data
					for id, account := range test.expectedAccs {
						account.Plan = test.plans[id]
						account.PortalApps = test.portalApps[id]
					}
					ts.Equal(test.expectedAccs, accountMap)

					accounts, err = ts.client2.GetAccountsByUser(context.Background(), test.userID)
					ts.Equal(test.err, err)
					accountMap = convertAccountsToMap(accounts)
					for id, account := range test.expectedAccs {
						account.Plan = test.plans[id]
						account.PortalApps = test.portalApps[id]
					}
					ts.Equal(test.expectedAccs, accountMap)
				}
			})
		}
	})

	/* ------ V2 User Read Tests ------ */

	ts.Run("Test_GetPortalUser", func() {
		tests := []struct {
			name               string
			userID             string
			expectedPortalUser *types.User
			err                error
		}{
			{
				name:               "Should get Portal User for provider user_1",
				userID:             "auth0|james_holden",
				expectedPortalUser: testdata.Users["user_1"],
			},
			{
				name:               "Should get Portal User for provider user_2",
				userID:             "auth0|paul_atreides",
				expectedPortalUser: testdata.Users["user_2"],
			},
			{
				name:               "Should get Portal User for provider user_3",
				userID:             "auth0|ellen_ripley",
				expectedPortalUser: testdata.Users["user_3"],
			},
			{
				name:               "Should get Portal User for provider user_4",
				userID:             "auth0|ulfric_stormcloak",
				expectedPortalUser: testdata.Users["user_4"],
			},
			{
				name:               "Should get Portal User for provider user_5",
				userID:             "auth0|chrisjen_avasarala",
				expectedPortalUser: testdata.Users["user_5"],
			},
			{
				name:               "Should get Portal User for provider user_6",
				userID:             "auth0|amos_burton",
				expectedPortalUser: testdata.Users["user_6"],
			},
			{
				name:               "Should get Portal User for provider user_7",
				userID:             "auth0|frodo_baggins",
				expectedPortalUser: testdata.Users["user_7"],
			},
			{
				name:               "Should get Portal User for portal user_1",
				userID:             "user_1",
				expectedPortalUser: testdata.Users["user_1"],
			},
			{
				name:               "Should get Portal User for portal user_2",
				userID:             "user_2",
				expectedPortalUser: testdata.Users["user_2"],
			},
			{
				name:               "Should get Portal User for portal user_3",
				userID:             "user_3",
				expectedPortalUser: testdata.Users["user_3"],
			},
			{
				name:               "Should get Portal User for portal user_4",
				userID:             "user_4",
				expectedPortalUser: testdata.Users["user_4"],
			},
			{
				name:               "Should get Portal User for portal user_5",
				userID:             "user_5",
				expectedPortalUser: testdata.Users["user_5"],
			},
			{
				name:               "Should get Portal User for portal user_6",
				userID:             "user_6",
				expectedPortalUser: testdata.Users["user_6"],
			},
			{
				name:               "Should get Portal User for portal user_7",
				userID:             "user_7",
				expectedPortalUser: testdata.Users["user_7"],
			},
			{
				name:   "Should error when user does not exist",
				userID: "facebook|ron_swanson",
				err:    fmt.Errorf("Response not OK. 404 Not Found: error in getPortalUser: user not found for ID: facebook|ron_swanson"),
			},
			{
				name:   "Should error when no user ID",
				userID: "",
				err:    fmt.Errorf("no user ID"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalUser, err := ts.client1.GetPortalUser(context.Background(), test.userID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedPortalUser, portalUser)

					portalUser, err = ts.client2.GetPortalUser(context.Background(), test.userID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedPortalUser, portalUser)
				}
			})
		}
	})

	ts.Run("Test_GetPortalUserID", func() {
		tests := []struct {
			name                 string
			userID               string
			expectedPortalUserID types.UserID
			err                  error
		}{
			{
				name:                 "Should get Portal User ID for provider user_1",
				userID:               "auth0|james_holden",
				expectedPortalUserID: "user_1",
			},
			{
				name:                 "Should get Portal User ID for provider user_2",
				userID:               "auth0|paul_atreides",
				expectedPortalUserID: "user_2",
			},
			{
				name:                 "Should get Portal User ID for provider user_3",
				userID:               "auth0|ellen_ripley",
				expectedPortalUserID: "user_3",
			},
			{
				name:                 "Should get Portal User ID for provider user_4",
				userID:               "auth0|ulfric_stormcloak",
				expectedPortalUserID: "user_4",
			},
			{
				name:                 "Should get Portal User ID for provider user_5",
				userID:               "auth0|chrisjen_avasarala",
				expectedPortalUserID: "user_5",
			},
			{
				name:                 "Should get Portal User ID for provider user_6",
				userID:               "auth0|amos_burton",
				expectedPortalUserID: "user_6",
			},
			{
				name:                 "Should get Portal User ID for provider user_7",
				userID:               "auth0|frodo_baggins",
				expectedPortalUserID: "user_7",
			},
			{
				name:                 "Should get Portal User ID for portal user_1",
				userID:               "user_1",
				expectedPortalUserID: "user_1",
			},
			{
				name:   "Should error when user does not exist",
				userID: "facebook|ron_swanson",
				err:    fmt.Errorf("Response not OK. 404 Not Found: error in getPortalUser: user not found for ID: facebook|ron_swanson"),
			},
			{
				name:   "Should error when no user ID",
				userID: "",
				err:    fmt.Errorf("no user ID"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalUserID, err := ts.client1.GetPortalUserID(context.Background(), test.userID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedPortalUserID, portalUserID)

					portalUserID, err = ts.client2.GetPortalUserID(context.Background(), test.userID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedPortalUserID, portalUserID)
				}
			})
		}
	})

	/* ------ V2 Plans Read Tests ------ */

	ts.Run("Test_GetAllPlans", func() {
		tests := []struct {
			name     string
			expected map[types.PayPlanType]*types.Plan
			err      error
		}{
			{
				name:     "Should get all plans",
				expected: testdata.PayPlans,
				err:      nil,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				plans, err := ts.client1.GetAllPlans(context.Background())
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(derefPlansMap(test.expected), convertPlansToMap(plans))

					plans, err = ts.client2.GetAllPlans(context.Background())
					ts.Equal(test.err, err)
					ts.Equal(derefPlansMap(test.expected), convertPlansToMap(plans))
				}
			})
		}
	})

	/* ------ V2 Blocked Contracts Read Tests ------ */

	ts.Run("Test_GetBlockedContracts", func() {
		tests := []struct {
			name               string
			expectedBlockedCon types.GlobalBlockedContracts
			err                error
		}{
			{
				name:               "Should get all blocked contracts",
				expectedBlockedCon: testdata.GlobalBlockedContracts,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				blockedContracts, err := ts.client1.GetBlockedContracts(context.Background())
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedBlockedCon, blockedContracts)

					blockedContracts, err = ts.client2.GetBlockedContracts(context.Background())
					ts.Equal(test.err, err)
					ts.Equal(test.expectedBlockedCon, blockedContracts)
				}
			})
		}
	})

	/* ------ Portal App ID Read Test ------ */

	ts.Run("Test_GetPortalUserID", func() {
		tests := []struct {
			name           string
			providerUserID string
			expectedUserID types.UserID
			err            error
		}{
			{
				name:           "Should fetch Portal User ID for provider user ID auth0|james_holden",
				providerUserID: "auth0|james_holden",
				expectedUserID: "user_1",
			},
			{
				name:           "Should fetch Portal User ID for provider user ID auth0|ellen_ripley",
				providerUserID: "auth0|ellen_ripley",
				expectedUserID: "user_3",
			},
			{
				name:           "Should fetch Portal User ID for provider user ID github|paul_atreides",
				providerUserID: "github|paul_atreides",
				expectedUserID: "user_2",
			},
			{
				name:           "Should fail if passed an invalid user ID",
				providerUserID: "auth0|george_carlin",
				err:            fmt.Errorf("Response not OK. 404 Not Found: error in getPortalUser: user not found for ID: auth0|george_carlin"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				userID, err := ts.client1.GetPortalUserID(context.Background(), test.providerUserID)
				ts.Equal(test.err, err)
				if test.err == nil {
					ts.Equal(test.expectedUserID, userID)
				}
			})
		}
	})

}

func chainsToMap(chains []*types.Chain) map[types.RelayChainID]*types.Chain {
	chainMap := make(map[types.RelayChainID]*types.Chain)
	for _, chain := range chains {
		chainMap[chain.ID] = chain
	}
	return chainMap
}

func gigastakeAppsToMap(apps []*types.GigastakeApp) map[types.GigastakeAppID]*types.GigastakeApp {
	appMap := make(map[types.GigastakeAppID]*types.GigastakeApp)
	for _, app := range apps {
		appMap[app.ID] = app
	}
	return appMap
}

func portalAppsToMap(apps []*types.PortalApp) map[types.PortalAppID]*types.PortalApp {
	result := make(map[types.PortalAppID]*types.PortalApp)
	for _, app := range apps {
		result[app.ID] = app
	}
	return result
}

func portalAppLitesToMap(portalAppLites []*types.PortalAppLite) map[types.PortalAppID]*types.PortalAppLite {
	portalAppLitesMap := make(map[types.PortalAppID]*types.PortalAppLite)
	for i := range portalAppLites {
		strKeys := make([]string, len(portalAppLites[i].PublicKeys))
		for j, key := range portalAppLites[i].PublicKeys {
			strKeys[j] = string(key)
		}

		sort.Strings(strKeys)

		for j, key := range strKeys {
			portalAppLites[i].PublicKeys[j] = types.PortalAppPublicKey(key)
		}

		portalAppLitesMap[portalAppLites[i].ID] = portalAppLites[i]
	}
	return portalAppLitesMap
}

func convertAccountsToMap(accounts []*types.Account) map[types.AccountID]*types.Account {
	accountMap := make(map[types.AccountID]*types.Account)
	for _, account := range accounts {
		accountMap[account.ID] = account
	}
	return accountMap
}

func convertPlansToMap(plans []types.Plan) map[types.PayPlanType]types.Plan {
	planMap := make(map[types.PayPlanType]types.Plan)
	for _, plan := range plans {
		planMap[plan.Type] = plan
	}
	return planMap
}

func derefPlansMap(plans map[types.PayPlanType]*types.Plan) map[types.PayPlanType]types.Plan {
	planMap := make(map[types.PayPlanType]types.Plan)
	for _, plan := range plans {
		planMap[plan.Type] = *plan
	}
	return planMap
}

func filterActiveChains(chains map[types.RelayChainID]*types.Chain) map[types.RelayChainID]*types.Chain {
	activeChains := make(map[types.RelayChainID]*types.Chain)

	for id, chain := range chains {
		if chain.Active {
			activeChains[id] = chain
		}
	}

	return activeChains
}

/* ---------- Test Suite Util Interfaces ---------- */

type phdE2EReadTestSuite struct {
	suite.Suite
	client1, client2 IDBClient
}

type phdE2EWriteTestSuite struct {
	suite.Suite
	client1, client2 IDBClient
}

const (
	phdPortOne = "8080"
	phdPortTwo = "8081"
)

func initDBClient(ts DBClientInitializer) error {
	baseConfig := Config{
		APIKey:  "test_api_key_6789",
		Retries: 1,
		Timeout: 10 * time.Second,
	}

	config1 := baseConfig
	config1.BaseURL = "http://localhost:8080"
	client1, err := NewDBClient(config1)
	if err != nil {
		return err
	}
	ts.SetClient1(client1)

	config2 := baseConfig
	config2.BaseURL = "http://localhost:8081"
	client2, err := NewDBClient(config2)
	if err != nil {
		return err
	}
	ts.SetClient2(client2)

	return nil
}

type DBClientInitializer interface {
	SetClient1(client IDBClient)
	SetClient2(client IDBClient)
	NoError(err error)
}

func (ts *phdE2EReadTestSuite) SetClient1(client IDBClient) {
	ts.client1 = client
}
func (ts *phdE2EReadTestSuite) SetClient2(client IDBClient) {
	ts.client2 = client
}
func (ts *phdE2EWriteTestSuite) SetClient1(client IDBClient) {
	ts.client1 = client
}
func (ts *phdE2EWriteTestSuite) SetClient2(client IDBClient) {
	ts.client2 = client
}
func (ts *phdE2EReadTestSuite) NoError(err error) {
	ts.Suite.NoError(err)
}
func (ts *phdE2EWriteTestSuite) NoError(err error) {
	ts.Suite.NoError(err)
}
