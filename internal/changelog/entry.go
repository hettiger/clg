package changelog

import "go.yaml.in/yaml/v3"

type Entry struct {
	Type    string `yaml:"type"`
	Message string `yaml:"message"`
}

func (e Entry) YAMLData() ([]byte, error) {
	return yaml.Marshal(e)
}

func (e Entry) YAML() (string, error) {
	data, err := e.YAMLData()
	if err != nil {
		return "", err
	}

	return string(data), nil
}
