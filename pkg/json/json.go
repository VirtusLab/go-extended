package json

import (
	"encoding/json"
	"io"
)

// ToInterface unmarshalls JSON reader to a "generic" interface
func ToInterface(reader io.Reader) (interface{}, error) {
	var result interface{}
	decoder := json.NewDecoder(reader)
	for {
		err := decoder.Decode(&result)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if result == nil {
		return map[string]interface{}{}, nil
	}
	return result, nil
}
