// Code generated by go-swagger; DO NOT EDIT.

package openstack

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new openstack API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for openstack API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	ListOpenstackNetworks(params *ListOpenstackNetworksParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackNetworksOK, error)

	ListOpenstackNetworksNoCredentials(params *ListOpenstackNetworksNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackNetworksNoCredentialsOK, error)

	ListOpenstackSecurityGroups(params *ListOpenstackSecurityGroupsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSecurityGroupsOK, error)

	ListOpenstackSecurityGroupsNoCredentials(params *ListOpenstackSecurityGroupsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSecurityGroupsNoCredentialsOK, error)

	ListOpenstackSizes(params *ListOpenstackSizesParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSizesOK, error)

	ListOpenstackSizesNoCredentials(params *ListOpenstackSizesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSizesNoCredentialsOK, error)

	ListOpenstackSubnets(params *ListOpenstackSubnetsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSubnetsOK, error)

	ListOpenstackSubnetsNoCredentials(params *ListOpenstackSubnetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSubnetsNoCredentialsOK, error)

	ListOpenstackTenants(params *ListOpenstackTenantsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackTenantsOK, error)

	ListOpenstackTenantsNoCredentials(params *ListOpenstackTenantsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackTenantsNoCredentialsOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  ListOpenstackNetworks Lists networks from openstack
*/
func (a *Client) ListOpenstackNetworks(params *ListOpenstackNetworksParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackNetworksOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackNetworksParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackNetworks",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/openstack/networks",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackNetworksReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackNetworksOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackNetworksDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackNetworksNoCredentials Lists networks from openstack
*/
func (a *Client) ListOpenstackNetworksNoCredentials(params *ListOpenstackNetworksNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackNetworksNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackNetworksNoCredentialsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackNetworksNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/networks",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackNetworksNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackNetworksNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackNetworksNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSecurityGroups Lists security groups from openstack
*/
func (a *Client) ListOpenstackSecurityGroups(params *ListOpenstackSecurityGroupsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSecurityGroupsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSecurityGroupsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSecurityGroups",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/openstack/securitygroups",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSecurityGroupsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSecurityGroupsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSecurityGroupsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSecurityGroupsNoCredentials Lists security groups from openstack
*/
func (a *Client) ListOpenstackSecurityGroupsNoCredentials(params *ListOpenstackSecurityGroupsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSecurityGroupsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSecurityGroupsNoCredentialsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSecurityGroupsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/securitygroups",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSecurityGroupsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSecurityGroupsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSecurityGroupsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSizes Lists sizes from openstack
*/
func (a *Client) ListOpenstackSizes(params *ListOpenstackSizesParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSizesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSizesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSizes",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/openstack/sizes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSizesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSizesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSizesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSizesNoCredentials Lists sizes from openstack
*/
func (a *Client) ListOpenstackSizesNoCredentials(params *ListOpenstackSizesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSizesNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSizesNoCredentialsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSizesNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/sizes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSizesNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSizesNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSizesNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSubnets Lists subnets from openstack
*/
func (a *Client) ListOpenstackSubnets(params *ListOpenstackSubnetsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSubnetsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSubnetsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSubnets",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/openstack/subnets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSubnetsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSubnetsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSubnetsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackSubnetsNoCredentials Lists subnets from openstack
*/
func (a *Client) ListOpenstackSubnetsNoCredentials(params *ListOpenstackSubnetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackSubnetsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackSubnetsNoCredentialsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackSubnetsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/subnets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackSubnetsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackSubnetsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackSubnetsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackTenants Lists tenants from openstack
*/
func (a *Client) ListOpenstackTenants(params *ListOpenstackTenantsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackTenantsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackTenantsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackTenants",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/openstack/tenants",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackTenantsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackTenantsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackTenantsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListOpenstackTenantsNoCredentials Lists tenants from openstack
*/
func (a *Client) ListOpenstackTenantsNoCredentials(params *ListOpenstackTenantsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter) (*ListOpenstackTenantsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListOpenstackTenantsNoCredentialsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listOpenstackTenantsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/tenants",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListOpenstackTenantsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListOpenstackTenantsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListOpenstackTenantsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
