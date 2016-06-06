package mark

import (
	"encoding/json"
)

// Op is an arbitrary operation
type Op struct {
	Op      string
	Body    interface{}
	Version int
}

// // UnmarshalJSON builds an op from json
// func (op *Op) UnmarshalJSON(b []byte) error {
//   err := json.Unmarshal(b, &op)
//   if err != nil {
//     return err
//   }
//
//   // TODO look this up in a registry
//   if op.Op == "eav" {
//
//   }
//
// }

// Feed is a sequence of operations
type Feed struct {
	Ops []Op `json:"ops"`
}

// FromBytes inflates a Feed object from binary
func FromBytes(bytes []byte) (*Feed, error) {
	var feed Feed
	if len(bytes) == 0 {
		var ops []Op
		return &Feed{Ops: ops}, nil
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
