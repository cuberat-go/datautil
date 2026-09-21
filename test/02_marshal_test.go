package datautil_test

import (
	// Built-in/core modules.
	"fmt"
	"testing"

	// Third-party modules.
	"github.com/stretchr/testify/assert"

	// First-party modules.
	"github.com/cuberat-go/datautil"
)

func TestByteSlice(t *testing.T) {
	data := []byte("hello world")

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := []byte{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestUint(t *testing.T) {
	data := uint64(42)

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := uint64(0)
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestUintPointer(t *testing.T) {
	data := new(uint64(42))

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := new(uint64(0))
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestInt(t *testing.T) {
	data := int64(-42)

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := int64(0)
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestFloat(t *testing.T) {
	data := 3.14159

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := 0.0
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestComplex(t *testing.T) {
	data := complex(3.14, -2.71)

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := complex(0, 0)
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestBool(t *testing.T) {
	data := true

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := false
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}

	data = false

	wireData, err = pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData = true
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestString(t *testing.T) {
	data := "hello world"

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := ""
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestIntSlice(t *testing.T) {
	data := []int64{1, 2, 3, 4, 5}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := []int64{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestNilIntSlice(t *testing.T) {
	var data []int64

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := []int64{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStringSlice(t *testing.T) {
	data := []string{"hello", "world", "foo", "bar"}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := []string{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestIntArray(t *testing.T) {
	data := [5]int64{1, 2, 3, 4, 5}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := [5]int64{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestMap(t *testing.T) {
	data := map[string]int64{"one": 1, "two": 2}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var newData map[string]int64
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestMapOfMaps(t *testing.T) {
	data := map[string]map[string]int64{
		"one": {"a": 1, "b": 2},
		"two": {"c": 3, "d": 4},
	}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := map[string]map[string]int64{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestMapOfSlices(t *testing.T) {
	data := map[string][]int64{
		"one": {1, 2, 3},
		"two": {4, 5, 6},
	}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := map[string][]int64{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithInts(t *testing.T) {
	type MyStruct struct {
		A int64
		B int64
	}

	data := MyStruct{A: 42, B: -42}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	assert.NoError(t, err)
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithUints(t *testing.T) {
	type MyStruct struct {
		A uint64
		B uint64
	}

	data := MyStruct{A: 42, B: 84}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithUintPointers(t *testing.T) {
	type MyStruct struct {
		A *uint64
		B *uint64
	}

	data := MyStruct{A: new(uint64(42)), B: new(uint64(84))}
	ptrData := MyStruct{A: data.A, B: data.B}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(ptrData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, ptrData, newData) {
		t.FailNow()
	}
}

func TestStructWithBools(t *testing.T) {
	type MyStruct struct {
		A bool
		B bool
	}

	data := MyStruct{A: true, B: false}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithStrings(t *testing.T) {
	type MyStruct struct {
		A string
		B string
	}

	data := MyStruct{A: "hello", B: "world"}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithSlices(t *testing.T) {
	type MyStruct struct {
		A []int64
		B []string
	}

	data := MyStruct{
		A: []int64{1, 2, 3},
		B: []string{"hello", "world"},
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithArrays(t *testing.T) {
	type MyStruct struct {
		A [3]int64
		B [2]string
	}

	data := MyStruct{
		A: [3]int64{1, 2, 3},
		B: [2]string{"hello", "world"},
	}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithMaps(t *testing.T) {
	type MyStruct struct {
		A map[string]int64
		B map[int64]string
	}

	data := MyStruct{
		A: map[string]int64{"one": 1, "two": 2},
		B: map[int64]string{1: "one", 2: "two"},
	}

	pfm := datautil.NewPrefixFramedMarshaler()
	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithStructs(t *testing.T) {
	type InnerStruct struct {
		X int
		Y string
	}

	type MyStruct struct {
		A InnerStruct
		B InnerStruct
	}

	data := MyStruct{
		A: InnerStruct{X: 1, Y: "one"},
		B: InnerStruct{X: 2, Y: "two"},
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithNullPointers(t *testing.T) {
	type InnerStruct struct {
		X int
		Y string
	}

	type MyStruct struct {
		A *InnerStruct
		B *InnerStruct
	}

	data := MyStruct{
		A: &InnerStruct{X: 1, Y: "one"},
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := MyStruct{}
	err = pfm.Unmarshal(wireData, &newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithNilSlice(t *testing.T) {
	type MyStruct struct {
		A []int64
		B []string
	}
	data := &MyStruct{
		A: []int64{42},
		B: nil,
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := &MyStruct{}
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithNilSlices(t *testing.T) {
	type MyStruct struct {
		A []int64
		B []string
	}

	data := &MyStruct{
		A: nil,
		B: nil,
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := &MyStruct{}
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithNilMap(t *testing.T) {
	type MyStruct struct {
		A map[string]int
		B map[string]string
	}

	data := &MyStruct{
		A: nil,
		B: nil,
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := &MyStruct{}
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

type CustomType struct {
	Num1   int
	Num2   int
	Float1 float64
}

func (c *CustomType) MarshalBinary() ([]byte, error) {
	return fmt.Appendf(nil, "%d,%d,%f", c.Num1, c.Num2, c.Float1), nil
}

func (c *CustomType) UnmarshalBinary(data []byte) error {
	_, err := fmt.Sscanf(string(data), "%d,%d,%f", &c.Num1, &c.Num2,
		&c.Float1)
	return err
}

func TestStructWithCustomMarshaler(t *testing.T) {

	type MyStruct struct {
		A *CustomType
		B *CustomType
	}

	data := &CustomType{
		Num1:   1,
		Num2:   2,
		Float1: 1.1,
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := &CustomType{}
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}

func TestStructWithCustomMarshalerFields(t *testing.T) {
	type MyStruct struct {
		A *CustomType
		B *CustomType
	}

	data := &MyStruct{
		A: &CustomType{
			Num1:   1,
			Num2:   2,
			Float1: 1.1,
		},
		B: &CustomType{
			Num1:   3,
			Num2:   4,
			Float1: 2.2,
		},
	}

	pfm := datautil.NewPrefixFramedMarshaler()

	wireData, err := pfm.Marshal(data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	newData := &MyStruct{}
	err = pfm.Unmarshal(wireData, newData)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	if !assert.Equal(t, data, newData) {
		t.FailNow()
	}
}
