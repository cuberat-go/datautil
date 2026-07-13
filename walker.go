package datautil

import (
	// Built-in modules.
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

// Function type for setting a value in a data structure. The function takes a
// value of any type and returns an error if the operation fails.
type SetValueFunc func(value any) error

// Function type for setting a string value in a data structure. The function
// takes a string value and returns an error if the operation fails.
type SetValueFuncString func(value string) error

// Function type for handling struct traversal. The handler will be called for
// each field in the struct. Use the set() function to modify the value in the
// struct if so desired.
//
// jsonPath is the JSONPath to the field (RFC 9535), structVal is the struct
// itself, fieldName is the name of the field, and fieldValue is the value of
// the field.
type StructWalkerHandler func(jsonPath string, structVal any, fieldName string,
	fieldValue any, set SetValueFunc) error

// Function type for handling slice traversal. The handler will be called for
// each element in the slice. Use the set() function to modify the value in
// the slice if so desired.
//
// jsonPath is the JSONPath to the element (RFC 9535), sliceVal is the slice
// itself, index is the index of the element, and elementValue is the value of
// the element.
type SliceWalkerHandler func(jsonPath string,
	sliceVal any, index int, elementValue any, set SetValueFunc) error

// Function type for handling map traversal. The handler will be called for each
// key-value pair in the map. Use the set() function to modify the value in
// the map if so desired.
//
// jsonPath is the JSONPath to the value (RFC 9535), mapVal is the map itself,
// key is the key in the map, and value is the value of the key-value pair.
type MapWalkerHandler func(jsonPath string, mapVal any, key any, value any,
	set SetValueFunc) error

// Function type for handling struct traversal with string values. The handler
// will be called for each field in the struct whose value is a string. Use
// the set() function to modify the value in the struct if so desired. As a
// special case, if the field value is a pointer to a string, the handler will
// be called with the dereferenced string value, and the set() function will
// modify the value in the struct by setting the pointer to a new string
// value.
//
// jsonPath is the JSONPath to the field (RFC 9535), structVal is the
// struct itself, fieldName is the name of the field, and fieldValue is the
// value of the field.
type StructWalkerHandlerString func(jsonPath string, structVal any,
	fieldName string, fieldValue string, set SetValueFuncString) error

// Function type for handling slice traversal with string values. The handler
// will be called for each element in the slice whose value is a string. Use the
// set() function to modify the value in the slice if so desired. As a special
// case, if the slice element is a pointer to a string, the handler will be
// called with the dereferenced string value, and the set() function will
// modify the value in the slice by setting the pointer to a new string value.
//
// jsonPath is the JSONPath to the element (RFC 9535), sliceVal is the slice
// itself, index is the index of the element, and elementValue is the value of
// the element.
type SliceWalkerHandlerString func(jsonPath string, sliceVal any, index int,
	elementValue string, set SetValueFuncString) error

// Function type for handling map traversal with string values. The handler will
// be called for each key-value pair in the map whose value is a string. Use the
// set() function to modify the value in the map if so desired. As a special
// case, if the map value is a pointer to a string, the handler will be
// called with the dereferenced string value, and the set() function will
// modify the value in the map by setting the pointer to a new string value.
//
// jsonPath is the JSONPath to the value (RFC 9535), mapVal is the map itself,
// key is the key in the map, and value is the value of the key-value pair.
type MapWalkerHandlerString func(jsonPath string, mapVal any, key any,
	value string, set SetValueFuncString) error

// Structure representing a walker that traverses data structures and invokes
// handlers for structs, slices, and maps.
type Walker struct {
	debug               bool
	logger              *slog.Logger
	logLevel            *slog.LevelVar
	structHandler       StructWalkerHandler
	structHandlerString StructWalkerHandlerString
	sliceHandler        SliceWalkerHandler
	sliceHandlerString  SliceWalkerHandlerString
	mapHandler          MapWalkerHandler
	mapHandlerString    MapWalkerHandlerString
}

// Creates a new instance of Walker.
func NewWalker() *Walker {
	w := &Walker{}
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
	w.logLevel = logLevel
	w.logger = slog.New(slog.NewTextHandler(
		os.Stderr, &slog.HandlerOptions{Level: logLevel}),
	)
	slog.SetDefault(w.logger)
	return w
}

// Used for debugging this module. Do not use, as this may change or go away
// anytime.
func (w *Walker) xWithDebug() *Walker {
	w.debug = true
	w.logLevel.Set(slog.LevelDebug)
	return w
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

	return w.doWalk(value, "$")
}

func (w *Walker) doWalk(value reflect.Value, jsonPath string) error {
	// value := reflect.ValueOf(data)
	kind := value.Kind()

	if kind == reflect.Pointer || kind == reflect.Interface {
		// Dereference the pointer until we reach a non-pointer value.
		for value = value.Elem(); value.Kind() == reflect.Pointer ||
			value.Kind() == reflect.Interface; value = value.Elem() {
		}
	}

	kind = value.Kind()

	w.logger.Debug("doWalk", "kind", kind, "jsonPath", jsonPath)

	switch kind {
	case reflect.Struct:
		return w.walkStruct(value, jsonPath)
	case reflect.Slice, reflect.Array:
		return w.walkSlice(value, jsonPath)
	case reflect.Map:
		return w.walkMap(value, jsonPath)
	default:
		return fmt.Errorf("unsupported data type: %s", kind)
	}
}

// Determines if the provided value is walkable (i.e., a struct, slice, or map).
func (w *Walker) isWalkable(value reflect.Value) bool {
	if value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		// Dereference the pointer until we reach a non-pointer value.
		for value = value.Elem(); value.Kind() == reflect.Pointer ||
			value.Kind() == reflect.Interface; value = value.Elem() {
			// if value.Kind() == reflect.Interface {
			// 	// FIXME: check the thing inside the interface.
			// 	if w.isCollection(value) {
			// 		return true
			// 	}
			// 	return false
			// }
		}
	}

	return w.isCollection(value)
}

func (w *Walker) isCollection(value reflect.Value) bool {
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
func (w *Walker) walkStruct(value reflect.Value, jsonPath string) error {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		newPath := w.addJSONPathName(jsonPath, fieldType.Name)

		if w.isWalkable(field) {
			if err := w.doWalk(field, newPath); err != nil {
				return err
			}
			continue
		}

		isString, origVal := w.isString(field)
		if isString && w.structHandlerString != nil {
			setFunc := func(newValue string) error {
				toSet, err := w.setString(field, newValue)
				if err != nil {
					return err
				}
				field.Set(toSet)
				return nil
			}
			err := w.structHandlerString(newPath, value.Interface(),
				fieldType.Name, origVal.String(), setFunc)
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
				if err := w.structHandler(newPath, value.Interface(),
					fieldType.Name, field.Interface(), setFunc); err != nil {
					return err
				}
			}
		}

	}
	return nil
}

// Traverses a slice or array, invoking the appropriate handlers for each
// element.
func (w *Walker) walkSlice(value reflect.Value, jsonPath string) error {
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)

		newPath := w.addJSONPathIndex(jsonPath, i)

		if w.isWalkable(elem) {
			if err := w.doWalk(elem, newPath); err != nil {
				return err
			}
			continue
		}

		isString, origVal := w.isString(elem)
		if isString && w.sliceHandlerString != nil {
			setFunc := func(newValue string) error {
				toSet, err := w.setString(elem, newValue)
				if err != nil {
					return err
				}
				elem.Set(toSet)
				return nil
			}

			err := w.sliceHandlerString(newPath, value.Interface(), i,
				origVal.String(), setFunc)
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
				if err := w.sliceHandler(newPath, value.Interface(), i,
					elem.Interface(), setFunc); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Traverses a map, invoking the appropriate handlers for each key-value pair.
func (w *Walker) walkMap(value reflect.Value, jsonPath string) error {
	for _, key := range value.MapKeys() {
		val := value.MapIndex(key)
		logger := w.logger.With(slog.String("jsonPath", jsonPath),
			slog.String("key", key.Interface().(string)))

		newPath := w.addJSONPathName(jsonPath, key.Interface().(string))

		if w.isWalkable(val) {
			logger.Debug("walkMap -- val is walkable")
			if err := w.doWalk(val, newPath); err != nil {
				return err
			}
			continue
		}

		isString, origVal := w.isString(val)

		logger.Debug("isString", slog.Bool("isString", isString))

		if isString && w.mapHandlerString != nil {
			setFunc := func(newValue string) error {
				toSet, err := w.setString(val, newValue)
				if err != nil {
					return err
				}
				value.SetMapIndex(key, toSet)

				return nil
			}
			err := w.mapHandlerString(newPath, value.Interface(),
				key.Interface(), origVal.String(), setFunc)
			if err != nil {
				return err
			}
		} else {
			if w.mapHandler != nil {
				setFunc := func(newValue any) error {
					value.SetMapIndex(key, reflect.ValueOf(newValue))
					return nil
				}
				if err := w.mapHandler(newPath, value.Interface(),
					key.Interface(), val.Interface(), setFunc); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Walks the chain of pointers and interfaces to determine if the final value is
// a string. Returns true and the final value if it is a string, otherwise
// returns false and the final value.
func (w *Walker) isString(value reflect.Value) (bool, reflect.Value) {
	if value.Kind() != reflect.Pointer && value.Kind() != reflect.Interface {
		if value.Kind() == reflect.String {
			return true, value
		}
		return false, value
	}

	// Dereference the pointer until we reach a non-pointer value.
	for value = value.Elem(); value.Kind() == reflect.Pointer ||
		value.Kind() == reflect.Interface; value = value.Elem() {
	}

	if value.Kind() == reflect.String {
		return true, value
	}
	return false, value
}

// Set the string value. `initVal` can be a string, a pointer to a string, an
// interface containing a string, a pointer to a pointer to a string, etc., in
// any combination. The return value `toSet` will be the string or initial
// pointer or interface to set in the original data structure. This method
// assumes that the caller has already checked that `initVal` is a string or a
// pointer to a string or an interface containing a string, etc.
func (w *Walker) setString(
	initVal reflect.Value,
	newVal string,
) (toSet reflect.Value, err error) {
	if initVal.Kind() != reflect.Pointer && initVal.Kind() != reflect.Interface {
		return reflect.ValueOf(newVal), nil
	}

	logger := w.logger.With(slog.String("initVal", initVal.String()))

	links := []reflect.Kind{}
	val := initVal
	done := false
	for {
		switch val.Kind() {
		case reflect.Pointer:
			links = append(links, reflect.Pointer)
		case reflect.Interface:
			links = append(links, reflect.Interface)
		default:
			done = true
		}
		if done {
			break
		}
		val = val.Elem()
	}

	val = reflect.ValueOf(newVal)
	for i := len(links) - 1; i >= 0; i-- {
		switch links[i] {
		case reflect.Pointer:
			v := reflect.New(val.Type())
			v.Elem().Set(val)
			val = v
		case reflect.Interface:
			v := reflect.New(val.Type()).Elem()
			v.Set(reflect.ValueOf(val.Interface()))
			val = v
		default:
			// should never happen
			logger.Debug("setString: unexpected kind", "kind", links[i])
			return reflect.Value{}, fmt.Errorf("unexpected kind: %s", links[i])
		}
	}

	toSet = val

	return toSet, nil
}

// Appends a JSON path member name to the existing JSON path. If the name is a
// valid JSON path member shorthand, it will be appended using dot notation.
// Otherwise, it will be appended using bracket notation with quotes.
func (w *Walker) addJSONPathName(jsonPath string, name string) string {
	if w.validJSONPathMemberShorthand(name) {
		return fmt.Sprintf("%s.%s", jsonPath, name)
	}
	return fmt.Sprintf("%s[%q]", jsonPath, name)
}

// Appends a JSON path index to the existing JSON path using bracket notation.
func (w *Walker) addJSONPathIndex(jsonPath string, index int) string {
	return fmt.Sprintf("%s[%d]", jsonPath, index)
}

// Determines if a rune is an invalid character for a JSON path member name.
func isInvalidJSONPathMemberCharsFunc(r rune) bool {
	if r > 'a' && r < 'z' || r > 'A' && r < 'Z' || r > '0' && r < '9' ||
		r == '_' {
		return false
	}
	if r >= 0x80 && r <= 0xd7ff || r >= 0xe000 && r <= 0x10ffff {
		return false
	}
	return true
}

// Determines if a JSON path member name is valid for shorthand notation. A
// valid JSON path member name can be used in dot notation without quotes.
func (w *Walker) validJSONPathMemberShorthand(name string) bool {
	return !strings.ContainsFunc(name, isInvalidJSONPathMemberCharsFunc)
}
