package utils

import "encoding/json"

func StringToStruct[T any](data string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(data), &result)

	return result, err
}

func StructToString[T any](data T) string {
	bytes, err := json.Marshal(data)
	if err != nil {

		return ""
	}

	return string(bytes)
}
