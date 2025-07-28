// Package vra8 provides primitives to interact with the openapi HTTP API.

package vra8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	// A struct embedding an interface allows for mock testing without
	// making any real n/w calls and also helps to pass a wrapper 
	// that logs, performs do() and then logs again etc
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {

	Login(ctx context.Context, params *LoginParams, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Authenticate(ctx context.Context, params *AuthenticationParams, body AuthenticationJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAccessTokenWithRefreshTokenWithBody request with any body
	GetAccessTokenWithRefreshTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetAccessTokenWithRefreshToken(ctx context.Context, body GetAccessTokenWithRefreshTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAccessTokenWithAuthorizationRequestWithBody request with any body
	GetAccessTokenWithAuthorizationRequestWithBody(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetAccessTokenWithAuthorizationRequest(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, body GetAccessTokenWithAuthorizationRequestJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	Logout(ctx context.Context, params *LogoutParams, body LogoutJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// LoginWithBody request with any body
	LoginWithBody(ctx context.Context, params *LoginParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

}


func (c *Client) Login(ctx context.Context, params *LoginParams, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}


// NewLoginRequest calls the generic Login builder with application/json body
func NewLoginRequest(server string, params *LoginParams, body LoginJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewLoginRequestWithBody(server, params, "application/json", bodyReader)
}


// NewLoginRequestWithBody generates requests for Login with any type of body
func NewLoginRequestWithBody(server string, params *LoginParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/csp/gateway/am/api/login")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.AccessToken != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "access_token", runtime.ParamLocationQuery, *params.AccessToken); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) GetAccessTokenWithRefreshTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAccessTokenWithRefreshTokenRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAccessTokenWithRefreshToken(ctx context.Context, body GetAccessTokenWithRefreshTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAccessTokenWithRefreshTokenRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAccessTokenWithAuthorizationRequestWithBody(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAccessTokenWithAuthorizationRequestRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAccessTokenWithAuthorizationRequest(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, body GetAccessTokenWithAuthorizationRequestJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAccessTokenWithAuthorizationRequestRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) LogoutWithBody(ctx context.Context, params *LogoutParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLogoutRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Logout(ctx context.Context, params *LogoutParams, body LogoutJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLogoutRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) LoginWithBody(ctx context.Context, params *LoginParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}


// NewGetAccessTokenWithRefreshTokenRequest calls the generic GetAccessTokenWithRefreshToken builder with application/json body
func NewGetAccessTokenWithRefreshTokenRequest(server string, body GetAccessTokenWithRefreshTokenJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetAccessTokenWithRefreshTokenRequestWithBody(server, "application/json", bodyReader)
}

// NewGetAccessTokenWithRefreshTokenRequestWithBody generates requests for GetAccessTokenWithRefreshToken with any type of body
func NewGetAccessTokenWithRefreshTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/csp/gateway/am/api/auth/api-tokens/authorize")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetAccessTokenWithAuthorizationRequestRequest calls the generic GetAccessTokenWithAuthorizationRequest builder with application/json body
func NewGetAccessTokenWithAuthorizationRequestRequest(server string, params *GetAccessTokenWithAuthorizationRequestParams, body GetAccessTokenWithAuthorizationRequestJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetAccessTokenWithAuthorizationRequestRequestWithBody(server, params, "application/json", bodyReader)
}

// NewGetAccessTokenWithAuthorizationRequestRequestWithBody generates requests for GetAccessTokenWithAuthorizationRequest with any type of body
func NewGetAccessTokenWithAuthorizationRequestRequestWithBody(server string, params *GetAccessTokenWithAuthorizationRequestParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/csp/gateway/am/api/auth/authorize")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.Authorization != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "authorization", runtime.ParamLocationHeader, *params.Authorization)
			if err != nil {
				return nil, err
			}

			req.Header.Set("authorization", headerParam0)
		}

	}

	return req, nil
}

// NewLogoutRequest calls the generic Logout builder with application/json body
func NewLogoutRequest(server string, params *LogoutParams, body LogoutJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewLogoutRequestWithBody(server, params, "application/json", bodyReader)
}

// NewLogoutRequestWithBody generates requests for Logout with any type of body
func NewLogoutRequestWithBody(server string, params *LogoutParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/csp/gateway/am/api/auth/logout")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.TheAccessTokenToBeInvalidated != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithLocation("simple", false, "The access token to be invalidated.", runtime.ParamLocationHeader, *params.TheAccessTokenToBeInvalidated)
			if err != nil {
				return nil, err
			}

			req.Header.Set("The access token to be invalidated.", headerParam0)
		}

	}

	return req, nil
}



func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetOpenidConfigurationWithResponse request
	GetOpenidConfigurationWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetOpenidConfigurationResponse, error)

	// GetAccessTokenWithRefreshTokenWithBodyWithResponse request with any body
	GetAccessTokenWithRefreshTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetAccessTokenWithRefreshTokenResponse, error)

	GetAccessTokenWithRefreshTokenWithResponse(ctx context.Context, body GetAccessTokenWithRefreshTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*GetAccessTokenWithRefreshTokenResponse, error)

	// GetAccessTokenWithAuthorizationRequestWithBodyWithResponse request with any body
	GetAccessTokenWithAuthorizationRequestWithBodyWithResponse(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetAccessTokenWithAuthorizationRequestResponse, error)

	GetAccessTokenWithAuthorizationRequestWithResponse(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, body GetAccessTokenWithAuthorizationRequestJSONRequestBody, reqEditors ...RequestEditorFn) (*GetAccessTokenWithAuthorizationRequestResponse, error)

	// GetKeysWithResponse request
	GetKeysWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetKeysResponse, error)

	// LogoutWithBodyWithResponse request with any body
	LogoutWithBodyWithResponse(ctx context.Context, params *LogoutParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LogoutResponse, error)

	LogoutWithResponse(ctx context.Context, params *LogoutParams, body LogoutJSONRequestBody, reqEditors ...RequestEditorFn) (*LogoutResponse, error)

	// GetAccessTokenPkceFlowWithBodyWithResponse request with any body
	GetAccessTokenPkceFlowWithBodyWithResponse(ctx context.Context, params *GetAccessTokenPkceFlowParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetAccessTokenPkceFlowResponse, error)

	GetAccessTokenPkceFlowWithFormdataBodyWithResponse(ctx context.Context, params *GetAccessTokenPkceFlowParams, body GetAccessTokenPkceFlowFormdataRequestBody, reqEditors ...RequestEditorFn) (*GetAccessTokenPkceFlowResponse, error)

	// GetPublicKeyWithResponse request
	GetPublicKeyWithResponse(ctx context.Context, params *GetPublicKeyParams, reqEditors ...RequestEditorFn) (*GetPublicKeyResponse, error)

	// SearchGroupsWithResponse request
	SearchGroupsWithResponse(ctx context.Context, params *SearchGroupsParams, reqEditors ...RequestEditorFn) (*SearchGroupsResponse, error)

	// GetLoggedInUserWithResponse request
	GetLoggedInUserWithResponse(ctx context.Context, params *GetLoggedInUserParams, reqEditors ...RequestEditorFn) (*GetLoggedInUserResponse, error)

	// GetUserDefaultOrgWithResponse request
	GetUserDefaultOrgWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetUserDefaultOrgResponse, error)

	// GetLoggedInUserDetailsWithResponse request
	GetLoggedInUserDetailsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetLoggedInUserDetailsResponse, error)

	// GetUserOrgs1WithResponse request
	GetUserOrgs1WithResponse(ctx context.Context, params *GetUserOrgs1Params, reqEditors ...RequestEditorFn) (*GetUserOrgs1Response, error)

	// GetLoggedInUserGroupsOnOrgWithResponse request
	GetLoggedInUserGroupsOnOrgWithResponse(ctx context.Context, orgId string, reqEditors ...RequestEditorFn) (*GetLoggedInUserGroupsOnOrgResponse, error)

	// GetUserOrgInfoWithResponse request
	GetUserOrgInfoWithResponse(ctx context.Context, orgId string, reqEditors ...RequestEditorFn) (*GetUserOrgInfoResponse, error)

	// GetUserOrgRolesWithResponse request
	GetUserOrgRolesWithResponse(ctx context.Context, orgId string, reqEditors ...RequestEditorFn) (*GetUserOrgRolesResponse, error)

	// GetUserOrgServiceRolesWithResponse request
	GetUserOrgServiceRolesWithResponse(ctx context.Context, orgId string, params *GetUserOrgServiceRolesParams, reqEditors ...RequestEditorFn) (*GetUserOrgServiceRolesResponse, error)

	// GetPrincipalUserProfileWithResponse request
	GetPrincipalUserProfileWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPrincipalUserProfileResponse, error)

	// UpdateUserProfileWithBodyWithResponse request with any body
	UpdateUserProfileWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserProfileResponse, error)

	UpdateUserProfileWithResponse(ctx context.Context, body UpdateUserProfileJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserProfileResponse, error)

	// UpdateUserPreferencesWithBodyWithResponse request with any body
	UpdateUserPreferencesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserPreferencesResponse, error)

	UpdateUserPreferencesWithResponse(ctx context.Context, body UpdateUserPreferencesJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserPreferencesResponse, error)

	// LoginWithBodyWithResponse request with any body
	LoginWithBodyWithResponse(ctx context.Context, params *LoginParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginResponse, error)

	LoginWithResponse(ctx context.Context, params *LoginParams, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginResponse, error)

	// LoginOauthWithBodyWithResponse request with any body
	LoginOauthWithBodyWithResponse(ctx context.Context, params *LoginOauthParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginOauthResponse, error)

	LoginOauthWithResponse(ctx context.Context, params *LoginOauthParams, body LoginOauthJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginOauthResponse, error)

	// GetByIdWithResponse request
	GetByIdWithResponse(ctx context.Context, orgId string, reqEditors ...RequestEditorFn) (*GetByIdResponse, error)

	// PatchOrgWithBodyWithResponse request with any body
	PatchOrgWithBodyWithResponse(ctx context.Context, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchOrgResponse, error)

	PatchOrgWithResponse(ctx context.Context, orgId string, body PatchOrgJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchOrgResponse, error)

	// RemoveGroupsFromOrganizationWithBodyWithResponse request with any body
	RemoveGroupsFromOrganizationWithBodyWithResponse(ctx context.Context, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RemoveGroupsFromOrganizationResponse, error)

	RemoveGroupsFromOrganizationWithResponse(ctx context.Context, orgId string, body RemoveGroupsFromOrganizationJSONRequestBody, reqEditors ...RequestEditorFn) (*RemoveGroupsFromOrganizationResponse, error)

	// GetOrganizationGroupsWithResponse request
	GetOrganizationGroupsWithResponse(ctx context.Context, orgId string, params *GetOrganizationGroupsParams, reqEditors ...RequestEditorFn) (*GetOrganizationGroupsResponse, error)

	// SearchOrgGroupsWithResponse request
	SearchOrgGroupsWithResponse(ctx context.Context, orgId string, params *SearchOrgGroupsParams, reqEditors ...RequestEditorFn) (*SearchOrgGroupsResponse, error)

	// GetNestedGroupsFromADGroupWithResponse request
	GetNestedGroupsFromADGroupWithResponse(ctx context.Context, orgId string, groupId string, params *GetNestedGroupsFromADGroupParams, reqEditors ...RequestEditorFn) (*GetNestedGroupsFromADGroupResponse, error)

	// GetGroupRolesOnOrganizationWithResponse request
	GetGroupRolesOnOrganizationWithResponse(ctx context.Context, orgId string, groupId string, reqEditors ...RequestEditorFn) (*GetGroupRolesOnOrganizationResponse, error)

	// UpdateGroupRolesOnOrganizationWithBodyWithResponse request with any body
	UpdateGroupRolesOnOrganizationWithBodyWithResponse(ctx context.Context, orgId string, groupId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateGroupRolesOnOrganizationResponse, error)

	UpdateGroupRolesOnOrganizationWithResponse(ctx context.Context, orgId string, groupId string, body UpdateGroupRolesOnOrganizationJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateGroupRolesOnOrganizationResponse, error)

	// GetPaginatedGroupUsersWithResponse request
	GetPaginatedGroupUsersWithResponse(ctx context.Context, orgId string, groupId string, params *GetPaginatedGroupUsersParams, reqEditors ...RequestEditorFn) (*GetPaginatedGroupUsersResponse, error)

	// DeleteOrgScopedOAuthClientWithBodyWithResponse request with any body
	DeleteOrgScopedOAuthClientWithBodyWithResponse(ctx context.Context, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*DeleteOrgScopedOAuthClientResponse, error)

	DeleteOrgScopedOAuthClientWithResponse(ctx context.Context, orgId string, body DeleteOrgScopedOAuthClientJSONRequestBody, reqEditors ...RequestEditorFn) (*DeleteOrgScopedOAuthClientResponse, error)

	// CreateOrgScopedOAuthClientWithBodyWithResponse request with any body
	CreateOrgScopedOAuthClientWithBodyWithResponse(ctx context.Context, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateOrgScopedOAuthClientResponse, error)

	CreateOrgScopedOAuthClientWithResponse(ctx context.Context, orgId string, body CreateOrgScopedOAuthClientJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateOrgScopedOAuthClientResponse, error)

	// GetOrgScopedOAuthClientWithResponse request
	GetOrgScopedOAuthClientWithResponse(ctx context.Context, orgId string, oauthAppId string, reqEditors ...RequestEditorFn) (*GetOrgScopedOAuthClientResponse, error)

	// GetOrgRolesWithResponse request
	GetOrgRolesWithResponse(ctx context.Context, orgId string, params *GetOrgRolesParams, reqEditors ...RequestEditorFn) (*GetOrgRolesResponse, error)

	// PatchOrgRolesWithBodyWithResponse request with any body
	PatchOrgRolesWithBodyWithResponse(ctx context.Context, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchOrgRolesResponse, error)

	PatchOrgRolesWithResponse(ctx context.Context, orgId string, body PatchOrgRolesJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchOrgRolesResponse, error)

	// GetRoleByOrgIdAndRoleIdWithResponse request
	GetRoleByOrgIdAndRoleIdWithResponse(ctx context.Context, orgId string, roleId string, reqEditors ...RequestEditorFn) (*GetRoleByOrgIdAndRoleIdResponse, error)

	// GetOrgSubOrgsWithResponse request
	GetOrgSubOrgsWithResponse(ctx context.Context, orgId string, reqEditors ...RequestEditorFn) (*GetOrgSubOrgsResponse, error)

	// GetPaginatedOrgUsersInfo1WithResponse request
	GetPaginatedOrgUsersInfo1WithResponse(ctx context.Context, orgId string, params *GetPaginatedOrgUsersInfo1Params, reqEditors ...RequestEditorFn) (*GetPaginatedOrgUsersInfo1Response, error)

	// SearchUsersWithResponse request
	SearchUsersWithResponse(ctx context.Context, orgId string, params *SearchUsersParams, reqEditors ...RequestEditorFn) (*SearchUsersResponse, error)

	// GetAccessTokenInfoWithResponse request
	GetAccessTokenInfoWithResponse(ctx context.Context, params *GetAccessTokenInfoParams, reqEditors ...RequestEditorFn) (*GetAccessTokenInfoResponse, error)

	// GetUserInAnyOrganization1WithResponse request
	GetUserInAnyOrganization1WithResponse(ctx context.Context, acct string, params *GetUserInAnyOrganization1Params, reqEditors ...RequestEditorFn) (*GetUserInAnyOrganization1Response, error)

	// GetUserInfoInOrganization1WithResponse request
	GetUserInfoInOrganization1WithResponse(ctx context.Context, acct string, orgId string, reqEditors ...RequestEditorFn) (*GetUserInfoInOrganization1Response, error)

	// GetUserRolesOnOrgWithGroupInfoWithResponse request
	GetUserRolesOnOrgWithGroupInfoWithResponse(ctx context.Context, userId string, orgId string, reqEditors ...RequestEditorFn) (*GetUserRolesOnOrgWithGroupInfoResponse, error)

	// GetUserRolesInOrganization1WithResponse request
	GetUserRolesInOrganization1WithResponse(ctx context.Context, userId string, orgId string, reqEditors ...RequestEditorFn) (*GetUserRolesInOrganization1Response, error)

	// PatchUserRolesInOrganizationWithBodyWithResponse request with any body
	PatchUserRolesInOrganizationWithBodyWithResponse(ctx context.Context, userId string, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchUserRolesInOrganizationResponse, error)

	PatchUserRolesInOrganizationWithResponse(ctx context.Context, userId string, orgId string, body PatchUserRolesInOrganizationJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchUserRolesInOrganizationResponse, error)

	// GetUserServiceRolesInOrganization1WithResponse request
	GetUserServiceRolesInOrganization1WithResponse(ctx context.Context, userId string, orgId string, params *GetUserServiceRolesInOrganization1Params, reqEditors ...RequestEditorFn) (*GetUserServiceRolesInOrganization1Response, error)

	// PatchUserServiceRolesInOrganizationWithBodyWithResponse request with any body
	PatchUserServiceRolesInOrganizationWithBodyWithResponse(ctx context.Context, userId string, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchUserServiceRolesInOrganizationResponse, error)

	PatchUserServiceRolesInOrganizationWithResponse(ctx context.Context, userId string, orgId string, body PatchUserServiceRolesInOrganizationJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchUserServiceRolesInOrganizationResponse, error)

	// GetUserShortInfoInOrganizationWithResponse request
	GetUserShortInfoInOrganizationWithResponse(ctx context.Context, userId string, orgId string, reqEditors ...RequestEditorFn) (*GetUserShortInfoInOrganizationResponse, error)

	// GetUserOrgsWithResponse request
	GetUserOrgsWithResponse(ctx context.Context, params *GetUserOrgsParams, reqEditors ...RequestEditorFn) (*GetUserOrgsResponse, error)

	// GetPaginatedOrgUsersInfoWithResponse request
	GetPaginatedOrgUsersInfoWithResponse(ctx context.Context, orgId string, params *GetPaginatedOrgUsersInfoParams, reqEditors ...RequestEditorFn) (*GetPaginatedOrgUsersInfoResponse, error)

	// GetUserInAnyOrganizationWithResponse request
	GetUserInAnyOrganizationWithResponse(ctx context.Context, userId string, params *GetUserInAnyOrganizationParams, reqEditors ...RequestEditorFn) (*GetUserInAnyOrganizationResponse, error)

	// GetUserInfoInOrganizationWithResponse request
	GetUserInfoInOrganizationWithResponse(ctx context.Context, userId string, orgId string, reqEditors ...RequestEditorFn) (*GetUserInfoInOrganizationResponse, error)

	// GetUserRolesInOrganizationWithResponse request
	GetUserRolesInOrganizationWithResponse(ctx context.Context, userId string, orgId string, reqEditors ...RequestEditorFn) (*GetUserRolesInOrganizationResponse, error)

	// GetUserServiceRolesInOrganizationWithResponse request
	GetUserServiceRolesInOrganizationWithResponse(ctx context.Context, userId string, orgId string, params *GetUserServiceRolesInOrganizationParams, reqEditors ...RequestEditorFn) (*GetUserServiceRolesInOrganizationResponse, error)

	// PatchUserRolesOnOrganizationWithBodyWithResponse request with any body
	PatchUserRolesOnOrganizationWithBodyWithResponse(ctx context.Context, userId string, orgId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchUserRolesOnOrganizationResponse, error)

	PatchUserRolesOnOrganizationWithResponse(ctx context.Context, userId string, orgId string, body PatchUserRolesOnOrganizationJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchUserRolesOnOrganizationResponse, error)

	// GetAllServiceDefinitionsWithResponse request
	GetAllServiceDefinitionsWithResponse(ctx context.Context, params *GetAllServiceDefinitionsParams, reqEditors ...RequestEditorFn) (*GetAllServiceDefinitionsResponse, error)

	// GetAllByOrgServiceDefinitions1WithResponse request
	GetAllByOrgServiceDefinitions1WithResponse(ctx context.Context, orgId string, params *GetAllByOrgServiceDefinitions1Params, reqEditors ...RequestEditorFn) (*GetAllByOrgServiceDefinitions1Response, error)

	// GetPagedServiceDefinitionOrgsWithResponse request
	GetPagedServiceDefinitionOrgsWithResponse(ctx context.Context, serviceDefinitionId string, params *GetPagedServiceDefinitionOrgsParams, reqEditors ...RequestEditorFn) (*GetPagedServiceDefinitionOrgsResponse, error)

	// GetAllByOrgServiceDefinitionsWithResponse request
	GetAllByOrgServiceDefinitionsWithResponse(ctx context.Context, orgId string, params *GetAllByOrgServiceDefinitionsParams, reqEditors ...RequestEditorFn) (*GetAllByOrgServiceDefinitionsResponse, error)

	// CheckIDTokenWithResponse request
	CheckIDTokenWithResponse(ctx context.Context, params *CheckIDTokenParams, reqEditors ...RequestEditorFn) (*CheckIDTokenResponse, error)
}

type GetAccessTokenWithRefreshTokenResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetAccessTokenWithRefreshTokenResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAccessTokenWithRefreshTokenResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAccessTokenWithAuthorizationRequestResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetAccessTokenWithAuthorizationRequestResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAccessTokenWithAuthorizationRequestResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type LogoutResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r LogoutResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r LogoutResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type LoginResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r LoginResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r LoginResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAccessTokenWithRefreshTokenWithBodyWithResponse request with arbitrary body returning *GetAccessTokenWithRefreshTokenResponse
func (c *ClientWithResponses) GetAccessTokenWithRefreshTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetAccessTokenWithRefreshTokenResponse, error) {
	rsp, err := c.GetAccessTokenWithRefreshTokenWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccessTokenWithRefreshTokenResponse(rsp)
}

func (c *ClientWithResponses) GetAccessTokenWithRefreshTokenWithResponse(ctx context.Context, body GetAccessTokenWithRefreshTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*GetAccessTokenWithRefreshTokenResponse, error) {
	rsp, err := c.GetAccessTokenWithRefreshToken(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccessTokenWithRefreshTokenResponse(rsp)
}

// GetAccessTokenWithAuthorizationRequestWithBodyWithResponse request with arbitrary body returning *GetAccessTokenWithAuthorizationRequestResponse
func (c *ClientWithResponses) GetAccessTokenWithAuthorizationRequestWithBodyWithResponse(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetAccessTokenWithAuthorizationRequestResponse, error) {
	rsp, err := c.GetAccessTokenWithAuthorizationRequestWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccessTokenWithAuthorizationRequestResponse(rsp)
}

func (c *ClientWithResponses) GetAccessTokenWithAuthorizationRequestWithResponse(ctx context.Context, params *GetAccessTokenWithAuthorizationRequestParams, body GetAccessTokenWithAuthorizationRequestJSONRequestBody, reqEditors ...RequestEditorFn) (*GetAccessTokenWithAuthorizationRequestResponse, error) {
	rsp, err := c.GetAccessTokenWithAuthorizationRequest(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccessTokenWithAuthorizationRequestResponse(rsp)
}

// LogoutWithBodyWithResponse request with arbitrary body returning *LogoutResponse
func (c *ClientWithResponses) LogoutWithBodyWithResponse(ctx context.Context, params *LogoutParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LogoutResponse, error) {
	rsp, err := c.LogoutWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLogoutResponse(rsp)
}

func (c *ClientWithResponses) LogoutWithResponse(ctx context.Context, params *LogoutParams, body LogoutJSONRequestBody, reqEditors ...RequestEditorFn) (*LogoutResponse, error) {
	rsp, err := c.Logout(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLogoutResponse(rsp)
}

// LoginWithBodyWithResponse request with arbitrary body returning *LoginResponse
func (c *ClientWithResponses) LoginWithBodyWithResponse(ctx context.Context, params *LoginParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginResponse, error) {
	rsp, err := c.LoginWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginResponse(rsp)
}

func (c *ClientWithResponses) LoginWithResponse(ctx context.Context, params *LoginParams, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginResponse, error) {
	rsp, err := c.Login(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginResponse(rsp)
}

// GetAccessTokenInfoWithResponse request returning *GetAccessTokenInfoResponse
func (c *ClientWithResponses) GetAccessTokenInfoWithResponse(ctx context.Context, params *GetAccessTokenInfoParams, reqEditors ...RequestEditorFn) (*GetAccessTokenInfoResponse, error) {
	rsp, err := c.GetAccessTokenInfo(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAccessTokenInfoResponse(rsp)
}

// ParseGetAccessTokenWithRefreshTokenResponse parses an HTTP response from a GetAccessTokenWithRefreshTokenWithResponse call
func ParseGetAccessTokenWithRefreshTokenResponse(rsp *http.Response) (*GetAccessTokenWithRefreshTokenResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAccessTokenWithRefreshTokenResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetAccessTokenWithAuthorizationRequestResponse parses an HTTP response from a GetAccessTokenWithAuthorizationRequestWithResponse call
func ParseGetAccessTokenWithAuthorizationRequestResponse(rsp *http.Response) (*GetAccessTokenWithAuthorizationRequestResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAccessTokenWithAuthorizationRequestResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}


// ParseLogoutResponse parses an HTTP response from a LogoutWithResponse call
func ParseLogoutResponse(rsp *http.Response) (*LogoutResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LogoutResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseLoginResponse parses an HTTP response from a LoginWithResponse call
func ParseLoginResponse(rsp *http.Response) (*LoginResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

