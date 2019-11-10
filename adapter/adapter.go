package adapter

type ConfigAdapter interface {

  ReadConfigFromFile(c interface{}, name string) []byte
  WriteConfigToFile(name string) interface{}

}