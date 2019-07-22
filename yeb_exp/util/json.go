package util

import (
	"encoding/json"
)

type TJson map[string]interface{}

func (t *TJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}
func (t *TJson) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*map[string]interface{})(t))
	return err
}
