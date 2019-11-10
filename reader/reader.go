package reader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadConfigFromFile (c interface{}, name string) interface{} {
	file, err := os.Open(name)

	if err != nil {
		panic("There ain't file")
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	xml.Unmarshal(byteValue, c)

	defer file.Close()

	return true
}