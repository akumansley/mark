package entities

import "strings"

const (
	Eq = "="
)

const (
	Ascending  = iota
	Descending = iota
)

type Filter struct {
	Attribute string
	Predicate string
	Value     interface{}
}

type Order struct {
	Attribute string
	Direction int
}

type Query struct {
	db      *DB
	filters []Filter
	order   []Order
	limit   int
	offset  int
}

func (q *Query) Filter(spec string, val interface{}) *Query {
	parts := strings.Split(spec, " ")
	attr, pred := parts[0], parts[1]
	q.filters = append(q.filters, Filter{Attribute: attr, Predicate: pred, Value: val})
	return q
}

func (q *Query) Order(spec string) *Query {
	direction := Ascending
	if strings.HasPrefix(spec, "-") {
		spec = spec[:len(spec)-1]
		direction = Descending
	}
	q.order = append(q.order, Order{Attribute: spec, Direction: direction})
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// TODO
func (q *Query) GetAll(dst interface{}) error {
	return nil
}
