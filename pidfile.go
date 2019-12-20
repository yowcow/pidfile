package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	syscall "golang.org/x/sys/unix"
)

// PIDFile is a state
type PIDFile struct {
	fd   int
	path string
}

// Create creates a pidfile at given path and locks
func Create(path string) (*PIDFile, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return nil, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if err := writeFile(path, ""); err != nil {
			return nil, err
		}
	} else if stat.IsDir() {
		return nil, fmt.Errorf("path %s is a dir", path)
	}

	fd, err := syscall.Open(path, syscall.O_RDONLY, 0000)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return nil, err
	}

	if err := writeFile(path, strconv.Itoa(os.Getpid())); err != nil {
		return nil, err
	}

	return &PIDFile{
		fd:   fd,
		path: path,
	}, nil
}

func writeFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), os.FileMode(0644))
}

// Remove unlocks a pidfile and removes
func (p *PIDFile) Remove() error {
	if err := syscall.Flock(p.fd, syscall.LOCK_UN); err != nil {
		return err
	}

	if err := syscall.Close(p.fd); err != nil {
		return err
	}

	if err := syscall.Unlink(p.path); err != nil {
		return err
	}

	return nil
}
