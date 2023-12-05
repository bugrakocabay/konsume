package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ProcessTemplate processes the template and returns the processed body
func ProcessTemplate(template map[string]interface{}, messageData map[string]interface{}) ([]byte, error) {
	processedBody := make(map[string]interface{})

	for key, templateValue := range template {
		if reflect.TypeOf(templateValue).Kind() == reflect.String {
			value := templateValue.(string)
			if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
				fieldName := strings.Trim(value, "{}")
				if v, ok := messageData[fieldName]; ok {
					processedBody[key] = v
				} else {
					return nil, fmt.Errorf("field %s not found in message", fieldName)
				}
			} else {
				processedBody[key] = templateValue
			}
		} else {
			processedBody[key] = templateValue
		}
	}

	bodyBytes, err := json.Marshal(processedBody)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
