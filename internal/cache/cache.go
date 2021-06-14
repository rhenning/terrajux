package cache

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Manager interface {
	Ensure() error
	Clear() error
	HasKey(key string) bool
	GetAbsKeyPath(key string) (abspath string, err error)
}

type Cache struct {
	Dir string
}

func New(dirpath string) *Cache {
	return &Cache{
		Dir: dirpath,
	}
}

func (c *Cache) Ensure() error {
	return os.MkdirAll(c.Dir, 0750)
}

func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.Dir); err != nil {
		return err
	}

	return c.Ensure()
}

func (c *Cache) HasKey(k string) bool {
	if _, err := os.Stat(path.Join(c.Dir, k)); os.IsNotExist(err) {
		return false
	}

	return true
}

func (c *Cache) GetAbsKeyPath(k string) (string, error) {
	var err error

	if !c.HasKey(k) {
		err = fmt.Errorf("cache key %q not found", k)
	}

	return filepath.Join(c.Dir, k), err
}
