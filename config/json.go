package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigAdapterJson struct {
	Input
}

func (adapter ConfigAdapterJson) ReadConfigFromFile (c *Config, name string) []byte {
	file, err := os.Open(name)

	if err != nil {
		panic("There ain't file")
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &c)

	defer file.Close()
	return byteValue
}

func (adapter ConfigAdapterJson) WriteConfigToFile (name string) interface{} {
	if _, err := os.Open(name); err == nil {
		fmt.Println("File's already exists")
		return nil
	}

	adapter.callUserInput()

	file, _ := json.MarshalIndent(adapter.Conf, "", " ")

	_ = ioutil.WriteFile(name, file, 0777)

	return "create config"
}
