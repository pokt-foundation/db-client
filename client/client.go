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
	v2Types "github.com/pokt-foundation/portal-db/v2/types"
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
		// GetBlockchains returns all blockchains in the DB - GET `<base URL>/<version>/blockchain`
		GetBlockchains(ctx context.Context) ([]*types.Blockchain, error)
		// GetBlockchainByID returns a single Blockchain by its relay chain ID - GET `<base URL>/<version>/blockchain/{id}`
		GetBlockchainByID(ctx context.Context, blockchainID string) (*types.Blockchain, error)
		// GetApplications returns all Applications in the DB - GET `<base URL>/<version>/application`
		GetApplications(ctx context.Context) ([]*types.Application, error)
		// GetApplicationByID returns a single Application by its ID - GET `<base URL>/<version>/application/{id}`
		GetApplicationByID(ctx context.Context, applicationID string) (*types.Application, error)
		// GetApplicationsByUserID returns all the Applications for a user - GET `<base URL>/<version>/user/{userID}/application`
		GetApplicationsByUserID(ctx context.Context, userID string) ([]*types.Application, error)
		// GetLoadBalancers returns all Load Balancers in the DB - GET `<base URL>/<version>/load_balancer`
		GetLoadBalancers(ctx context.Context) ([]*types.LoadBalancer, error)
		// GetLoadBalancerByID returns a single Load Balancer by its ID - GET `<base URL>/<version>/load_balancer/{id}`
		GetLoadBalancerByID(ctx context.Context, loadBalancerID string) (*types.LoadBalancer, error)
		// GetLoadBalancersByUserID returns all the load balancers for a user - GET `<base URL>/<version>/user/{userID}/load_balancer`.*/
		// This method can be filtered by the user's role for a given LB. To return all LBs for the user pass nil for the roleNameFilter param.
		GetLoadBalancersByUserID(ctx context.Context, userID string, roleNameFilter *types.RoleName) ([]*types.LoadBalancer, error)
		// GetPendingLoadBalancersByUserID returns all the pending load balancers for an userID - GET `<base URL>/<version>/user/{userID}/load_balancer/pending`.*/
		GetPendingLoadBalancersByUserID(ctx context.Context, userID string) ([]*types.LoadBalancer, error)
		// GetLoadBalancersCountByUserID returns the number of loadbalancers owned by an userID - GET `<base URL>/<version>/user/{userID}/load_balancer/count`.`
		GetLoadBalancersCountByUserID(ctx context.Context, userID string) (int, error)
		// GetPayPlans returns all Pay Plans in the DB - GET `<base URL>/<version>/pay_plan`
		GetPayPlans(ctx context.Context) ([]*types.PayPlan, error)
		// GetPayPlanByType returns a single Pay Plan by its type - GET `<base URL>/<version>/pay_plan/{type}`
		GetPayPlanByType(ctx context.Context, payPlanType types.PayPlanType) (*types.PayPlan, error)
		// GetUserPermissionsByUserID returns all load balancer UserPermissions for a given User ID - GET `<base URL>/<version>/user/{userID}/permission`
		GetUserPermissionsByUserID(ctx context.Context, userID types.UserID) (*types.UserPermissions, error)
		// GetUserByUserID returns a single user bases on his User ID - GET `<base URL>/<version>/user/{userID}`
		GetUserByUserID(ctx context.Context, userID types.UserID) (*v2Types.User, error)
	}
	// IDBWriter interface contains write methods for interacting with the Pocket HTTP DB
	IDBWriter interface {
		// CreateBlockchain creates a single Blockchain in the DB - POST `<base URL>/<version>/blockchain`
		CreateBlockchain(ctx context.Context, blockchain types.Blockchain) (*types.Blockchain, error)
		// CreateBlockchainRedirect creates a single Blockchain Redirect in the DB - POST `<base URL>/<version>/blockchain/redirect`
		CreateBlockchainRedirect(ctx context.Context, redirect types.Redirect) (*types.Redirect, error)
		// CreateLoadBalancer creates a single Load Balancer in the DB - POST `<base URL>/<version>/load_balancer`
		CreateLoadBalancer(ctx context.Context, loadBalancer types.LoadBalancer) (*types.LoadBalancer, error)
		// CreateLoadBalancerUser adds a single User to a single Load Balancer in the DB - POST `<base URL>/<version>/load_balancer/{id}/user`
		CreateLoadBalancerUser(ctx context.Context, loadBalancerID string, user types.UserAccess) (*types.LoadBalancer, error)
		// CreatePortalUser adds a single User to the database and create a new account - POST `<base URL>/<version>/user`
		CreatePortalUser(ctx context.Context, userInput v2Types.CreateUser) (*v2Types.User, error)
		// CreateLoadBalancerIntegration adds account integrations to a single Load Balancer - POST `<base URL>/<version>/load_balancer/{id}/integration`
		CreateLoadBalancerIntegration(ctx context.Context, loadBalancerID string, integrationsInput types.AccountIntegrations) (*types.LoadBalancer, error)
		// ActivateBlockchain toggles a single Blockchain's `active` field` - PUT `<base URL>/<version>/blockchain/{id}/activate`
		ActivateBlockchain(ctx context.Context, blockchainID string, active bool) (bool, error)
		// UpdateAppFirstDateSurpassed updates a slice of Applications' FirstDateSurpassed fields in the DB - POST `<base URL>/<version>/first_date_surpassed`
		UpdateAppFirstDateSurpassed(ctx context.Context, updateInput types.UpdateFirstDateSurpassed) ([]*types.Application, error)
		// RemoveApplication removes a single Application by updating its status field - PUT `<base URL>/<version>/application/{id}` with Remove: true
		RemoveApplication(ctx context.Context, id string) (*types.Application, error)
		// UpdateBlockchain updates a single LoadBalancer in the DB - PUT `<base URL>/<version>/blockchain/{id}`
		UpdateBlockchain(ctx context.Context, blockchainID string, chainUpdate types.UpdateBlockchain) (*types.Blockchain, error)
		// UpdateLoadBalancer updates a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}`
		UpdateLoadBalancer(ctx context.Context, id string, lbUpdate types.UpdateApplication) (*types.LoadBalancer, error)
		// UpdateLoadBalancerUserRole updates a single User's role for a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}/user`
		UpdateLoadBalancerUserRole(ctx context.Context, loadBalancerID string, update types.UpdateUserAccess) (*types.LoadBalancer, error)
		// AcceptLoadBalancerUser updates a single User's UserID and Accepted fields for a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}/user/accept`
		AcceptLoadBalancerUser(ctx context.Context, loadBalancerID, userID string) (*types.LoadBalancer, error)
		// RemoveLoadBalancer removes a single LoadBalancer by updating its user field to null - PUT `<base URL>/<version>/load_balancer/{id}` with Remove: true
		RemoveLoadBalancer(ctx context.Context, id string) (*types.LoadBalancer, error)
		// DeleteLoadBalancerUser deletes a single User from a single Load Balancer  - DELETE `<base URL>/<version>/load_balancer/{id}/user/{userID}`
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
	permissionPath         subPath = "permission"
	acceptPath             subPath = "accept"
	pendingPath            subPath = "pending"
	countPath              subPath = "count"
)

// New API versions should be added to both the APIVersion enum and ValidAPIVersions map
const (
	V1 APIVersion = "v1"
)

var ValidAPIVersions = map[APIVersion]bool{
	V1: true,
}

var (
	errBaseURLNotProvided       error = errors.New("base URL not provided")
	errAPIKeyNotProvided        error = errors.New("API key not provided")
	errVersionNotProvided       error = errors.New("version not provided")
	errInvalidVersionProvided   error = errors.New("invalid version provided")
	errNoUserID                 error = errors.New("no user ID")
	errNoBlockchainID           error = errors.New("no blockchain ID")
	errNoApplicationID          error = errors.New("no application ID")
	errNoLoadBalancerID         error = errors.New("no load balancer ID")
	errNoPayPlanType            error = errors.New("no pay plan type")
	errInvalidBlockchainJSON    error = errors.New("invalid blockchain JSON")
	errInvalidAppJSON           error = errors.New("invalid application JSON")
	errInvalidLoadBalancerJSON  error = errors.New("invalid load balancer JSON")
	errInvalidIntegrationsJSON  error = errors.New("invalid integrations JSON")
	errInvalidActivationJSON    error = errors.New("invalid active field JSON")
	errInvalidRoleName          error = errors.New("invalid role name")
	errOwnerRequiresUpdateEmail error = errors.New("transferring ownership requires providing the updater's email")
	errInvalidRoleNameFilter    error = errors.New("invalid role name filter")
	errResponseNotOK            error = errors.New("Response not OK")
	errInvalidCreateUserJSON    error = errors.New("invalid create user JSON")
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

// GetBlockchainByID returns a single Blockchain by its relay chain ID - GET `<base URL>/<version>/blockchain/{id}`
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
// This method can be filtered by the user's role for a given LB. To return all LBs for the user pass nil for the roleNameFilter param.
func (db *DBClient) GetLoadBalancersByUserID(ctx context.Context, userID string, roleNameFilter *types.RoleName) ([]*types.LoadBalancer, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(userPath), userID, loadBalancerPath)

	if roleNameFilter != nil {
		filter := *roleNameFilter

		if !types.ValidRoleNames[filter] {
			return nil, errInvalidRoleNameFilter
		}

		endpoint = fmt.Sprintf("%s?filter=%s", endpoint, filter)
	}

	return get[[]*types.LoadBalancer](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetPendingLoadBalancersByUserID returns all the pending load balancers for an userID - GET `<base URL>/<version>/user/{portalID}/load_balancer/pending`.*/
func (db *DBClient) GetPendingLoadBalancersByUserID(ctx context.Context, userID string) ([]*types.LoadBalancer, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s/%s", db.versionedBasePath(userPath), userID, loadBalancerPath, pendingPath)

	return get[[]*types.LoadBalancer](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetLoadBalancersCountByUserID returns all the pending load balancers for an userID - GET `<base URL>/<version>/user/{portalID}/load_balancer/count`.*/
func (db *DBClient) GetLoadBalancersCountByUserID(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s/%s", db.versionedBasePath(userPath), userID, loadBalancerPath, countPath)

	return get[int](endpoint, db.getAuthHeaderForRead(), db.httpClient)
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

// GetUserPermissionsByUserID returns all load balancer UserPermissions for a given User ID - GET `<base URL>/<version>/user/{userID}/permission`
func (db *DBClient) GetUserPermissionsByUserID(ctx context.Context, userID types.UserID) (*types.UserPermissions, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(userPath), userID, permissionPath)

	return get[*types.UserPermissions](endpoint, db.getAuthHeaderForRead(), db.httpClient)
}

// GetUserByUserID returns a single user bases on his User ID - GET `<base URL>/<version>/user/{userID}`
func (db *DBClient) GetUserByUserID(ctx context.Context, userID types.UserID) (*v2Types.User, error) {
	if userID == "" {
		return nil, errNoUserID
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(userPath), userID)

	return get[*v2Types.User](endpoint, db.getAuthHeaderForRead(), db.httpClient)
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

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(blockchainPath), redirectPath)

	return post[*types.Redirect](endpoint, db.getAuthHeaderForWrite(), redirectJSON, db.httpClient)
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

// CreatePortalUser adds a single User to the database and create a new account - POST `<base URL>/<version>/user`
func (db *DBClient) CreatePortalUser(ctx context.Context, userInput v2Types.CreateUser) (*v2Types.User, error) {
	portalUserJSON, err := json.Marshal(userInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidCreateUserJSON, err)
	}

	return post[*v2Types.User](db.versionedBasePath(userPath), db.getAuthHeaderForWrite(), portalUserJSON, db.httpClient)
}

// CreateLoadBalancerIntegration adds account integrations to a single Load Balancer - POST `<base URL>/<version>/load_balancer/{id}/integration`
func (db *DBClient) CreateLoadBalancerIntegration(ctx context.Context, loadBalancerID string, integrationsInput types.AccountIntegrations) (*types.LoadBalancer, error) {
	integrationsJSON, err := json.Marshal(integrationsInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidIntegrationsJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID, "integration")

	return post[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), integrationsJSON, db.httpClient)
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

// UpdateAppFirstDateSurpassed updates a slice of Applications' FirstDateSurpassed fields in the DB - POST `<base URL>/<version>/first_date_surpassed`
func (db *DBClient) UpdateAppFirstDateSurpassed(ctx context.Context, updateInput types.UpdateFirstDateSurpassed) ([]*types.Application, error) {
	firstDateSurpassedJSON, err := json.Marshal(updateInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAppJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(applicationPath), firstDateSurpassedPath)

	return post[[]*types.Application](endpoint, db.getAuthHeaderForWrite(), firstDateSurpassedJSON, db.httpClient)
}

// UpdateBlockchain updates a single LoadBalancer in the DB - PUT `<base URL>/<version>/blockchain/{id}`
func (db *DBClient) UpdateBlockchain(ctx context.Context, blockchainID string, chainUpdate types.UpdateBlockchain) (*types.Blockchain, error) {
	if blockchainID == "" {
		return nil, errNoBlockchainID
	}

	blockchainUpdateJSON, err := json.Marshal(chainUpdate)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidBlockchainJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(blockchainPath), blockchainID)

	return put[*types.Blockchain](endpoint, db.getAuthHeaderForWrite(), blockchainUpdateJSON, db.httpClient)
}

// UpdateLoadBalancer updates a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}`
// NOTE: It is intended that the UpdateAppliation struct be used here as part of the V2 changes
func (db *DBClient) UpdateLoadBalancer(ctx context.Context, id string, lbUpdate types.UpdateApplication) (*types.LoadBalancer, error) {
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
func (db *DBClient) UpdateLoadBalancerUserRole(ctx context.Context, loadBalancerID string, update types.UpdateUserAccess) (*types.LoadBalancer, error) {
	if loadBalancerID == "" {
		return nil, errNoLoadBalancerID
	}
	if update.UserID == "" {
		return nil, errNoUserID
	}
	if update.RoleName == types.RoleName("") || !types.ValidRoleNames[update.RoleName] {
		return nil, errInvalidRoleName
	}
	if update.RoleName == types.RoleOwner && update.Email == "" {
		return nil, errOwnerRequiresUpdateEmail
	}

	updateStruct := types.UpdateUserAccess{
		UserID:   update.UserID,
		RoleName: update.RoleName,
		Email:    update.Email,
	}
	/* 	if update.RoleName == types.RoleOwner {
		updateStruct.UpdaterEmail = update.UpdaterEmail
	} */

	loadBalancerUserUpdateJSON, err := json.Marshal(updateStruct)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID, userPath)

	return put[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerUserUpdateJSON, db.httpClient)
}

// AcceptLoadBalancerUser updates a single User's UserID and Accepted fields for a single LoadBalancer in the DB - PUT `<base URL>/<version>/load_balancer/{id}/user/accept`
func (db *DBClient) AcceptLoadBalancerUser(ctx context.Context, loadBalancerID, userID string) (*types.LoadBalancer, error) {
	if userID == "" {
		return nil, errNoUserID
	}
	if loadBalancerID == "" {
		return nil, errNoLoadBalancerID
	}
	if userID == "" {
		return nil, errNoUserID
	}

	loadBalancerAcceptUserJSON, err := json.Marshal(types.UpdateUserAccess{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidLoadBalancerJSON, err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s/%s", db.versionedBasePath(loadBalancerPath), loadBalancerID, userPath, acceptPath)

	return put[*types.LoadBalancer](endpoint, db.getAuthHeaderForWrite(), loadBalancerAcceptUserJSON, db.httpClient)
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

	endpoint := fmt.Sprintf("%s/%s", db.versionedBasePath(loadBalancerPath), id)

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

// DeleteLoadBalancerUser deletes a single User from a single Load Balancer by user userID - DELETE `<base URL>/<version>/load_balancer/{id}/user/{userID}`
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
