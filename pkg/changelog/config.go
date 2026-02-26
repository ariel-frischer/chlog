package changelog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile   = ".chlog.yaml"
	DefaultPublicFile   = "CHANGELOG.md"
	DefaultInternalFile = "CHANGELOG-internal.md"
)

// Config holds project-specific chlog settings.
type Config struct {
	RepoURL         string `yaml:"repo_url,omitempty"`
	IncludeInternal bool   `yaml:"include_internal,omitempty"`
	PublicFile      string `yaml:"public_file,omitempty"`
	InternalFile    string `yaml:"internal_file,omitempty"`
}

// PublicFilePath returns PublicFile if set, otherwise the default.
func (c *Config) PublicFilePath() string {
	if c.PublicFile != "" {
		return c.PublicFile
	}
	return DefaultPublicFile
}

// InternalFilePath returns InternalFile if set, otherwise the default.
func (c *Config) InternalFilePath() string {
	if c.InternalFile != "" {
		return c.InternalFile
	}
	return DefaultInternalFile
}

// LoadConfig reads a config file from the given path.
// Returns an empty config (not an error) if the file doesn't exist.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

// ResolveRepoURL returns the repo URL from config, falling back to git remote detection.
func ResolveRepoURL(cfg *Config) string {
	if cfg != nil && cfg.RepoURL != "" {
		return cfg.RepoURL
	}
	url, err := DetectRepoURL()
	if err != nil {
		return ""
	}
	return url
}
