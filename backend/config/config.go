package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Bucket string `yaml:"bucket"`
}

func (c *Config) Conf() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "beta"
	}
	// 0 means this file itself
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename))
	confpath := filepath.Join(rootDir, "configs", fmt.Sprintf("%s.yaml", env))
	yamlFile, err := ioutil.ReadFile(confpath)
	if err != nil {
		log.Fatalf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

var c = &Config{}
var CONFIG = c.Conf()
