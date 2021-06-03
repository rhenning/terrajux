package terrajux

import (
	"fmt"
	"os/user"
	"path/filepath"
)

var (
	Name       = "terrajux"
	Version    = "dev"
	ProjectURL = "https://github.com/rhenning/terrajux"
)

type Config struct {
	Name       string
	Version    string
	ProjectURL string

	DataDir    string
	CacheDir   string
	CacheClear bool
	DiffTool   string

	GitURL     string
	GitRefV1   string
	GitRefV2   string
	GitSubpath string
}

func NewConfig() *Config {
	config := &Config{
		Name:       Name,
		Version:    Version,
		ProjectURL: ProjectURL,
	}

	config.setDefaultDataDir()
	config.setDefaultCacheDir()

	return config
}

func (c *Config) setDefaultDataDir() {
	u, err := user.Current()

	if err != nil {
		panic("couldn't lookup current user x.x")
	}

	c.DataDir = filepath.Join(u.HomeDir, fmt.Sprintf(".%s", Name))
}

func (c *Config) setDefaultCacheDir() {
	c.CacheDir = filepath.Join(c.DataDir, "cache")
}
