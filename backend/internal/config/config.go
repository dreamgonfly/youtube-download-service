package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

const DefaultEnv = "beta"

var Conf, RootDir, Env = NewConfig()

type Config struct {
	Bucket string `yaml:"bucket"`
}

func NewConfig() (conf *Config, rootDir string, env string) {
	var c *Config = &Config{}
	env = os.Getenv("ENV")
	if env == "" {
		log.Printf("ENV is not set. use default (%s)", DefaultEnv)
		env = DefaultEnv
	}
	var _, filename, _, _ = runtime.Caller(0) // 0 means this file itself
	rootDir = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	confpath := filepath.Join(rootDir, "configs", fmt.Sprintf("%s.yaml", env))
	yamlFile, err := os.ReadFile(confpath)
	if err != nil {
		log.Fatalf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c, rootDir, env
}
