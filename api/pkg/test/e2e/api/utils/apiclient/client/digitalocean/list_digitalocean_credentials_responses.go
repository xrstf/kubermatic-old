// Code generated by go-swagger; DO NOT EDIT.

package digitalocean

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/models"
)

// ListDigitaloceanCredentialsReader is a Reader for the ListDigitaloceanCredentials structure.
type ListDigitaloceanCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListDigitaloceanCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListDigitaloceanCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewListDigitaloceanCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListDigitaloceanCredentialsOK creates a ListDigitaloceanCredentialsOK with default headers values
func NewListDigitaloceanCredentialsOK() *ListDigitaloceanCredentialsOK {
	return &ListDigitaloceanCredentialsOK{}
}

/*ListDigitaloceanCredentialsOK handles this case with default header values.

CredentialList
*/
type ListDigitaloceanCredentialsOK struct {
	Payload *models.CredentialList
}

func (o *ListDigitaloceanCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/providers/digitalocean/credentials][%d] listDigitaloceanCredentialsOK  %+v", 200, o.Payload)
}

func (o *ListDigitaloceanCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CredentialList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListDigitaloceanCredentialsDefault creates a ListDigitaloceanCredentialsDefault with default headers values
func NewListDigitaloceanCredentialsDefault(code int) *ListDigitaloceanCredentialsDefault {
	return &ListDigitaloceanCredentialsDefault{
		_statusCode: code,
	}
}

/*ListDigitaloceanCredentialsDefault handles this case with default header values.

ErrorResponse is the default representation of an error
*/
type ListDigitaloceanCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorDetails
}

// Code gets the status code for the list digitalocean credentials default response
func (o *ListDigitaloceanCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *ListDigitaloceanCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/providers/digitalocean/credentials][%d] listDigitaloceanCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ListDigitaloceanCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorDetails)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
