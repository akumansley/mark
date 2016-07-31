package entities

import (
	"io"
	"sort"
)

type orderIterator struct {
	o          *order
	inner      queryIterator
	workingSet []string
	db         *DB
	returned   int
}

type sortPair struct {
	eid string
	val string
}
type sorts []sortPair

func (a sorts) Len() int           { return len(a) }
func (a sorts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sorts) Less(i, j int) bool { return a[i].val < a[j].val }

func newOrderIterator(o *order, db *DB, inner queryIterator) *orderIterator {
	i := orderIterator{o: o, db: db, inner: inner}
	return &i
}

func (i *orderIterator) Next() (string, error) {
	if i.returned == 0 {
		// TODO errs
		i.fill()
		i.sort()
	}
	if i.returned == len(i.workingSet) {
		return "", io.EOF
	}
	i.returned++
	return i.workingSet[i.returned-1], nil
}

func (i *orderIterator) fill() {
	// TODO assert err is EOF
	for eid, err := i.inner.Next(); err == nil; eid, err = i.inner.Next() {
		i.workingSet = append(i.workingSet, eid)
	}
}

func (i *orderIterator) sort() error {
	var s sorts
	for _, v := range i.workingSet {
		k := NewKey("eav", v, i.o.Attribute)
		val, err := i.db.store.Get(k.ToBytes())
		if err != nil {
			return err
		}
		s = append(s, sortPair{eid: v, val: string(val)})
	}
	if i.o.Direction == Descending {
		sort.Sort(sort.Reverse(s))
	} else {
		sort.Sort(s)
	}
	var eids []string

	for _, v := range s {
		eids = append(eids, v.eid)
	}

	i.workingSet = eids
	return nil
}
