package datautil_test

import (
	// Built-in modules.
	"testing"

	// External modules.

	// Internal modules.
	"github.com/cuberat-go/datautil"
	"github.com/stretchr/testify/assert"
)

func TestJSONPathStruct(t *testing.T) {
	// Test the JSON path generation for a struct with nested fields.
	type NestedStruct struct {
		Field1 string
		Field2 int
	}
	type TestStruct struct {
		Nested NestedStruct
	}

	testStruct := TestStruct{
		Nested: NestedStruct{
			Field1: "value1",
			Field2: 42,
		},
	}

	expectedPath := "$.Nested.Field1"
	gotPath := ""

	walker := datautil.NewWalker().
		WithStructHandlerString(func(jsonPath string, structVal any, fieldName string, fieldValue string, set datautil.SetValueFuncString) error {
			if fieldName == "Field1" {
				gotPath = jsonPath
			}

			return nil
		})

	err := walker.Walk(&testStruct)
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, gotPath)
}

func TestJSONPathSlice(t *testing.T) {
	// Test the JSON path generation for a slice of structs.
	type TestStruct struct {
		Field1 string
	}

	testSlice := []TestStruct{
		{Field1: "value1"},
		{Field1: "value2"},
	}

	expectedPath := "$[0].Field1"
	gotPath := ""

	walker := datautil.NewWalker().
		WithStructHandlerString(func(jsonPath string, structVal any, fieldName string, fieldValue string, set datautil.SetValueFuncString) error {
			if fieldName == "Field1" && fieldValue == "value1" {
				gotPath = jsonPath
			}

			return nil
		})

	err := walker.Walk(&testSlice)
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, gotPath)
}

func TestJSONPATHMap(t *testing.T) {
	// Test the JSON path generation for a map with string keys and values.
	testMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	expectedPath := "$.key1"
	gotPath := ""

	walker := datautil.NewWalker().
		WithMapHandlerString(func(jsonPath string, mapVal any, key any, value string, set datautil.SetValueFuncString) error {
			if key == "key1" {
				gotPath = jsonPath
			}

			return nil
		})

	err := walker.Walk(&testMap)
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, gotPath)
}

func TestJSONPATHNestedMapWithSlice(t *testing.T) {
	// Test the JSON path generation for a nested map containing a slice.
	testData := map[string]any{
		"outerKey": []any{
			map[string]string{
				"innerKey": "innerValue",
			},
		},
	}

	expectedPath := "$.outerKey[0].innerKey"
	gotPath := ""

	walker := datautil.NewWalker().
		WithMapHandlerString(func(jsonPath string, mapVal any, key any, value string, set datautil.SetValueFuncString) error {
			if key == "innerKey" {
				gotPath = jsonPath
			}
			return nil
		})

	err := walker.Walk(&testData)
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, gotPath)
}
