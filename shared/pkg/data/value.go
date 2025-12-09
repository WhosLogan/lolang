package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"shared/pkg/types"
)

type Value struct {
	Type types.TypeCode
	Data []byte
}

func NewValue(v any) (Value, error) {
	switch x := v.(type) {
	case int64:
		data := make([]byte, 8)
		binary.BigEndian.PutUint64(data, uint64(x))
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
	}

	return Value{}, errors.New("invalid value type")
}

func (v *Value) GetInt() int64 {
	return int64(v.Data[0])
}

func (v *Value) GetBool() bool {
	return v.Data[0] == 1
}

func (v *Value) Equal(other *Value) bool {
	if v.Type != other.Type || len(v.Data) != len(other.Data) {
		return false
	}

	return bytes.Equal(v.Data, other.Data)
}
