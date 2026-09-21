package datautil

import (
	// Built-in/core modules.
	"io"

	// First-party modules.
	"github.com/cuberat-go/datautil/internal/marshal"
)

// Object structure for prefix-framed marshaling of data.
//
// Datatypes are marshaled into bytes according to their kind, with support for
// structs, maps, slices, arrays, booleans, integers, unsigned integers,
// floats, strings, and complex numbers. Custom types implementing the
// encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interfaces are also
// supported.
type PrefixFramedMarshaler struct {
	marshaler *marshal.PrefixFramedMarshaler
}

// Returns a new prefix-framed data marshaler.
func NewPrefixFramedMarshaler() *PrefixFramedMarshaler {
	return &PrefixFramedMarshaler{
		marshaler: &marshal.PrefixFramedMarshaler{},
	}
}

// Marshals the provided data to the provided writer using prefix-framed
// encoding.
func (pfm *PrefixFramedMarshaler) MarshalTo(data any, w io.Writer) error {
	return pfm.marshaler.MarshalTo(data, w)
}

// Marshals the provided data and returns the resulting byte slice using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) Marshal(data any) ([]byte, error) {
	return pfm.marshaler.Marshal(data)
}

// Unmarshals data from the provided reader into the given value using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) UnmarshalFrom(
	r io.Reader,
	val any,
) error {
	return pfm.marshaler.UnmarshalFrom(r, val)
}

// Unmarshals data from the provided byte slice into the given value using
// prefix-framed encoding.
func (pfm *PrefixFramedMarshaler) Unmarshal(
	data []byte,
	val any,
) error {
	return pfm.marshaler.Unmarshal(data, val)
}
