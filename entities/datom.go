package entities

import "fmt"

// Datom is an entity-attribute-value statement
type Datom struct {
	EntityID  string
	Attribute string
	Value     interface{}
	Added     bool
}

// EAVKey returns the key in the EAV index for this datom
func (d *Datom) EAVKey() []byte {
	return NewKey("eav", d.EntityID, d.Attribute).ToBytes()
}

// AEVKey returns the key in the AEV index for this datom
func (d *Datom) AEVKey() []byte {
	return NewKey("aev", d.Attribute, d.EntityID).ToBytes()
}

// AVEKey returns the key in the AVE index for this datom
// Has to include the entity ID for uniqueness
func (d *Datom) AVEKey() []byte {
	return NewKey("ave", d.Attribute, fmt.Sprint(d.Value), d.EntityID).ToBytes()
}

// VAEKey returns the key in the VAE index for this datom
// Has to include the entity ID for uniqueness
func (d *Datom) VAEKey() []byte {
	return NewKey("vae", fmt.Sprint(d.Value), d.Attribute, d.EntityID).ToBytes()
}
