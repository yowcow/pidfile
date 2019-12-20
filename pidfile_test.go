package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	syscall "golang.org/x/sys/unix"
)

func TestCreateRemove(t *testing.T) {
	dir, err := ioutil.TempDir("", "pidfile-test-")
	if err != nil {
		t.Fatal("failed creating a tmpdir:", err)
	}

	path := filepath.Join(dir, "test.pid")
	fmt.Println("pidfile at", path)

	// Create
	p, err := Create(path)
	if err != nil {
		t.Fatal("expected no error but got", err)
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		t.Error("expected", path, "to exist but does not")
	}

	r, err := os.Open(path)
	if err != nil {
		t.Error("expected no error but got", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error("expected no error but got", err)
	}

	expectedPid := syscall.Getpid()

	var actualPid int
	fmt.Sscanf(string(b), "%d", &actualPid)

	if expectedPid != actualPid {
		t.Error("expected pid", expectedPid, "but got", actualPid)
	}

	// Remove
	err = p.Remove()
	if err != nil {
		t.Error("expected no error but got", err)
	}

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		t.Error("expected", path, "NOT to exist but it does")
	}
}
