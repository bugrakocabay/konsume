package util

import "encoding/json"

// ParseJSONToMap parses the JSON message and returns the map
func ParseJSONToMap(msg []byte) (map[string]interface{}, error) {
	var messageData map[string]interface{}
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		return nil, err
	}
	return messageData, nil
}
