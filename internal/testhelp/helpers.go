package testhelp

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	d, err := os.Open(filepath.Clean(dirpath))

	if err != nil {
		return false, err
	}

	defer func() {
		if derr := d.Close(); derr != nil {
			err = derr
			fmt.Printf("%+v\n", err)
		}
	}()

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

	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
			fmt.Printf("%+v\n", err)
		}
	}()

	_, err = f.WriteString(content)

	return err
}
