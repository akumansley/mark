package entities

type kindIterator struct {
	kind string
	db   *DB
	iter Iterator
}

func newKindIterator(db *DB, kind string) queryIterator {
	i := kindIterator{db: db, kind: kind}
	i.init()
	return &i
}

func (i *kindIterator) init() error {
	prefix := NewKey("ave", "db/Kind", i.kind)
	iter, err := i.db.store.Prefix(prefix.ToBytes())
	i.iter = iter
	return err
}

func (i *kindIterator) Next() (string, error) {
	_, v, err := i.iter.Next()
	return string(v), err
}
