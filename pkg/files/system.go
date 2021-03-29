package files

import (
	"crypto/md5"
	"io"
	"os"
)

const (
	passwd = "/etc/passwd"
	group  = "/etc/group"
	shadow = "/etc/shadow"
)

var systemFiles = []string{
	passwd,
	group,
	shadow,
}

// hashFile hashes the file at the given path.
func hashFile(path string) ([]byte, error) {
	h := md5.New()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// ProcessSystemFiles
func ProcessSystemFiles() error {
	for _, sf := range systemFiles {
		b, err := hashFile(sf)
		if err != nil {
			return err
		}
		_ = b
	}

	return nil
}
