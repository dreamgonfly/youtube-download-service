package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var Config, RootDir, Env = NewConfig()

const DefaultEnv = "beta"

type Configuration struct {
	Bucket    string `yaml:"bucket"`
	LogServer string `yaml:"log_server"`
}

func NewConfig() (config *Configuration, rootDir string, env string) {
	var c *Configuration = &Configuration{}

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
		log.Fatalf("could not read config yaml file: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
	}
	return c, rootDir, env
}
