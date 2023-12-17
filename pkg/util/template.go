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

// ProcessGraphQLTemplate processes the graphql template and returns the processed body
func ProcessGraphQLTemplate(graphqlTemplate string, messageData map[string]interface{}) (string, error) {
	processedQuery := graphqlTemplate

	for key, value := range messageData {
		placeholder := fmt.Sprintf("{{%s}}", key)
		var valueStr string

		// Convert different types to their string representations
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

	return processedQuery, nil
}
