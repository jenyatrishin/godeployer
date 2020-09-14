package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"dep2go/tools"
)

type ConfigAdapterXml struct {
	Input
}

func (adapter ConfigAdapterXml) ReadConfigFromFile (c *Config, name string) []byte {
	file, err := os.Open(name)

	if err != nil {
		fmt.Println(tools.UserError("There ain't file"))
		tools.WriteLog("There ain't file")
		os.Exit(1)
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	xml.Unmarshal(byteValue, &c)

	defer file.Close()

	return byteValue
}

func (adapter ConfigAdapterXml) WriteConfigToFile (name string) interface{} {

	if _, err := os.Open(name); err == nil {
		fmt.Println("File already exists")
		return nil
	}

	adapter.callUserInput()

	file, _ := xml.MarshalIndent(adapter.Conf, "", " ")

	_ = ioutil.WriteFile(name, file, 0777)

	return "create config"
}
