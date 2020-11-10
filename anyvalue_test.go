package anyvalue

import (
	"io/ioutil"
	"testing"
)

func LoadConfigYaml(config string) (*AnyValue, error) {
	dataBytes, err := ioutil.ReadFile(config)

	if err != nil {
		return nil, err
	}
	return NewFromYaml(dataBytes)
}

func LoadConfigJson(config string) (*AnyValue, error) {
	dataBytes, err := ioutil.ReadFile(config)

	if err != nil {
		return nil, err
	}
	return NewFromJson(dataBytes)
}

func TestAnyvalueYaml(t *testing.T) {

	config, err := LoadConfigYaml("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config.Get("listen").AsStr())
	t.Log(config.GetValue("redis.addr").AsStr())
	t.Log(config.GetValue("redis.max_conn").AsInt())
}

func TestAnyvalueJson(t *testing.T) {

	config, err := LoadConfigJson("./config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config.Get("listen").AsStr())
	t.Log(config.GetValue("redis.addr").AsStr())
	t.Log(config.GetValue("redis.max_conn").AsInt())
}
