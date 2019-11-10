package config

import (
	"../adapter"
	"encoding/xml"
	"reflect"
)

type Config struct {
	//deprecated - gonna be removed
	ConfigFileType string
	XMLName xml.Name `xml:"config"`
	Envs []Env `xml:"environment"`
	Version string `xml:"version"`
}

func (c *Config) ReadConfig (name string, adapterIns adapter.ConfigAdapter) *Config {
	x := adapterIns.ReadConfigFromFile(c, name)
	xml.Unmarshal(x, &c)
	return c
}

func (c *Config) WriteConfig (name string, adapterIns adapter.ConfigAdapter) {
	//adapterIns := GetAdapterByType(c.ConfigFileType)
	adapterIns.WriteConfigToFile(name)
}

func (c *Config) GetEnvs () []Env {
	return c.Envs
}

func (c *Config) GetEnvByType (envType string) Env {
	envs := c.Envs
	var out Env

	for _, envIns := range envs {
		if envIns.EnvType == envType {
			out = envIns
			break
		}
	}

	return out
}

type Env struct {
	XMLName xml.Name `xml:"environment"`
	EnvType string `xml:"type,attr"`
	Server string `xml:"server"`
	Login string `xml:"login"`
	Password string `xml:"password"`
	HomeDir string `xml:"homeDir"`
	BeforeDeploy []Command `xml:"beforeDeploy"`
	AfterDeploy []Command `xml:"afterDeploy"`
}

func (e *Env) SetParam (name string, value string) *Env {

	return e
}

func (e *Env) GetParam (name string) interface{} {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(name)
	return f
}

type Command struct {
	Item string `xml:"command"`
}