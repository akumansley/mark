package mark

import (
  "encoding/json"
)

// Op is a serialized single operation
type Op struct {
	Op string
	Body interface{}
	Version int
}

// Feed is a sequence of operations
type Feed struct {
  Ops []Op  `json:"ops"`
}

// FromBytes inflates a Feed object from binary
func FromBytes(bytes []byte) (*Feed, error) {
  var feed Feed;
  if len(bytes) == 0 {
    var ops []Op
    return &Feed{Ops:ops}, nil
  }
  err := json.Unmarshal(bytes, &feed)
  if err != nil {
    return nil, err
  }
  return &feed, nil
}

// Append adds an Op to the end of a feed
func (feed *Feed) Append(op Op) {
  feed.Ops = append(feed.Ops, op)
}
