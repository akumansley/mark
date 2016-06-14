package mark

import (
	"crypto/rsa"
	"encoding/json"
	"reflect"
)

// DB is the access point to the db
type DB struct {
	store        Store
	typeRegistry map[string]reflect.Type
}

// Entity is an envelope for a model class
type Entity struct {
	ID         string
	EntityType string
	Body       interface{}
}

// EAV is an entity-attribute-value statement
type EAV struct {
	EntityID  string
	Attribute string
	Value     interface{}
	Added     bool
}

// Register tells the DB about a given type of struct
func (db *DB) Register(nilOfType interface{}) {
	typeObject := reflect.ValueOf(nilOfType).Elem().Type()
	db.typeRegistry[typeObject.Name()] = typeObject
}

// Inflate turns a list of EAVs into a list of entities
// later EAVs (in the slice) override earlier ones
func (db *DB) Inflate(eavs []EAV) []Entity {
	entitiesMap := make(map[string]map[string]interface{})
	var entities []Entity

	for _, eav := range eavs {
		entry, exists := entitiesMap[eav.EntityID]
		if !exists {
			entitiesMap[eav.EntityID] = make(map[string]interface{})
			entry = entitiesMap[eav.EntityID]
		}
		if eav.Added {
			entry[eav.Attribute] = eav.Value
		} else {
			delete(entry, eav.Attribute)
		}
	}

	for i, entityMap := range entitiesMap {
		entityTypeName := entityMap["Type"].(string) // if not present, fall back to something
		entityType := db.typeRegistry[entityTypeName]
		entity := reflect.New(entityType).Interface()
		for k, v := range entityMap {
			field := reflect.ValueOf(entity).Elem().FieldByName(k)
			if field.IsValid() {
				field.Set(reflect.ValueOf(v))
			}
		}
		entities = append(entities, Entity{ID: i, Body: entity, EntityType: entityTypeName})
	}
	return entities
}

// GetAll returns all entities of a given type
func (db *DB) GetAll(key *rsa.PublicKey, dst interface{}) error {
	feed, err := db.FeedForPubKey(key)
	if err != nil {
		return err
	}
	var allEavs []EAV
	for _, op := range feed.Ops {
		if op.Op == "eav" {
			var eavs []EAV
			err = json.Unmarshal(op.RawBody, &eavs)
			if err != nil {
				return nil
			}
			allEavs = append(allEavs, eavs...)
		}
	}

	entities := db.Inflate(allEavs)

	dv := reflect.ValueOf(dst).Elem() // dv is a Value(sliceInstance)
	dstTypeName := dv.Type().Elem().Name()

	for _, entity := range entities {
		if entity.EntityType == dstTypeName {
			dv.Set(reflect.Append(dv, reflect.ValueOf(entity.Body).Elem()))
		}
	}

	return nil
}

// Close closes the db's underlying store
func (db *DB) Close() error {
	return db.store.Close()
}

// DBFromStore is a constructor for a db
func DBFromStore(store Store) *DB {
	db := new(DB)
	db.store = store
	db.typeRegistry = make(map[string]reflect.Type, 16)
	return db
}

// MarshalJSON maps an EAV to a JSON list
func (eav EAV) MarshalJSON() ([]byte, error) {
	ls := []interface{}{eav.EntityID, eav.Attribute, eav.Value, eav.Added}
	return json.Marshal(ls)
}

// UnmarshalJSON maps a JSON list to an EAV
func (eav *EAV) UnmarshalJSON(bytes []byte) error {
	var ls []interface{}
	err := json.Unmarshal(bytes, &ls)
	if err != nil {
		return err
	}
	eav.EntityID = ls[0].(string)
	eav.Attribute = ls[1].(string)
	eav.Value = ls[2]
	eav.Added = ls[3].(bool)
	return nil
}

// ToEAV turns an entity into a series of EAV statements
func (entity *Entity) ToEAV() ([]EAV, error) {
	entityID := entity.ID

	elem := reflect.ValueOf(entity.Body).Elem()
	var eavs []EAV

	typeEav := EAV{
		EntityID:  entityID,
		Attribute: "Type",
		Value:     elem.Type().Name(),
		Added:     true,
	}
	eavs = append(eavs, typeEav)

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
	op := Op{Op: "eav", Body: eavs}

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
	return FromBytes(key, feedBytes)
}
