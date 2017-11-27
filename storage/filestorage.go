package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
)

// FileKeyValueStore is a key-value store implemented in a directory
// the keys are names of files and values are written in the files
type FileKeyValueStore struct {
	dataDir string
	tmpDir  string
	n       int // number of files (k/v pairs)
}

func (s *FileKeyValueStore) getPathInTmpDir(key string) string {
	return path.Join(s.tmpDir, key)
}

// EditResource opens an existing file and writes the new value to the file
func (s *FileKeyValueStore) EditResource(t Tuple) error {
	t.Lock()
	defer t.Unlock()

	path := t.GetKey()
	value := t.GetValue()

	return ioutil.WriteFile(s.getPathInTmpDir(path), []byte(value), 0644)
}

// CreateResource creates a new file with the appropriate value
func (s *FileKeyValueStore) CreateResource(t Tuple) error {
	path := t.GetKey()
	value := t.GetValue()

	path = s.tmpDir + "/" + path
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
	return os.Remove(s.getPathInTmpDir(path))
}

// ReadResource reads the value of the file at the given path
func (s *FileKeyValueStore) ReadResource(path string) (string, error) {
	val, err := ioutil.ReadFile(s.getPathInTmpDir(path))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// CommitWrite TODO
func (s *FileKeyValueStore) CommitWrite(t Tuple) error {
	return nil
}

// NewFileKeyValueStore creates a directory at the given path and
// returns a *FileKeyValueStore configured to write at that path
func NewFileKeyValueStore(dataDir, tmpDir string) (*FileKeyValueStore, error) {

	logrus.Infof("creating data directory at %v", dataDir)
	err := os.MkdirAll(dataDir, 0777)
	if err != nil {
		logrus.Info(err)
	}

	logrus.Infof("creating temporary directory at %v", tmpDir)
	err = os.MkdirAll(tmpDir, 0777)
	if err != nil {
		logrus.Info(err)
	}
	return &FileKeyValueStore{
		dataDir,
		tmpDir,
		0,
	}, nil
}
