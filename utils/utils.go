package utils

import "encoding/json"

func StructToMap(v any) (map[string]any, error) {
    data, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }

    var result map[string]any
    err = json.Unmarshal(data, &result)
    return result, err
}
