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
	t.Log(config.Get("redis.addr").AsStr())
	t.Log(config.Get("redis.max_conn").AsInt())
}

func TestAnyvalueCheck(t *testing.T) {

	config, err := LoadConfigYaml("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config.Get("listen").IsArray())
	t.Log(config.Get("listen").IsMap())
	t.Log(config.Get("listen").IsBool())
	t.Log(config.Get("listen").IsStr())
	t.Log(config.Get("listen").IsNumber())

	t.Log(config.Get("redis").IsMap())
	t.Log(config.Get("redis.max_conn").IsNumber())
}

func TestAnyvalueCheck2(t *testing.T) {

	config, err := LoadConfigJson("./config.json")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config.Get("listen").IsArray())
	t.Log(config.Get("listen").IsMap())
	t.Log(config.Get("listen").IsBool())
	t.Log(config.Get("listen").IsStr())
	t.Log(config.Get("listen").IsNumber())

	t.Log(config.Get("redis").IsMap())
	t.Log(config.Get("redis.max_conn").IsNumber())
}

func TestAnyvalueJson(t *testing.T) {

	config, err := LoadConfigJson("./config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config.Get("listen").AsStr())
	t.Log(config.Get("redis.addr").AsStr())
	t.Log(config.Get("redis.max_conn").AsInt())
}

func TestExist(t *testing.T) {
	config, err := LoadConfigJson("./config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config.Exist("listen"))
	t.Log(config.Exist("redis.max_conn"))
	t.Log(config.Exist("redis.max_conn1"))

}

func TestSet(t *testing.T) {
	av := New().Set("a", "hello").Set("b", 100).Set("c", "haha").
		Set("data.name", "starjiang").Set("data.age", 100)
	out, err := av.EncodeJson()

	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}

func TestMsgPack(t *testing.T) {
	av := New().Set("a", "hello").Set("b", 100).Set("c", "haha").
		Set("data.name", "starjiang").Set("data.age", 100)
	out, err := av.EncodeMsgPack()

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("out=%v", string(out))
	av, err = NewFromMsgPack(out)
	out, err = av.EncodeJson()
	t.Logf("out=%v", string(out))
}
