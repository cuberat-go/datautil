package marshal

import (
	// Built-in/core modules.
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

var ErrInvalidData = fmt.Errorf("invalid data")
var ErrUnsupportedType = fmt.Errorf("unsupported type")
var ErrNullValue = fmt.Errorf("null value")

type Reader interface {
	io.Reader
	io.ByteReader
}

type fieldConf struct {
	Kind reflect.Kind
}

// Object structure for prefix-framed marshaling of data.
//
// Datatypes are marshaled into bytes according to their kind, with support for
// structs, maps, slices, arrays, booleans, integers, unsigned integers,
// floats, strings, and complex numbers. Custom types implementing the
// encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interfaces are also
// supported.
type PrefixFramedMarshaler struct {
}

// Marshals the given data into the provided writer using prefix-framed
// encoding.
func (pfm *PrefixFramedMarshaler) MarshalTo(data any, w io.Writer) error {
	return pfm.marshalTo(reflect.ValueOf(data), w)
}

// Marshals the given data and returns the resulting byte slice using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) Marshal(data any) ([]byte, error) {
	w := &bytes.Buffer{}
	err := pfm.marshalTo(reflect.ValueOf(data), w)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Unmarshals data from the provided reader into the given value using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) UnmarshalFrom(
	r io.Reader,
	val any,
) error {
	var mr Reader
	var ok bool
	if mr, ok = r.(Reader); !ok {
		mr = bufio.NewReader(r)
	}
	err := pfm.unmarshalFrom(mr, reflect.ValueOf(val))
	if err == ErrNullValue {
		return nil
	}
	return err
}

// Unmarshals data from the provided byte slice into the given value using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) Unmarshal(
	data []byte,
	val any,
) error {
	err := pfm.unmarshalFrom(bytes.NewReader(data), reflect.ValueOf(val))
	if err == ErrNullValue {
		return nil
	}
	return err
}

func (pfm *PrefixFramedMarshaler) marshalCustomValueTo(
	conf fieldConf,
	val encoding.BinaryMarshaler,
	w io.Writer,
) error {
	data, err := val.MarshalBinary()
	if err != nil {
		return err
	}
	return pfm.marshalByteFieldTo(conf, data, w)
}

func (pfm *PrefixFramedMarshaler) unmarshalCustomValueFrom(
	val encoding.BinaryUnmarshaler,
	r Reader,
) error {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	data := make([]byte, size)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}
	return val.UnmarshalBinary(data)
}

func (pfm *PrefixFramedMarshaler) marshalTo(
	val reflect.Value,
	w io.Writer,
) error {
	if !val.IsValid() {
		return ErrInvalidData
	}

	conf := fieldConf{
		Kind: val.Kind(),
	}

	if mType, ok := val.Interface().(encoding.BinaryMarshaler); ok {
		return pfm.marshalCustomValueTo(conf, mType, w)
	}

	kind := val.Kind()
	for ; kind == reflect.Pointer ||
		kind == reflect.Interface; kind = val.Kind() {
		if val.IsNil() {
			return pfm.marshalByteFieldTo(conf, nil, w)
		}
		val = val.Elem()
	}

	if !val.IsValid() {
		return fmt.Errorf("%w: after dereferencing in marshalTo",
			ErrInvalidData)
	}

	switch kind {
	case reflect.Struct:
		return pfm.marshalStructValueTo(conf, val, w)
	case reflect.Map:
		return pfm.marshalMapValueTo(conf, val, w)
	case reflect.Slice, reflect.Array:
		return pfm.marshalSliceValueTo(conf, val, w)
	case reflect.Bool:
		return pfm.marshalBoolValueTo(conf, val, w)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return pfm.marshalIntValueTo(conf, val, w)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return pfm.marshalUintValueTo(conf, val, w)
	case reflect.Float32, reflect.Float64:
		return pfm.marshalFloatValueTo(conf, val, w)
	case reflect.String:
		return pfm.marshalStringValueTo(conf, val, w)
	case reflect.Complex64, reflect.Complex128:
		return pfm.marshalComplexValueTo(conf, val, w)
	default:
		return fmt.Errorf("%w: %s in marshalTo", ErrUnsupportedType, kind)
	}
}

func (pfm *PrefixFramedMarshaler) unmarshalFrom(
	r Reader,
	val reflect.Value,
) error {
	if !val.IsValid() {
		return ErrInvalidData
	}

	kind := val.Kind()
	if kind != reflect.Pointer {
		return fmt.Errorf("expected pointer type, got: %s", kind)
	}

	val = pfm.ensureValue(val)
	for val.IsValid() {
		if mType, ok := val.Interface().(encoding.BinaryUnmarshaler); ok {
			return pfm.unmarshalCustomValueFrom(mType, r)
		}
		if val.Kind() != reflect.Pointer && val.Kind() != reflect.Interface {
			break
		}
		if val.IsNil() {
			return fmt.Errorf("%w: nil pointer in unmarshalFrom", ErrInvalidData)
		}
		val = val.Elem()
	}

	if !val.IsValid() {
		return fmt.Errorf("%w: after dereferencing in unmarshalFrom",
			ErrInvalidData)
	}

	conf := fieldConf{
		Kind: kind,
	}

	kind = val.Kind()
	switch kind {
	case reflect.Struct:
		return pfm.unmarshalStructValueFrom(r, pfm.ensureStruct(val))
	case reflect.Map:
		return pfm.unmarshalMapValueFrom(r, pfm.ensureMap(val))
	case reflect.Slice, reflect.Array:
		return pfm.unmarshalSliceValueFrom(r, pfm.ensureSlice(val))
	case reflect.Bool:
		return pfm.unmarshalBoolValueFrom(conf, r, val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return pfm.unmarshalIntValueFrom(conf, r, val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return pfm.unmarshalUintValueFrom(conf, r, val)
	case reflect.Float32, reflect.Float64:
		return pfm.unmarshalFloatValueFrom(conf, r, val)
	case reflect.String:
		return pfm.unmarshalStringValueFrom(conf, r, val)
	case reflect.Complex64, reflect.Complex128:
		return pfm.unmarshalComplexValueFrom(conf, r, val)
	case reflect.Chan:
		return fmt.Errorf("%w: %s in unmarshalFrom", ErrUnsupportedType, kind)
	case reflect.Func:
		return fmt.Errorf("%w: %s in unmarshalFrom", ErrUnsupportedType, kind)
	default:
		return fmt.Errorf("%w: %s in unmarshalFrom", ErrUnsupportedType, kind)
	}
}

func (pfm *PrefixFramedMarshaler) ensureValue(
	value reflect.Value,
) reflect.Value {
	switch value.Kind() {
	case reflect.Pointer:
		for thisVal := value; thisVal.Kind() == reflect.Pointer; thisVal = thisVal.Elem() {
			if thisVal.IsNil() || !thisVal.Elem().IsValid() {
				thisVal.Set(pfm.ensureValue(reflect.New(thisVal.Type().Elem())))
			}
		}
		return value
	case reflect.Struct:
		return pfm.ensureStruct(value)
	case reflect.Map:
		return pfm.ensureMap(value)
	case reflect.Slice:
		return pfm.ensureSlice(value)
	}

	return value
}

func (pfm *PrefixFramedMarshaler) ensureMap(value reflect.Value) reflect.Value {
	if value.IsNil() {
		value.Set(reflect.MakeMap(value.Type()))
	}

	return value
}

func (pfm *PrefixFramedMarshaler) ensureSlice(
	value reflect.Value,
) reflect.Value {
	if value.Kind() == reflect.Array {
		return value
	}

	if value.IsNil() {
		value.Set(reflect.MakeSlice(value.Type(), 0, 0))
	}

	return value
}

func (pfm *PrefixFramedMarshaler) ensureStruct(value reflect.Value) reflect.Value {
	if !value.IsValid() {
		value.Set(reflect.New(value.Type()).Elem())
	}

	return value
}

func (pfm *PrefixFramedMarshaler) marshalMapValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	mapOut := &bytes.Buffer{}
	iter := value.MapRange()
	for iter.Next() {
		err := pfm.marshalTo(iter.Key(), mapOut)
		if err != nil {
			return err
		}
		err = pfm.marshalTo(iter.Value(), mapOut)
		if err != nil {
			return err
		}
	}

	return pfm.marshalByteFieldTo(conf, mapOut.Bytes(), w)
}

func (pfm *PrefixFramedMarshaler) unmarshalMapValueFrom(
	r Reader,
	value reflect.Value,
) error {
	keyType := value.Type().Key()
	valType := value.Type().Elem()

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	if size == 0 {
		value.SetZero()
		return nil
	}

	limitedReader := bufio.NewReader(io.LimitReader(r, int64(size)))

	for err == nil {
		keyPtr := reflect.New(keyType)
		err = pfm.unmarshalFrom(limitedReader, keyPtr)
		if err != nil {
			break
		}
		valPtr := reflect.New(valType)
		err = pfm.unmarshalFrom(limitedReader, valPtr)
		if err != nil {
			break
		}
		value.SetMapIndex(keyPtr.Elem(), valPtr.Elem())
	}

	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (pfm *PrefixFramedMarshaler) marshalStructValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	if value.IsZero() {
		return pfm.marshalByteFieldTo(conf, nil, w)
	}
	b := &bytes.Buffer{}
	for field, value := range value.Fields() {
		if !field.IsExported() {
			continue
		}
		err := pfm.marshalTo(value, b)
		if err != nil {
			return err
		}
	}

	return pfm.marshalByteFieldTo(conf, b.Bytes(), w)
}

func (pfm *PrefixFramedMarshaler) unmarshalStructValueFrom(
	r Reader,
	value reflect.Value,
) error {
	if !value.IsValid() {
		return ErrInvalidData
	}

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if size == 0 {
		value.SetZero()
		return ErrNullValue
	}
	limitedReader := bufio.NewReader(io.LimitReader(r, int64(size)))

	valType := value.Type()

	for i := 0; i < valType.NumField(); i++ {
		typeField := valType.Field(i)
		if !typeField.IsExported() {
			continue
		}
		fieldType := typeField.Type

		newValPtr := reflect.New(fieldType)
		newValPtr = pfm.ensureValue(newValPtr)
		err = pfm.unmarshalFrom(limitedReader, newValPtr)
		if err == ErrNullValue {
			value.Field(i).SetZero()
			continue
		}
		if err != nil {
			return err
		}
		value.Field(i).Set(newValPtr.Elem())
	}

	return nil
}

func (pfm *PrefixFramedMarshaler) marshalStringValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	return pfm.marshalByteFieldTo(conf, []byte(value.String()), w)
}

func (pfm *PrefixFramedMarshaler) marshalSliceValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	sliceType := value.Type()
	elemType := sliceType.Elem()

	if elemType == reflect.TypeFor[byte]() {
		return pfm.marshalByteSliceValue(conf, value, w)
	}

	sliceOut := &bytes.Buffer{}
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		err := pfm.marshalTo(elem, sliceOut)
		if err != nil {
			return err
		}
	}
	return pfm.marshalByteFieldTo(conf, sliceOut.Bytes(), w)
}

func (pfm *PrefixFramedMarshaler) unmarshalSliceValueFrom(
	r Reader,
	value reflect.Value,
) error {
	sliceType := value.Type()
	elemType := sliceType.Elem()
	byteType := reflect.TypeFor[byte]()

	if elemType == byteType {
		return pfm.unmarshalByteFieldFrom(r, value)
	}

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	if size == 0 {
		value.SetZero()
		return nil
	}

	limitedReader := bufio.NewReader(io.LimitReader(r, int64(size)))
	for i := 0; err == nil; i++ {
		elem := reflect.New(elemType)
		err = pfm.unmarshalFrom(limitedReader, elem)
		if err == nil {
			if value.Kind() == reflect.Array {
				if i >= value.Len() {
					continue
				}
				e := value.Index(i)
				e.Set(elem.Elem())
				continue
			}
			value.Set(reflect.Append(value, elem.Elem()))
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func (pfm *PrefixFramedMarshaler) marshalByteSliceValue(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	data := value.Bytes()
	return pfm.marshalByteFieldTo(conf, data, w)
}

func (pfm *PrefixFramedMarshaler) unmarshalByteFieldFrom(
	r Reader,
	value reflect.Value,
) error {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if size == 0 {
		value.SetBytes(nil)
		return nil
	}
	data := make([]byte, size)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}
	value.SetBytes(data)
	return nil
}

func (pfm *PrefixFramedMarshaler) unmarshalByteFieldToBytes(
	r Reader,
) ([]byte, error) {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return nil, ErrNullValue
	}
	data := make([]byte, size)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (pfm *PrefixFramedMarshaler) writeSize(size int, w io.Writer) error {
	sizeBuf := make([]byte, binary.MaxVarintLen32)
	sizeLen := binary.PutUvarint(sizeBuf, uint64(size))
	_, err := w.Write(sizeBuf[:sizeLen])
	return err
}

func (pfm *PrefixFramedMarshaler) marshalByteFieldTo(
	conf fieldConf,
	data []byte,
	w io.Writer,
) error {
	err := pfm.writeSize(len(data), w)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func (pfm *PrefixFramedMarshaler) marshalUIntField(
	conf fieldConf,
	value uint64,
	w io.Writer,
) error {
	// var err error
	// fieldBuf := make([]byte, binary.MaxVarintLen64)
	// fieldSize := binary.PutUvarint(fieldBuf, value)
	// return pfm.marshalByteFieldTo(fieldBuf[:fieldSize], w)

	var fieldBuf []byte

	switch conf.Kind {
	case reflect.Uint8:
		fieldBuf = make([]byte, 1)
		fieldBuf[0] = byte(value)
	case reflect.Uint16:
		fieldBuf = make([]byte, 2)
		binary.BigEndian.PutUint16(fieldBuf, uint16(value))
	case reflect.Uint32:
		fieldBuf = make([]byte, 4)
		binary.BigEndian.PutUint32(fieldBuf, uint32(value))
	case reflect.Uint64:
		fieldBuf = make([]byte, 8)
		binary.BigEndian.PutUint64(fieldBuf, value)
	default:
		return fmt.Errorf("unsupported uint kind: %v", conf.Kind)
	}

	return pfm.marshalByteFieldTo(conf, fieldBuf, w)
}

func (pfm *PrefixFramedMarshaler) marshalIntValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	var fieldBuf []byte
	switch value.Kind() {
	case reflect.Int8:
		fieldBuf = make([]byte, 1)
		fieldBuf[0] = byte(value.Int())
	case reflect.Int16:
		fieldBuf = make([]byte, 2)
		binary.BigEndian.PutUint16(fieldBuf, uint16(value.Int()))
	case reflect.Int32:
		fieldBuf = make([]byte, 4)
		binary.BigEndian.PutUint32(fieldBuf, uint32(value.Int()))
	case reflect.Int64, reflect.Int:
		fieldBuf = make([]byte, 8)
		binary.BigEndian.PutUint64(fieldBuf, uint64(value.Int()))
	default:
		return fmt.Errorf("unsupported int kind: %v", value.Kind())
	}

	return pfm.marshalByteFieldTo(conf, fieldBuf, w)
}

func (pfm *PrefixFramedMarshaler) marshalComplexValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	fieldBuf := make([]byte, 16)
	binary.BigEndian.PutUint64(fieldBuf[:8],
		math.Float64bits(real(value.Complex())))
	binary.BigEndian.PutUint64(fieldBuf[8:],
		math.Float64bits(imag(value.Complex())))

	return pfm.marshalByteFieldTo(conf, fieldBuf, w)
}

func (pfm *PrefixFramedMarshaler) unmarshalComplexValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	if len(fieldData) != 16 {
		return fmt.Errorf("invalid complex field size")
	}
	realBits := binary.BigEndian.Uint64(fieldData[:8])
	imagBits := binary.BigEndian.Uint64(fieldData[8:])
	value.SetComplex(complex(math.Float64frombits(realBits),
		math.Float64frombits(imagBits)))

	return nil
}

func (pfm *PrefixFramedMarshaler) marshalFloatValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	fieldBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(fieldBuf, math.Float64bits(value.Float()))

	return pfm.marshalByteFieldTo(conf, fieldBuf, w)
}

func (pfm *PrefixFramedMarshaler) unmarshalFloatValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	if len(fieldData) != 8 {
		return fmt.Errorf("invalid float field size")
	}
	bits := binary.BigEndian.Uint64(fieldData)
	value.SetFloat(math.Float64frombits(bits))

	return nil
}

func (pfm *PrefixFramedMarshaler) marshalBoolValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	fieldBuf := []byte{0}
	if value.Bool() {
		fieldBuf[0] = 1
	}

	return pfm.marshalByteFieldTo(conf, fieldBuf, w)
}

func (pfm *PrefixFramedMarshaler) unmarshalBoolValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	if len(fieldData) != 1 {
		return fmt.Errorf("invalid boolean field size")
	}
	value.SetBool(fieldData[0] != 0)

	return nil
}

func (pfm *PrefixFramedMarshaler) unmarshalStringValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	value.SetString(string(fieldData))

	return nil
}

func (pfm *PrefixFramedMarshaler) unmarshalUintValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	intVal := uint64(0)
	switch len(fieldData) {
	case 8:
		intVal = binary.BigEndian.Uint64(fieldData)
	case 4:
		intVal = uint64(binary.BigEndian.Uint32(fieldData))
	case 2:
		intVal = uint64(binary.BigEndian.Uint16(fieldData))
	case 1:
		intVal = uint64(fieldData[0])
	default:
		return fmt.Errorf("invalid uint field size")
	}

	// fieldValue, _ := binary.Uvarint(fieldData)
	// value.SetUint(fieldValue)
	value.SetUint(intVal)

	return nil
}

func (pfm *PrefixFramedMarshaler) unmarshalIntValueFrom(
	conf fieldConf,
	r Reader,
	value reflect.Value,
) error {
	fieldData, err := pfm.unmarshalByteFieldToBytes(r)
	if err != nil {
		return err
	}

	switch len(fieldData) {
	case 8:
		value.SetInt(int64(binary.BigEndian.Uint64(fieldData)))
	case 4:
		value.SetInt(int64(binary.BigEndian.Uint32(fieldData)))
	case 2:
		value.SetInt(int64(binary.BigEndian.Uint16(fieldData)))
	case 1:
		value.SetInt(int64(fieldData[0]))
	default:
		return fmt.Errorf("invalid int field size")
	}

	return nil
}

func (pfm *PrefixFramedMarshaler) marshalUintValueTo(
	conf fieldConf,
	value reflect.Value,
	w io.Writer,
) error {
	conf.Kind = value.Kind()
	return pfm.marshalUIntField(conf, value.Uint(), w)
}
