package testhelp

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func CreateTempDir(t *testing.T) string {
	tdir, err := ioutil.TempDir("", "terrajux-test")
	t.Logf("Created temp directory: %v", tdir)

	if err != nil {
		t.Fatalf("Error creating temp directory: %v", err)
	}

	return tdir
}

func DirIsEmpty(dirpath string) (bool, error) {
	d, err := os.Open(dirpath)
	if err != nil {
		return false, err
	}
	defer d.Close()

	_, err = d.Readdirnames(1)

	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func WriteFile(t *testing.T, fpath string, content string) error {
	t.Logf("Creating file: %v", fpath)
	f, err := os.Create(fpath)

	if err != nil {
		return err
	}

	if _, err = f.WriteString(content); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
