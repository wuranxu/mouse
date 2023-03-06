package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	MouseToken      = "mouse_token"
	MouseUser       = "mouse_user"
	MouseParameters = "mouse_parameters"
)

var Conf = new(Config)

type SqlConfig struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DBName   string `json:"dbName" yaml:"dbName"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Config struct {
	Database SqlConfig `json:"database"`
}

func ParseConfig(filepath string, cfg interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

func Init(filepath string) error {
	return ParseConfig(filepath, Conf)
}
