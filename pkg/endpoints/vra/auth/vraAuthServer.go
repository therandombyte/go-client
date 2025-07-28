// Package vra8 provides primitives to interact with the openapi HTTP API.

package vra8

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Exchange organization scoped API-token for user access token.
	// (POST /csp/gateway/am/api/auth/api-tokens/authorize)
	GetAccessTokenWithRefreshToken(ctx echo.Context) error
	// Get an access token.
	// (POST /csp/gateway/am/api/auth/authorize)
	GetAccessTokenWithAuthorizationRequest(ctx echo.Context, params GetAccessTokenWithAuthorizationRequestParams) error
	// Performs logout.
	// (POST /csp/gateway/am/api/auth/logout)
	Logout(ctx echo.Context, params LogoutParams) error
	// Login.
	// (POST /csp/gateway/am/api/login)
	Login(ctx echo.Context, params LoginParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAccessTokenWithRefreshToken converts echo context to params.
func (w *ServerInterfaceWrapper) GetAccessTokenWithRefreshToken(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAccessTokenWithRefreshToken(ctx)
	return err
}

// GetAccessTokenWithAuthorizationRequest converts echo context to params.
func (w *ServerInterfaceWrapper) GetAccessTokenWithAuthorizationRequest(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAccessTokenWithAuthorizationRequestParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for authorization, got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "authorization", valueList[0], &Authorization, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorization: %s", err))
		}

		params.Authorization = &Authorization
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAccessTokenWithAuthorizationRequest(ctx, params)
	return err
}

// Logout converts echo context to params.
func (w *ServerInterfaceWrapper) Logout(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params LogoutParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "The access token to be invalidated." -------------
	if valueList, found := headers[http.CanonicalHeaderKey("The access token to be invalidated.")]; found {
		var TheAccessTokenToBeInvalidated string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for The access token to be invalidated., got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "The access token to be invalidated.", valueList[0], &TheAccessTokenToBeInvalidated, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter The access token to be invalidated.: %s", err))
		}

		params.TheAccessTokenToBeInvalidated = &TheAccessTokenToBeInvalidated
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Logout(ctx, params)
	return err
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params LoginParams
	// ------------- Optional query parameter "access_token" -------------

	err = runtime.BindQueryParameter("form", true, false, "access_token", ctx.QueryParams(), &params.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter access_token: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx, params)
	return err
}


// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/csp/gateway/am/api/auth/api-tokens/authorize", wrapper.GetAccessTokenWithRefreshToken)
	router.POST(baseURL+"/csp/gateway/am/api/auth/authorize", wrapper.GetAccessTokenWithAuthorizationRequest)
	router.POST(baseURL+"/csp/gateway/am/api/auth/logout", wrapper.Logout)
	router.POST(baseURL+"/csp/gateway/am/api/login", wrapper.Login)
}
