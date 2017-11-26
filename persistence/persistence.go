package persistence

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

// Tuple is a single key-value pair
type Tuple struct {
	key   string
	value string
}

// NewTuple returns a newly created tuple
func NewTuple(key, value string) Tuple {
	return Tuple{key, value}
}

func (t *Tuple) getKey() string {
	return t.key
}

func (t *Tuple) getValue() string {
	return t.value
}

// KeyValueReadWriter is a driver for a read/write key-value store.
// The methods are very strict by design because I want them to break easily.
type KeyValueReadWriter interface {
	CreateResource(Tuple) error
	ReadResource(string) (string, error)
	EditResource(Tuple) error
	DeleteResource(string) error
}

// CommitState describes the state of a distributed commit
type CommitState int

const (
	initial CommitState = 1 + iota
	pending
	applied
	failed
)

// NewKeyValueStore returns an instance of a k/v store
// driven by whatever is in options["driver"]
// available values: filesystem
func NewKeyValueStore(driver string, options map[string]string) (KeyValueReadWriter, error) {
	var filepath string
	var KVRW KeyValueReadWriter
	var err error

	switch driver {
	case "file":
		filepath = options["filepath"]
		if filepath == "" {
			return nil, errors.New("no filepath provided")
		}
		KVRW, err = NewFileKeyValueStore(filepath)
		if err != nil {
			return nil, err
		}
	default:
		logrus.Errorf("unknown driver: %s", driver)
	}

	return KVRW, nil
}

// TwoPhaseCommitResource implements functions necessary to
// perform 2-phase commits as a Resource (Slave)
type TwoPhaseCommitResource interface {
	ProcessXactRequest() error
}

// TwoPhaseCommitManager implements functions necessary to
// perform 2-phase commits as a transaction manager
type TwoPhaseCommitManager interface {
}

// XATransaction is a transaction that can be performed in a distributed system
type XATransaction struct {
	Xid   int
	State CommitState
}

// func (tpcr *TwoPhaseCommitResource) ProcessXactRequest() (bool, error) {
// 	// ProcessXactRequest
// 	// write to disk
// 	// return ok if transaction can go thru

// 	return false, nil
// }
