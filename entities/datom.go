package entities

import (
	"encoding/json"
	"fmt"
)

// Datom is an entity-attribute-value statement
type Datom struct {
	EntityID  string
	Attribute string
	Value     interface{}
	Added     bool
}

func (d *Datom) MarshalJSON() ([]byte, error) {
	var ary []interface{}
	ary = append(ary, d.EntityID)
	ary = append(ary, d.Attribute)
	ary = append(ary, d.Value)
	ary = append(ary, d.Added)
	return json.Marshal(ary)
}

func (d *Datom) UnmarshalJSON(data []byte) error {
	var ary []interface{}
	err := json.Unmarshal(data, &ary)
	if err != nil {
		return err
	}
	d.EntityID = ary[0].(string)
	d.Attribute = ary[1].(string)

	switch ary[2].(type) {
	case float64:
		v := ary[2].(float64)
		d.Value = int(v)
	case string:
		d.Value = ary[2]
	case bool:
		d.Value = ary[2]
	case nil:
		d.Value = nil
	default:
		panic("Invalid value")
	}
	d.Added = ary[3].(bool)
	return nil
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
	return NewKey("ave", d.Attribute, fmt.Sprintf("%v", d.Value), d.EntityID).ToBytes()
}

// VAEKey returns the key in the VAE index for this datom
// Has to include the entity ID for uniqueness
func (d *Datom) VAEKey() []byte {
	return NewKey("vae", fmt.Sprintf("%v", d.Value), d.Attribute, d.EntityID).ToBytes()
}
