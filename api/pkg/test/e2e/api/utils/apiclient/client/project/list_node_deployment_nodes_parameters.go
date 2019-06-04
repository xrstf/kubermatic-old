// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListNodeDeploymentNodesParams creates a new ListNodeDeploymentNodesParams object
// with the default values initialized.
func NewListNodeDeploymentNodesParams() *ListNodeDeploymentNodesParams {
	var ()
	return &ListNodeDeploymentNodesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListNodeDeploymentNodesParamsWithTimeout creates a new ListNodeDeploymentNodesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListNodeDeploymentNodesParamsWithTimeout(timeout time.Duration) *ListNodeDeploymentNodesParams {
	var ()
	return &ListNodeDeploymentNodesParams{

		timeout: timeout,
	}
}

// NewListNodeDeploymentNodesParamsWithContext creates a new ListNodeDeploymentNodesParams object
// with the default values initialized, and the ability to set a context for a request
func NewListNodeDeploymentNodesParamsWithContext(ctx context.Context) *ListNodeDeploymentNodesParams {
	var ()
	return &ListNodeDeploymentNodesParams{

		Context: ctx,
	}
}

// NewListNodeDeploymentNodesParamsWithHTTPClient creates a new ListNodeDeploymentNodesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListNodeDeploymentNodesParamsWithHTTPClient(client *http.Client) *ListNodeDeploymentNodesParams {
	var ()
	return &ListNodeDeploymentNodesParams{
		HTTPClient: client,
	}
}

/*ListNodeDeploymentNodesParams contains all the parameters to send to the API endpoint
for the list node deployment nodes operation typically these are written to a http.Request
*/
type ListNodeDeploymentNodesParams struct {

	/*ClusterID*/
	ClusterID string
	/*Dc*/
	Dc string
	/*HideInitialConditions*/
	HideInitialConditions *bool
	/*NodedeploymentID*/
	NodedeploymentID string
	/*ProjectID*/
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithTimeout(timeout time.Duration) *ListNodeDeploymentNodesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithContext(ctx context.Context) *ListNodeDeploymentNodesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithHTTPClient(client *http.Client) *ListNodeDeploymentNodesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithClusterID(clusterID string) *ListNodeDeploymentNodesParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithDc adds the dc to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithDc(dc string) *ListNodeDeploymentNodesParams {
	o.SetDc(dc)
	return o
}

// SetDc adds the dc to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetDc(dc string) {
	o.Dc = dc
}

// WithHideInitialConditions adds the hideInitialConditions to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithHideInitialConditions(hideInitialConditions *bool) *ListNodeDeploymentNodesParams {
	o.SetHideInitialConditions(hideInitialConditions)
	return o
}

// SetHideInitialConditions adds the hideInitialConditions to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetHideInitialConditions(hideInitialConditions *bool) {
	o.HideInitialConditions = hideInitialConditions
}

// WithNodedeploymentID adds the nodedeploymentID to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithNodedeploymentID(nodedeploymentID string) *ListNodeDeploymentNodesParams {
	o.SetNodedeploymentID(nodedeploymentID)
	return o
}

// SetNodedeploymentID adds the nodedeploymentId to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetNodedeploymentID(nodedeploymentID string) {
	o.NodedeploymentID = nodedeploymentID
}

// WithProjectID adds the projectID to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) WithProjectID(projectID string) *ListNodeDeploymentNodesParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the list node deployment nodes params
func (o *ListNodeDeploymentNodesParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *ListNodeDeploymentNodesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param dc
	if err := r.SetPathParam("dc", o.Dc); err != nil {
		return err
	}

	if o.HideInitialConditions != nil {

		// query param hideInitialConditions
		var qrHideInitialConditions bool
		if o.HideInitialConditions != nil {
			qrHideInitialConditions = *o.HideInitialConditions
		}
		qHideInitialConditions := swag.FormatBool(qrHideInitialConditions)
		if qHideInitialConditions != "" {
			if err := r.SetQueryParam("hideInitialConditions", qHideInitialConditions); err != nil {
				return err
			}
		}

	}

	// path param nodedeployment_id
	if err := r.SetPathParam("nodedeployment_id", o.NodedeploymentID); err != nil {
		return err
	}

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}