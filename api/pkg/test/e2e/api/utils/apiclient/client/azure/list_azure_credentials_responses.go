// Code generated by go-swagger; DO NOT EDIT.

package azure

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kubermatic/kubermatic/api/pkg/test/e2e/api/utils/apiclient/models"
)

// ListAzureCredentialsReader is a Reader for the ListAzureCredentials structure.
type ListAzureCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAzureCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListAzureCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewListAzureCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAzureCredentialsOK creates a ListAzureCredentialsOK with default headers values
func NewListAzureCredentialsOK() *ListAzureCredentialsOK {
	return &ListAzureCredentialsOK{}
}

/*ListAzureCredentialsOK handles this case with default header values.

CredentialList
*/
type ListAzureCredentialsOK struct {
	Payload *models.CredentialList
}

func (o *ListAzureCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/providers/azure/credentials][%d] listAzureCredentialsOK  %+v", 200, o.Payload)
}

func (o *ListAzureCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CredentialList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAzureCredentialsDefault creates a ListAzureCredentialsDefault with default headers values
func NewListAzureCredentialsDefault(code int) *ListAzureCredentialsDefault {
	return &ListAzureCredentialsDefault{
		_statusCode: code,
	}
}

/*ListAzureCredentialsDefault handles this case with default header values.

ErrorResponse is the default representation of an error
*/
type ListAzureCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorDetails
}

// Code gets the status code for the list azure credentials default response
func (o *ListAzureCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *ListAzureCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/providers/azure/credentials][%d] listAzureCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ListAzureCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorDetails)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
