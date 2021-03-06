// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PacketDatacenterSpec PacketDatacenterSpec specifies a datacenter of Packet.
//
// swagger:model PacketDatacenterSpec
type PacketDatacenterSpec struct {

	// facilities
	Facilities []string `json:"facilities"`
}

// Validate validates this packet datacenter spec
func (m *PacketDatacenterSpec) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PacketDatacenterSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PacketDatacenterSpec) UnmarshalBinary(b []byte) error {
	var res PacketDatacenterSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
