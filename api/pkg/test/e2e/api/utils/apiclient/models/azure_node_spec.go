// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AzureNodeSpec AzureNodeSpec describes settings for an Azure node
//
// swagger:model AzureNodeSpec
type AzureNodeSpec struct {

	// should the machine have a publicly accessible IP address
	AssignPublicIP bool `json:"assignPublicIP,omitempty"`

	// Data disk size in GB
	DataDiskSize int32 `json:"dataDiskSize,omitempty"`

	// OS disk size in GB
	OSDiskSize int32 `json:"osDiskSize,omitempty"`

	// VM size
	// Required: true
	Size *string `json:"size"`

	// Additional metadata to set
	Tags map[string]string `json:"tags,omitempty"`
}

// Validate validates this azure node spec
func (m *AzureNodeSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSize(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AzureNodeSpec) validateSize(formats strfmt.Registry) error {

	if err := validate.Required("size", "body", m.Size); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AzureNodeSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AzureNodeSpec) UnmarshalBinary(b []byte) error {
	var res AzureNodeSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
