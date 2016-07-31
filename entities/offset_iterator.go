package entities

type offsetIterator struct {
	offset   int
	inner    queryIterator
	consumed int
}

func newOffsetIterator(offset int, inner queryIterator) *offsetIterator {
	i := offsetIterator{offset: offset, inner: inner}
	return &i
}

func (i *offsetIterator) Next() (string, error) {
	for i.consumed < i.offset {
		i.inner.Next()
		i.consumed++
	}
	return i.inner.Next()
}
