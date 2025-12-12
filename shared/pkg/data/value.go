package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"shared/pkg/types"
)

type Value struct {
	Type   types.TypeCode
	Data   []byte
	Fields map[string]Value
}

func NewValue(v any) (Value, error) {
	switch x := v.(type) {
	case int:
		return NewValue(int64(x))
	case int32:
		return NewValue(int64(x))
	case int64:
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(x))
		return Value{
			Type: types.LoInt,
			Data: data,
		}, nil
	case bool:
		data := make([]byte, 1)
		if x {
			data[0] = 1
		}

		return Value{
			Type: types.LoBool,
			Data: data,
		}, nil
	case string:
		return Value{
			Type: types.LoString,
			Data: []byte(x),
		}, nil
	}

	return Value{}, errors.New("invalid value type")
}

func NewStructValue(typeCode types.TypeCode) Value {
	return Value{
		Type:   typeCode,
		Fields: make(map[string]Value),
	}
}

func (v *Value) SetField(name string, val Value) {
	if v.Fields == nil {
		v.Fields = make(map[string]Value)
	}
	v.Fields[name] = val
}

func (v *Value) GetField(name string) (Value, bool) {
	if v.Fields == nil {
		return Value{}, false
	}
	val, ok := v.Fields[name]
	return val, ok
}

func MustNewValue(v any) Value {
	val, err := NewValue(v)
	if err != nil {
		panic(err)
	}

	return val
}

func (v *Value) GetInt() int64 {
	return int64(binary.LittleEndian.Uint64(v.Data))
}

func (v *Value) GetBool() bool {
	return v.Data[0] == 1
}

func (v *Value) GetString() string {
	return string(v.Data)
}

func (v *Value) Equal(other *Value) bool {
	if (v.Type == types.LoBool && other.Type == types.LoInt) || (v.Type == types.LoInt && other.Type == types.LoBool) {
		return v.Data[0] == other.Data[0]
	}

	if v.Type != other.Type || len(v.Data) != len(other.Data) {
		return false
	}

	return bytes.Equal(v.Data, other.Data)
}
