package util

import (
	"encoding/json"
	"io"
	"net/http"
)

// ParseJSONToMap parses the JSON message and returns the map
func ParseJSONToMap(msg []byte) (map[string]interface{}, error) {
	var messageData map[string]interface{}
	err := json.Unmarshal(msg, &messageData)
	if err != nil {
		return nil, err
	}
	return messageData, nil
}

// ReadRequestBody reads the response body and returns it as a byte slice
func ReadRequestBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
