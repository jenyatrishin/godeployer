package reader

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"dep2go/tools"
)

func ReadConfigFromFile (c interface{}, name string) interface{} {
	file, err := os.Open(name)

	if err != nil {
		fmt.Println(tools.UserError("There ain't file"))
		tools.WriteLog("There ain't file")
		os.Exit(1)
	}

	fmt.Println("Config file is opened")

	byteValue, _ := ioutil.ReadAll(file)

	xml.Unmarshal(byteValue, c)

	defer file.Close()

	return true
}