package datautil

import (
	// Built-in/core modules.

	"io"
	"iter"

	// First-party modules.
	"github.com/cuberat-go/datautil/internal/marshal"
)

// Object structure for length-prefix framed marshaling of data.
//
// Datatypes are marshaled into bytes according to their kind, with support for
// structs, maps, slices, arrays, booleans, integers, unsigned integers,
// floats, strings, and complex numbers. Custom types implementing the
// encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interfaces are also
// supported. The length of each data element is prefixed to its byte
// representation, in the form of a varint, enabling efficient parsing and
// unmarshaling.
//
// ### Structs
//
// A struct value is marshaled by prefixing the length of its byte
// representation, followed by the serialized fields of the struct, in field
// order, excluding unexported fields. Note this means that if you add new
// fields in the middle of the struct, it will affect backward compatibility.
// Adding a field to the end of a struct is generally safe and maintains
// backward compatibility.
//
// #### Struct Field Tags

// Struct field tags can be used to customize the marshaling behavior of
// individual fields. Use `datautilpfm` as the key, e.g.,
//
//	type MyStruct struct {
//	    A uint64 `datautilpfm:"varint"`
//	}
//
// Valid values for the `datautilpfm` tag:
// - "varint": Indicates that an unsigned integer (uint, uint8,uint16, uint32, uint64) field should be marshaled using varint encoding.
//
// ### Maps
//
// A map value is marshaled by prefixing the length of its byte
// representation, followed by the serialized key-value pairs. The order of
// key-value pairs when serialized depends on the underlying map implementation.
//
// ### Slices and Arrays
//
// A slice or array value is marshaled by prefixing the length of its byte
// representation, followed by the serialized elements in order.
// ### Booleans
//
// A boolean value is marshaled by prefixing the length of its byte
// representation, followed by a single byte representing the boolean value
// (0 for false, 1 for true).
//
// ### Integers
//
// An integer value is marshaled by prefixing the length of its byte
// representation, followed by the serialized integer in big-endian order.
//
// ### Unsigned Integers
//
// An unsigned integer value is marshaled by prefixing the length of its byte
// representation, followed by the serialized unsigned integer in big-endian order.
//
// ### Floats
//
// A float value is marshaled by prefixing the length of its byte
// representation, followed by the serialized float in big-endian order.
//
// ### Strings
//
// A string value is marshaled by prefixing the length of its byte
// representation, followed by the serialized string bytes.
//
// ### Complex Numbers
//
// A complex number value is marshaled by prefixing the length of its byte
// representation, followed by the serialized real and imaginary parts in
// big-endian order.
//
// ### BinaryMarshaler
//
// A value implementing the BinaryMarshaler interface is marshaled by
// prefixing the length of its byte representation, followed by the serialized
// bytes returned by the MarshalBinary method.
type PrefixFramedMarshaler struct {
	marshaler *marshal.PrefixFramedMarshaler
}

// Returns a new prefix-framed data marshaler.
func NewPrefixFramedMarshaler() *PrefixFramedMarshaler {
	return &PrefixFramedMarshaler{
		marshaler: &marshal.PrefixFramedMarshaler{},
	}
}

// Marshals the provided data to the provided writer using length-prefix framed
// encoding.
func (pfm *PrefixFramedMarshaler) MarshalTo(data any, w io.Writer) error {
	return pfm.marshaler.MarshalTo(data, w)
}

// Marshals the provided data and returns the resulting byte slice using
// length-prefix framing.
func (pfm *PrefixFramedMarshaler) Marshal(data any) ([]byte, error) {
	return pfm.marshaler.Marshal(data)
}

// Unmarshals data from the provided reader into the given value using
// length-prefix framing.
func (pfm *PrefixFramedMarshaler) UnmarshalFrom(
	r io.Reader,
	val any,
) error {
	return pfm.marshaler.UnmarshalFrom(r, val)
}

// Unmarshals data from the provided byte slice into the given value using
// length-prefix framing.
func (pfm *PrefixFramedMarshaler) Unmarshal(
	data []byte,
	val any,
) error {
	return pfm.marshaler.Unmarshal(data, val)
}

// Returns a sequence iterator for reading multiple values of the specified type from the provided reader using length-prefix framing.
func (pfm *PrefixFramedMarshaler) SeqFrom[T any](
	template T,
	r io.Reader,
) iter.Seq2[T, error] {
	return pfm.marshaler.SeqFrom(template, r)
}
