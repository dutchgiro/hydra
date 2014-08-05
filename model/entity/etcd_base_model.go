package entity

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
)

type EtcdModelizer interface {
	ExportEtcdOperations() map[string]string
}

type EtcdBaseModel map[string]interface{}

type EtcdBaseModels []EtcdBaseModel

func NewModelFromEvent(event *store.Event) (*EtcdBaseModel, error) {
	model := make(map[string]interface{})
	if err := proccessStruct(event, model); err != nil {
		return nil, err
	}
	m := EtcdBaseModel(model)
	return &m, nil
}

func NewModelsFromEvent(event *store.Event) (*EtcdBaseModels, error) {
	models := make([]EtcdBaseModel, 0)
	nodes := []*store.NodeExtern(reflect.ValueOf(event).Elem().FieldByName("Node").Elem().FieldByName("Nodes").Interface().(store.NodeExterns))
	for _, node := range nodes {
		model := make(map[string]interface{})
		if err := proccessStruct(node, model); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	m := EtcdBaseModels(models)
	return &m, nil
}

func ExtractJsonKeyFromEtcdKey(s string) (string, error) {
	lastIndex := strings.LastIndex(s, "/")
	if lastIndex == -1 {
		// TODO: improve error
		return "", errors.New("Bad etcd key")
	}
	key := s[lastIndex+1:]
	if len(key) == 0 {
		// TODO: improve error
		return "", errors.New("Bad etcd key")
	}
	return key, nil
}

func proccessStruct(s interface{}, m map[string]interface{}) error {
	if node := reflect.ValueOf(s).Elem().FieldByName("Node"); node.IsValid() {
		proccessStruct(node.Interface(), m)
	} else if exists := reflect.ValueOf(s).Elem().FieldByName("Nodes"); exists.IsValid() && !exists.IsNil() {
		key, _ := ExtractJsonKeyFromEtcdKey(reflect.ValueOf(s).Elem().FieldByName("Key").Interface().(string))
		nodes := []*store.NodeExtern(reflect.ValueOf(s).Elem().FieldByName("Nodes").Interface().(store.NodeExterns))
		m[key] = make(map[string]interface{})
		for _, node := range nodes {
			proccessStruct(node, m[key].(map[string]interface{}))
		}
	} else {
		value, err := CastInterfaceToString(reflect.ValueOf(s).Elem().FieldByName("Value").Interface())
		if err != nil {
			return err
		}
		key, _ := ExtractJsonKeyFromEtcdKey(reflect.ValueOf(s).Elem().FieldByName("Key").Interface().(string))
		m[key] = value
	}
	return nil
}

func CastInterfaceToString(v interface{}) (string, error) {
	var str string
	switch v.(type) {
	case nil:
		str = ""
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = strconv.Itoa(v.(int))
	case float32, float64:
		str = strconv.FormatFloat(v.(float64), 'f', 2, 64)
	case bool:
		str = strconv.FormatBool(v.(bool))
	case string:
		str = v.(string)
	default:
		// TODO: improve error
		return "", errors.New("Bad interface")
	}
	return str, nil
}

func (e EtcdBaseModel) ExportEtcdOperations() (map[string]string, error) {
	var operations map[string]string
	operations = make(map[string]string)

	var processInterface func(interface{}, string) error
	var processMap func(map[string]interface{}, string) error
	var processSlice func([]interface{}, string) error

	processInterface = func(in interface{}, key string) error {
		switch reflect.ValueOf(in).Kind() {
		case reflect.Map:
			processMap(in.(map[string]interface{}), key)
		case reflect.Slice:
			processSlice(in.([]interface{}), key)
		default:
			valueString, err := CastInterfaceToString(in)
			if err != nil {
				return err
			}
			operations[key] = valueString
		}
		return nil
	}

	processSlice = func(s []interface{}, parentKey string) error {
		for key, value := range s {
			if err := processInterface(value, parentKey+"/"+strconv.Itoa(key)); err != nil {
				return err
			}
		}
		return nil
	}

	processMap = func(mp map[string]interface{}, parentKey string) error {
		for key, value := range mp {
			if err := processInterface(value, parentKey+"/"+key); err != nil {
				return err
			}
		}
		return nil
	}

	// Process entity
	if err := processMap(e, ""); err != nil {
		return nil, err
	}

	return operations, nil
}

func CastStringToInterface(s string) (interface{}, error) {
	var i interface{}
	var err error

	i, err = strconv.Atoi(s)
	if err == nil {
		return i, nil
	}
	i, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return i, nil
	}
	i, err = strconv.ParseBool(s)
	if err == nil {
		return i, nil
	}
	if s == "" {
		return nil, nil
	}
	return s, nil
}

func (e *EtcdBaseModel) ToJsonableMap() (map[string]interface{}, error) {
	var checkMapKeysAsInteger func(map[string]interface{}) bool
	var processSlice func(map[string]interface{}) ([]interface{}, error)
	var processInterface func(interface{}) (interface{}, error)
	var processMap func(map[string]interface{}) (map[string]interface{}, error)

	checkMapKeysAsInteger = func(mp map[string]interface{}) bool {
		for key, _ := range mp {
			_, err := strconv.Atoi(key)
			if err != nil {
				return false
			}
		}
		return true
	}

	processSlice = func(mapToProcess map[string]interface{}) ([]interface{}, error) {
		finalSliceSize := len(mapToProcess)
		finalSlice := make([]interface{}, finalSliceSize, finalSliceSize)
		for key, value := range mapToProcess {
			v, err := processInterface(value)
			if err != nil {
				return nil, err
			}
			i, _ := strconv.Atoi(key)
			finalSlice[i] = v
		}
		return finalSlice, nil
	}

	processInterface = func(in interface{}) (interface{}, error) {
		var intrfac interface{}
		var err error

		switch reflect.ValueOf(in).Kind() {
		case reflect.Map:
			if checkMapKeysAsInteger(in.(map[string]interface{})) {
				intrfac, err = processSlice(in.(map[string]interface{}))
			} else {
				intrfac, err = processMap(in.(map[string]interface{}))
			}
		default:
			intrfac, err = CastStringToInterface(in.(string))
		}
		return intrfac, err
	}

	processMap = func(mapToProcess map[string]interface{}) (map[string]interface{}, error) {
		finalMap := make(map[string]interface{})
		for key, value := range mapToProcess {
			v, err := processInterface(value)
			if err != nil {
				return nil, err
			}
			finalMap[key] = v
		}
		return finalMap, nil
	}

	eP := map[string]interface{}(*e)
	jsonMap, err := processMap(eP)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func (e *EtcdBaseModel) Explode() (string, map[string]interface{}, error) {
	eP := map[string]interface{}(*e)
	if len(eP) != 1 {
		return "", nil, errors.New("This is not a valid Model, it contains more than one super key")
	}
	for key, value := range eP {
		return key, value.(map[string]interface{}), nil
	}
	return "", nil, nil
}

func (e *EtcdBaseModel) ToJson() ([]byte, error) {
	mp, err := e.ToJsonableMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(mp)
}

func (e *EtcdBaseModels) ToJson() ([]byte, error) {
	sliceSize := len([]EtcdBaseModel(*e))
	s := make([]map[string]interface{}, sliceSize, sliceSize)
	i := 0
	for _, value := range []EtcdBaseModel(*e) {
		v, err := value.ToJsonableMap()
		if err != nil {
			return nil, err
		}
		s[i] = v
		i++
	}
	return json.Marshal(s)
}
