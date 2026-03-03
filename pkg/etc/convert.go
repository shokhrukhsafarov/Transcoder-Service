package etc

import "encoding/json"

func ConvertStructToStruct[T any, U any](from T) (U, error) {
	var to U

	jsonData, err := json.Marshal(from)
	if err != nil {
		return to, err
	}

	err = json.Unmarshal(jsonData, &to)
	if err != nil {
		return to, err
	}

	return to, nil
}
