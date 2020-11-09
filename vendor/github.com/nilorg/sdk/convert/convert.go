package convert

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// ToString 转换成string
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

// ToBool 转换成Bool
func ToBool(i interface{}) bool {
	switch b := i.(type) {
	case bool:
		return b
	case nil:
		return false
	case int:
		if i.(int) != 0 {
			return true
		}
		return false
	case string:
		v, err := strconv.ParseBool(ToString(i))
		if err != nil {
			return false
		}
		return v
	default:
		return false
	}
}

// ToInt 转换成int
func ToInt(value interface{}) int {
	return int(ToInt64(value))
}

// ToInt32 转换成int32
func ToInt32(value interface{}) int32 {
	return int32(ToInt64(value))
}

// ToInt64 转换成int64
func ToInt64(value interface{}) int64 {
	num, err := strconv.ParseInt(ToString(value), 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// ToUint 转换成uint
func ToUint(value interface{}) uint {
	return uint(ToUint64(value))
}

// ToUint32 转换成uint32
func ToUint32(value interface{}) uint32 {
	return uint32(ToUint64(value))
}

// ToUint64 转换成uint64
func ToUint64(value interface{}) uint64 {
	num, err := strconv.ParseUint(ToString(value), 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// ToFloat32 转换成float32
func ToFloat32(value interface{}) float32 {
	return float32(ToFloat64(value))
}

// ToFloat64 转换成float64
func ToFloat64(value interface{}) float64 {
	num, err := strconv.ParseFloat(ToString(value), 64)
	if err != nil {
		return 0
	}
	return num
}

// BytesToInt32 字节转Int32
func BytesToInt32(data []byte) int32 {
	var num int32
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &num)
	return num
}

// BytesToInt 字节转Int
func BytesToInt(data []byte) int {
	return int(BytesToInt32(data))
}

// BytesToInt64 字节转Int64
func BytesToInt64(data []byte) int64 {
	var num int64
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.BigEndian, &num)
	return num
}
