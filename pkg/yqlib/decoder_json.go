package yqlib

import (
	"encoding/json"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type jsonDecoder struct {
	decoder json.Decoder
}

func NewJSONDecoder() Decoder {
	return &jsonDecoder{}
}

func (dec *jsonDecoder) Init(reader io.Reader) {
	dec.decoder = *json.NewDecoder(reader)
}

func (dec *jsonDecoder) Decode(rootYamlNode *yaml.Node) error {

	var dataBucket interface{}
	log.Debug("going to decode")
	err := dec.decoder.Decode(&dataBucket)
	if err != nil {
		return err
	}
	node, err := dec.convertToYamlNode(dataBucket)

	if err != nil {
		return err
	}
	rootYamlNode.Kind = yaml.DocumentNode
	rootYamlNode.Content = []*yaml.Node{node}
	return nil
}

func (dec *jsonDecoder) convertToYamlNode(data interface{}) (*yaml.Node, error) {
	switch data.(type) {
	case float64, float32:
		// json decoder returns ints as float.
		return parseSnippet(fmt.Sprintf("%v", data))
	case int, int64, int32, string, bool, nil:
		return createScalarNode(data, fmt.Sprintf("%v", data)), nil
	case map[string]interface{}:
		return dec.parseMap(data.(map[string]interface{}))
	case []interface{}:
		return dec.parseArray(data.([]interface{}))
	default:
		return nil, fmt.Errorf("unrecognised type :(")
	}
}

func (dec *jsonDecoder) parseMap(dataMap map[string]interface{}) (*yaml.Node, error) {

	var yamlMap = &yaml.Node{Kind: yaml.MappingNode}

	for key, value := range dataMap {
		yamlValue, err := dec.convertToYamlNode(value)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, createScalarNode(key, fmt.Sprintf("%v", key)), yamlValue)
	}
	return yamlMap, nil
}

func (dec *jsonDecoder) parseArray(dataArray []interface{}) (*yaml.Node, error) {

	var yamlMap = &yaml.Node{Kind: yaml.SequenceNode}

	for _, value := range dataArray {
		yamlValue, err := dec.convertToYamlNode(value)
		if err != nil {
			return nil, err
		}
		yamlMap.Content = append(yamlMap.Content, yamlValue)
	}
	return yamlMap, nil
}
