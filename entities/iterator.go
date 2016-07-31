package entities

// queryIterator is the interface for query execution
type queryIterator interface {
	Next() (string, error)
}
