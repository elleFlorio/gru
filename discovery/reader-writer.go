package discovery

import (
	"encoding/json"
)

func WriteData(path string, src interface{}) error {
	var err error
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = service().Set(path, string(data), Options{})
	if err != nil {
		return err
	}

	return nil
}

func ReadData(path string, dest interface{}) error {
	var err error
	resp, err := service().Get(path, Options{})
	if err != nil {
		return err
	}
	data := resp[path]
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return err
	}

	return nil
}
