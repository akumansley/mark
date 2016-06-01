package mark

import (
	"crypto/rsa"
	"encoding/json"
	"reflect"
)

// Bookmark is a model class representing a bookmark
type Bookmark struct {
	URL  string
	Note string
}

// DB is the access point to the dece
type DB struct {
	store Store
}

// Close closes the db's underlying store
func (db *DB) Close() error {
	return db.store.Close()
}

// DBFromStore is a constructor for a db
func DBFromStore(store Store) *DB {
	db := new(DB)
	db.store = store
	return db
}

// Entity is an envelope for a model class
type Entity struct {
	ID   string
	Body interface{}
}

// EAV is an entity-attribute-value statement
type EAV struct {
	EntityID  string
	Attribute string
	Value     interface{}
	Added     bool
}

// MarshalJSON map an EAV to a JSON list
func (eav EAV) MarshalJSON() ([]byte, error) {
	ls := []interface{}{eav.EntityID, eav.Attribute, eav.Value, eav.Added}
	return json.Marshal(ls)
}

// ToEAV turns an entity into a series of EAV statements
func (entity *Entity) ToEAV() ([]EAV, error) {
	entityID := entity.ID

	elem := reflect.ValueOf(entity.Body).Elem()
	var eavs []EAV

	for i := 0; i < elem.NumField(); i++ {
		valueField := elem.Field(i)
		typeField := elem.Type().Field(i)

		eav := EAV{
			EntityID:  entityID,
			Attribute: typeField.Name,
			Value:     valueField.Interface(),
			Added:     true,
		}
		eavs = append(eavs, eav)
	}

	return eavs, nil
}

// Add adds an element to the DB under a given key's space
func (db *DB) Add(key *rsa.PrivateKey, entity Entity) error {
	feed, err := db.FeedForPubKey(&key.PublicKey)
	if err != nil {
		return err
	}

	eavs, err := entity.ToEAV()
	op := Op{Op: "eav", Version: 0, Body: eavs}

	feed.Append(op)

	bytes, err := json.Marshal(feed)
	if err != nil {
		return err
	}

	keyFingerprint, err := Fingerprint(&key.PublicKey)
	if err != nil {
		return err
	}

	err = db.store.Set(keyFingerprint, bytes)
	if err != nil {
		return err
	}

	return nil
}

// FeedForPubKey returns the feed for a given pub key
func (db *DB) FeedForPubKey(key *rsa.PublicKey) (*Feed, error) {
	keyFingerprint, err := Fingerprint(key)
	if err != nil {
		return nil, err
	}

	feedBytes, err := db.store.Get(keyFingerprint)
	if err != nil {
		return nil, err
	}
	return FromBytes(feedBytes)
}
