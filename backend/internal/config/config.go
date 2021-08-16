package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var Config, RootDir = NewConfig()

const DefaultEnv = "beta"

type Configuration struct {
	Bucket string `yaml:"bucket"`
}

func NewConfig() (config *Configuration, rootDir string) {
	var c *Configuration = &Configuration{}
	
	var env = os.Getenv("ENV")
	if env == "" {
		log.Printf("ENV is not set. use default (%s)", DefaultEnv)
		env = DefaultEnv
	}
	log.Printf("%s environment is set", env)

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
	return c, rootDir
}
