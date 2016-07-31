package entities

import "io"

type limitIterator struct {
	limit    int
	inner    queryIterator
	returned int
}

func newLimitIterator(limit int, inner queryIterator) *limitIterator {
	i := limitIterator{limit: limit, inner: inner}
	return &i
}

func (i *limitIterator) Next() (string, error) {
	if i.returned == i.limit {
		return "", io.EOF
	}
	i.returned++
	return i.inner.Next()
}
