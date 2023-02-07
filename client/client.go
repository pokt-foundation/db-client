package dbclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojektech/heimdall"
	"github.com/pokt-foundation/portal-db/types"
)

type (
	// DBClient struct contains all the possible methods to interact
	// with the Pocket HTTP DB amd satisfies the IDBClient interface
	DBClient struct {
		IDBClient
		httpClient *httpclient.Client
		config     Config
	}
	// Config struct to provide config options
	Config struct {
		BaseURL, APIKey string
		Version         APIVersion
		Retries         int
		Timeout         time.Duration
	}

	// IDBClient interface contains all read & write methods to interact with the Pocket HTTP DB
	IDBClient interface {
		IDBReader
		IDBWriter
	}
	// IDBReader interface contains read-only methods for interacting with the Pocket HTTP DB
	IDBReader interface {
		GetBlockchains(ctx context.Context) ([]*types.Blockchain, error)
		GetBlockchainByID(ctx context.Context, blockchainID string) (*types.Blockchain, error)
		GetApplications(ctx context.Context) ([]*types.Application, error)
		GetApplicationByID(ctx context.Context, applicationID string) (*types.Application, error)
		GetApplicationsByUserID(ctx context.Context, userID string) ([]*types.Application, error)
		GetLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error)
		GetLoadBalancerByID(ctx context.Context, loadBalancerID string) (*types.LoadBalancer, error)
		GetLoadBalancersByUserID(ctx context.Context, userID string) ([]*types.LoadBalancer, error)
		GetPayPlans(ctx context.Context) ([]*types.PayPlan, error)
		GetPayPlanByType(ctx context.Context, payPlanType types.PayPlanType) (*types.PayPlan, error)
	}
	// IDBWriter interface contains write methods for interacting with the Pocket HTTP DB
	IDBWriter interface {
		CreateBlockchain(ctx context.Context, blockchain types.Blockchain) (*types.Blockchain, error)
		CreateBlockchainRedirect(ctx context.Context, redirect types.Redirect) (*types.Redirect, error)
		CreateApplication(ctx context.Context, application types.Application) (*types.Application, error)
		CreateLoadBalancer(ctx context.Context, loadBalancer types.LoadBalancer) (*types.LoadBalancer, error)
		CreateLoadBalancerUser(ctx context.Context, loadBalancerID string, user types.UserAccess) (*types.LoadBalancer, error)
		ActivateBlockchain(ctx context.Context, blockchainID string, active bool) (bool, error)
		UpdateApplication(ctx context.Context, id string, update types.UpdateApplication) (*types.Application, error)
		UpdateAppFirstDateSurpassed(ctx context.Context, updateInput types.UpdateFirstDateSurpassed) ([]*types.Application, error)
		RemoveApplication(ctx context.Context, id string) (*types.Application, error)
		UpdateLoadBalancer(ctx context.Context, id string, lbUpdate types.UpdateLoadBalancer) (*types.LoadBalancer, error)
		UpdateLoadBalancerUserRole(ctx context.Context, id string, userUpdate types.UpdateUserAccess) (*types.LoadBalancer, error)
		RemoveLoadBalancer(ctx context.Context, id string) (*types.LoadBalancer, error)
		DeleteLoadBalancerUser(ctx context.Context, loadBalancerID, userID string) (*types.LoadBalancer, error)
	}

	basePath   string
	subPath    string
	APIVersion string
)

const (
	blockchainPath   basePath = "blockchain"
	applicationPath  basePath = "application"
	loadBalancerPath basePath = "load_balancer"
	payPlanPath      basePath = "pay_plan"
	userPath         basePath = "user"

	redirectPath           subPath = "redirect"
	activatePath           subPath = "activate"
	firstDateSurpassedPath subPath = "first_date_surpassed"
)

// New API versions should be added to both the APIVersion enum and ValidAPIVersions map
const (
	V0 APIVersion = "v0" // TODO remove when dropping v0 support
	V1 APIVersion = "v1"
)

var ValidAPIVersions = map[APIVersion]bool{
	V0: true, // TODO remove when dropping v0 support
	V1: true,
}

var (
	errBaseURLNotProvided      error = errors.New("base URL not provided")
	errAPIKeyNotProvided       error = errors.New("API key not provided")
	errVersionNotProvided      error = errors.New("version not provided")
	errInvalidVersionProvided  error = errors.New("invalid version provided")
	errNoUserID                error = errors.New("no user ID")
	errNoBlockchainID          error = errors.New("no blockchain ID")
	errNoApplicationID         error = errors.New("no application ID")
	errNoLoadBalancerID        error = errors.New("no load balancer ID")
	errNoPayPlanType           error = errors.New("no pay plan type")
	errInvalidBlockchainJSON   error = errors.New("invalid blockchain JSON")
	errInvalidAppJSON          error = errors.New("invalid application JSON")
	errInvalidLoadBalancerJSON error = errors.New("invalid load balancer JSON")
	errInvalidActivationJSON   error = errors.New("invalid active field JSON")
	errResponseNotOK           error = errors.New("Response not OK")
)

// NewDBClient returns a read-write HTTP client to use the Pocket HTTP DB - https://github.com/pokt-foundation/pocket-http-db
func NewDBClient(config Config) (IDBClient, error) {
	if err := config.validateConfig(); err != nil {
		return nil, err
	}

	return &DBClient{httpClient: newHTTPClient(config), config: config}, nil
}

// NewReadOnlyDBClient returns a read-only HTTP client to use the Pocket HTTP DB - https://github.com/pokt-foundation/pocket-http-db
func NewReadOnlyDBClient(config Config) (IDBReader, error) {
	if err := config.validateConfig(); err != nil {
		return nil, err
	}

	return &DBClient{httpClient: newHTTPClient(config), config: config}, nil
}

// newHTTPClient creates a new heimdall HTTP client with retries and exponential backoff
func newHTTPClient(config Config) *httpclient.Client {
	backoff := heimdall.NewExponentialBackoff(2*time.Millisecond, 9*time.Millisecond, 2, 2*time.Millisecond)
	retrier := heimdall.NewRetrier(backoff)

	return httpclient.NewClient(
		httpclient.WithHTTPTimeout(config.Timeout),
		httpclient.WithRetryCount(config.Retries),
		httpclient.WithRetrier(retrier),
	)
}

// validateConfig ensures that a valid configuration is provided to the DB client
func (c Config) validateConfig() error {
	if c.BaseURL == "" {
		return errBaseURLNotProvided
	}
	if c.APIKey == "" {
		return errAPIKeyNotProvided
	}
	if c.Version == "" {
		return errVersionNotProvided
	}
	if !ValidAPIVersions[c.Version] {
		return errInvalidVersionProvided
	}
	return nil
}

// versionedBasePath returns the base path for a given data type eg. `https://pocket.http-db-url.com/v1/application`
func (db *DBClient) versionedBasePath(dataTypePath basePath) string {
	if db.config.Version == V0 { // TODO remove when dropping v0 support
		return fmt.Sprintf("%s/%s", db.config.BaseURL, dataTypePath)
	}

	return fmt.Sprintf("%s/%s/%s", db.config.BaseURL, db.config.Version, dataTypePath)
}

// getAuthHeaderForRead returns the auth header for read operations to PHD
func (db *DBClient) getAuthHeaderForRead() http.Header {
	return http.Header{"Authorization": {db.config.APIKey}}
}

// getAuthHeaderForRead returns the auth header for write operations to PHD
func (db *DBClient) getAuthHeaderForWrite() http.Header {
	return http.Header{"Authorization": {db.config.APIKey}, "Content-Type": {"application/json"}}
}

/* -- Read Methods -- */

// GetBlockchains returns all blockchains in the DB - GET `<base URL>/<version>/blockchain`
func (db *DBClient) GetBlockchains(ctx context.Context) ([]*types.Blockchain, error) {
	endpoint := db.versionedBasePath(blockchainPath)

	return get[[]*types.Blockchain](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetBlockchain returns a single Blockchain by its relay chain ID - GET `<base URL>/<version>/blockchain/{id}`
func (db *DBClient) GetBlockchainByID(ctx context.Context, blockchainID string) (*types.Blockchain, error) {
	if blockchainID == "" {
		return nil, errNoBlockchainID
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(blockchainPath), blockchainID)

	return get[*types.Blockchain](endpoint, db.getAuthHeaderForRead(), db.httpClient)

}

// GetApplications returns all Applications in the DB - GET `<base URL>/<version>/application`
func (db *DBClient) GetApplications(ctx context.Context) ([]*types.Application, error) {
	endpoint := db.versionedBasePath(applicationPath)

	return get[[]*types.Application](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetApplicationByID returns a single Application by its ID - GET `<base URL>/<version>/application/{id}`
func (db *DBClient) GetApplicationByID(ctx context.Context, applicationID string) (*types.Application, error) {
	if applicationID == "" {
		return nil, errNoApplicationID
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(applicationPath), applicationID)

	return get[*types.Application](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetApplicationsByUserID returns all the Applications for a user - GET `<base URL>/<version>/user/{userID}/application`
func (db *DBClient) GetApplicationsByUserID(ctx context.Context, userID string) ([]*types.Application, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(userPath), userID, applicationPath)

	return get[[]*types.Application](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetLoadBalancers returns all Load Balancers in the DB - GET `<base URL>/<version>/load_balancer`
func (db *DBClient) GetLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error) {
	endpoint := db.versionedBasePath(loadBalancerPath)

	return get[[]*types.LoadBalancer](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetLoadBalancerByID returns a single Load Balancer by its ID - GET `<base URL>/<version>/load_balancer/{id}`
func (db *DBClient) GetLoadBalancerByID(ctx context.Context, loadBalancerID string) (*types.LoadBalancer, error) {
	if loadBalancerID == "" {
		return nil, errNoLoadBalancerID
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID)

	return get[*types.LoadBalancer](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetLoadBalancersByUserID returns all the load balancers for a user - GET `<base URL>/<version>/user/{userID}/load_balancer`
func (db *DBClient) GetLoadBalancersByUserID(ctx context.Context, userID string) ([]*types.LoadBalancer, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(userPath), userID, loadBalancerPath)

	return get[[]*types.LoadBalancer](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetPayPlans returns all Pay Plans in the DB - GET `<base URL>/<version>/pay_plan`
func (db *DBClient) GetPayPlans(ctx context.Context) ([]*types.PayPlan, error) {
	endpoint := db.versionedBasePath(payPlanPath)

	return get[[]*types.PayPlan](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetPayPlanByType returns a single Pay Plan by its type - GET `<base URL>/<version>/pay_plan/{type}`
func (db *DBClient) GetPayPlanByType(ctx context.Context, payPlanType types.PayPlanType) (*types.PayPlan, error) {
	if payPlanType == "" {
		return nil, errNoPayPlanType
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(payPlanPath), payPlanType)

	return get[*types.PayPlan](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

/* -- Create Methods -- */

// CreateBlockchain creates a single Blockchain in the DB - POST `<base URL>/<version>/blockchain`
func (db *DBClient) CreateBlockchain(ctx context.Context, blockchain types.Blockchain) (*types.Blockchain, error) {
	blockchainJSON, err := json.Marshal(blockchain)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidBlockchainJSON, err)
	}

	endpoint := db.versionedBasePath(blockchainPath)

	return post[*types.Blockchain](endpoint, db.getAuthHeaderForWrite(), blockchainJSON, db.httpClient)
}

// CreateBlockchainRedirect creates a single Blockchain Redirect in the DB - POST `<base URL>/<version>/blockchain/redirect`
func (db *DBClient) CreateBlockchainRedirect(ctx context.Context, redirect types.Redirect) (*types.Redirect, error) {
	redirectJSON, err := json.Marshal(redirect)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	var endpoint string
	if db.config.Version == V0 { // TODO remove when dropping v0 support
		endpoint = db.versionedBasePath(basePath(redirectPath))
	} else {
		endpoint = fmt.Sprintf("%s/%s", db.versionedBasePath(blockchainPath), redirectPath)
	}

	return post[*types.Redirect](endpoint, db.getAuthHeaderForWrite(), redirectJSON, db.httpClient)
}

// CreateApplication creates a single Application in the DB - POST `<base URL>/<version>/application`
func (db *DBClient) CreateApplication(ctx context.Context, app types.Application) (*types.Application, error) {
	appJSON, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	endpoint := db.versionedBasePath(applicationPath)

	return post[*types.Application](endpoint, db.getAuthHeaderForWrite(), appJSON, db.httpClient)
}

// CreateLoadBalancer creates a single Load Balancer in the DB - POST `<base URL>/<version>/load_balancer`
func (db *DBClient) CreateLoadBalancer(ctx context.Context, loadBalancer types.LoadBalancer) (*types.LoadBalancer, error) {
	loadBalancerJSON, err := json.Marshal(loadBalancer)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	return post[*types.LoadBalancer](db.versionedBasePath(loadBalancerPath), db.getAuthHeaderForWrite(), loadBalancerJSON, db.httpClient)
}

// CreateLoadBalancerUser adds a single User to a single Load Balancer in the DB - POST `<base URL>/<version>/load_balancer/{id}/user`
func (db *DBClient) CreateLoadBalancerUser(ctx context.Context, loadBalancerID string, user types.UserAccess) (*types.LoadBalancer, error) {
	loadBalancerUserJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID, userPath)

	return post[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerUserJSON, db.httpClient)
}

/* -- Update Methods -- */

// ActivateBlockchain toggles a single Blockchain's `active` field` - PUT `<base URL>/<version>/blockchain/{id}/activate`
func (db *DBClient) ActivateBlockchain(ctx context.Context, blockchainID string, active bool) (bool, error) {
	if blockchainID == "" {
		return false, errNoBlockchainID
	}

	activeJSON, err := json.Marshal(active)
	if err != nil {
		return false, fmt.Errorf("%w: %s", errInvalidActivationJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(blockchainPath), blockchainID, activatePath)

	return post[bool](endpoint, db.getAuthHeaderForWrite(), activeJSON, db.httpClient)
}

// UpdateApplication updates a single Application in the DB - PUT `<base URL>/<version>/application/{id}`
func (db *DBClient) UpdateApplication(ctx context.Context, id string, appUpdate types.UpdateApplication) (*types.Application, error) {
	if id == "" {
		return nil, errNoApplicationID
	}

	appUpdateJSON, err := json.Marshal(appUpdate)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(applicationPath), id)

	return put[*types.Application](endpoint, db.getAuthHeaderForWrite(), appUpdateJSON, db.httpClient)
}

// UpdateAppFirstDateSurpassed updates a slice of Applications' FirstDateSurpassed fields in the DB - PUT `<base URL>/<version>/first_date_surpassed`
func (db *DBClient) UpdateAppFirstDateSurpassed(ctx context.Context, updateInput types.UpdateFirstDateSurpassed) ([]*types.Application, error) {
	firstDateSurpassedJSON, err := json.Marshal(updateInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(applicationPath), firstDateSurpassedPath)

	return post[[]*types.Application](endpoint, db.getAuthHeaderForWrite(), firstDateSurpassedJSON, db.httpClient)
}

// UpdateLoadBalancer updates a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}`
func (db *DBClient) UpdateLoadBalancer(ctx context.Context, id string, lbUpdate types.UpdateLoadBalancer) (*types.LoadBalancer, error) {
	if id == "" {
		return nil, errNoLoadBalancerID
	}

	loadBalancerUpdateJSON, err := json.Marshal(lbUpdate)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(loadBalancerPath), id)

	return put[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerUpdateJSON, db.httpClient)
}

// UpdateLoadBalancerUserRole updates a single User's role for a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}/user`
func (db *DBClient) UpdateLoadBalancerUserRole(ctx context.Context, id string, userUpdate types.UpdateUserAccess) (*types.LoadBalancer, error) {
	if id == "" {
		return nil, errNoLoadBalancerID
	}

	loadBalancerUserUpdateJSON, err := json.Marshal(userUpdate)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(loadBalancerPath), id, userPath)

	return put[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerUserUpdateJSON, db.httpClient)
}

// RemoveApplication removes a single Application by updating its status field - PUT `<base URL>/<version>/application/{id}` with Remove: true
func (db *DBClient) RemoveApplication(ctx context.Context, id string) (*types.Application, error) {
	if id == "" {
		return nil, errNoApplicationID
	}

	appRemoveJSON, err := json.Marshal(types.UpdateApplication{Remove: true})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(applicationPath), id)

	return put[*types.Application](endpoint, db.getAuthHeaderForWrite(), appRemoveJSON, db.httpClient)
}

// RemoveLoadBalancer removes a single LoadBalancer by updating its user field to null - PUT `<base URL>/<version>/load_balancer/{id}` with Remove: true
func (db *DBClient) RemoveLoadBalancer(ctx context.Context, id string) (*types.LoadBalancer, error) {
	if id == "" {
		return nil, errNoLoadBalancerID
	}

	loadBalancerRemoveJSON, err := json.Marshal(types.UpdateLoadBalancer{Remove: true})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(loadBalancerPath), id)

	return put[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerRemoveJSON, db.httpClient)
}

/* -- Delete Methods -- */

// DeleteLoadBalancerUser deletes a single User from a single Load Balancer  - DELETE `<base URL>/<version>/load_balancer/{id}/user/{userID}` with Remove: true
func (db *DBClient) DeleteLoadBalancerUser(ctx context.Context, loadBalancerID, userID string) (*types.LoadBalancer, error) {
	if loadBalancerID == "" {
		return nil, errNoLoadBalancerID
	}
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID, userPath, userID)

	return delete[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), db.httpClient)
}

/* -- PHD Client HTTP Funcs -- */

// Generic HTTP GET request
func get[T any](endpoint string, header http.Header, httpClient *httpclient.Client) (T, error) {
	var data T

	response, err := httpClient.Get(endpoint, header)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, parseErrorResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Generic HTTP POST request
func post[T any](endpoint string, header http.Header, postData []byte, httpClient *httpclient.Client) (T, error) {
	var data T

	postBody := bytes.NewBufferString(string(postData))

	response, err := httpClient.Post(endpoint, postBody, header)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, parseErrorResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Generic HTTP PUT request
func put[T any](endpoint string, header http.Header, postData []byte, httpClient *httpclient.Client) (T, error) {
	var data T

	postBody := bytes.NewBufferString(string(postData))

	response, err := httpClient.Put(endpoint, postBody, header)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, parseErrorResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Generic HTTP DELETE request
func delete[T any](endpoint string, header http.Header, httpClient *httpclient.Client) (T, error) {
	var data T

	response, err := httpClient.Delete(endpoint, header)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, parseErrorResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Parses the error reponse and returns the status code and error message
func parseErrorResponse(errResponse *http.Response) error {
	code := errResponse.StatusCode
	text := http.StatusText(code)

	errString := fmt.Errorf("%s. %d %s", errResponseNotOK, code, text)

	body, err := io.ReadAll(errResponse.Body)
	if err != nil {
		return errString
	}

	var errorMap map[string]string
	err = json.Unmarshal(body, &errorMap)
	if err != nil {
		return errString
	}

	if errorMessage, ok := errorMap["error"]; ok {
		errString = fmt.Errorf("%s: %s", errString, errorMessage)
	}

	return errString
}
