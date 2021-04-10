package settings

import "go-common/utils"

type Settings struct {
	Debug bool   `json:"debug"`
	Host  string `json:"host"`
	Cert  string `json:"cert"`
	Key   string `json:"key"`
}

func LoadSettings(filename string) (*Settings, error) {
	settings := &Settings{
		Debug: true,
		Host:  "0.0.0.0:80",
		Cert:  "",
		Key:   "",
	}

	if err := utils.LoadSettings(settings, filename); err != nil {
		return nil, err
	}

	return settings, nil
}
