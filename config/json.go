package config

type ConfigAdapterJson struct {
	Input
}

func (adapter ConfigAdapterJson) ReadConfigFromFile (c *Config, name string) []byte {

}

func (adapter ConfigAdapterJson) WriteConfigToFile (name string) interface{} {

}
