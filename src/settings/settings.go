package settings

import "gopkg.in/go-mixed/go-common.v1/conf.v1"

type Settings struct {
	Debug bool   `yaml:"debug"`
	Host  string `yaml:"host"`
	Cert  string `yaml:"cert"`
	Key   string `yaml:"key"`
}

func LoadSettings(filename string) (*Settings, error) {
	settings := &Settings{
		Debug: true,
		Host:  "0.0.0.0:80",
		Cert:  "",
		Key:   "",
	}

	if err := conf.LoadSettings(settings, filename); err != nil {
		return nil, err
	}

	return settings, nil
}
