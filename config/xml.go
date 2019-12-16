package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigAdapterXml struct {
	Input
}

func (adapter ConfigAdapterXml) ReadConfigFromFile (c *Config, name string) []byte {
	file, err := os.Open(name)

	if err != nil {
		panic("There ain't file")
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	xml.Unmarshal(byteValue, &c)

	defer file.Close()
	return byteValue
}

func (adapter ConfigAdapterXml) WriteConfigToFile (name string) interface{} {

	if _, err := os.Open(name); err == nil {
		fmt.Println("File's already exists")
		return nil
	}

	adapter.callUserInput()

	file, _ := xml.MarshalIndent(adapter.Conf, "", " ")

	_ = ioutil.WriteFile(name, file, 0777)

	return "create config"
}
