package datautil

import (
	// Built-in modules.
	"fmt"
	"reflect"
)

// Function type for setting a value in a data structure. The function takes a
// value of any type and returns an error if the operation fails.
type SetValueFunc func(value any) error

// Function type for setting a string value in a data structure. The function
// takes a string value and returns an error if the operation fails.
type SetValueFuncString func(value string) error

// Function type for handling struct traversal.
type StructWalkerHandler func(structVal any, fieldName string,
	fieldValue any, set SetValueFunc) error

// Function type for handling slice traversal.
type SliceWalkerHandler func(sliceVal any, index int, elementValue any,
	set SetValueFunc) error

// Function type for handling map traversal.
type MapWalkerHandler func(mapVal any, key any, value any,
	set SetValueFunc) error

// Function type for handling struct traversal with string values. The handler
// will not be called if the field value is not a string.
type StructWalkerHandlerString func(structVal any, fieldName string,
	fieldValue string, set SetValueFuncString) error

// Function type for handling slice traversal with string values. The handler
// will not be called if the element value is not a string.
type SliceWalkerHandlerString func(sliceVal any, index int,
	elementValue string, set SetValueFuncString) error

// Function type for handling map traversal with string values. The handler will
// not be called if the value is not a string.
type MapWalkerHandlerString func(mapVal any, key any, value string,
	set SetValueFuncString) error

// Structure representing a walker that traverses data structures and invokes
// handlers for structs, slices, and maps.
type Walker struct {
	structHandler       StructWalkerHandler
	structHandlerString StructWalkerHandlerString
	sliceHandler        SliceWalkerHandler
	sliceHandlerString  SliceWalkerHandlerString
	mapHandler          MapWalkerHandler
	mapHandlerString    MapWalkerHandlerString
}

// Creates a new instance of Walker.
func NewWalker() *Walker {
	return &Walker{}
}

// Sets the handler for struct traversal and returns the Walker instance.
func (w *Walker) WithStructHandler(handler StructWalkerHandler) *Walker {
	w.structHandler = handler
	return w
}

// Sets the handler for struct traversal with string values and returns the
// Walker instance. The handler will not be called if the field value is not a
// string.
func (w *Walker) WithStructHandlerString(
	handler StructWalkerHandlerString,
) *Walker {
	w.structHandlerString = handler
	return w
}

// Sets the handler for slice traversal and returns the Walker instance.
func (w *Walker) WithSliceHandler(handler SliceWalkerHandler) *Walker {
	w.sliceHandler = handler
	return w
}

// Sets the handler for slice traversal with string values and returns the
// Walker instance. The handler will not be called if the element value is not a
// string.
func (w *Walker) WithSliceHandlerString(
	handler SliceWalkerHandlerString,
) *Walker {
	w.sliceHandlerString = handler
	return w
}

// Sets the handler for map traversal and returns the Walker instance.
func (w *Walker) WithMapHandler(handler MapWalkerHandler) *Walker {
	w.mapHandler = handler
	return w
}

// Sets the handler for map traversal with string values and returns the
// Walker instance. The handler will not be called if the value is not a string.
func (w *Walker) WithMapHandlerString(
	handler MapWalkerHandlerString,
) *Walker {
	w.mapHandlerString = handler
	return w
}

// Traverses the provided data structure, invoking the appropriate handlers for
// structs, slices, and maps.
func (w *Walker) Walk(data any) error {
	value := reflect.ValueOf(data)
	kind := value.Kind()

	if kind != reflect.Pointer {
		return fmt.Errorf("data must be a pointer, got %s", kind)
	}

	// Dereference the pointer until we reach a non-pointer value.
	for value = value.Elem(); value.Kind() == reflect.Pointer; value = value.Elem() {
	}

	kind = value.Kind()

	switch kind {
	case reflect.Struct:
		return w.walkStruct(value)
	case reflect.Slice, reflect.Array:
		return w.walkSlice(value)
	case reflect.Map:
		return w.walkMap(value)
	default:
		return fmt.Errorf("unsupported data type: %s", kind)
	}
}

// Determines if the provided value is walkable (i.e., a struct, slice, or map).
func (w *Walker) isWalkable(value reflect.Value) bool {
	if value.Kind() == reflect.Pointer {
		// Dereference the pointer until we reach a non-pointer value.
		for value = value.Elem(); value.Kind() == reflect.Pointer; value = value.Elem() {
		}
	}

	kind := value.Kind()
	switch kind {
	case reflect.Struct:
		return true
	case reflect.Slice, reflect.Array:
		return true
	case reflect.Map:
		return true
	default:
		return false
	}

}

// Traverses a struct, invoking the appropriate handlers for each field.
func (w *Walker) walkStruct(value reflect.Value) error {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		if w.isWalkable(field) {
			if err := w.Walk(field); err != nil {
				return err
			}
			continue
		}

		dereferenced, ptrs := w.getPtrList(field)

		if dereferenced.Kind() == reflect.String &&
			w.structHandlerString != nil {
			setFunc := func(newValue string) error {
				if len(ptrs) > 0 {
					if !field.CanSet() || !dereferenced.CanSet() {
						return fmt.Errorf("cannot set value for field %s",
							fieldType.Name)
					}
					field.Set(ptrs[len(ptrs)-1])
					ptrs[0].Elem().SetString(newValue)
				} else {
					if !field.CanSet() {
						return fmt.Errorf("cannot set value for field %s",
							fieldType.Name)
					}
					field.SetString(newValue)
				}
				return nil
			}
			err := w.structHandlerString(value.Interface(), fieldType.Name,
				dereferenced.String(), setFunc)
			if err != nil {
				return err
			}
		} else {
			if w.structHandler != nil {
				setFunc := func(newValue any) error {
					if !field.CanSet() {
						return fmt.Errorf("cannot set value for field %s", fieldType.Name)
					}
					field.Set(reflect.ValueOf(newValue))
					return nil
				}
				if err := w.structHandler(value.Interface(), fieldType.Name, field.Interface(), setFunc); err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (w *Walker) getPtrList(
	value reflect.Value,
) (dereferenced reflect.Value, ptrs []reflect.Value) {
	if value.Kind() != reflect.Pointer {
		return value, nil
	}

	ptrCount := 1
	// Dereference the pointer until we reach a non-pointer value.
	for value = value.Elem(); value.Kind() == reflect.Pointer; value = value.Elem() {
		ptrCount++
	}

	dereferenced = value
	// Collect the pointers in reverse order.
	ptrs = make([]reflect.Value, 0, ptrCount)

	val := dereferenced
	for range ptrCount {
		newVal := reflect.New(val.Type())
		ptrs = append(ptrs, newVal)
		val = newVal
	}

	return dereferenced, ptrs
}

// Traverses a slice or array, invoking the appropriate handlers for each
// element.
func (w *Walker) walkSlice(value reflect.Value) error {
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)

		if w.isWalkable(elem) {
			if err := w.Walk(elem); err != nil {
				return err
			}
			continue
		}

		dereferenced, ptrs := w.getPtrList(elem)
		if dereferenced.Kind() == reflect.String &&
			w.sliceHandlerString != nil {
			setFunc := func(newValue string) error {
				if len(ptrs) > 0 {
					if !elem.CanSet() || !dereferenced.CanSet() {
						return fmt.Errorf("cannot set value for index %d", i)
					}
					elem.Set(ptrs[len(ptrs)-1])
					ptrs[0].Elem().SetString(newValue)
				} else {
					if !elem.CanSet() {
						return fmt.Errorf("cannot set value for index %d", i)
					}
					elem.SetString(newValue)
				}
				return nil
			}
			err := w.sliceHandlerString(value.Interface(), i, dereferenced.String(),
				setFunc)
			if err != nil {
				return err
			}
		} else {
			if w.sliceHandler != nil {
				setFunc := func(newValue any) error {
					if !elem.CanSet() {
						return fmt.Errorf("cannot set value for index %d", i)
					}
					elem.Set(reflect.ValueOf(newValue))
					return nil
				}
				if err := w.sliceHandler(value.Interface(), i,
					elem.Interface(), setFunc); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Traverses a map, invoking the appropriate handlers for each key-value pair.
func (w *Walker) walkMap(value reflect.Value) error {
	for _, key := range value.MapKeys() {
		val := value.MapIndex(key)

		if w.isWalkable(val) {
			if err := w.Walk(val); err != nil {
				return err
			}
			continue
		}

		dereferenced, ptrs := w.getPtrList(val)

		if dereferenced.Kind() == reflect.String && w.mapHandlerString != nil {
			setFunc := func(newValue string) error {
				if len(ptrs) > 0 {
					value.SetMapIndex(key, ptrs[len(ptrs)-1])
					ptrs[0].Elem().SetString(newValue)
				} else {
					value.SetMapIndex(key, reflect.ValueOf(newValue))
				}

				return nil
			}
			err := w.mapHandlerString(value.Interface(), key.Interface(),
				dereferenced.String(), setFunc)
			if err != nil {
				return err
			}
		} else {
			if w.mapHandler != nil {
				setFunc := func(newValue any) error {
					value.SetMapIndex(key, reflect.ValueOf(newValue))
					return nil
				}
				if err := w.mapHandler(value.Interface(), key.Interface(),
					val.Interface(), setFunc); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
