package dbclient

import (
	"context"
	"fmt"
	"net/http"
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
				err:     fmt.Errorf("Response not OK. 404 Not Found: blockchain not found"),
			},
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				got, err := ts.client1.GetChainByID(testCtx, test.chainID)
				ts.Equal(test.err, err)

				if err == nil {
					// Assign GigastakeApp to the chain's GigastakeApps
					test.expectedChain.GigastakeApps = make(map[types.GigastakeAppID]*types.GigastakeApp)
					test.expectedChain.GigastakeApps[test.gigastakeApp.ID] = test.gigastakeApp

					// Compare the expectedChain and actual
					ts.Equal(test.expectedChain, got)
				}
			})
		}
	})
}

// Runs all the write endpoint tests after the read tests
// This ensures the write tests do not modify the seed data expected by the read tests
func (ts *phdE2EWriteTestSuite) Test_WriteTests() {

	/* ------ V2 Chain Create/Update Tests ------ */

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
				createdChainResp, err := ts.client1.CreateChainAndGigastakeApps(testCtx, test.newChainInput)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)
					createdChain := createdChainResp.Chain
					createdGigastakeApps := createdChainResp.GigastakeApps
					timestamp := createdChain.CreatedAt

					test.newChainInput.Chain.CreatedAt = timestamp
					test.newChainInput.Chain.UpdatedAt = timestamp

					ts.Equal(test.newChainInput.Chain, createdChain)
					for _, expectedApp := range test.newChainInput.GigastakeApps {
						expectedApp.CreatedAt = timestamp
						expectedApp.UpdatedAt = timestamp
						for _, createdApp := range createdGigastakeApps {
							expectedApp.ID = createdApp.ID
							expectedApp.ChainIDs = createdApp.ChainIDs
							ts.Equal(test.newChainInput.GigastakeApps, createdGigastakeApps)
						}
					}

					createdChainByID, err := ts.client1.GetChainByID(testCtx, createdChain.ID)
					ts.NoError(err)
					createdChainByID.CreatedAt = timestamp
					createdChainByID.UpdatedAt = timestamp
					ts.Len(createdChainByID.GigastakeApps, 1)
					createdChainByID.GigastakeApps = nil
					ts.Equal(createdChain, createdChainByID)

					createdChainByID, err = ts.client2.GetChainByID(testCtx, createdChain.ID)
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
		}

		for _, test := range tests {
			ts.Run(test.name, func() {
				createdGigastakeApp, err := ts.client1.CreateGigastakeApp(testCtx, test.gigastakeAppInput)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)
					timestamp := createdGigastakeApp.CreatedAt

					// Ensure timestamps are the same before comparing
					test.expected.CreatedAt = timestamp
					test.expected.UpdatedAt = timestamp
					test.expected.ID = createdGigastakeApp.ID
					test.expected.PrivateKey = ""

					ts.Equal(test.expected, createdGigastakeApp)
				}
			})
		}
	})

	ts.Run("Test_UpdateChain", func() {
		tests := []struct {
			name        string
			chainUpdate types.Chain
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
				chainUpdateResponse, err := ts.client1.UpdateChain(testCtx, test.chainUpdate)
				ts.Equal(test.err, err)

				if test.err == nil {
					<-time.After(50 * time.Millisecond)

					ts.NotEmpty(chainUpdateResponse)

					timestamp := chainUpdateResponse.CreatedAt

					test.chainUpdate.CreatedAt = timestamp
					test.chainUpdate.UpdatedAt = timestamp

					ts.Equal(test.chainUpdate, *chainUpdateResponse)

					updatedChainByID, err := ts.client1.GetChainByID(testCtx, chainUpdateResponse.ID)
					ts.NoError(err)
					if test.noSubtables {
						test.chainUpdate.Altruists = updatedChainByID.Altruists
						test.chainUpdate.Checks = updatedChainByID.Checks
						test.chainUpdate.AliasDomains = updatedChainByID.AliasDomains
					}
					updatedChainByID.CreatedAt = timestamp
					updatedChainByID.UpdatedAt = timestamp
					ts.NotEmpty(updatedChainByID.GigastakeApps, 1)
					updatedChainByID.GigastakeApps = nil
					ts.Equal(test.chainUpdate, *updatedChainByID)

					updatedChainByID, err = ts.client2.GetChainByID(testCtx, chainUpdateResponse.ID)
					if test.noSubtables {
						test.chainUpdate.Altruists = updatedChainByID.Altruists
						test.chainUpdate.Checks = updatedChainByID.Checks
						test.chainUpdate.AliasDomains = updatedChainByID.AliasDomains
					}
					ts.NoError(err)
					updatedChainByID.CreatedAt = timestamp
					updatedChainByID.UpdatedAt = timestamp
					ts.NotEmpty(updatedChainByID.GigastakeApps, 1)
					updatedChainByID.GigastakeApps = nil
					ts.Equal(test.chainUpdate, *updatedChainByID)
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
				chainActive, err := ts.client1.ActivateChain(testCtx, test.chainID, test.active)
				ts.Equal(test.err, err)

				if err == nil {
					<-time.After(50 * time.Millisecond)
					ts.Equal(test.active, chainActive)

					fetchedChain, err := ts.client1.GetChainByID(testCtx, test.chainID)
					ts.NoError(err)
					ts.Equal(test.active, fetchedChain.Active)

					fetchedChain, err = ts.client2.GetChainByID(testCtx, test.chainID)
					ts.NoError(err)
					ts.Equal(test.active, fetchedChain.Active)
				}
			})
		}
	})
}

/* ---------- Test Suite Util Interfaces ---------- */
var testCtx = context.Background()

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
