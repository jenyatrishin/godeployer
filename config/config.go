package config

import (
	"encoding/xml"
	"reflect"
)

const (
	XML string = "xml"
	JSON string = "json"
	AUTH_PASSWORD string  = "password"
	AUTH_KEY string = "key"
)

type Config struct {
	XMLName xml.Name `xml:"config" json:"config"`
	Envs []Env `xml:"environment" json:"environment"`
	ProjectName string `xml:"project_name" json:"projectName"`
//	Version string `xml:"version"`
}

func (c *Config) ReadConfig(name string, ext string) *Config {
	adapterIns := getAdapterByType(ext)
	adapterIns.ReadConfigFromFile(c, name)

	return c
}

func (c *Config) WriteConfig(name string, ext string) {
	adapterIns := getAdapterByType(ext)
	adapterIns.WriteConfigToFile(name)
}

func (c *Config) GetEnvs() []Env {
	return c.Envs
}

func (c *Config) GetEnvByType(envType string) Env {
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

func (c *Config) ValidateConfig(name string, ext string) bool {
	adapterIns := getAdapterByType(ext)
	adapterIns.ReadConfigFromFile(c, name)
	result := true
	allowed := [3]string{"developer","staging","production"}

	if len(c.GetEnvs()) < 1 || c.ProjectName == "" {
		result = false
	}

	for _, envIns := range c.GetEnvs() {
		if envInArray(envIns.EnvType, allowed) == false || envIns.Server == "" ||
			envIns.Login == "" || envIns.HomeDir == "" {
			result = false
			break
		}
		if envIns.AuthType == AUTH_PASSWORD && string(envIns.Password) == "" ||
			envIns.AuthType == AUTH_KEY && envIns.KeyFile == "" {
			result = false
			break
		}
		if envIns.GitConfig.Repository == "" || envIns.GitConfig.User == "" ||
			string(envIns.GitConfig.Password) == "" || envIns.GitConfig.Branch == "" {
			result = false
			break
		}
	}

	return result
}

type Env struct {
	XMLName xml.Name `xml:"environment" json:"environment"`
	EnvType string `xml:"type,attr" json:"type"`
	Server string `xml:"server" json:"server"`
	Login string `xml:"login" json:"login"`
	AuthType string `xml:"auth_type" json:"authType"`
	Password []byte `xml:"password" json:"password"`
	KeyFile string `xml:"key" json:"key"`
	HomeDir string `xml:"homeDir" json:"homeDir"`
	BeforeDeploy []Command `xml:"beforeDeploy" json:"beforeDeploy"`
	AfterDeploy []Command `xml:"afterDeploy" json:"afterDeploy"`
	GitConfig GitConfig `xml:"git" json:"git"`
}

//deprecated
//gonna be removed
func (e *Env) SetParam(name string, value string) *Env {
	return e
}

//deprecated
//gonna be removed
func (e *Env) GetParam(name string) interface{} {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(name)
	return f
}

type Command struct {
	Item string `xml:"command" json:"command"`
}

type GitConfig struct {
	Repository string `xml:"repository" json:"repository"`
	User string `xml:"user" json:"user"`
	Password []byte `xml:"password" json:"password"`
	Branch string `xml:"branch" json:"branch"`
}

func getAdapterByType(ext string) ConfigAdapter {
	if ext == XML {
		return ConfigAdapterXml{}
	} else if ext == JSON {
		return ConfigAdapterJson{}
	}

	return ConfigAdapterXml{}
}

func envInArray(env string, envs [3]string) bool {
	output := false
	for _, v := range envs {
		if env == v {
			output = true
			break;
		}
	}

	return output
}