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

// Tests the Walker with handlers that accept string values. The handlers will
// only be called for fields, elements, or values that are of type string.
func TestWalkerString(t *testing.T) {
	expectedMap := map[string]string{
		"key1": "newValue1",
		"key2": "value2",
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
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

// Tests the Walker with handlers that accept any type of value.
func TestWalkerAny(t *testing.T) {
	structHandler := func(
		jsonPath string,
		structVal any,
		fieldName string,
		fieldValue any,
		set datautil.SetValueFunc,
	) error {
		if fieldValue == 1 {
			// Set a new value for the field in the struct.
			err := set(2)
			if err != nil {
				return err
			}
		}
		t.Logf("Struct field: %s, Value: %v\n", fieldName, fieldValue)
		return nil
	}

	sliceHandler := func(
		jsonPath string,
		sliceVal any,
		index int,
		elementValue any,
		set datautil.SetValueFunc,
	) error {
		if elementValue == "element3" {
			// Set a new value for the element in the slice.
			err := set("element3Mod")
			if err != nil {
				return err
			}
		}
		t.Logf("Slice index: %d, Value: %v\n", index, elementValue)
		return nil
	}

	expectedMap := map[string]string{
		"key1": "newValue1",
		"key2": "value2",
	}

	mapHandler := func(
		jsonPath string,
		mapVal any,
		key any,
		value any,
		set datautil.SetValueFunc,
	) error {
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
		WithMapHandler(mapHandler)

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

// Tests for values that are pointers to string (*string). This test ensures
// that the Walker can correctly dereference pointers and modify the underlying
// string values.
func TestWalkerStringPtr(t *testing.T) {
	var err error

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	expectedMap := map[string]*string{
		"key1": new("newValue1"),
		"key2": new("value2"),
	}
	map1 := map[string]*string{
		"key1": new("value1"),
		"key2": new("value2"),
	}

	err = walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
		return
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	slice1 := []*string{new("element1"), new("element2"), new("element3")}
	expectedSlice := []*string{new("element1"), new("element2"),
		new("element3Mod")}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
		return
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct struct {
		Field1 string
		Field2 *string
	}

	expectedStruct := myStruct{
		Field1: "structValue1",
		Field2: new("structValue2Mod"),
	}

	struct1 := myStruct{
		Field1: "structValue1",
		Field2: new("structValue2"),
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
		return
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %#v\n", struct1)
}

// Test for values that are pointers to pointers to string (**string). This test
// ensures that the Walker can correctly dereference multiple levels of
// pointers and modify the underlying string values.
func TestWalkerStringPtrPtr(t *testing.T) {
	var err error

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	expectedMap := map[string]**string{
		"key1": new(new("newValue1")),
		"key2": new(new("value2")),
	}
	map1 := map[string]**string{
		"key1": new(new("value1")),
		"key2": new(new("value2")),
	}

	err = walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
		return
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	slice1 := []**string{new(new("element1")), new(new("element2")),
		new(new("element3"))}
	expectedSlice := []**string{new(new("element1")), new(new("element2")),
		new(new("element3Mod"))}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
		return
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct struct {
		Field1 string
		Field2 **string
	}

	expectedStruct := myStruct{
		Field1: "structValue1",
		Field2: new(new("structValue2Mod")),
	}

	struct1 := myStruct{
		Field1: "structValue1",
		Field2: new(new("structValue2")),
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
		return
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %#v\n", struct1)
}

// Tests the Walker with handlers that accept string values, but the input data
// is provided as a pointer to interface{}. This test ensures that the Walker
// can correctly handle interface{} types and still invoke the appropriate
// handlers for string values.
func TestWalkerStringInterface(t *testing.T) {
	expectedMap := map[string]string{
		"key1": "newValue1",
		"key2": "value2",
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	var map1 any = map[string]string{
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

	var slice1 any = []string{"element1", "element2", "element3"}
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

	var expectedStruct any = myStruct{
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

// Tests the Walker with handlers that accept string values, but the input data
// contains interface{} types. This test ensures that the Walker
// can correctly handle interface{} types and still invoke the appropriate
// handlers for string values.
func TestWalkerStringInterface2(t *testing.T) {
	expectedMap := map[string]any{
		"key1": "newValue1",
		"key2": "value2",
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	var map1 any = map[string]any{
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

	var slice1 any = []any{"element1", "element2", "element3"}
	expectedSlice := []any{"element1", "element2", "element3Mod"}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct struct {
		Field1 string
		Field2 any
	}

	var expectedStruct any = myStruct{
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

func TestWalkerStringRecursive(t *testing.T) {
	expectedMap := map[string]map[string]string{
		"key1": {"subKey1": "newValue1"},
		"key2": {"subKey2": "value2"},
	}

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	var map1 any = map[string]map[string]string{
		"key1": {"subKey1": "value1"},
		"key2": {"subKey2": "value2"},
	}

	err := walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	var slice1 any = [][]string{{"element1", "element2", "element3"}}
	expectedSlice := [][]string{{"element1", "element2", "element3Mod"}}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct2 struct {
		Field1 string
		Field2 string
	}

	type myStruct struct {
		Field1 string
		Field2 *myStruct2
	}

	var expectedStruct any = myStruct{
		Field1: "structValue1",
		Field2: &myStruct2{
			Field1: "structValue1",
			Field2: "structValue2Mod",
		},
	}

	struct1 := myStruct{
		Field1: "structValue1",
		Field2: &myStruct2{
			Field1: "structValue1",
			Field2: "structValue2",
		},
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %+v\n", struct1)

	type myStruct3 struct {
		Field1 string
		Field2 myStruct2
	}

	expectedMapOfStructs := map[string]*myStruct3{
		"key1": {
			Field1: "structValue1",
			Field2: myStruct2{
				Field1: "structValue1",
				Field2: "structValue2Mod",
			},
		},
		"key2": {
			Field1: "structValue1",
			Field2: myStruct2{
				Field1: "structValue1",
				Field2: "structValue3",
			},
		},
	}
	mapOfStructs := map[string]*myStruct3{
		"key1": {
			Field1: "structValue1",
			Field2: myStruct2{
				Field1: "structValue1",
				Field2: "structValue2",
			},
		},
		"key2": {
			Field1: "structValue1",
			Field2: myStruct2{
				Field1: "structValue1",
				Field2: "structValue3",
			},
		},
	}

	err = walker.Walk(&mapOfStructs)
	if err != nil {
		t.Errorf("Error walking map of structs: %v", err)
	}

	assert.Equal(t, expectedMapOfStructs, mapOfStructs,
		"The map of structs should have been modified correctly.")
	t.Logf("Modified map of structs: %+v\n", mapOfStructs)
}

func TestWalkerStringNonPointers(t *testing.T) {
	expectedMap := map[string]map[string]string{
		"key1": {"subKey1": "newValue1"},
		"key2": {"subKey2": "value2"},
	}

	// mapHandlerString := func(mapVal any, key any, value string,
	// 	set datautil.SetValueFuncString) error {
	// 	if value == "value1" {
	// 		// Set a new value for the key in the map.
	// 		err := set("newValue1")
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	t.Logf("Map key: %v, Value: %s\n", key, value)
	// 	return nil
	// }

	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithStructHandlerString(structHandlerString).
		WithSliceHandlerString(sliceHandlerString).
		WithMapHandlerString(mapHandlerString)

	var map1 any = map[string]map[string]string{
		"key1": {"subKey1": "value1"},
		"key2": {"subKey2": "value2"},
	}

	err := walker.Walk(&map1)
	if err != nil {
		t.Errorf("Error walking map: %v", err)
	}

	assert.Equal(t, expectedMap, map1,
		"The map should have been modified correctly.")
	t.Logf("Modified map: %+v\n", map1)

	var slice1 any = [][]string{{"element1", "element2", "element3"}}
	expectedSlice := [][]string{{"element1", "element2", "element3Mod"}}
	err = walker.Walk(&slice1)
	if err != nil {
		t.Errorf("Error walking slice: %v", err)
	}

	assert.Equal(t, expectedSlice, slice1,
		"The slice should have been modified correctly.")
	t.Logf("Modified slice: %+v\n", slice1)

	type myStruct2 struct {
		SubField1 string
		SubField2 string
	}

	type myStruct struct {
		Field1 string
		Field2 myStruct2
	}

	var expectedStruct any = myStruct{
		Field1: "structValue1",
		Field2: myStruct2{
			SubField1: "structValue1",
			SubField2: "structValue2Mod",
		},
	}

	struct1 := myStruct{
		Field1: "structValue1",
		Field2: myStruct2{
			SubField1: "structValue1",
			SubField2: "structValue2",
		},
	}

	err = walker.Walk(&struct1)
	if err != nil {
		t.Errorf("Error walking struct: %v", err)
	}

	assert.Equal(t, expectedStruct, struct1,
		"The struct should have been modified correctly.")
	t.Logf("Modified struct: %+v\n", struct1)
}

func ExampleWalker_Walk_trivial() {
	sliceHandler := func(
		jsonPath string,
		sliceVal any,
		index int,
		elementValue string,
		set datautil.SetValueFuncString,
	) error {
		if elementValue == "element3" {
			// Set a new value for the element in the slice.
			err := set("element3Mod")
			if err != nil {
				return err
			}
		}

		return nil
	}
	// Create a new Walker instance.
	walker := datautil.NewWalker().
		WithSliceHandlerString(sliceHandler)
	slice := []string{"element1", "element2", "element3"}
	err := walker.Walk(&slice)
	if err != nil {
		fmt.Printf("Error walking slice: %v", err)
	}
	fmt.Printf("Modified slice: %+v\n", slice)
	// Output:
	// Modified slice: [element1 element2 element3Mod]
}

func mapHandlerString(
	jsonPath string,
	mapVal any,
	key any,
	value string,
	set datautil.SetValueFuncString,
) error {
	if value == "value1" {
		// Set a new value for the key in the map.
		err := set("newValue1")
		if err != nil {
			return err
		}
	}

	return nil
}

func sliceHandlerString(
	jsonPath string,
	sliceVal any,
	index int,
	elementValue string,
	set datautil.SetValueFuncString,
) error {
	if elementValue == "element3" {
		// Set a new value for the element in the slice.
		err := set("element3Mod")
		if err != nil {
			return err
		}
	}
	return nil
}

func structHandlerString(
	jsonPath string,
	structVal any,
	fieldName string,
	fieldValue string,
	set datautil.SetValueFuncString,
) error {
	if fieldValue == "structValue2" {
		// Set a new value for the field in the struct.
		err := set("structValue2Mod")
		if err != nil {
			return err
		}
	}
	return nil
}
