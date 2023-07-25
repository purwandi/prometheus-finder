package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

type FinderMode string

const (
	ModeFull  FinderMode = "full_search"
	ModeFirst FinderMode = "first_of_line"
	ModeLast  FinderMode = "last_of_line"
)

func (e *FinderMode) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FinderMode(s)
	case string:
		*e = FinderMode(s)
	default:
		return fmt.Errorf("unsupported scan type for FinderMode: %T", src)
	}
	return nil
}

type Config struct {
	Name   string            `yaml:"metric_name"`
	Path   string            `yaml:"path"`
	Query  string            `yaml:"query"`
	Mode   FinderMode        `yaml:"mode"`
	Labels prometheus.Labels `yaml:"labels"`
}

func NewConfig(configPath string) (*[]Config, error) {
	configs := &[]Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
