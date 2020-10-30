package servicemanager

import "encoding/json"

type config struct {
	Services map[string]json.RawMessage `json:"services"`
}

func readConfig(raw []byte) (*config, error) {
	c := new(config)
	if raw == nil {
		c.Services = make(map[string]json.RawMessage)
		return c, nil
	}
	if err := json.Unmarshal(raw, c); err != nil {
		return nil, err
	}
	return c, nil
}
