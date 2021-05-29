package cache

import "os"

type Cache struct {
	Dir string
}

func New(dirpath string) *Cache {
	return &Cache{
		Dir: dirpath,
	}
}

func (c *Cache) Ensure() error {
	return os.MkdirAll(c.Dir, 0755)
}

func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.Dir); err != nil {
		return err
	}

	return c.Ensure()
}
