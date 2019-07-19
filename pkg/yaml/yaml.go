package yaml

import (
	"gopkg.in/yaml.v3"
	"io"
)

// ToInterface unmarshalls YAML reader to a "generic" interface
func ToInterface(reader io.Reader) (interface{}, error) {
	var values []interface{}
	decoder := yaml.NewDecoder(reader)
	for {
		var value map[string]interface{}
		err := decoder.Decode(&value)
		if err == io.EOF {
			break
		} else if err != nil {
			return []interface{}{}, err
		}
		values = append(values, value)
	}
	var result interface{}
	if len(values) == 0 {
		// empty
	} else if len(values) == 1 {
		result = values[0].(map[string]interface{})
	} else {
		result = values
	}
	return result, nil
}
