package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/awans/mark/feed"
	"github.com/nu7hatch/gouuid"
	"github.com/square/go-jose"
)

// DB is the access point to the entity DB
type DB struct {
	store Store
	fp    []byte
	c     *feed.Coder
}

// ConvertDatoms implements Converter
func ConvertDatoms(bytes []byte) (interface{}, error) {
	var ds []Datom
	err := json.Unmarshal(bytes, &ds)
	return ds, err
}

// ConvertJWK implements Converter
func ConvertJWK(bytes []byte) (interface{}, error) {
	var jwk jose.JsonWebKey
	err := json.Unmarshal(bytes, &jwk)
	return &jwk, err
}

// NewDB is a constructor for a db
func NewDB(store Store, fp []byte) *DB {
	c := feed.NewCoder()
	c.RegisterOp("eav", ConvertDatoms)
	c.RegisterOp("declare-key", ConvertJWK)

	return &DB{store: store, fp: fp, c: c}
}

// Close closes the db
func (db *DB) Close() {
	db.store.Close()
}

// RebuildIndexes deletes all keys in the eav indexes and then loads each feed
func (db *DB) RebuildIndexes() error {
	for _, index := range []string{"eav", "aev", "ave", "vae"} {
		p, err := db.store.Prefix(NewKey(index).ToBytes())
		if err != nil {
			return err
		}
		for k, _, err := p.Next(); err == nil; k, _, err = p.Next() {
			db.store.Delete(k)
		}
	}
	fi, err := db.store.Prefix(NewKey("user").ToBytes())
	if err != nil {
		return err
	}
	for _, v, err := fi.Next(); err == nil; _, v, err = fi.Next() {
		feed, err := db.c.Decode(v)
		if err != nil {
			return err
		}
		db.LoadFeed(feed)
	}
	return nil
}

// LoadFeed applies each op to the db in turn and saves it under the user/feed key
func (db *DB) LoadFeed(feed *feed.Feed) {
	for _, op := range feed.Ops {
		db.applyOp(op)
	}
}

func (db *DB) applyOp(op feed.Op) {
	if op.Op != "eav" {
		return
	}
	datoms := op.Body.([]Datom)
	for _, datom := range datoms {
		db.applyDatom(datom)
	}
}

// UserFeed loads the feed for the user in this session
func (db *DB) UserFeed() (*feed.Feed, error) {
	feedK := Key{path: [][]byte{[]byte("user"), db.fp, []byte("feed")}}
	feedBytes, err := db.store.Get(feedK.ToBytes())
	if err != nil {
		return nil, err
	}
	return db.c.Decode(feedBytes)
}

// PutFeed sets a feed in the db
func (db *DB) PutFeed(f *feed.Feed) error {
	fp, err := f.Fingerprint()
	if err != nil {
		return err
	}
	feedBytes, err := db.c.Encode(f)
	if err != nil {
		return err
	}
	feedK := Key{path: [][]byte{[]byte("user"), fp, []byte("feed")}}
	db.store.Set(feedK.ToBytes(), feedBytes)
	return nil
}

func (db *DB) applyDatom(d Datom) {
	// eav, aev, ave, vae
	// we probably don't need all of these..
	if d.Added {
		// TODO value shouldn't always have to be a string
		// maybe store it as gob?
		db.store.Set(d.EAVKey(), []byte(fmt.Sprint(d.Value)))
		db.store.Set(d.AEVKey(), []byte(fmt.Sprint(d.Value)))
		db.store.Set(d.AVEKey(), []byte(d.EntityID))
		db.store.Set(d.VAEKey(), []byte(d.EntityID))
	} else {
		db.store.Delete(d.EAVKey())
		db.store.Delete(d.AEVKey())
		db.store.Delete(d.AVEKey())
		db.store.Delete(d.VAEKey())
	}
}

func getKindFromSlicePtr(slice interface{}) string {
	return reflect.ValueOf(slice).Elem().Type().Elem().Name()
}

func getKindFromInstance(instance interface{}) string {
	return reflect.ValueOf(instance).Type().Name()
}

// GetAll returns all entities of a given type
func (db *DB) GetAll(dst interface{}) error {
	kind := getKindFromSlicePtr(dst)
	prefix := NewKey("ave", "db/kind", kind)
	i, err := db.store.Prefix(prefix.ToBytes())
	if err != nil {
		return err
	}

	var entityIDs []string
	for _, v, err := i.Next(); err == nil; _, v, err = i.Next() {
		entityIDs = append(entityIDs, string(v))
	}
	db.GetMulti(entityIDs, dst)
	return nil
}

// Get returns a single entity by id
func (db *DB) Get(id string, dst interface{}) error {
	prefix := NewKey("eav", id)
	i, err := db.store.Prefix(prefix.ToBytes())
	if err != nil {
		return err
	}

	entityType := reflect.ValueOf(dst).Elem().Type()
	entity := reflect.New(entityType).Interface()

	for k, v, err := i.Next(); err == nil; k, v, err = i.Next() {
		components := bytes.Split(k, Separator)
		// eav/123/user/name = Andrew
		attr := string(components[3])
		field := reflect.ValueOf(entity).Elem().FieldByName(attr)
		sv := string(v)
		if field.IsValid() {
			if attr == "ID" {
				field.Set(reflect.ValueOf(id))
			} else {
				field.Set(reflect.ValueOf(sv))
			}
		}
	}
	reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(entity).Elem())
	return nil
}

// GetMulti fetches many keys
// dst is a pointer to a slice
func (db *DB) GetMulti(ids []string, dst interface{}) {
	v := reflect.ValueOf(dst).Elem() // v is a Value(sliceInstance)
	entityType := v.Type().Elem()    // v is a V(sliceInstance)->T(sliceType)->T(inner type)

	for _, id := range ids {
		entity := reflect.New(entityType).Interface()
		db.Get(id, entity)
		v.Set(reflect.Append(v, reflect.ValueOf(entity).Elem()))
	}
}

func eavOp(datoms []Datom) feed.Op {
	op := feed.Op{Op: "eav", Body: datoms}
	return op
}

// Put sets src at id
// TODO load it first and store the delta
func (db *DB) Put(id string, src interface{}) error {
	kind := getKindFromInstance(src)
	c := reflect.ValueOf(src)
	cType := c.Type()

	var datoms []Datom
	d := Datom{
		EntityID:  id,
		Attribute: "db/kind",
		Value:     kind,
		Added:     true,
	}
	datoms = append(datoms, d)

	for i := 0; i < cType.NumField(); i++ {
		valueField := c.Field(i)
		typeField := cType.Field(i)

		attrName := kind + "/" + typeField.Name

		d = Datom{
			EntityID:  id,
			Attribute: attrName,
			Value:     valueField.Interface(),
			Added:     true,
		}
		datoms = append(datoms, d)
	}

	feed, err := db.UserFeed()
	if err != nil {
		return err
	}

	op := eavOp(datoms)
	feed.Append(op)

	err = db.PutFeed(feed)
	if err != nil {
		return err
	}
	db.applyOp(op)
	return nil
}

// Add adds a new entity to the db
func (db *DB) Add(src interface{}) error {
	u, err := uuid.NewV4()
	if err != nil {
		return err
	}
	id := u.String()
	db.Put(id, src)
	return nil
}