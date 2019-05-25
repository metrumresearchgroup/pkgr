package configlib

import (
	"gopkg.in/yaml.v3"
)

// Format provides a standard format for yaml
func Format(b []byte) ([]byte, error) {
	cfg := yaml.Node{}
	err := yaml.Unmarshal([]byte(b), &cfg)
	if err != nil {
		return b, err
	}

	d, err := yaml.Marshal(&cfg)
	if err != nil {
		return b, err
	}
	return d, nil
}
