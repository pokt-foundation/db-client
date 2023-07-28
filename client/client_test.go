package dbclient

import (
	"context"
	"fmt"
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
				err:     fmt.Errorf("Response not OK. 404 Not Found: chain not found"),
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

	ts.Run("Test_GetAllChains", func() {
		tests := []struct {
			name           string
			expectedChains map[types.RelayChainID]*types.Chain
			gigastakeApps  map[types.GigastakeAppID]*types.GigastakeApp
			options        ChainsOptions
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
				options:        ChainsOptions{IncludeInactive: true},
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

				if test.options.IncludeInactive {
					chains, err = ts.client1.GetAllChains(context.Background(), test.options)
				} else {
					chains, err = ts.client1.GetAllChains(context.Background())
				}
				ts.Equal(test.err, err)

				if test.err == nil {
					ts.Equal(test.expectedChains, chainsToMap(chains))

					if test.options.IncludeInactive {
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
				err:         fmt.Errorf("Response not OK. 404 Not Found: portal app not found"),
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
			options      PortalAppsOptions
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
				name:    "Should get portal apps where user_1 is OWNER",
				userID:  "user_1",
				options: PortalAppsOptions{RoleNameFilter: types.RoleOwner},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_1": testdata.PortalApps["test_app_1"],
				},
			},
			{
				name:    "Should get portal apps where user_6 is ADMIN",
				userID:  "user_6",
				options: PortalAppsOptions{RoleNameFilter: types.RoleAdmin},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_3": testdata.PortalApps["test_app_3"],
				},
			},
			{
				name:    "Should get portal apps where user_7 is MEMBER",
				userID:  "user_7",
				options: PortalAppsOptions{RoleNameFilter: types.RoleMember},
				expectedApps: map[types.PortalAppID]*types.PortalApp{
					"test_app_3": testdata.PortalApps["test_app_3"],
				},
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				var portalApps []*types.PortalApp
				var err error

				if test.options.RoleNameFilter != "" {
					portalApps, err = ts.client1.GetPortalAppsByUser(context.Background(), test.userID, test.options)
				} else {
					portalApps, err = ts.client1.GetPortalAppsByUser(context.Background(), test.userID)
				}
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedApps, portalAppsToMap(portalApps))

					if test.options.RoleNameFilter != "" {
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

	ts.Run("Test_GetUserPermissionByUserID", func() {
		tests := []struct {
			name                string
			providerUserID      types.ProviderUserID
			expectedPermissions *types.UserPermissions
			err                 error
		}{
			{
				name:                "Should get user permissions for user_1",
				providerUserID:      "auth0|james_holden",
				expectedPermissions: testdata.UserPermissions["user_1"],
			},
			{
				name:                "Should get user permissions for user_2",
				providerUserID:      "auth0|paul_atreides",
				expectedPermissions: testdata.UserPermissions["user_2"],
			},
			{
				name:                "Should get user permissions for user_3",
				providerUserID:      "auth0|ellen_ripley",
				expectedPermissions: testdata.UserPermissions["user_3"],
			},
			{
				name:                "Should get user permissions for user_4",
				providerUserID:      "auth0|ulfric_stormcloak",
				expectedPermissions: testdata.UserPermissions["user_4"],
			},
			{
				name:                "Should get user permissions for user_5",
				providerUserID:      "auth0|chrisjen_avasarala",
				expectedPermissions: testdata.UserPermissions["user_5"],
			},
			{
				name:                "Should get user permissions for user_6",
				providerUserID:      "auth0|amos_burton",
				expectedPermissions: testdata.UserPermissions["user_6"],
			},
			{
				name:                "Should get user permissions for user_7",
				providerUserID:      "auth0|frodo_baggins",
				expectedPermissions: testdata.UserPermissions["user_7"],
			},
			{
				name:           "Should error when no user ID",
				providerUserID: "",
				err:            fmt.Errorf("no user ID"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				userPermissions, err := ts.client1.GetUserPermissionByUserID(context.Background(), test.providerUserID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedPermissions, userPermissions)

					userPermissions, err = ts.client2.GetUserPermissionByUserID(context.Background(), test.providerUserID)
					ts.Equal(test.err, err)
					ts.Equal(test.expectedPermissions, userPermissions)
				}
			})
		}
	})

	ts.Run("Test_GetPortalUserIDFromProviderUserID", func() {
		tests := []struct {
			name                 string
			providerUserID       types.ProviderUserID
			expectedPortalUserID types.UserID
			err                  error
		}{
			{
				name:                 "Should get Portal User ID for provider user_1",
				providerUserID:       "auth0|james_holden",
				expectedPortalUserID: "user_1",
			},
			{
				name:                 "Should get Portal User ID for provider user_2",
				providerUserID:       "auth0|paul_atreides",
				expectedPortalUserID: "user_2",
			},
			{
				name:                 "Should get Portal User ID for provider user_3",
				providerUserID:       "auth0|ellen_ripley",
				expectedPortalUserID: "user_3",
			},
			{
				name:                 "Should get Portal User ID for provider user_4",
				providerUserID:       "auth0|ulfric_stormcloak",
				expectedPortalUserID: "user_4",
			},
			{
				name:                 "Should get Portal User ID for provider user_5",
				providerUserID:       "auth0|chrisjen_avasarala",
				expectedPortalUserID: "user_5",
			},
			{
				name:                 "Should get Portal User ID for provider user_6",
				providerUserID:       "auth0|amos_burton",
				expectedPortalUserID: "user_6",
			},
			{
				name:                 "Should get Portal User ID for provider user_7",
				providerUserID:       "auth0|frodo_baggins",
				expectedPortalUserID: "user_7",
			},
			{
				name:           "Should error when user does not exist",
				providerUserID: "facebook|ron_swanson",
				err:            fmt.Errorf("Response not OK. 404 Not Found"),
			},
			{
				name:           "Should error when no user ID",
				providerUserID: "",
				err:            fmt.Errorf("no user ID"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				portalUserID, err := ts.client1.GetPortalUserIDFromProviderUserID(context.Background(), test.providerUserID)
				ts.Equal(test.err, err)

				if err == nil {
					ts.Equal(test.expectedPortalUserID, portalUserID)

					portalUserID, err = ts.client2.GetPortalUserIDFromProviderUserID(context.Background(), test.providerUserID)
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
}

// Runs all the write endpoint tests after the read tests
// This ensures the write tests do not modify the seed data expected by the read tests
func (ts *phdE2EWriteTestSuite) Test_WriteTests() {

	/* ------ V2 Chain Write Tests ------ */

	ts.Run("Test_CreateChainAndGigastakeApps", func() {
		tests := []struct {
			name          string
			newChainInput types.NewChainInput
			err           error
		}{
			{
				name:          "Should create a new blockchain and its Gigastake apps in the DB",
				newChainInput: testdata.TestCreateNewChainInput,
			},
			{
				name:          "Should fail if Chain is missing",
				newChainInput: types.NewChainInput{},
				err:           fmt.Errorf("Response not OK. 400 Bad Request: error chain cannot be nil"),
			},
			{
				name: "Should fail if GigastakeApp is missing",
				newChainInput: types.NewChainInput{
					Chain: &types.Chain{ID: "1234"},
				},
				err: fmt.Errorf("Response not OK. 400 Bad Request: error gigastakeApps slice cannot be empty"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdChainResp, err := ts.client1.CreateChainAndGigastakeApps(context.Background(), test.newChainInput)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)
					ts.NotEmpty(createdChainResp.GigastakeApps[0].ID)

					createdChain := createdChainResp.Chain
					createdGigastakeApps := createdChainResp.GigastakeApps
					timestamp := createdChain.CreatedAt

					test.newChainInput.Chain.CreatedAt = timestamp
					test.newChainInput.Chain.UpdatedAt = timestamp

					ts.Equal(test.newChainInput.Chain, createdChain)
					for _, expectedApp := range test.newChainInput.GigastakeApps {
						expectedApp.ID = createdChainResp.GigastakeApps[0].ID
						expectedApp.ChainIDs = map[types.RelayChainID]struct{}{
							test.newChainInput.Chain.ID: {},
						}
						expectedApp.CreatedAt = timestamp
						expectedApp.UpdatedAt = timestamp
						ts.Equal(test.newChainInput.GigastakeApps, createdGigastakeApps)
					}

					createdChainByID, err := ts.client1.GetChainByID(context.Background(), createdChain.ID)
					ts.NoError(err)
					createdChainByID.CreatedAt = timestamp
					createdChainByID.UpdatedAt = timestamp
					ts.Len(createdChainByID.GigastakeApps, 1)
					createdChainByID.GigastakeApps = nil
					ts.Equal(createdChain, createdChainByID)

					createdChainByID, err = ts.client2.GetChainByID(context.Background(), createdChain.ID)
					ts.NoError(err)
					createdChainByID.CreatedAt = timestamp
					createdChainByID.UpdatedAt = timestamp
					ts.Len(createdChainByID.GigastakeApps, 1)
					createdChainByID.GigastakeApps = nil
					ts.Equal(createdChain, createdChainByID)
				}
			})
		}
	})

	ts.Run("Test_CreateGigastakeApp", func() {
		tests := []struct {
			name              string
			gigastakeAppInput types.GigastakeApp
			err               error
			expected          *types.GigastakeApp
		}{
			{
				name:              "Should create a new Gigastake app in the DB",
				gigastakeAppInput: testdata.TestCreateGigastakeApp,
				expected:          &testdata.TestCreateGigastakeApp,
			},
			{
				name: "Should return an error if no name provided",
				gigastakeAppInput: types.GigastakeApp{
					Name: "",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: gigastake app name cannot be empty"),
			},
			{
				name: "Should return an error for non-existent chain ID",
				gigastakeAppInput: types.GigastakeApp{
					Name:     "whatever",
					ChainIDs: map[types.RelayChainID]struct{}{"0666": {}},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error chain does not exist for chain ID '0666'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdGigastakeApp, err := ts.client1.CreateGigastakeApp(context.Background(), test.gigastakeAppInput)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(createdGigastakeApp.ID)

					timestamp := createdGigastakeApp.CreatedAt

					// Ensure ID and timestamps are the same before comparing
					test.expected.ID = createdGigastakeApp.ID
					test.expected.CreatedAt = timestamp
					test.expected.UpdatedAt = timestamp
					test.expected.PrivateKey = ""

					ts.Equal(test.expected, createdGigastakeApp)

					// Check the GigastakeApp is included in each chain
					for chainID := range test.gigastakeAppInput.ChainIDs {
						chain, err := ts.client1.GetChainByID(context.Background(), chainID)
						ts.NoError(err)
						chain.GigastakeApps[test.expected.ID].CreatedAt = timestamp
						chain.GigastakeApps[test.expected.ID].UpdatedAt = timestamp
						ts.Equal(test.expected, chain.GigastakeApps[test.expected.ID])

						chain, err = ts.client2.GetChainByID(context.Background(), chainID)
						ts.NoError(err)
						chain.GigastakeApps[test.expected.ID].CreatedAt = timestamp
						chain.GigastakeApps[test.expected.ID].UpdatedAt = timestamp
						ts.Equal(test.expected, chain.GigastakeApps[test.expected.ID])
					}
				}
			})
		}
	})

	ts.Run("Test_UpdateChain", func() {
		tests := []struct {
			name        string
			chainUpdate types.UpdateChain
			noSubtables bool
			err         error
		}{
			{
				name:        "Should update the blockchain in the DB",
				chainUpdate: testdata.UpdateChainOne,
			},
			{
				name:        "Should update the blockchain again in the DB",
				chainUpdate: testdata.UpdateChainTwo,
			},
			{
				name:        "Should update the blockchain a third time in the DB without removing any subtables",
				chainUpdate: testdata.UpdateChainThree,
				noSubtables: true, // When no subtables are passed in the update do not modify the subtables of the expected chain
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				chainUpdateResponse, err := ts.client1.UpdateChain(context.Background(), test.chainUpdate)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(chainUpdateResponse)

					updatedChainByID1, err := ts.client1.GetChainByID(context.Background(), chainUpdateResponse.ID)
					ts.NoError(err)
					updatedChainByID2, err := ts.client2.GetChainByID(context.Background(), chainUpdateResponse.ID)
					ts.NoError(err)

					updatedChains := []*types.Chain{updatedChainByID1, updatedChainByID2}

					for _, updatedChain := range updatedChains {
						ts.Equal(chainUpdateResponse, updatedChain)

						expectedChain := &types.Chain{
							GigastakeApps: chainUpdateResponse.GigastakeApps,
						}

						// Only compare the fields present in the update struct
						if test.chainUpdate.Blockchain != nil {
							expectedChain.Blockchain = *test.chainUpdate.Blockchain
						}

						if test.chainUpdate.Description != nil {
							expectedChain.Description = *test.chainUpdate.Description
						}

						if test.chainUpdate.EnforceResult != nil {
							expectedChain.EnforceResult = *test.chainUpdate.EnforceResult
						}

						if test.chainUpdate.Path != nil {
							expectedChain.Path = *test.chainUpdate.Path
						}

						if test.chainUpdate.Ticker != nil {
							expectedChain.Ticker = *test.chainUpdate.Ticker
						}

						if test.chainUpdate.AllowedMethods != nil {
							expectedChain.AllowedMethods = test.chainUpdate.AllowedMethods
						}

						if test.chainUpdate.LogLimitBlocks != nil {
							expectedChain.LogLimitBlocks = *test.chainUpdate.LogLimitBlocks
						}

						if test.chainUpdate.RequestTimeout != nil {
							expectedChain.RequestTimeout = *test.chainUpdate.RequestTimeout
						}

						if test.chainUpdate.Active != nil {
							expectedChain.Active = *test.chainUpdate.Active
						}

						if !test.noSubtables {
							expectedChain.Altruists = *test.chainUpdate.Altruists
							expectedChain.Checks = *test.chainUpdate.Checks
							expectedChain.AliasDomains = *test.chainUpdate.AliasDomains
						} else {
							expectedChain.Altruists = updatedChain.Altruists
							expectedChain.Checks = updatedChain.Checks
							expectedChain.AliasDomains = updatedChain.AliasDomains
						}

						expectedChain.ID = updatedChain.ID
						expectedChain.CreatedAt = updatedChain.CreatedAt
						expectedChain.UpdatedAt = updatedChain.UpdatedAt

						ts.Equal(expectedChain, updatedChain)
					}
				}
			})
		}
	})

	ts.Run("Test_UpdateGigastakeApp", func() {
		tests := []struct {
			name               string
			gigastakeAppUpdate types.UpdateGigastakeApp
			err                error
			expected           *types.UpdateGigastakeApp
		}{
			{
				name: "Should update GigastakeApp ChainIDs in the database",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake",
					ChainIDs: []types.RelayChainID{"0001", "0040"},
				},
				err: nil,
				expected: &types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake",
					ChainIDs: []types.RelayChainID{"0001", "0040"},
				},
			},
			{
				name: "Should update both GigastakeApp name and ChainIDs in the database",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake_updated",
					ChainIDs: []types.RelayChainID{"0001", "0040", "0053"},
				},
				err: nil,
				expected: &types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake_updated",
					ChainIDs: []types.RelayChainID{"0001", "0040", "0053"},
				},
			},
			{
				name: "Should update both GigastakeApp name and ChainIDs in the database back to original values",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake",
					ChainIDs: []types.RelayChainID{"0001"},
				},
				err: nil,
				expected: &types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "pokt_gigastake",
					ChainIDs: []types.RelayChainID{"0001"},
				},
			},
			{
				name: "Should return an error if the GigastakeApp name is empty",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:   "test_gigastake_app_1",
					Name: "",
				},
				err:      fmt.Errorf("Response not OK. 500 Internal Server Error: gigastake app name cannot be empty"),
				expected: nil,
			},
			{
				name: "Should return an error if the GigastakeApp ChainIDs is empty",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "whatever",
					ChainIDs: []types.RelayChainID{},
				},
				err:      fmt.Errorf("Response not OK. 500 Internal Server Error: chainIDs cannot be empty for gigastake app update"),
				expected: nil,
			},
			{
				name: "Should return an error if the chain doesn't exist",
				gigastakeAppUpdate: types.UpdateGigastakeApp{
					ID:       "test_gigastake_app_1",
					Name:     "whatever",
					ChainIDs: []types.RelayChainID{"0666"},
				},
				err:      fmt.Errorf("Response not OK. 500 Internal Server Error: error chain does not exist for chain ID '0666'"),
				expected: nil,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				updatedGigastakeApp, err := ts.client1.UpdateGigastakeApp(context.Background(), test.gigastakeAppUpdate)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(updatedGigastakeApp.ID)
					ts.Equal(test.expected, updatedGigastakeApp)

					// Check the GigastakeApp is included in each chain
					for _, chainID := range test.gigastakeAppUpdate.ChainIDs {
						chain, err := ts.client1.GetChainByID(context.Background(), chainID)
						ts.NoError(err)
						ts.Equal(test.gigastakeAppUpdate.Name, chain.GigastakeApps[updatedGigastakeApp.ID].Name)
						ts.Contains(chain.GigastakeApps, updatedGigastakeApp.ID)

						chain, err = ts.client2.GetChainByID(context.Background(), chainID)
						ts.NoError(err)
						ts.Equal(test.gigastakeAppUpdate.Name, chain.GigastakeApps[updatedGigastakeApp.ID].Name)
						ts.Contains(chain.GigastakeApps, updatedGigastakeApp.ID)
					}
				}
			})
		}
	})

	ts.Run("Test_ActivateChain", func() {
		tests := []struct {
			name    string
			chainID types.RelayChainID
			active  bool
			err     error
		}{
			{
				name:    "Should activate a blockchain in the DB",
				chainID: "0064",
				active:  true,
			},
			{
				name:    "Should deactivate a blockchain in the DB",
				chainID: "0064",
				active:  false,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				chainActive, err := ts.client1.ActivateChain(context.Background(), test.chainID, test.active)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)
					ts.Equal(test.active, chainActive)

					fetchedChain, err := ts.client1.GetChainByID(context.Background(), test.chainID)
					ts.NoError(err)
					ts.Equal(test.active, fetchedChain.Active)

					fetchedChain, err = ts.client2.GetChainByID(context.Background(), test.chainID)
					ts.NoError(err)
					ts.Equal(test.active, fetchedChain.Active)
				}
			})
		}
	})

	/* ------ V2 Portal App Write Tests ------ */

	ts.Run("Test_CreatePortalApp", func() {
		tests := []struct {
			name           string
			portalAppInput *types.PortalApp
			aatInput       types.AAT
			err            error
			expected       *types.PortalApp
		}{
			{
				name:           "Should create a new Portal app in the DB",
				portalAppInput: testdata.TestCreatePortalApp,
				aatInput:       testdata.TestCreatePortalAppAAT,
				expected: &types.PortalApp{
					Name:      "create_pokt_app_1",
					AccountID: "account_4",
					Settings: types.Settings{
						Environment:       "production",
						SecretKey:         "test_3e3fb7949c9e3b193cfba5348f93bb2f",
						SecretKeyRequired: true,
					},
					AATs: map[types.ProtocolAppID]types.AAT{
						"test_protocol_app_5": testdata.TestCreatePortalAppAAT,
					},
					Notifications: map[types.NotificationType]types.AppNotification{
						types.NotificationTypeEmail: {
							Type:        types.NotificationTypeEmail,
							Active:      true,
							Destination: "ulfric.stormcloak123@test.com",
							Events: map[types.NotificationEvent]bool{
								types.NotificationEventFull:          true,
								types.NotificationEventSignedUp:      true,
								types.NotificationEventThreeQuarters: true,
							},
						},
					},
					FirstDateSurpassed: testdata.MockTimestamp,
					CreatedAt:          testdata.MockTimestamp,
					UpdatedAt:          testdata.MockTimestamp,
					LegacyFields: types.LegacyFields{
						PlanType:       types.FreetierV0,
						DailyLimit:     250_000,
						RequestTimeout: 15_000,
					},
				},
			},
			{
				name: "Should return an error if no name provided",
				portalAppInput: &types.PortalApp{
					Name:      "",
					AccountID: "account_4",
					LegacyFields: types.LegacyFields{
						PlanType: "basic_plan",
					},
					Settings: types.Settings{
						Environment: "production",
					},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: portal app name cannot be empty"),
			},
			{
				name: "Should return an error if invalid environment provided",
				portalAppInput: &types.PortalApp{
					Name:      "Test App",
					AccountID: "account_4",
					LegacyFields: types.LegacyFields{
						PlanType: "basic_plan",
					},
					Settings: types.Settings{
						Environment: "cascadia",
					},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: invalid portal app environment provided: cascadia"),
			},
			{
				name: "Should return an error for non-existent account ID",
				portalAppInput: &types.PortalApp{
					Name:      "Test App",
					AccountID: "non_existing_account_id",
					LegacyFields: types.LegacyFields{
						PlanType: "basic_plan",
					},
					Settings: types.Settings{
						Environment: "production",
					},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error account does not exist for account ID 'non_existing_account_id'"),
			},
			{
				name: "Should return an error for non-existent plan",
				portalAppInput: &types.PortalApp{
					Name:      "Test App",
					AccountID: "account_4",
					LegacyFields: types.LegacyFields{
						PlanType: "non_existing_plan",
					},
					Settings: types.Settings{
						Environment: "production",
					},
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error pay plan 'non_existing_plan' does not exist"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				test.portalAppInput.AATs = map[types.ProtocolAppID]types.AAT{
					test.aatInput.ID: test.aatInput,
				}
				createdPortalApp, err := ts.client1.CreatePortalApp(context.Background(), *test.portalAppInput)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(createdPortalApp.ID)

					timestamp := testdata.MockTimestamp

					// Ensure ID and timestamps are the same before comparing
					test.expected.ID = createdPortalApp.ID
					createdPortalApp.CreatedAt = timestamp
					createdPortalApp.UpdatedAt = timestamp
					aat := test.expected.AATs["test_protocol_app_5"]
					var testAATID types.ProtocolAppID
					for aatID := range createdPortalApp.AATs {
						testAATID = aatID
						break
					}
					aat.ID = testAATID
					aat.PrivateKey = ""
					test.expected.AATs[testAATID] = aat
					delete(test.expected.AATs, "test_protocol_app_5")

					ts.Equal(test.expected, createdPortalApp)

					portalApp, err := ts.client1.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.NoError(err)
					portalApp.CreatedAt = timestamp
					portalApp.UpdatedAt = timestamp
					ts.Equal(test.expected, portalApp)

					portalApp, err = ts.client2.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.NoError(err)
					portalApp.CreatedAt = timestamp
					portalApp.UpdatedAt = timestamp
					ts.Equal(test.expected, portalApp)
				}
			})
		}
	})

	ts.Run("Test_UpdatePortalApp", func() {
		tests := []struct {
			name                     string
			updatePortalApp          types.UpdatePortalApp
			testUpdateTime           time.Time
			testUpdatedName          string
			testUpdatedSettings      types.Settings
			testUpdatedNotifications map[types.NotificationType]types.AppNotification
			testUpdatedWhitelists    types.Whitelists
			testUpdatedLegacyFields  types.LegacyFields
			err                      error
		}{
			{
				name: "Should update a new PortalApp in the database with all fields",
				updatePortalApp: types.UpdatePortalApp{
					Name:          testdata.UpdatePortalAppName,
					Settings:      testdata.UpdatePortalAppSettings,
					Notifications: testdata.UpdatePortalAppNotifications,
					Whitelists:    testdata.UpdatePortalAppWhitelists,
					PlanType:      testdata.UpdatePortalAppPlan.PlanType,
					DailyLimit:    testdata.UpdatePortalAppPlan.DailyLimit,
				},
				testUpdateTime:  testdata.MockTimestamp,
				testUpdatedName: "portal-app-updated",
				testUpdatedSettings: types.Settings{
					Environment:       types.EnvironmentProduction,
					SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
					SecretKeyRequired: true,
					MonthlyRelayLimit: 2_500_000,
					FavoritedChainIDs: map[types.RelayChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
				},
				testUpdatedNotifications: map[types.NotificationType]types.AppNotification{
					types.NotificationTypeEmail: {
						Type:        types.NotificationTypeEmail,
						Active:      true,
						Destination: "user@example.com",
						Trigger:     "daily",
						Events: map[types.NotificationEvent]bool{
							types.NotificationEventSignedUp:      true,
							types.NotificationEventHalf:          true,
							types.NotificationEventQuarter:       true,
							types.NotificationEventThreeQuarters: true,
							types.NotificationEventFull:          true,
						},
					},
					types.NotificationTypeWebhook: {
						Type:        types.NotificationTypeWebhook,
						Active:      true,
						Destination: "https://example.com/webhook",
						Trigger:     "hourly",
						Events: map[types.NotificationEvent]bool{
							types.NotificationEventHalf: true,
							types.NotificationEventFull: true,
						},
					},
				},
				testUpdatedWhitelists: types.Whitelists{
					Origins:     map[types.Origin]struct{}{"https://portalgun.io": {}, "https://subdomain.example.com": {}, "https://www.example.com": {}},
					UserAgents:  map[types.UserAgent]struct{}{"Brave": {}, "Google Chrome": {}, "Mozilla Firefox": {}, "Netscape Navigator": {}, "Safari": {}},
					Blockchains: map[types.RelayChainID]struct{}{"0001": {}, "0002": {}, "003E": {}, "0056": {}},
					Contracts: map[types.RelayChainID]map[types.Contract]struct{}{
						"0001": {"0xtest_2f78db6436527729929aaf6c616361de0f7": {}, "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}},
						"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
						"003E": {"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}, "0xtest_f958d2ee523a2206206994597c13d831ec7": {}},
						"0056": {"0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}, "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}},
					},
					Methods: map[types.RelayChainID]map[types.Method]struct{}{
						"0001": {"GET": {}, "POST": {}, "PUT": {}},
						"0002": {"DELETE": {}, "GET": {}, "POST": {}, "PUT": {}},
						"003E": {"GET": {}},
						"0056": {"GET": {}, "POST": {}},
					},
				},
				testUpdatedLegacyFields: types.LegacyFields{
					PlanType:   types.PayAsYouGoV0,
					DailyLimit: 0,
				},
				err: nil,
			},
			{
				name: "Should update a new PortalApp in the database with only a new Name",
				updatePortalApp: types.UpdatePortalApp{
					Name: testdata.UpdatePortalAppName,
				},
				testUpdateTime:  testdata.MockTimestamp,
				testUpdatedName: "portal-app-updated",
				err:             nil,
			},
			{
				name: "Should update a new PortalApp in the database with only new Settings",
				updatePortalApp: types.UpdatePortalApp{
					Settings: testdata.UpdatePortalAppSettings,
				},
				testUpdateTime: testdata.MockTimestamp,
				testUpdatedSettings: types.Settings{
					Environment:       types.EnvironmentProduction,
					SecretKey:         "test_9d07c8a96ad53e7c288b0e86f37c5680",
					SecretKeyRequired: true,
					MonthlyRelayLimit: 2_500_000,
					FavoritedChainIDs: map[types.RelayChainID]struct{}{"0003": {}, "0009": {}, "00H3": {}},
				},
				err: nil,
			},
			{
				name: "Should update a new PortalApp in the database with only new Notifications",
				updatePortalApp: types.UpdatePortalApp{
					Notifications: testdata.UpdatePortalAppNotifications,
				},
				testUpdateTime: testdata.MockTimestamp,
				testUpdatedNotifications: map[types.NotificationType]types.AppNotification{
					types.NotificationTypeEmail: {
						Type:        types.NotificationTypeEmail,
						Active:      true,
						Destination: "user@example.com",
						Trigger:     "daily",
						Events: map[types.NotificationEvent]bool{
							types.NotificationEventSignedUp:      true,
							types.NotificationEventHalf:          true,
							types.NotificationEventQuarter:       true,
							types.NotificationEventThreeQuarters: true,
							types.NotificationEventFull:          true,
						},
					},
					types.NotificationTypeWebhook: {
						Type:        types.NotificationTypeWebhook,
						Active:      true,
						Destination: "https://example.com/webhook",
						Trigger:     "hourly",
						Events: map[types.NotificationEvent]bool{
							types.NotificationEventHalf: true,
							types.NotificationEventFull: true,
						},
					},
				},
				err: nil,
			},
			{
				name: "Should update a new PortalApp in the database with only new Whitelists",
				updatePortalApp: types.UpdatePortalApp{
					Whitelists: testdata.UpdatePortalAppWhitelists,
				},
				testUpdateTime: testdata.MockTimestamp,
				testUpdatedWhitelists: types.Whitelists{
					Origins:     map[types.Origin]struct{}{"https://portalgun.io": {}, "https://subdomain.example.com": {}, "https://www.example.com": {}},
					UserAgents:  map[types.UserAgent]struct{}{"Brave": {}, "Google Chrome": {}, "Mozilla Firefox": {}, "Netscape Navigator": {}, "Safari": {}},
					Blockchains: map[types.RelayChainID]struct{}{"0001": {}, "0002": {}, "003E": {}, "0056": {}},
					Contracts: map[types.RelayChainID]map[types.Contract]struct{}{
						"0001": {"0xtest_2f78db6436527729929aaf6c616361de0f7": {}, "0xtest_5fbfe3e9af3971dd833d26ba9b5c936f0be": {}},
						"0002": {"0xtest_1111117dc0aa78b770fa6a738034120c302": {}, "0xtest_a39b223fe8d0a0e5c4f27ead9083c756cc2": {}},
						"003E": {"0xtest_0a85d5af5bf1d1762f925bdaddc4201f984": {}, "0xtest_f958d2ee523a2206206994597c13d831ec7": {}},
						"0056": {"0xtest_00000f279d81a1d3cc75430faa017fa5a2e": {}, "0xtest_5068778dd592e39a122f4f5a5cf09c90fe2": {}},
					},
					Methods: map[types.RelayChainID]map[types.Method]struct{}{
						"0001": {"GET": {}, "POST": {}, "PUT": {}},
						"0002": {"DELETE": {}, "GET": {}, "POST": {}, "PUT": {}},
						"003E": {"GET": {}},
						"0056": {"GET": {}, "POST": {}},
					},
				},
				err: nil,
			},
			{
				name: "Should update a new PortalApp in the database with a new plan",
				updatePortalApp: types.UpdatePortalApp{
					PlanType:   testdata.UpdatePortalAppPlan.PlanType,
					DailyLimit: testdata.UpdatePortalAppPlan.DailyLimit,
				},
				testUpdateTime: testdata.MockTimestamp,
				testUpdatedLegacyFields: types.LegacyFields{
					PlanType:   types.PayAsYouGoV0,
					DailyLimit: 0,
				},
				err: nil,
			},
			{
				name: "Should update a new PortalApp in the database with an Enterprise plan",
				updatePortalApp: types.UpdatePortalApp{
					PlanType:   testdata.UpdatePortalAppEnterprisePlan.PlanType,
					DailyLimit: testdata.UpdatePortalAppEnterprisePlan.DailyLimit,
				},
				testUpdateTime: testdata.MockTimestamp,
				testUpdatedLegacyFields: types.LegacyFields{
					PlanType:    types.Enterprise,
					CustomLimit: 5_600_000,
				},
				err: nil,
			},
		}

		for i, test := range tests {

			ts.Run(test.name, func() {
				// Create new portal app for test case
				createApp := *testdata.TestUpdatePortalApp
				createApp.Name = fmt.Sprintf("test-update-portal-app-%d", i+1)
				createdPortalApp, err := ts.client1.CreatePortalApp(context.Background(), createApp)
				ts.NoError(err)

				checkUpdatedPortalApp := func(updatedPortalApp *types.PortalApp) {
					if test.testUpdatedName != "" {
						ts.Equal(test.testUpdatedName, updatedPortalApp.Name)
					} else {
						ts.Equal(createdPortalApp.Name, updatedPortalApp.Name)
					}

					if test.testUpdatedSettings.Environment != "" {
						ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
					} else {
						ts.Equal(createdPortalApp.Settings, updatedPortalApp.Settings)
					}

					if len(test.testUpdatedNotifications) != 0 {
						ts.Equal(test.testUpdatedNotifications, updatedPortalApp.Notifications)
					} else {
						ts.Equal(createdPortalApp.Notifications, updatedPortalApp.Notifications)
					}

					if len(test.testUpdatedWhitelists.Origins) != 0 {
						ts.Equal(test.testUpdatedWhitelists, updatedPortalApp.Whitelists)
					} else {
						ts.Equal(createdPortalApp.Whitelists, updatedPortalApp.Whitelists)
					}

					if test.updatePortalApp.PlanType != "" {
						ts.Equal(test.testUpdatedLegacyFields.PlanType, updatedPortalApp.LegacyFields.PlanType)
					}
					if test.updatePortalApp.DailyLimit != 0 {
						ts.Equal(test.testUpdatedLegacyFields.DailyLimit, updatedPortalApp.LegacyFields.DailyLimit)
					}
					if test.updatePortalApp.CustomLimit != 0 {
						ts.Equal(test.testUpdatedLegacyFields.CustomLimit, updatedPortalApp.LegacyFields.CustomLimit)
					}
				}

				// Update created portal app for test case
				updateApp := test.updatePortalApp
				updateApp.AppID = createdPortalApp.ID
				_, err = ts.client1.UpdatePortalApp(context.Background(), updateApp)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)

					updatedPortalApp, err := ts.client1.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.NoError(err)
					checkUpdatedPortalApp(updatedPortalApp)

					updatedPortalApp, err = ts.client2.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.NoError(err)
					checkUpdatedPortalApp(updatedPortalApp)
				}
			})
		}
	})

	ts.Run("Test_DeletePortalApp", func() {
		tests := []struct {
			name     string
			expected map[string]string
			err      error
		}{
			{
				name:     "Should delete the Portal App in the DB",
				expected: map[string]string{"status": "deleted"},
				err:      fmt.Errorf("Response not OK. 404 Not Found: portal app not found"),
			},
		}

		for i, test := range tests {
			// Create new portal app for test case
			createApp := *testdata.TestUpdatePortalApp
			createApp.Name = fmt.Sprintf("test-delete-portal-app-%d", i+1)
			createdPortalApp, err := ts.client1.CreatePortalApp(context.Background(), createApp)
			ts.NoError(err)

			<-time.After(50 * time.Millisecond)

			// Ensure the Portal App exists in both clients
			portalApp, err := ts.client1.GetPortalAppByID(context.Background(), createdPortalApp.ID)
			ts.NoError(err)
			ts.NotEmpty(portalApp)
			portalApp, err = ts.client2.GetPortalAppByID(context.Background(), createdPortalApp.ID)
			ts.NoError(err)
			ts.NotEmpty(portalApp)

			ts.Run(test.name, func() {
				response, err := ts.client1.DeletePortalApp(context.Background(), createdPortalApp.ID)
				ts.NoError(err)

				if err == nil {
					ts.Equal(test.expected, response)

					<-time.After(50 * time.Millisecond)

					// Ensure the Portal App is deleted for both clients
					portalApp, err := ts.client1.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.Equal(test.err, err)
					ts.Nil(portalApp)

					portalApp, err = ts.client2.GetPortalAppByID(context.Background(), createdPortalApp.ID)
					ts.Equal(test.err, err)
					ts.Nil(portalApp)
				}
			})
		}
	})

	ts.Run("Test_UpdatePortalAppsFirstDateSurpassed", func() {
		tests := []struct {
			name                     string
			firstDateSurpassedUpdate types.UpdateFirstDateSurpassed
			err                      error
			expected                 map[string]string
		}{
			{
				name: "Should update FirstDateSurpassed field for a given list of Portal apps in the DB",
				firstDateSurpassedUpdate: types.UpdateFirstDateSurpassed{
					PortalAppIDs:       []types.PortalAppID{"test_app_1", "test_app_2"},
					FirstDateSurpassed: testdata.MockTimestamp,
				},
				expected: map[string]string{"status": "updated"},
			},
			{
				name: "Should return an error if no app IDs provided",
				firstDateSurpassedUpdate: types.UpdateFirstDateSurpassed{
					PortalAppIDs:       []types.PortalAppID{},
					FirstDateSurpassed: testdata.MockTimestamp,
				},
				err: fmt.Errorf("Response not OK. 400 Bad Request"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				result, err := ts.client1.UpdatePortalAppsFirstDateSurpassed(context.Background(), test.firstDateSurpassedUpdate)
				ts.Equal(test.err, err)

				if test.err == nil {
					ts.Equal(test.expected, result)

					for _, appID := range test.firstDateSurpassedUpdate.PortalAppIDs {
						portalApp, err := ts.client1.GetPortalAppByID(context.Background(), appID)
						ts.NoError(err)
						ts.Equal(test.firstDateSurpassedUpdate.FirstDateSurpassed, portalApp.FirstDateSurpassed)

						portalApp, err = ts.client2.GetPortalAppByID(context.Background(), appID)
						ts.NoError(err)
						ts.Equal(test.firstDateSurpassedUpdate.FirstDateSurpassed, portalApp.FirstDateSurpassed)
					}
				}
			})
		}
	})

	/* ------ V2 Account Write Tests ------ */

	ts.Run("Test_CreateAccount", func() {
		tests := []struct {
			name         string
			ownerID      types.UserID
			accountInput *types.Account
			err          error
			expected     *types.Account
			expectedPlan *types.Plan
		}{
			{
				name:         "Should create a new Account in the DB",
				ownerID:      "user_1",
				accountInput: testdata.TestCreateAccount,
				expected: &types.Account{
					PlanType:  types.PayPlanType("developer_plan"),
					CreatedAt: testdata.MockTimestamp,
					UpdatedAt: testdata.MockTimestamp,
					Users: map[types.UserID]types.AccountUserAccess{
						"user_1": {
							UserID:   "user_1",
							Email:    "james.holden123@test.com",
							Owner:    true,
							Accepted: true,
						},
					},
				},
				expectedPlan: &types.Plan{
					Type:              types.PayPlanType("developer_plan"),
					ChainIDs:          map[types.RelayChainID]struct{}{"0001": {}, "0053": {}, "0021": {}, "0034": {}},
					MonthlyRelayLimit: 500_000,
					ThroughputLimit:   500,
					AppLimit:          1,
					LegacyDailyLimit:  100,
					CreatedAt:         testdata.MockTimestamp,
				},
			},
			{
				name:         "Should fail if input does not have a User ID set",
				ownerID:      "",
				accountInput: &types.Account{},
				err:          fmt.Errorf("no user ID"),
			},
			{
				name:         "Should fail if input Account does not have a PayPlanType set",
				ownerID:      "user_1",
				accountInput: &types.Account{PlanType: ""},
				err:          fmt.Errorf("no plan type set"),
			},
			{
				name:         "Should fail if input Account has an invalid plan type",
				ownerID:      "user_1",
				accountInput: &types.Account{PlanType: types.PayPlanType("turbo_ultra_mega_plan")},
				err:          fmt.Errorf("Response not OK. 500 Internal Server Error: error pay plan 'turbo_ultra_mega_plan' does not exist"),
			},
			{
				name:         "Should fail if input User does not exist in the db",
				ownerID:      "user_451",
				accountInput: testdata.Accounts[types.AccountID("account_5")],
				err:          fmt.Errorf("Response not OK. 500 Internal Server Error: error user does not exist for portal ID 'user_451'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdAccount, err := ts.client1.CreateAccount(context.Background(), test.ownerID, *test.accountInput, time.Now())
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)
					ts.NotEmpty(createdAccount.ID)

					// Ensure ID and timestamps are the same before comparing
					test.expected.ID = createdAccount.ID
					createdAccount.CreatedAt = testdata.MockTimestamp
					createdAccount.UpdatedAt = testdata.MockTimestamp
					ts.Equal(test.expected, createdAccount)

					test.expected.Plan = test.expectedPlan
					expectedUsers := test.expected.Users[test.ownerID]
					expectedUsers.PortalAppRoles = map[types.PortalAppID]types.RoleName{}
					test.expected.Users[test.ownerID] = expectedUsers

					account, err := ts.client1.GetAccountByID(context.Background(), createdAccount.ID)
					ts.NoError(err)
					account.CreatedAt = testdata.MockTimestamp
					account.UpdatedAt = testdata.MockTimestamp
					ts.Equal(test.expected, account)

					account, err = ts.client2.GetAccountByID(context.Background(), createdAccount.ID)
					ts.NoError(err)
					account.CreatedAt = testdata.MockTimestamp
					account.UpdatedAt = testdata.MockTimestamp
					ts.Equal(test.expected, account)
				}
			})
		}
	})

	ts.Run("Test_UpdateAccount", func() {
		tests := []struct {
			name                string
			accountBeforeUpdate *types.Account
			update              types.UpdateAccount
			err                 error
			expected            *types.Account
		}{
			{
				name:                "Should update the account's PlanType field",
				accountBeforeUpdate: testdata.Accounts["account_1"],
				update: types.UpdateAccount{
					AccountID: "account_1",
					PlanType:  types.Enterprise,
				},
				expected: &types.Account{
					ID:       "account_1",
					PlanType: types.Enterprise,
				},
			},
			{
				name: "Should fail if an invalid account ID is provided",
				update: types.UpdateAccount{
					AccountID: "account_8823",
					PlanType:  types.Enterprise,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error account does not exist for account ID 'account_8823'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				if test.err == nil {
					accountBeforeUpdate, err := ts.client1.GetAccountByID(context.Background(), test.accountBeforeUpdate.ID)
					ts.NoError(err)
					ts.Equal(test.accountBeforeUpdate.PlanType, accountBeforeUpdate.PlanType)
				}

				updatedAccount, err := ts.client1.UpdateAccount(context.Background(), test.update)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.Equal(test.expected.PlanType, updatedAccount.PlanType)

					accountAfterUpdate, err := ts.client1.GetAccountByID(context.Background(), updatedAccount.ID)
					ts.NoError(err)
					ts.Equal(test.expected.PlanType, accountAfterUpdate.PlanType)
				}
			})
		}
	})

	ts.Run("Test_CreateAccountIntegration", func() {
		tests := []struct {
			name                       string
			accountIntegrationInput    *types.AccountIntegrations
			expectedAccountIntegration types.AccountIntegrations
			err                        error
		}{
			{
				name: "Should create a new account integration",
				accountIntegrationInput: &types.AccountIntegrations{
					AccountID:          "account_5",
					CovalentAPIKeyFree: "created_covalent_api_key_1",
				},
				expectedAccountIntegration: types.AccountIntegrations{
					AccountID:          "account_5",
					CovalentAPIKeyFree: "created_covalent_api_key_1",
				},
				err: nil,
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdAccountIntegration, err := ts.client1.CreateAccountIntegration(context.Background(), test.accountIntegrationInput.AccountID, *test.accountIntegrationInput)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(createdAccountIntegration.AccountID)

					// Ensure AccountID is the same before comparing
					test.expectedAccountIntegration.AccountID = createdAccountIntegration.AccountID

					ts.Equal(&test.expectedAccountIntegration, createdAccountIntegration)

					account, err := ts.client1.GetAccountByID(context.Background(), createdAccountIntegration.AccountID)
					ts.NoError(err)
					ts.Equal(test.expectedAccountIntegration, account.Integrations)

					account, err = ts.client2.GetAccountByID(context.Background(), createdAccountIntegration.AccountID)
					ts.NoError(err)
					ts.Equal(test.expectedAccountIntegration, account.Integrations)
				}
			})
		}
	})

	ts.Run("Test_UpdateAccountIntegration", func() {
		tests := []struct {
			name                    string
			accountIntegrationInput types.AccountIntegrations
			err                     error
			expected                *types.AccountIntegrations
		}{
			{
				name: "Should update existing Account integration in the DB",
				accountIntegrationInput: types.AccountIntegrations{
					AccountID:          "account_5",
					CovalentAPIKeyFree: "updated_covalent_api_key_1",
				},
				err: nil,
				expected: &types.AccountIntegrations{
					AccountID:          "account_5",
					CovalentAPIKeyFree: "updated_covalent_api_key_1",
				},
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				updatedAccountIntegration, err := ts.client1.UpdateAccountIntegration(context.Background(), test.accountIntegrationInput.AccountID, test.accountIntegrationInput)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(updatedAccountIntegration.AccountID)
					test.expected.AccountID = updatedAccountIntegration.AccountID
					ts.Equal(test.expected, updatedAccountIntegration)

					account, err := ts.client1.GetAccountByID(context.Background(), updatedAccountIntegration.AccountID)
					ts.NoError(err)
					ts.Equal(test.expected.CovalentAPIKeyFree, account.Integrations.CovalentAPIKeyFree)

					account, err = ts.client2.GetAccountByID(context.Background(), updatedAccountIntegration.AccountID)
					ts.NoError(err)
					ts.Equal(test.expected.CovalentAPIKeyFree, account.Integrations.CovalentAPIKeyFree)
				}
			})
		}
	})

	ts.Run("Test_DeleteAccount", func() {
		tests := []struct {
			name     string
			ownerID  types.UserID
			expected map[string]string
			err      error
		}{
			{
				name:     "Should delete the Account in the DB",
				ownerID:  "user_7",
				expected: map[string]string{"status": "deleted"},
				err:      fmt.Errorf("Response not OK. 404 Not Found: account not found"),
			},
		}

		for i, test := range tests {
			// Create new account for test case
			createAcc := types.Account{
				ID:                     types.AccountID(fmt.Sprintf("test-delete-account-%d", i+1)),
				PlanType:               types.FreetierV0,
				PartnerChainIDs:        map[types.RelayChainID]struct{}{"chain_1": {}},
				PartnerThroughputLimit: 1000,
				PartnerAppLimit:        5,
			}
			createdAccount, err := ts.client1.CreateAccount(context.Background(), test.ownerID, createAcc, testdata.MockTimestamp)
			ts.NoError(err)

			<-time.After(50 * time.Millisecond)

			// Ensure the Account exists in both clients
			account, err := ts.client1.GetAccountByID(context.Background(), createdAccount.ID)
			ts.NoError(err)
			ts.NotEmpty(account)
			account, err = ts.client2.GetAccountByID(context.Background(), createdAccount.ID)
			ts.NoError(err)
			ts.NotEmpty(account)

			ts.Run(test.name, func() {
				response, err := ts.client1.DeleteAccount(context.Background(), createdAccount.ID)
				ts.NoError(err)

				if err == nil {
					ts.Equal(test.expected, response)

					<-time.After(50 * time.Millisecond)

					// Ensure the Account is deleted for both clients
					account, err := ts.client1.GetAccountByID(context.Background(), createdAccount.ID)
					ts.Equal(test.err, err)
					ts.Nil(account)

					account, err = ts.client2.GetAccountByID(context.Background(), createdAccount.ID)
					ts.Equal(test.err, err)
					ts.Nil(account)
				}
			})
		}
	})

	/* ------ V2 Account User Write Tests ------ */

	ts.Run("Test_WriteAccountUser", func() {
		tests := []struct {
			name                   string
			createAccountUserInput types.CreateAccountUserAccess
			expected               *types.AccountUserAccess
			err                    error
		}{
			{
				name: "Should create a new AccountUserAccess in the DB",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_4",
					PortalAppID: "test_app_1",
					Email:       "bernard.marx@test.com",
					RoleName:    types.RoleMember,
				},
				expected: &types.AccountUserAccess{
					AccountID:      "account_4",
					UserID:         "user_11",
					Email:          "bernard.marx@test.com",
					Accepted:       false,
					PortalAppRoles: map[types.PortalAppID]types.RoleName{"test_app_1": types.RoleMember},
				},
			},
			{
				name: "Should create a new AccountUserAccess in the DB for a user that hasn't signed up yet",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					Email:       "winston.smith@test.com",
					RoleName:    types.RoleAdmin,
				},
				expected: &types.AccountUserAccess{
					AccountID:      "account_3",
					UserID:         "", // UserID created when user created
					Email:          "winston.smith@test.com",
					Accepted:       false,
					PortalAppRoles: map[types.PortalAppID]types.RoleName{"test_app_3": types.RoleAdmin},
				},
			},
			{
				name: "Should fail if an invalid email provided",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_4",
					PortalAppID: "test_app_3",
					Email:       "winston.smith",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error email input is not a valid email address 'winston.smith'"),
			},
			{
				name: "Should fail if account does not exist",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_674",
					PortalAppID: "test_app_222",
					Email:       "winston.smith@test.com",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error account does not exist for account ID 'account_674'"),
			},
			{
				name: "Should fail if an empty email string is provided",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_1",
					PortalAppID: "test_app_3",
					Email:       "",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("no email"),
			},
			{
				name: "Should fail if an empty PortalAppID string is provided",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_1",
					PortalAppID: "",
					Email:       "valid.email@test.com",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("no portal app ID"),
			},
			{
				name: "Should fail if an empty AccountID string is provided",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "",
					PortalAppID: "test_app_1",
					Email:       "valid.email@test.com",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("no account ID"),
			},
			{
				name: "Should fail if the AccountID provided does not exist",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "non_existent_account",
					PortalAppID: "test_app_1",
					Email:       "valid.email@test.com",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error account does not exist for account ID 'non_existent_account'"),
			},
			{
				name: "Should fail if the PortalAppID provided does not exist",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_1",
					PortalAppID: "non_existent_app",
					Email:       "valid.email@test.com",
					RoleName:    types.RoleAdmin,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error portal app does not exist for ID 'non_existent_app'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				userIDResp, err := ts.client1.WriteAccountUser(context.Background(), test.createAccountUserInput, testdata.MockTimestamp)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(userIDResp["userID"])
					accountID := test.createAccountUserInput.AccountID
					test.expected.AccountID = ""

					if test.expected.UserID == "" {
						test.expected.UserID = userIDResp["userID"]
					}

					account, err := ts.client1.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(*test.expected, account.Users[test.expected.UserID])

					account, err = ts.client2.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(*test.expected, account.Users[test.expected.UserID])

					// Clean up created account user
					_, err = ts.client1.RemoveAccountUser(context.Background(), types.UpdateRemoveAccountUser{
						AccountID:   test.createAccountUserInput.AccountID,
						PortalAppID: test.createAccountUserInput.PortalAppID,
						UserID:      test.expected.UserID,
					})
					ts.NoError(err)
				}

			})
		}
	})

	ts.Run("Test_SetAccountUserRole", func() {
		tests := []struct {
			name                    string
			updateAccountUserRole   types.UpdateAccountUserRole
			accountUsersAfterUpdate map[types.UserID]types.AccountUserAccess
			testCreatedTime         time.Time
			err                     error
		}{
			{
				name: "Should update an existing AccountUserAccess row's role to non-OWNER role",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "test_app_3",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    types.RoleAdmin,
				},
				accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
					"user_5":  testdata.AccountUserAccess[5],
					"user_6":  testdata.AccountUserAccess[6],
					"user_10": testdata.AccountUserAccess[12],
					"user_7": {
						UserID:   "user_7",
						Email:    "frodo.baggins123@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_3": types.RoleAdmin,
						},
					},
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             nil,
			},
			{
				name: "Should update an existing AccountUserAccess row's role back to original role",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "test_app_3",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    types.RoleMember,
				},
				accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
					"user_5":  testdata.AccountUserAccess[5],
					"user_6":  testdata.AccountUserAccess[6],
					"user_10": testdata.AccountUserAccess[12],
					"user_7": {
						UserID:   "user_7",
						Email:    "frodo.baggins123@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_3": types.RoleMember,
						},
					},
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             nil,
			},
			{
				name: "Should transfer the OWNER of an Account",
				updateAccountUserRole: types.UpdateAccountUserRole{
					AccountID:   "account_2",
					PortalAppID: "test_app_2",
					UserID:      "user_4",
					RoleName:    types.RoleOwner,
				},
				accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
					"user_9": testdata.AccountUserAccess[9],
					"user_2": testdata.AccountUserAccess[10],
					"user_3": {
						UserID:   "user_3",
						Email:    "ellen.ripley789@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_2": types.RoleAdmin,
						},
					},
					"user_4": {
						Owner:    true,
						UserID:   "user_4",
						Email:    "ulfric.stormcloak123@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_2": types.RoleOwner,
						},
					},
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             nil,
			},
			{
				name: "Should revert the transfer of the OWNER of an Account",
				updateAccountUserRole: types.UpdateAccountUserRole{
					AccountID:   "account_2",
					PortalAppID: "test_app_2",
					UserID:      "user_3",
					RoleName:    types.RoleOwner,
				},
				accountUsersAfterUpdate: map[types.UserID]types.AccountUserAccess{
					"user_9": testdata.AccountUserAccess[9],
					"user_2": testdata.AccountUserAccess[10],
					"user_4": {
						UserID:   "user_4",
						Email:    "ulfric.stormcloak123@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_2": types.RoleAdmin,
						},
					},
					"user_3": {
						Owner:    true,
						UserID:   "user_3",
						Email:    "ellen.ripley789@test.com",
						Accepted: true,
						PortalAppRoles: map[types.PortalAppID]types.RoleName{
							"test_app_2": types.RoleOwner,
						},
					},
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             nil,
			},
			{
				name: "Should fail if attempting to transfer ownership to user that has not accepted their invite yet",
				updateAccountUserRole: types.UpdateAccountUserRole{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					UserID:      "user_10",
					RoleName:    types.RoleOwner,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("Response not OK. 500 Internal Server Error: error cannot transfer ownership to user ID 'user_10' for account ID 'account_3' because the user has not accepted their invite"),
			},
			{
				name: "Should fail if User is not a member of an Account",
				updateAccountUserRole: types.UpdateAccountUserRole{
					AccountID:   "account_2",
					PortalAppID: "test_app_2",
					UserID:      "user_512",
					RoleName:    types.RoleMember,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("Response not OK. 500 Internal Server Error: error user ID 'user_512' does not exist for portal app ID 'account_2'"),
			},
			{
				name: "Should fail if RoleName is empty",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "test_app_3",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    "",
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("no role name"),
			},
			{
				name: "Should fail if RoleName is invalid",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "test_app_3",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    "INVALID_ROLE_NAME",
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("Response not OK. 500 Internal Server Error: error invalid role name set"),
			},
			{
				name: "Should fail if PortalAppID is empty",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    types.RoleAdmin,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("no portal app ID"),
			},
			{
				name: "Should fail if AccountID is empty",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "test_app_3",
					AccountID:   "",
					UserID:      "user_7",
					RoleName:    types.RoleAdmin,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("no account ID"),
			},
			{
				name: "Should fail if PortalAppID does not exist",
				updateAccountUserRole: types.UpdateAccountUserRole{
					PortalAppID: "non_existent_app",
					AccountID:   "account_3",
					UserID:      "user_7",
					RoleName:    types.RoleAdmin,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("Response not OK. 500 Internal Server Error: error portal app does not exist for ID 'non_existent_app'"),
			},
			{
				name: "Should fail if User is not a member of an Account",
				updateAccountUserRole: types.UpdateAccountUserRole{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					UserID:      "non_member_user",
					RoleName:    types.RoleMember,
				},
				testCreatedTime: testdata.MockTimestamp,
				err:             fmt.Errorf("Response not OK. 500 Internal Server Error: error user ID 'non_member_user' does not exist for portal app ID 'account_3'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.SetAccountUserRole(context.Background(), test.updateAccountUserRole, test.testCreatedTime)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					accountID := test.updateAccountUserRole.AccountID

					account, err := ts.client1.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(test.accountUsersAfterUpdate, account.Users)

					account, err = ts.client2.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(test.accountUsersAfterUpdate, account.Users)
				}

			})
		}
	})

	ts.Run("Test_UpdateAcceptAccountUser", func() {
		tests := []struct {
			name                   string
			accountID              types.AccountID
			userID                 string
			acceptAccountUserInput types.UpdateAcceptAccountUser
			expected               *types.AccountUserAccess
			err                    error
		}{
			{
				name:      "Should create a new UserAuthProvider for an existing user in the DB",
				accountID: "account_3",
				userID:    "user_10",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "test_app_3",
					UserID:           "user_10",
					AuthProviderType: types.AuthTypeAuth0Username,
					ProviderUserID:   "auth0|daenerys_targaryen",
				},
				expected: &types.AccountUserAccess{
					AccountID:      "account_3",
					UserID:         "user_10",
					Email:          "daenerys.targaryen123@test.com",
					Accepted:       true,
					PortalAppRoles: map[types.PortalAppID]types.RoleName{"test_app_3": types.RoleMember},
				},
				err: nil,
			},
			{
				name: "Should fail if no provider user ID provided",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "account_3",
					UserID:           "user_10",
					AuthProviderType: types.AuthType("ask_jeeves"),
				},
				err: fmt.Errorf("no provider user ID"),
			},
			{
				name: "Should fail if an invalid auth provider type provided",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "account_3",
					UserID:           "user_10",
					AuthProviderType: types.AuthType("ask_jeeves"),
					ProviderUserID:   "auth0|daenerys_targaryen",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error invalid auth provider type 'ask_jeeves'"),
			},
			{
				name: "Should fail if AuthProviderType is not provided",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "account_3",
					UserID:           "user_10",
					AuthProviderType: "",
				},
				err: fmt.Errorf("no auth provider type"),
			},
			{
				name: "Should fail if PortalAppID is not provided",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "",
					UserID:           "user_10",
					AuthProviderType: types.AuthTypeAuth0Username,
				},
				err: fmt.Errorf("no portal app ID"),
			},
			{
				name: "Should fail if user does not exist",
				acceptAccountUserInput: types.UpdateAcceptAccountUser{
					PortalAppID:      "account_3",
					UserID:           "user_123",
					AuthProviderType: types.AuthTypeAuth0Username,
					ProviderUserID:   "auth0|who_dis",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error user ID 'user_123' does not exist for portal app ID 'account_3'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.UpdateAcceptAccountUser(context.Background(), test.acceptAccountUserInput, testdata.MockTimestamp)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					accountID := test.accountID
					test.expected.AccountID = ""

					account, err := ts.client1.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(*test.expected, account.Users[test.expected.UserID])

					account, err = ts.client2.GetAccountByID(context.Background(), accountID)
					ts.NoError(err)
					ts.Equal(*test.expected, account.Users[test.expected.UserID])
				}
			})
		}
	})

	ts.Run("Test_RemoveAccountUser", func() {
		tests := []struct {
			name                    string
			createAccountUserInput  types.CreateAccountUserAccess
			updateRemoveAccountUser types.UpdateRemoveAccountUser
			accountUsersAfterDelete map[types.UserID]types.AccountUserAccess
			err                     error
		}{
			{
				name: "Should delete a single AccountUserAccess row",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					Email:       "user_5@test.com",
					RoleName:    types.RoleMember,
				},
				updateRemoveAccountUser: types.UpdateRemoveAccountUser{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
				},
				accountUsersAfterDelete: map[types.UserID]types.AccountUserAccess{
					"user_1": testdata.AccountUserAccess[1],
					"user_2": testdata.AccountUserAccess[2],
					"user_8": testdata.AccountUserAccess[8],
				},
				err: nil,
			},
			{
				name: "Should fail if provided a UserID that doesn't exist for the Account",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					Email:       "nonexistent_user@test.com",
					RoleName:    types.RoleMember,
				},
				updateRemoveAccountUser: types.UpdateRemoveAccountUser{
					AccountID:   "account_3",
					PortalAppID: "test_app_3",
					UserID:      "user_789",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error user ID 'user_789' does not exist for portal app ID 'test_app_3'"),
			},
			{
				name: "Should fail if attempting to delete the current Account OWNER",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_1",
					PortalAppID: "test_app_1",
					Email:       "user_1@test.com",
					RoleName:    types.RoleOwner,
				},
				updateRemoveAccountUser: types.UpdateRemoveAccountUser{
					AccountID:   "account_1",
					PortalAppID: "test_app_1",
					UserID:      "user_1",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error cannot delete user ID 'user_1' for account ID 'account_1' because this user is the current account owner"),
			},
			{
				name: "Should fail if provided a UserID that doesn't exist for the Account",
				createAccountUserInput: types.CreateAccountUserAccess{
					AccountID:   "account_1",
					PortalAppID: "test_app_1",
					Email:       "nonexistent_user@test.com",
					RoleName:    types.RoleMember,
				},
				updateRemoveAccountUser: types.UpdateRemoveAccountUser{
					AccountID:   "account_1",
					PortalAppID: "test_app_1",
					UserID:      "user_nonexistent",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error user ID 'user_nonexistent' does not exist for portal app ID 'test_app_1'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				var userID types.UserID

				if test.updateRemoveAccountUser.UserID == types.UserID("") {
					userIDResp, err := ts.client1.WriteAccountUser(context.Background(), test.createAccountUserInput, testdata.MockTimestamp)
					ts.NoError(err)

					<-time.After(50 * time.Millisecond)
					ts.NotEmpty(userIDResp["userID"])

					userID = userIDResp["userID"]

					test.updateRemoveAccountUser.UserID = userID
					account, err := ts.client1.GetAccountByID(context.Background(), test.createAccountUserInput.AccountID)
					ts.NoError(err)
					ts.Contains(account.Users, userID)
				}

				_, err := ts.client1.RemoveAccountUser(context.Background(), test.updateRemoveAccountUser)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					account, err := ts.client1.GetAccountByID(context.Background(), test.createAccountUserInput.AccountID)
					ts.NoError(err)
					ts.NotContains(account.Users, userID)

					account, err = ts.client2.GetAccountByID(context.Background(), test.createAccountUserInput.AccountID)
					ts.NoError(err)
					ts.NotContains(account.Users, userID)
				}
			})
		}
	})

	/* ------ V2 User Write Tests ------ */

	ts.Run("Test_CreateUser", func() {
		tests := []struct {
			name             string
			userInput        types.CreateUser
			expectedStatus   int
			expectedResponse *types.CreateUserResponse
			err              error
		}{
			{
				name: "Should create a single user in the DB",
				userInput: types.CreateUser{
					Email:          "jiminy.cricket@test.com",
					ProviderUserID: "auth0|jiminy_cricket",
				},
				expectedResponse: &types.CreateUserResponse{
					User: types.User{
						Email: "jiminy.cricket@test.com",
						AuthProviders: map[types.AuthType]types.UserAuthProvider{
							types.AuthTypeAuth0Username: {
								Type:           types.AuthTypeAuth0Username,
								ProviderUserID: "auth0|jiminy_cricket",
								Provider:       types.AuthProviderAuth0,
								Federated:      false,
							},
						},
						SignedUp: true,
					},
					AccountID: "",
				},
			},
			{
				name: "Should fail if there's no email",
				userInput: types.CreateUser{
					ProviderUserID: "auth0|test",
				},
				err: fmt.Errorf("no email"),
			},
			{
				name: "Should fail if invalid email provided",
				userInput: types.CreateUser{
					Email:          "jiminy.cricket",
					ProviderUserID: "auth0|test",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error email input is not a valid email address 'jiminy.cricket'"),
			},
			{
				name: "Should fail if there's no provider type",
				userInput: types.CreateUser{
					Email:          "email@test.com",
					ProviderUserID: "wtf|test",
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error invalid auth provider type 'wtf'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdUser, err := ts.client1.CreateUser(context.Background(), test.userInput)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					test.expectedResponse.User.ID = createdUser.User.ID
					test.expectedResponse.User.UpdatedAt = createdUser.User.UpdatedAt
					test.expectedResponse.User.CreatedAt = createdUser.User.CreatedAt
					test.expectedResponse.AccountID = createdUser.AccountID
					ts.Equal(test.expectedResponse, createdUser)
					ts.NotEmpty(createdUser.AccountID, "Should have an accountID")

					// If the user was created, it should have permissions
					providerID := createdUser.User.AuthProviders[types.AuthTypeAuth0Username].ProviderUserID
					permission, err := ts.client1.GetUserPermissionByUserID(context.Background(), providerID)
					ts.NoError(err)
					ts.NotNil(permission)

					permission, err = ts.client2.GetUserPermissionByUserID(context.Background(), providerID)
					ts.NoError(err)
					ts.NotNil(permission)
				}
			})
		}
	})

	ts.Run("Test_DeleteUser", func() {
		tests := []struct {
			name           string
			userID         types.UserID
			providerUserID types.ProviderUserID
			expectedErr    error
		}{
			{
				name:           "Should delete a User",
				userID:         "user_11",
				providerUserID: "auth0|bernard_marx",
				expectedErr:    nil,
			},
			{
				name:           "Should fail to delete a User if they are on the team of any accounts",
				userID:         "user_1",
				providerUserID: "auth0|james_holden",
				expectedErr:    fmt.Errorf("Response not OK. 500 Internal Server Error: error cannot delete user because they are still on an account team"),
			},
			{
				name:           "Should fail if the user does not exist in the database",
				userID:         "user_42",
				providerUserID: "auth0|gengelspiel",
				expectedErr:    fmt.Errorf("Response not OK. 500 Internal Server Error: error user does not exist for portal ID 'user_42'"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.DeleteUser(context.Background(), test.userID)
				ts.Equal(test.expectedErr, err)

				if test.expectedErr == nil {
					accounts, err := ts.client1.GetAccountsByUser(context.Background(), test.userID)
					ts.Error(err)
					ts.Nil(accounts)

					userPermissions, err := ts.client1.GetUserPermissionByUserID(context.Background(), test.providerUserID)
					ts.NoError(err)
					ts.Empty(userPermissions)

					accounts, err = ts.client2.GetAccountsByUser(context.Background(), test.userID)
					ts.Error(err)
					ts.Nil(accounts)

					userPermissions, err = ts.client2.GetUserPermissionByUserID(context.Background(), test.providerUserID)
					ts.NoError(err)
					ts.Empty(userPermissions)
				}
			})
		}
	})

	/* ------ V2 Blocked Contract Write Tests ------ */

	ts.Run("Test_WriteBlockedContract", func() {
		tests := []struct {
			name                     string
			blockedContract          types.BlockedContract
			expectedBlockedContracts types.GlobalBlockedContracts
			err                      error
		}{
			{
				name: "Should add a new blocked address to the global blocked contracts table",
				blockedContract: types.BlockedContract{
					BlockedAddress: "0xtest_newabcdef0123456789abcdef01234567",
					Active:         true,
				},
				expectedBlockedContracts: types.GlobalBlockedContracts{
					BlockedAddresses: map[types.BlockedAddress]struct{}{
						"0xtest_6789abcdef0123456789abcdef01234567":   {},
						"0xtest_f0123456789abcdef0123456789abcdef01":  {},
						"0xtest_cdef0123456789abcdef0123456789abcdef": {},
						"0xtest_56789abcdef0123456789abcdef01234567":  {},
						"0xtest_789abcdef0123456789abcdef0123456789":  {},
						"0xtest_newabcdef0123456789abcdef01234567":    {},
					},
				},
				err: nil,
			},
			{
				name: "Should return an error if the address is empty",
				blockedContract: types.BlockedContract{
					BlockedAddress: "",
					Active:         true,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error blockchain address must be provided"),
			},
			{
				name: "Should return an error if the address is a duplicate",
				blockedContract: types.BlockedContract{
					BlockedAddress: "0xtest_cdef0123456789abcdef0123456789abcdef",
					Active:         true,
				},
				err: fmt.Errorf("Response not OK. 500 Internal Server Error: error blockchain address 0xtest_cdef0123456789abcdef0123456789abcdef is already blocked"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.WriteBlockedContract(context.Background(), test.blockedContract)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					globalBlockedContracts, err := ts.client1.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)

					globalBlockedContracts, err = ts.client2.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)
				}
			})
		}
	})

	ts.Run("Test_UpdateBlockedContractActive", func() {
		tests := []struct {
			name                     string
			blockedAddress           types.BlockedAddress
			active                   bool
			expectedBlockedContracts types.GlobalBlockedContracts
			err                      error
		}{
			{
				name:           "Should deactivate a blocked address in the global blocked contracts table",
				blockedAddress: "0xtest_cdef0123456789abcdef0123456789abcdef",
				active:         false,
				expectedBlockedContracts: types.GlobalBlockedContracts{
					BlockedAddresses: map[types.BlockedAddress]struct{}{
						"0xtest_6789abcdef0123456789abcdef01234567":  {},
						"0xtest_f0123456789abcdef0123456789abcdef01": {},
						"0xtest_56789abcdef0123456789abcdef01234567": {},
						"0xtest_789abcdef0123456789abcdef0123456789": {},
						"0xtest_newabcdef0123456789abcdef01234567":   {},
					},
				},
				err: nil,
			},
			{
				name:           "Should reactivate a blocked address in the global blocked contracts table",
				blockedAddress: "0xtest_cdef0123456789abcdef0123456789abcdef",
				active:         true,
				expectedBlockedContracts: types.GlobalBlockedContracts{
					BlockedAddresses: map[types.BlockedAddress]struct{}{
						"0xtest_6789abcdef0123456789abcdef01234567":   {},
						"0xtest_f0123456789abcdef0123456789abcdef01":  {},
						"0xtest_cdef0123456789abcdef0123456789abcdef": {},
						"0xtest_56789abcdef0123456789abcdef01234567":  {},
						"0xtest_789abcdef0123456789abcdef0123456789":  {},
						"0xtest_newabcdef0123456789abcdef01234567":    {},
					},
				},
				err: nil,
			},
			{
				name:           "Should return an error if the address is empty",
				blockedAddress: "",
				err:            fmt.Errorf("no blocked address provided"),
			},
			{
				name:           "Should return an error if the address doesn't exist in the database",
				blockedAddress: "0xtest_34095u439fh49fh30fj239ru923kf3f09823fk",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: error blockchain address 0xtest_34095u439fh49fh30fj239ru923kf3f09823fk does not exist"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.UpdateBlockedContractActive(context.Background(), test.blockedAddress, test.active)
				ts.Equal(test.err, err)
				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					globalBlockedContracts, err := ts.client1.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)

					globalBlockedContracts, err = ts.client2.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)
				}
			})
		}
	})

	ts.Run("Test_RemoveBlockedContract", func() {
		tests := []struct {
			name                     string
			blockedAddress           types.BlockedAddress
			expectedBlockedContracts types.GlobalBlockedContracts
			err                      error
		}{
			{
				name:           "Should delete a blocked address in the global blocked contracts table",
				blockedAddress: "0xtest_789abcdef0123456789abcdef0123456789",
				expectedBlockedContracts: types.GlobalBlockedContracts{
					BlockedAddresses: map[types.BlockedAddress]struct{}{
						"0xtest_6789abcdef0123456789abcdef01234567":   {},
						"0xtest_f0123456789abcdef0123456789abcdef01":  {},
						"0xtest_cdef0123456789abcdef0123456789abcdef": {},
						"0xtest_56789abcdef0123456789abcdef01234567":  {},
						"0xtest_newabcdef0123456789abcdef01234567":    {},
					},
				},
				err: nil,
			},
			{
				name:           "Should return an error if the address is empty",
				blockedAddress: "",
				err:            fmt.Errorf("no blocked address provided"),
			},
			{
				name:           "Should return an error if the address doesn't exist in the database",
				blockedAddress: "0xtest_34095u439fh49fh30fj239ru923kf3f09823fk",
				err:            fmt.Errorf("Response not OK. 500 Internal Server Error: error blockchain address 0xtest_34095u439fh49fh30fj239ru923kf3f09823fk does not exist"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				_, err := ts.client1.RemoveBlockedContract(context.Background(), test.blockedAddress)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)

					globalBlockedContracts, err := ts.client1.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)

					globalBlockedContracts, err = ts.client2.GetBlockedContracts(context.Background())
					ts.NoError(err)
					ts.Equal(test.expectedBlockedContracts, globalBlockedContracts)
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
