package entities

type filterIterator struct {
	f     *filter
	inner queryIterator
	db    *DB
}

func newFilterIterator(f *filter, db *DB, inner queryIterator) *filterIterator {
	i := filterIterator{f: f, db: db, inner: inner}
	return &i
}

func (i *filterIterator) Next() (string, error) {
	var eid string
	var err error
	for eid, err = i.inner.Next(); err == nil && !i.match(eid); eid, err = i.inner.Next() {
	}
	return eid, err
}

// match checks the predicate for a given entity id
// Currently assumes the predicate is =
func (i *filterIterator) match(eid string) bool {
	k := NewKey("eav", eid, i.f.Attribute)
	v, err := i.db.store.Get(k.ToBytes())
	if err != nil {
		return false
	}
	return string(v) == i.f.Value.(string) // TODO :(
}
