package datautil_test

import (
	// Built-in modules.
	"testing"

	// External modules.
	"fmt"
	// Internal modules.
	"github.com/cuberat-go/datautil"
	"github.com/stretchr/testify/assert"
)

func TestWalkerString(t *testing.T) {
	structHandler := func(structVal any, fieldName string, fieldValue string,
		set datautil.SetValueFuncString) error {
		if fieldValue == "structValue2" {
			// Set a new value for the field in the struct.
			err := set("structValue2Mod")
			if err != nil {
				return err
			}
		}
		fmt.Printf("Struct field: %s, Value: %s\n", fieldName, fieldValue)
		return nil
	}

	sliceHandler := func(sliceVal any, index int, elementValue string,
		set datautil.SetValueFuncString) error {
		if elementValue == "element3" {
			// Set a new value for the element in the slice.
			err := set("element3Mod")
			if err != nil {
				return err
			}
		}
		fmt.Printf("Slice index: %d, Value: %s\n", index, elementValue)
		return nil
	}

	expectedMap := map[string]string{
		"key1": "newValue1",
		"key2": "value2",
	}

	mapHandlerString := func(mapVal any, key any, value string,
		set datautil.SetValueFuncString) error {
		if value == "value1" {
			// Set a new value for the key in the map.
			err := set("newValue1")
			if err != nil {
				return err
			}
		}
		t.Logf("Map key: %v, Value: %s\n", key, value)
		return nil
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandler).
		WithSliceHandlerString(sliceHandler).
		WithMapHandlerString(mapHandlerString)

	map1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	slice1 := []string{"element1", "element2", "element3"}
	expectedSlice := []string{"element1", "element2", "element3Mod"}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct struct {
		Field1 string
		Field2 string
	}

	expectedStruct := myStruct{
		Field1: "structValue1",
		Field2: "structValue2Mod",
	}

	struct1 := myStruct{
		Field1: "structValue1",
		Field2: "structValue2",
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %+v\n", struct1)
}

func TestWalkerAny(t *testing.T) {
	structHandler := func(structVal any, fieldName string, fieldValue any,
		set datautil.SetValueFunc) error {
		if fieldValue == 1 {
			// Set a new value for the field in the struct.
			err := set(2)
			if err != nil {
				return err
			}
		}
		fmt.Printf("Struct field: %s, Value: %v\n", fieldName, fieldValue)
		return nil
	}

	sliceHandler := func(sliceVal any, index int, elementValue any,
		set datautil.SetValueFunc) error {
		if elementValue == "element3" {
			// Set a new value for the element in the slice.
			err := set("element3Mod")
			if err != nil {
				return err
			}
		}
		fmt.Printf("Slice index: %d, Value: %v\n", index, elementValue)
		return nil
	}

	expectedMap := map[string]string{
		"key1": "newValue1",
		"key2": "value2",
	}

	mapHandlerString := func(mapVal any, key any, value any,
		set datautil.SetValueFunc) error {
		if value == "value1" {
			// Set a new value for the key in the map.
			err := set("newValue1")
			if err != nil {
				return err
			}
		}
		t.Logf("Map key: %v, Value: %v\n", key, value)
		return nil
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandler(structHandler).
		WithSliceHandler(sliceHandler).
		WithMapHandler(mapHandlerString)

	map1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	slice1 := []string{"element1", "element2", "element3"}
	expectedSlice := []string{"element1", "element2", "element3Mod"}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct struct {
		Field1 int
		Field2 string
	}

	expectedStruct := myStruct{
		Field1: 2,
		Field2: "structValue2",
	}

	struct1 := myStruct{
		Field1: 1,
		Field2: "structValue2",
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %+v\n", struct1)
}
