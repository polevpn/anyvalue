package anyvalue

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v2"
)

// returns the current implementation version
func Version() string {
	return "0.5.0"
}

var AVNil = AnyValue{nil}

type AnyValue struct {
	data interface{}
}

// Implements the json.Unmarshaler interface.
func (j *AnyValue) UnmarshalJSON(p []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(p))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

// Implements the json.Unmarshaler interface.
func (j *AnyValue) UnmarshalMsgPack(p []byte) error {

	dec := msgpack.NewDecoder(bytes.NewReader(p))
	dec.SetCustomStructTag("json")
	return dec.Decode(&j.data)
}

// Implements the yaml.Unmarshaler interface.
func (j *AnyValue) UnmarshalYAML(p []byte) error {
	dec := yaml.NewDecoder(bytes.NewBuffer(p))
	return dec.Decode(&j.data)
}

// NewFromReader returns a *AnyValue by decoding from an io.Reader
func NewFromJsonReader(r io.Reader) (*AnyValue, error) {
	j := new(AnyValue)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&j.data)
	return j, err
}

// NewFromReader returns a *AnyValue by decoding from an io.Reader
func NewFromMsgPackReader(r io.Reader) (*AnyValue, error) {
	j := new(AnyValue)
	dec := msgpack.NewDecoder(r)
	dec.SetCustomStructTag("json")
	err := dec.Decode(&j.data)
	return j, err
}

// NewFromReader returns a *AnyValue by decoding from an io.Reader
func NewFromYamlReader(r io.Reader) (*AnyValue, error) {
	j := new(AnyValue)
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&j.data)
	return j, err
}

// Float64 coerces into a float64
func (j *AnyValue) Float64() (float64, error) {
	switch j.data.(type) {
	case json.Number:
		return j.data.(json.Number).Float64()
	case float32, float64:
		return reflect.ValueOf(j.data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int coerces into an int
func (j *AnyValue) Int() (int, error) {
	switch j.data.(type) {
	case json.Number:
		i, err := j.data.(json.Number).Int64()
		return int(i), err
	case float32, float64:
		return int(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int64 coerces into an int64
func (j *AnyValue) Int64() (int64, error) {
	switch j.data.(type) {
	case json.Number:
		return j.data.(json.Number).Int64()
	case float32, float64:
		return int64(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(j.data).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Uint64 coerces into an uint64
func (j *AnyValue) Uint64() (uint64, error) {
	switch j.data.(type) {
	case json.Number:
		return strconv.ParseUint(j.data.(json.Number).String(), 10, 64)
	case float32, float64:
		return uint64(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(j.data).Uint(), nil
	}
	return 0, errors.New("invalid value type")
}

// NewFromJson returns a pointer to a new `AnyValue` object
// after unmarshaling `body` bytes
func NewFromJson(body []byte) (*AnyValue, error) {
	j := new(AnyValue)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// NewFromMsgPack returns a pointer to a new `AnyValue` object
// after unmarshaling `body` bytes
func NewFromMsgPack(body []byte) (*AnyValue, error) {
	j := new(AnyValue)
	err := j.UnmarshalMsgPack(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// NewFromYaml returns a pointer to a new `AnyValue` object
// after unmarshaling `body` bytes
func NewFromYaml(body []byte) (*AnyValue, error) {
	j := new(AnyValue)
	err := j.UnmarshalYAML(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// New returns a pointer to a new, empty `AnyValue` object
func New() *AnyValue {
	return &AnyValue{
		data: make(map[string]interface{}),
	}
}

// New returns a pointer to a new, empty `AnyValue` object
func NewFromInf(data interface{}) *AnyValue {
	return &AnyValue{
		data: data,
	}
}

// Interface returns the underlying data
func (j *AnyValue) Interface() interface{} {
	return j.data
}

// Encode returns its marshaled data as `[]byte`
func (j *AnyValue) EncodeJson() ([]byte, error) {
	return j.MarshalJSON()
}

// Encode returns its marshaled data as `[]byte`
func (j *AnyValue) EncodeMsgPack() ([]byte, error) {
	return j.MarshalMsgPack()
}

// EncodePretty returns its marshaled data as `[]byte` with indentation
func (j *AnyValue) EncodeJsonPretty() ([]byte, error) {
	return json.MarshalIndent(&j.data, "", "  ")
}

// Implements the json.Marshaler interface.
func (j *AnyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// Implements the msgpack.Marshaler interface.
func (j *AnyValue) MarshalMsgPack() ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetCustomStructTag("json")
	err := enc.Encode(&j.data)
	return buf.Bytes(), err
}

// Encode returns its marshaled data as `[]byte`
func (j *AnyValue) EncodeYaml() ([]byte, error) {
	return j.MarshalYAML()
}

// Implements the yaml.Marshaler interface.
func (j *AnyValue) MarshalYAML() ([]byte, error) {

	return yaml.Marshal(&j.data)
}

func (j *AnyValue) Set(path string, val interface{}) *AnyValue {
	branch := strings.Split(path, ".")
	return j.SetPath(branch, val)
}

// SetPath modifies `AnyValue`, recursively checking/creating map keys for the supplied path,
// and then finally writing in the value
func (j *AnyValue) SetPath(branch []string, val interface{}) *AnyValue {
	if len(branch) == 0 {
		j.data = val
		return j
	}

	// in order to insert our branch, we need map[string]interface{}
	if _, ok := (j.data).(map[string]interface{}); !ok {
		// have to replace with something suitable
		j.data = make(map[string]interface{})
	}
	curr := j.data.(map[string]interface{})

	for i := 0; i < len(branch)-1; i++ {
		b := branch[i]
		// key exists?
		if _, ok := curr[b]; !ok {
			n := make(map[string]interface{})
			curr[b] = n
			curr = n
			continue
		}

		// make sure the value is the right sort of thing
		if _, ok := curr[b].(map[string]interface{}); !ok {
			// have to replace with something suitable
			n := make(map[string]interface{})
			curr[b] = n
		}

		curr = curr[b].(map[string]interface{})
	}

	// add remaining k/v
	curr[branch[len(branch)-1]] = val
	return j
}

// Del modifies `AnyValue` map by deleting `key` if it is present.
func (j *AnyValue) Del(key string) {
	m, err := j.Map()
	if err != nil {
		return
	}
	delete(m, key)
}

// Get returns a pointer to a new `AnyValue` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *AnyValue) getValue(key string) *AnyValue {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &AnyValue{val}
		}
	}
	return &AVNil
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (j *AnyValue) GetPath(branch ...string) *AnyValue {
	jin := j
	for _, p := range branch {
		jin = jin.getValue(p)
	}
	return jin
}

// GetValue searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetValue("top_level.dict")
func (j *AnyValue) Get(path string) *AnyValue {
	branch := strings.Split(path, ".")
	jin := j
	for _, p := range branch {
		jin = jin.getValue(p)
	}
	return jin
}

// CheckGet returns a pointer to a new `AnyValue` object and
// a `bool` identifying success or failure
//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").Exist("inner"); ok {
//        log.Println(data)
//    }
func (j *AnyValue) Exist(path string) (*AnyValue, bool) {
	branch := strings.Split(path, ".")
	jin := j
	for _, p := range branch {

		jin = jin.getValue(p)
	}

	if jin != &AVNil {
		return jin, true
	} else {
		return jin, false
	}
}

// GetIndex returns a pointer to a new `AnyValue` object
// for `index` in its `array` representation
//
// this is the analog to Get when accessing elements of
// a json array instead of a json object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (j *AnyValue) GetIndex(index int) *AnyValue {
	a, err := j.Array()
	if err == nil {
		if len(a) > index {
			return &AnyValue{a[index]}
		}
	}
	return &AVNil
}

// Map type asserts to `map`
func (j *AnyValue) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	} else if m, ok := (j.data).(map[interface{}]interface{}); ok {
		var mn map[string]interface{}
		mn = make(map[string]interface{})
		for k, v := range m {
			kn, ok := k.(string)
			if ok {
				mn[kn] = v
			}
		}
		return mn, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// Array type asserts to an `array`
func (j *AnyValue) Array() ([]interface{}, error) {
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// Bool type asserts to `bool`
func (j *AnyValue) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (j *AnyValue) Str() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Bytes type asserts to `[]byte`
func (j *AnyValue) Bytes() ([]byte, error) {
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

// StrArr type asserts to an `array` of `string`
func (j *AnyValue) StrArr() ([]string, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, "")
			continue
		}
		s, ok := a.(string)
		if !ok {
			continue
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}

// IntArr type asserts to an `array` of `int`
func (j *AnyValue) Int64Arr() ([]int64, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]int64, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, 0)
			continue
		}
		var t int64 = 0
		switch a.(type) {
		case json.Number:
			t, _ = a.(json.Number).Int64()
		case float32, float64:
			t = int64(reflect.ValueOf(a).Float())
		case int, int8, int16, int32, int64:
			t = reflect.ValueOf(a).Int()
		case uint, uint8, uint16, uint32, uint64:
			t = int64(reflect.ValueOf(a).Uint())
		}
		retArr = append(retArr, t)
	}
	return retArr, nil
}

// IntArr type asserts to an `array` of `int`
func (j *AnyValue) UInt64Arr() ([]uint64, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]uint64, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, 0)
			continue
		}
		var t uint64 = 0
		switch a.(type) {
		case json.Number:
			t, _ = strconv.ParseUint(a.(json.Number).String(), 10, 64)
		case float32, float64:
			t = uint64(reflect.ValueOf(a).Float())
		case int, int8, int16, int32, int64:
			t = uint64(reflect.ValueOf(a).Int())
		case uint, uint8, uint16, uint32, uint64:
			t = reflect.ValueOf(a).Uint()
		}
		retArr = append(retArr, t)
	}
	return retArr, nil
}

// IntArr type asserts to an `array` of `int`
func (j *AnyValue) Float64Arr() ([]float64, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]float64, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, 0)
			continue
		}
		var t float64 = 0
		switch a.(type) {
		case json.Number:
			t, _ = a.(json.Number).Float64()
		case float32, float64:
			t = reflect.ValueOf(a).Float()
		case int, int8, int16, int32, int64:
			t = float64(reflect.ValueOf(a).Int())
		case uint, uint8, uint16, uint32, uint64:
			t = float64(reflect.ValueOf(a).Uint())
		}
		retArr = append(retArr, t)
	}
	return retArr, nil
}

// AsArray guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, v := range js.Get("results").AsArray() {
//			fmt.Println(i, v)
//		}
func (j *AnyValue) AsArray(args ...[]interface{}) []interface{} {
	var def []interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsArray() received too many arguments %d", len(args))
	}

	a, err := j.Array()
	if err == nil {
		return a
	}

	return def
}

// AsMap guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").AsMap() {
//			fmt.Println(k, v)
//		}
func (j *AnyValue) AsMap(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsMap() received too many arguments %d", len(args))
	}

	a, err := j.Map()
	if err == nil {
		return a
	}

	return def
}

// AsStr guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").AsStr(), js.Get("optional_param").AsStr("Asmy_ault"))
func (j *AnyValue) AsStr(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsString() received too many arguments %d", len(args))
	}

	s, err := j.Str()
	if err == nil {
		return s
	}

	return def
}

// AsStrArr guarantees the return of a `[]string` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, s := range js.Get("results").AsStrArr() {
//			fmt.Println(i, s)
//		}
func (j *AnyValue) AsStrArr(args ...[]string) []string {
	var def []string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsStrArr() received too many arguments %d", len(args))
	}

	a, err := j.StrArr()
	if err == nil {
		return a
	}

	return def
}

func (j *AnyValue) AsInt64Arr(args ...[]int64) []int64 {
	var def []int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsStrArr() received too many arguments %d", len(args))
	}

	a, err := j.Int64Arr()
	if err == nil {
		return a
	}

	return def
}

func (j *AnyValue) AsUInt64Arr(args ...[]uint64) []uint64 {
	var def []uint64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsStrArr() received too many arguments %d", len(args))
	}

	a, err := j.UInt64Arr()
	if err == nil {
		return a
	}

	return def
}

func (j *AnyValue) AsFloat64Arr(args ...[]float64) []float64 {
	var def []float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsStrArr() received too many arguments %d", len(args))
	}

	a, err := j.Float64Arr()
	if err == nil {
		return a
	}

	return def
}

// AsInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").AsInt(), js.Get("optional_param").AsInt(5150))
func (j *AnyValue) AsInt(args ...int) int {
	var def int

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsInt() received too many arguments %d", len(args))
	}

	i, err := j.Int()
	if err == nil {
		return i
	}

	return def
}

// AsFloat64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.Get("param1").AsFloat64(), js.Get("optional_param").AsFloat64(5.150))
func (j *AnyValue) AsFloat64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsFloat64() received too many arguments %d", len(args))
	}

	f, err := j.Float64()
	if err == nil {
		return f
	}

	return def
}

// AsBool guarantees the return of a `bool` (with optional default)
//
// useful when you explicitly want a `bool` in a single value return context:
//     myFunc(js.Get("param1").AsBool(), js.Get("optional_param").AsBool(true))
func (j *AnyValue) AsBool(args ...bool) bool {
	var def bool

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsBool() received too many arguments %d", len(args))
	}

	b, err := j.Bool()
	if err == nil {
		return b
	}

	return def
}

// AsInt64 guarantees the return of an `int64` (with optional default)
//
// useful when you explicitly want an `int64` in a single value return context:
//     myFunc(js.Get("param1").AsInt64(), js.Get("optional_param").AsInt64(5150))
func (j *AnyValue) AsInt64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsInt64() received too many arguments %d", len(args))
	}

	i, err := j.Int64()
	if err == nil {
		return i
	}

	return def
}

// AsUInt64 guarantees the return of an `uint64` (with optional default)
//
// useful when you explicitly want an `uint64` in a single value return context:
//     myFunc(js.Get("param1").AsUint64(), js.Get("optional_param").AsUint64(5150))
func (j *AnyValue) AsUint64(args ...uint64) uint64 {
	var def uint64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("AsUint64() received too many arguments %d", len(args))
	}

	i, err := j.Uint64()
	if err == nil {
		return i
	}

	return def
}
