// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DatacenterSpec DatacenterSpec specifies the data for a datacenter.
//
// swagger:model DatacenterSpec
type DatacenterSpec struct {

	// country
	Country string `json:"country,omitempty"`

	// EnforceAuditLogging enforces audit logging on every cluster within the DC,
	// ignoring cluster-specific settings.
	EnforceAuditLogging bool `json:"enforceAuditLogging,omitempty"`

	// location
	Location string `json:"location,omitempty"`

	// provider
	Provider string `json:"provider,omitempty"`

	// Deprecated. Automatically migrated to the RequiredEmailDomains field.
	RequiredEmailDomain string `json:"requiredEmailDomain,omitempty"`

	// required email domains
	RequiredEmailDomains []string `json:"requiredEmailDomains"`

	// seed
	Seed string `json:"seed,omitempty"`

	// alibaba
	Alibaba *AlibabaDatacenterSpec `json:"alibaba,omitempty"`

	// aws
	Aws *AWSDatacenterSpec `json:"aws,omitempty"`

	// azure
	Azure *AzureDatacenterSpec `json:"azure,omitempty"`

	// bringyourown
	Bringyourown BringYourOwnDatacenterSpec `json:"bringyourown,omitempty"`

	// digitalocean
	Digitalocean *DigitialoceanDatacenterSpec `json:"digitalocean,omitempty"`

	// gcp
	Gcp *GCPDatacenterSpec `json:"gcp,omitempty"`

	// hetzner
	Hetzner *HetznerDatacenterSpec `json:"hetzner,omitempty"`

	// kubevirt
	Kubevirt KubevirtDatacenterSpec `json:"kubevirt,omitempty"`

	// openstack
	Openstack *OpenstackDatacenterSpec `json:"openstack,omitempty"`

	// packet
	Packet *PacketDatacenterSpec `json:"packet,omitempty"`

	// vsphere
	Vsphere *VSphereDatacenterSpec `json:"vsphere,omitempty"`
}

// Validate validates this datacenter spec
func (m *DatacenterSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlibaba(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAws(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAzure(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDigitalocean(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGcp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHetzner(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOpenstack(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePacket(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVsphere(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DatacenterSpec) validateAlibaba(formats strfmt.Registry) error {

	if swag.IsZero(m.Alibaba) { // not required
		return nil
	}

	if m.Alibaba != nil {
		if err := m.Alibaba.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("alibaba")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateAws(formats strfmt.Registry) error {

	if swag.IsZero(m.Aws) { // not required
		return nil
	}

	if m.Aws != nil {
		if err := m.Aws.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("aws")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateAzure(formats strfmt.Registry) error {

	if swag.IsZero(m.Azure) { // not required
		return nil
	}

	if m.Azure != nil {
		if err := m.Azure.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("azure")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateDigitalocean(formats strfmt.Registry) error {

	if swag.IsZero(m.Digitalocean) { // not required
		return nil
	}

	if m.Digitalocean != nil {
		if err := m.Digitalocean.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("digitalocean")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateGcp(formats strfmt.Registry) error {

	if swag.IsZero(m.Gcp) { // not required
		return nil
	}

	if m.Gcp != nil {
		if err := m.Gcp.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("gcp")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateHetzner(formats strfmt.Registry) error {

	if swag.IsZero(m.Hetzner) { // not required
		return nil
	}

	if m.Hetzner != nil {
		if err := m.Hetzner.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hetzner")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateOpenstack(formats strfmt.Registry) error {

	if swag.IsZero(m.Openstack) { // not required
		return nil
	}

	if m.Openstack != nil {
		if err := m.Openstack.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("openstack")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validatePacket(formats strfmt.Registry) error {

	if swag.IsZero(m.Packet) { // not required
		return nil
	}

	if m.Packet != nil {
		if err := m.Packet.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("packet")
			}
			return err
		}
	}

	return nil
}

func (m *DatacenterSpec) validateVsphere(formats strfmt.Registry) error {

	if swag.IsZero(m.Vsphere) { // not required
		return nil
	}

	if m.Vsphere != nil {
		if err := m.Vsphere.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("vsphere")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DatacenterSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DatacenterSpec) UnmarshalBinary(b []byte) error {
	var res DatacenterSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
