package config

type ConfigAdapter interface {
	ReadConfigFromFile(c *Config, name string) []byte
	WriteConfigToFile(name string) interface{}
}
