package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Credentail  string `yaml:"credential"`
	LogFile     string `yaml:"logfile"`
	SpreadSheet struct {
		Id        string `yaml:"id"`
		ReadRange string `yaml:"readrange"`
	} `yaml:"spreadsheet"`
	Calendar struct {
		Id    string `yaml:"id"`
		Event string `yaml:"event"`
	} `yaml:"calendar"`
	Target struct {
		URL    string `yaml:"url"`
		ApiKey string `yaml:"apikey"`
	} `yaml:"target"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

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

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")

	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}
