// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/models"
)

// CreateRoleBindingReader is a Reader for the CreateRoleBinding structure.
type CreateRoleBindingReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateRoleBindingReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewCreateRoleBindingCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewCreateRoleBindingUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 403:
		result := NewCreateRoleBindingForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewCreateRoleBindingDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateRoleBindingCreated creates a CreateRoleBindingCreated with default headers values
func NewCreateRoleBindingCreated() *CreateRoleBindingCreated {
	return &CreateRoleBindingCreated{}
}

/*CreateRoleBindingCreated handles this case with default header values.

RoleBinding
*/
type CreateRoleBindingCreated struct {
	Payload *models.RoleBinding
}

func (o *CreateRoleBindingCreated) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings][%d] createRoleBindingCreated  %+v", 201, o.Payload)
}

func (o *CreateRoleBindingCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RoleBinding)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateRoleBindingUnauthorized creates a CreateRoleBindingUnauthorized with default headers values
func NewCreateRoleBindingUnauthorized() *CreateRoleBindingUnauthorized {
	return &CreateRoleBindingUnauthorized{}
}

/*CreateRoleBindingUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type CreateRoleBindingUnauthorized struct {
}

func (o *CreateRoleBindingUnauthorized) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings][%d] createRoleBindingUnauthorized ", 401)
}

func (o *CreateRoleBindingUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateRoleBindingForbidden creates a CreateRoleBindingForbidden with default headers values
func NewCreateRoleBindingForbidden() *CreateRoleBindingForbidden {
	return &CreateRoleBindingForbidden{}
}

/*CreateRoleBindingForbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type CreateRoleBindingForbidden struct {
}

func (o *CreateRoleBindingForbidden) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings][%d] createRoleBindingForbidden ", 403)
}

func (o *CreateRoleBindingForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateRoleBindingDefault creates a CreateRoleBindingDefault with default headers values
func NewCreateRoleBindingDefault(code int) *CreateRoleBindingDefault {
	return &CreateRoleBindingDefault{
		_statusCode: code,
	}
}

/*CreateRoleBindingDefault handles this case with default header values.

errorResponse
*/
type CreateRoleBindingDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the create role binding default response
func (o *CreateRoleBindingDefault) Code() int {
	return o._statusCode
}

func (o *CreateRoleBindingDefault) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings][%d] createRoleBinding default  %+v", o._statusCode, o.Payload)
}

func (o *CreateRoleBindingDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
