package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

// Ternary Realize ternary operations
func ternary[T interface{}](operation bool, a, b T) T {
	if operation {
		return a
	}
	return b
}

// Tap Call the given Closure with the given value then return the value.
func tap[T interface{}](value T, callback ...func(value T)) T {
	if len(callback) > 0 {
		callback[0](value)
	}

	return value
}

// VarDef Obtain the first value of a variable parameter or give a default value
func varDef[T interface{}](value []T, def ...T) T {
	if len(value) > 0 {
		return value[0]
	}
	if len(def) > 0 {
		return def[0]
	}
	var d T
	return d
}

// InArray Determine whether a one-dimensional array contains an element
func inArray[V comparable, T []V](needle V, haystack T) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}

	return false
}

// Unique Removing duplicate values from an array
func unique[V comparable, T []V](arr T) T {
	var (
		size   = len(arr)
		result = make(T, 0, size)
	)

	for i := 0; i < size; i++ {
		if !inArray(arr[i], result) {
			result = append(result, arr[i])
		}
	}

	return result
}

// filter 数组过滤
func filter[V interface{}, T []V](arr T, callback func(v V) bool) T {
	var items T

	for _, v := range arr {
		if callback(v) {
			items = append(items, v)
		}
	}

	return items
}

// Map 遍历修改数组
func arrMap[V interface{}, T []V](arr T, callback func(v V) V) T {
	var items T

	for _, v := range arr {
		items = append(items, callback(v))
	}

	return items
}

// ReplaceByArray returns a copy of `origin`,
// which is replaced by a slice in order, case-sensitively.
func replaceByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = strings.Replace(origin, array[i], array[i+1], -1)
	}
	return origin
}

// String converts `any` to string.
// It's most commonly used converting function.
func convString(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return convString(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

// Trim strips whitespace (or other characters) from the beginning and end of a string.
// The optional parameter `characterMask` specifies the additional stripped characters.
func trim(str string, characterMask ...string) string {
	trimChars := varDef(characterMask, DefaultCharset)
	return strings.Trim(str, trimChars)
}

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
func ucFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if isLetterLower(s[0]) {
		return string(s[0]-32) + s[1:]
	}
	return s
}

// isLetterLower checks whether the given byte b is in lower case.
func isLetterLower(b byte) bool {
	if b >= byte('a') && b <= byte('z') {
		return true
	}
	return false
}
