package files

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// BinaryStateStore contains the previously recorded
// state and the current binary state.
type BinaryStateStore struct {
	Previous *BinaryState
	Current  *BinaryState
}

// Compare compares the previous state to the current state and returns
// a slice of states found to be different.
func (b *BinaryStateStore) Compare() ([]BinaryStateStore, error) {
	if reflect.DeepEqual(b.Current, b.Previous) {
		return nil, nil
	}

	// not equal
	currentKeys := make([]string, len(b.Current.fileData))
	for k := range b.Current.fileData {
		currentKeys = append(currentKeys, k)
	}
	sort.Strings(currentKeys)

	previousKeys := make([]string, len(b.Previous.fileData))
	for k := range b.Previous.fileData {
		previousKeys = append(previousKeys, k)
	}
	sort.Strings(previousKeys)

	bsd := make([]BinaryStateStore, 0)
	for k := range b.Current.fileData {
		if diff := bytes.Compare(b.Current.fileData[k], b.Previous.fileData[k]); diff != 0 {
			bsd = append(bsd, BinaryStateStore{
				Previous: b.Previous,
				Current:  b.Current,
			})
		}
	}

	return bsd, nil
}

// BinaryState
type BinaryState struct {
	mu       *sync.Mutex
	fileData map[string][]byte
}

// NewBinaryState
func NewBinaryState() *BinaryState {
	return &BinaryState{
		mu:       &sync.Mutex{},
		fileData: make(map[string][]byte),
	}
}

// GenerateBinaryHashes
func (b *BinaryState) GenerateBinaryHashes() error {
	var wg sync.WaitGroup
	for _, p := range strings.Split(os.Getenv("PATH"), ":") {
		wg.Add(1)
		go func(bs *BinaryState, p string, wg *sync.WaitGroup) {
			defer wg.Done()

			if err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					if _, ok := err.(*os.PathError); ok {
						// suppressing the error
						return nil
					}
					return err
				}
				// ignore symlinks
				if info.Mode()&os.ModeSymlink == os.ModeSymlink || info.IsDir() {
					return nil
				}

				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()

				h := md5.New()
				if _, err := io.Copy(h, f); err != nil {
					return err
				}

				bs.mu.Lock()
				defer bs.mu.Unlock()

				if _, ok := bs.fileData[path]; !ok {
					bs.fileData[path] = h.Sum(nil)
				}

				return nil
			}); err != nil {
				logrus.Println(err)
			}
		}(b, p, &wg)
	}

	wg.Wait()

	return nil
}
