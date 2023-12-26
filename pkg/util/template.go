package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ProcessTemplate processes the template and returns the processed body
func ProcessTemplate(template map[string]interface{}, messageData map[string]interface{}) ([]byte, error) {
	processedBody, err := process(template, messageData)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := json.Marshal(processedBody)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

// process checks each value in the template recursively and replaces the values with the values in messageData
func process(template map[string]interface{}, messageData map[string]interface{}) (map[string]interface{}, error) {
	processedBody := make(map[string]interface{})

	for key, templateValue := range template {
		switch v := templateValue.(type) {
		case string:
			if strings.Contains(v, "{{") && strings.Contains(v, "}}") {
				fieldName := strings.Trim(v, "{}")
				if value, ok := messageData[fieldName]; ok {
					processedBody[key] = value
				} else {
					return nil, fmt.Errorf("field %s not found in message", fieldName)
				}
			} else {
				processedBody[key] = v
			}
		case map[string]interface{}:
			processedMap, err := process(v, messageData)
			if err != nil {
				return nil, err
			}
			processedBody[key] = processedMap
		default:
			processedBody[key] = v
		}
	}

	return processedBody, nil
}

// ProcessGraphQLTemplate processes the graphql template and returns the processed body
func ProcessGraphQLTemplate(graphqlTemplate string, messageData map[string]interface{}) (string, error) {
	processedQuery := graphqlTemplate

	if strings.Contains(graphqlTemplate, "{{") && strings.Contains(graphqlTemplate, "}}") {
		for key, value := range messageData {
			placeholder := fmt.Sprintf("{{%s}}", key)
			var valueStr string

			switch v := value.(type) {
			case string:
				valueStr = fmt.Sprintf("\"%s\"", v)
			case int, float64, bool:
				valueStr = fmt.Sprintf("%v", v)
			default:
				return "", fmt.Errorf("unsupported type for key %s", key)
			}

			processedQuery = strings.Replace(processedQuery, placeholder, valueStr, -1)
		}
	}

	return processedQuery, nil
}
