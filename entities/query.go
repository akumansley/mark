package entities

import "strings"

// Predicates
const (
	Eq = "="
)

// Sort directions
const (
	Ascending  = iota
	Descending = iota
)

type filter struct {
	Attribute string
	Predicate string
	Value     interface{}
}

type order struct {
	Attribute string
	Direction int
}

// Query is a db query
type Query struct {
	db      *DB
	filters []filter
	order   []order
	kind    string
	limit   int
	offset  int
}

// Filter adds a filter to the query
func (q *Query) Filter(spec string, val interface{}) *Query {
	parts := strings.Split(spec, " ")
	attr, pred := parts[0], parts[1]
	if attr == "FeedID" || attr == "ID" {
		attr = "db/" + attr
	} else {
		attr = q.kind + "/" + attr
	}
	q.filters = append(q.filters, filter{Attribute: attr, Predicate: pred, Value: val})
	return q
}

// Order adds a sort order to the query
func (q *Query) Order(spec string) *Query {
	direction := Ascending
	if strings.HasPrefix(spec, "-") {
		spec = spec[1:]
		direction = Descending
	}
	attr := q.kind + "/" + spec
	q.order = append(q.order, order{Attribute: attr, Direction: direction})
	return q
}

// Limit caps the number of rows
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Offset skips offset rows
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// GetAll returns the results of the query
func (q *Query) GetAll(dst interface{}) error {
	i := newKindIterator(q.db, q.kind)

	for _, f := range q.filters {
		i = newFilterIterator(&f, q.db, i)
	}

	for _, o := range q.order {
		i = newOrderIterator(&o, q.db, i)
	}

	if q.offset != -1 {
		i = newOffsetIterator(q.offset, i)
	}

	if q.limit != -1 {
		i = newLimitIterator(q.limit, i)
	}

	var eids []string
	for eid, err := i.Next(); err == nil; eid, err = i.Next() {
		eids = append(eids, eid)
	}
	return q.db.GetMulti(eids, dst)
}
