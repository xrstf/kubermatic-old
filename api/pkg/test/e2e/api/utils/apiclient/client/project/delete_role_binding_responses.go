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

// DeleteRoleBindingReader is a Reader for the DeleteRoleBinding structure.
type DeleteRoleBindingReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteRoleBindingReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDeleteRoleBindingOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewDeleteRoleBindingUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 403:
		result := NewDeleteRoleBindingForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewDeleteRoleBindingDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteRoleBindingOK creates a DeleteRoleBindingOK with default headers values
func NewDeleteRoleBindingOK() *DeleteRoleBindingOK {
	return &DeleteRoleBindingOK{}
}

/*DeleteRoleBindingOK handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteRoleBindingOK struct {
}

func (o *DeleteRoleBindingOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}][%d] deleteRoleBindingOK ", 200)
}

func (o *DeleteRoleBindingOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteRoleBindingUnauthorized creates a DeleteRoleBindingUnauthorized with default headers values
func NewDeleteRoleBindingUnauthorized() *DeleteRoleBindingUnauthorized {
	return &DeleteRoleBindingUnauthorized{}
}

/*DeleteRoleBindingUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteRoleBindingUnauthorized struct {
}

func (o *DeleteRoleBindingUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}][%d] deleteRoleBindingUnauthorized ", 401)
}

func (o *DeleteRoleBindingUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteRoleBindingForbidden creates a DeleteRoleBindingForbidden with default headers values
func NewDeleteRoleBindingForbidden() *DeleteRoleBindingForbidden {
	return &DeleteRoleBindingForbidden{}
}

/*DeleteRoleBindingForbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteRoleBindingForbidden struct {
}

func (o *DeleteRoleBindingForbidden) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}][%d] deleteRoleBindingForbidden ", 403)
}

func (o *DeleteRoleBindingForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteRoleBindingDefault creates a DeleteRoleBindingDefault with default headers values
func NewDeleteRoleBindingDefault(code int) *DeleteRoleBindingDefault {
	return &DeleteRoleBindingDefault{
		_statusCode: code,
	}
}

/*DeleteRoleBindingDefault handles this case with default header values.

errorResponse
*/
type DeleteRoleBindingDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete role binding default response
func (o *DeleteRoleBindingDefault) Code() int {
	return o._statusCode
}

func (o *DeleteRoleBindingDefault) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/roles/{namespace}/{role_id}/bindings/{binding_id}][%d] deleteRoleBinding default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteRoleBindingDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
