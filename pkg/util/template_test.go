package util

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	template := map[string]interface{}{
		"name":   "John",
		"age":    "{{age}}",
		"nested": map[string]interface{}{"city": "{{city}}"},
	}
	messageData := map[string]interface{}{
		"age":  30,
		"city": "New York",
	}

	expectedProcessedBody := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"nested": map[string]interface{}{"city": "New York"},
	}
	expectedBytes, _ := json.Marshal(expectedProcessedBody)

	result, err := ProcessTemplate(template, messageData)
	if err != nil {
		t.Errorf("ProcessTemplate() error = %v", err)
	}

	if !reflect.DeepEqual(result, expectedBytes) {
		t.Errorf("ProcessTemplate() got = %v, want %v", result, expectedBytes)
	}
}

func TestProcessTemplate_FieldNotFound(t *testing.T) {
	template := map[string]interface{}{
		"name": "{{name}}",
	}
	messageData := map[string]interface{}{
		"age": 30,
	}

	_, err := ProcessTemplate(template, messageData)
	if err == nil {
		t.Errorf("ProcessTemplate() expected error, got nil")
	}
}

func TestProcessGraphQLTemplate_UnsupportedType(t *testing.T) {
	graphqlTemplate := "query { user(id: \"{{id}}\") { name }}"
	messageData := map[string]interface{}{
		"id": []int{1, 2, 3},
	}

	_, err := ProcessGraphQLTemplate(graphqlTemplate, messageData)
	if err == nil {
		t.Errorf("ProcessGraphQLTemplate() expected error, got nil")
	}
}
