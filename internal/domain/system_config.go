package domain

type SystemConfig struct {
	Base
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ValueType   string `json:"value_type"`
	Module      string `json:"module"`
	Description string `json:"description"`
}
