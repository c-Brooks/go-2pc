package persistence

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

// FileKeyValueStore is a key-value store implemented in a directory
// the keys are names of files and values are written in the files
type FileKeyValueStore struct {
	basePath string // path to base directory
	n        int    // number of files (k/v pairs)
}

func (s *FileKeyValueStore) getPath(key string) string {
	return path.Join(s.basePath, key)
}

// EditResource opens an existing file and writes the new value to the file
func (s *FileKeyValueStore) EditResource(t Tuple) error {
	path := t.getKey()
	value := t.getValue()

	return ioutil.WriteFile(s.getPath(path), []byte(value), 0644)
}

// CreateResource creates a new file with the appropriate value
func (s *FileKeyValueStore) CreateResource(t Tuple) error {
	path := t.getKey()
	value := t.getValue()

	path = s.basePath + "/" + path
	logrus.Infof("creating resource: %s: %s", path, value)
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logrus.Warnf("could not create file at %s: %v", path, err)
		return err
	}

	_, err = f.WriteString(value)
	return err
}

// DeleteResource removes a file located at the given path
func (s *FileKeyValueStore) DeleteResource(path string) error {
	return os.Remove(s.getPath(path))
}

// ReadResource reads the value of the file at the given path
func (s *FileKeyValueStore) ReadResource(path string) (string, error) {
	val, err := ioutil.ReadFile(s.getPath(path))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// NewFileKeyValueStore creates a directory at the given path and
// returns a *FileKeyValueStore configured to write at that path
func NewFileKeyValueStore(path string) (*FileKeyValueStore, error) {

	logrus.Infof("creating db at %v", path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		logrus.Info(err)
	}

	return &FileKeyValueStore{path, 0}, nil
}
