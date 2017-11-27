package storage

import (
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
)

// Tuple is a single key-value pair.
// It's the only data that can be committed
// and includes a RWMutex to lock the tuple
type Tuple struct {
	*sync.RWMutex
	key   string
	value string
}

// NewTuple returns a newly created tuple
func NewTuple(key, value string) *Tuple {
	var m sync.RWMutex
	return &Tuple{&m, key, value}
}

// GetKey returns the tuple's key
func (t *Tuple) GetKey() string {
	return t.key
}

// GetValue returns the tuple's value
func (t *Tuple) GetValue() string {
	return t.value
}

// TwoPhaseReadWriter is a driver for a read/write key-value store.
// The methods are very strict by design because I want them to break easily.
type TwoPhaseReadWriter interface {
	ReadResource(string) (string, error)
	CreateResource(Tuple) error
	EditResource(Tuple) error
	DeleteResource(string) error

	CommitWrite(Tuple) error // TODO
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
// driven by whatever is in options["driver"].
// Available values: files
func NewKeyValueStore(driver string, options map[string]string) (TwoPhaseReadWriter, error) {
	var dataDir string
	var tmpDir string
	var tprw TwoPhaseReadWriter
	var err error

	switch driver {
	case "file":
		dataDir = options["dataDir"]
		if dataDir == "" {
			return nil, errors.New("no dataDir provided")
		}
		tmpDir = options["tmpDir"]
		if tmpDir == "" {
			return nil, errors.New("no tmpDir provided")
		}
		tprw, err = NewFileKeyValueStore(dataDir, tmpDir)
		if err != nil {
			return nil, err
		}

	default:
		logrus.Fatalf("unknown database driver: %s", driver)
	}

	return tprw, nil
}

// XATransaction is a transaction that can be performed in a distributed system
type XATransaction struct {
	Xid   int
	State CommitState
}
