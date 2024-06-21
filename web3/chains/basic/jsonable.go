package basic

import "encoding/json"

type Jsonable interface {
	JsonString() (*OptionalString, error)
}

func JsonString(o interface{}) (*OptionalString, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return &OptionalString{Value: string(data)}, nil
}

func FromJsonString(jsonStr string, out interface{}) error {
	bytes := []byte(jsonStr)
	return json.Unmarshal(bytes, out)
}
