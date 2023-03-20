package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Config struct {
	App `yaml:"app"`
}

type Env string

const (
	EnvProduction Env = "prod"
	EnvDev        Env = "dev"
	EnvStaging    Env = "staging"
)

type (
	App struct {
		// just app name
		Name string `yaml:"name"`
		// prod/dev/staging
		Env    Env    `yaml:"env" validate:"oneof=prod dev staging"`
		Listen Listen `yaml:"listen"`
		JWT    JWT    `yaml:"jwt"`

		// todo
		// Tracing *Tracing `yaml:"tracing"`
		// Metrics *Metrics `yaml:"metrics"`
		// Debug   *Debug   `yaml:"debug"`
		// Logging *Logging `yaml:"logging"`
	}

	// Listen configuration for address and port
	Listen struct {
		HTTP string `yaml:"http"`
	}
	JWT struct {
		Issuer     string        `yaml:"issuer"`
		PublicKey  string        `yaml:"public_key"`
		PrivateKey string        `yaml:"private_key"`
		ExpiresIn  time.Duration `yaml:"expires_in"`
	}
)

func ReadYaml(path string, target interface{}) error {
	confBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config from path=%s: %w", path, err)
	}
	if len(confBytes) == 0 {
		return fmt.Errorf("empty config file")
	}

	if err := UnmarshalYaml(confBytes, target); err != nil {
		return fmt.Errorf(`parse config: %w`, err)
	}

	validate := validator.New()

	err = validate.Struct(target)
	if err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	return nil
}

func UnmarshalYaml(content []byte, target interface{}) error {
	return yaml.UnmarshalStrict([]byte(os.ExpandEnv(string(content))), target)
}
