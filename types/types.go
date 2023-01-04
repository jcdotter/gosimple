// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package types

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// paramTypeError formats and returns an error
// using a template for type mismatches on func params
// function: the name of the function
// typ: the expected types for the param
// value: the provided value for the param
func paramTypeError(function string, typ string, value any) error {
	return typeError(function, "  expected %v type,\n  received %T type", typ, value)
}

func typeError(function string, format string, a ...any) error {
	return fmt.Errorf("failed call to utils.types.%v:\n"+format, append([]any{function}, a...)...)
}

// mustBe validates that 'a' must be of Type 't'
// if not, it throughs an error
func mustBe(a any, t ...Type) error {
	typeA := TypeOf(a)
	var vTypes string
	for i, typeT := range t {
		if typeA == typeT {
			return nil
		}
		vTypes += " "
		if len(t) > 1 && i == len(t)-1 {
			vTypes += "or "
		}
		vTypes += typeNames[typeT]
	}
	return paramTypeError(source(), vTypes, a)
}

func source() string {
	pc, _, _, _ := runtime.Caller(2)
	fnc := runtime.FuncForPC(pc)
	_, ln := fnc.FileLine(pc)
	return fmt.Sprintf("%s line:%v", fnc.Name(), ln)
}

// Type represents the abstract data type
// TYPE ABSTRACTIONS:
// Type: String, Bool, Int, Float, Uint, Time, Array, Map, Struct
// String: string
// Bool: bool
// Int: int, int8, int16, int32, int64
// Float: float32, float64
// Uint: uint, uint8, uint16, uint32, uint64
// Number: int, float, uint
// Basic: string, int, float, uint, bool
// Array: array, slice
// Map: map
// Struct: struct
// Value: string, int, float, uint, bool, time, slice, map

type Type uint

const (
	Invalid Type = iota
	String
	Int
	Float
	Uint
	Bool
	Time
	Array
	Map
	Struct
	Func
	Ptr
	Any
)

var typeNames = []string{
	Invalid: "invalid",
	String:  "string",
	Int:     "int",
	Float:   "float",
	Uint:    "uint",
	Bool:    "bool",
	Time:    "time",
	Array:   "array",
	Map:     "map",
	Struct:  "struct",
	Func:    "function",
	Ptr:     "pointer",
	Any:     "any",
}

func (t Type) String() string {
	return typeNames[uint(t)]
}

// TypeOf returns the abstracted data type of 'a':
//   String, Bool, Int, Uint, Float,
//   Time, Array, Map or Struct
func TypeOf(a any) Type {
	switch {
	case IsString(a):
		return String
	case IsBool(a):
		return Bool
	case IsInt(a):
		return Int
	case IsFloat(a):
		return Float
	case IsUint(a):
		return Uint
	case IsTime(a):
		return Time
	case IsArray(a):
		return Array
	case IsMap(a):
		return Map
	case IsStruct(a):
		return Struct
	case IsPtr(a):
		return Ptr
	case IsFunc(a):
		return Func
	default:
		return Invalid
	}
}

// TypeByName returns the Type using the string name of the type
func TypeByName(s string) (Type, error) {
	s = strings.ToLower(s)
	found := false
	typ := Invalid
	for i, v := range typeNames {
		if v == s {
			found = true
			typ = Type(i)
			break
		}
	}
	if found {
		return typ, nil
	} else {
		return typ, paramTypeError("TypeByName", "string of a valid", s)
	}
}

// To converts any value 'a' to type 't' by returning
// a single value map with a key of 't', and
// panic if can't convert 'a' to type 't'
// convertables types include String, Bool, Int, Float, Uint, Time
// example: To(t,a)[t]
func To(t Type, a any) (map[Type]any, error) {
	m := map[Type]any{}
	var err error
	switch t {
	case String:
		m[String], err = ToString(a)
		break
	case Bool:
		m[Bool], err = ToBool(a)
		break
	case Int:
		m[Int], err = ToInt(a)
		break
	case Float:
		m[Float], err = ToFloat(a)
		break
	case Uint:
		m[Uint], err = ToUint(a)
		break
	case Time:
		m[Time], err = ToTime(a)
		break
	default:
		err = fmt.Errorf("")
	}
	if err != nil {
		return map[Type]any{}, typeError("To", " could not convert type %T to %s", a, t)
	}
	return m, err
}

// To converts any value 'a' to type 't' by returning
// a single value map with a key of 't', and
// panic if can't convert 'a' to type 't'
// convertable types include String, Bool, Int#, Float#, Uint#
// example: StrictlyTo(t,a)[reflect.TypeOf(t).Kind()]
func StrictlyTo(t any, a any) (map[reflect.Kind]any, error) {
	m := map[reflect.Kind]any{}
	var err error
	switch TypeOf(t) {
	case String:
		m[reflect.String], err = ToString(a)
		break
	case Bool:
		m[reflect.Bool], err = ToBool(a)
		break
	case Int, Float, Uint:
		v, err := ToFloat(a)
		if err != nil {
			break
		}
		k := reflect.TypeOf(t).Kind()
		if ConversionOverflow(k, a) {
			err = typeError("To", " overflow error")
			break
		}
		switch k {
		case reflect.Int:
			m[reflect.Int] = int(v)
			break
		case reflect.Int8:
			m[reflect.Int8] = int8(v)
			break
		case reflect.Int16:
			m[reflect.Int16] = int16(v)
			break
		case reflect.Int32:
			m[reflect.Int32] = int32(v)
			break
		case reflect.Int64:
			m[reflect.Int64] = int64(v)
			break
		case reflect.Uint:
			m[reflect.Int] = uint(v)
			break
		case reflect.Uint8:
			m[reflect.Uint8] = uint8(v)
			break
		case reflect.Uint16:
			m[reflect.Uint16] = uint16(v)
			break
		case reflect.Uint32:
			m[reflect.Int32] = uint32(v)
			break
		case reflect.Uint64:
			m[reflect.Int64] = uint64(v)
			break
		case reflect.Float32:
			m[reflect.Float32] = float32(v)
			break
		case reflect.Float64:
			m[reflect.Float64] = float64(v)
			break
		}
		break
	default:
		err = fmt.Errorf("")
	}
	if err != nil || len(m) == 0 {
		return map[reflect.Kind]any{}, typeError("To", " could not convert type %T to %s", a, t)
	}
	return m, nil
}

// TypeOverflowLimit returns the value limit for numeric types
// 't' is the reflect package Kind of the numeric type
// returns and error if the Kind is not numeric
func TypeOverflowLimit(t reflect.Kind) (float64, error) {
	l := map[reflect.Kind]float64{
		reflect.Int:     float64(math.MaxInt),
		reflect.Int8:    float64(math.MaxInt8),
		reflect.Int16:   float64(math.MaxInt16),
		reflect.Int32:   float64(math.MaxInt32),
		reflect.Int64:   float64(math.MaxInt64),
		reflect.Uint:    float64(math.MaxUint),
		reflect.Uint8:   float64(math.MaxUint8),
		reflect.Uint16:  float64(math.MaxUint16),
		reflect.Uint32:  float64(math.MaxUint32),
		reflect.Uint64:  float64(math.MaxUint64),
		reflect.Float32: float64(math.MaxFloat32),
		reflect.Float64: float64(math.MaxFloat64),
	}
	r, ok := l[t]
	if !ok {
		return 0, fmt.Errorf("not a numberic value type")
	}
	return r, nil
}

// ValueOverflowLimit returns the value limit for numeric types
// 'a' is a numeric value from which the type is determined
// returns and error if 'a' is not numeric
func ValueOverflowLimit(a any) (float64, error) {
	return TypeOverflowLimit(reflect.TypeOf(a).Kind())
}

// ConversionOverflow evaluates whether 'a' will overflow
// if converted to type 't', which is
// the reflect.Kind of a data type
// returns true if value is not convertable
func ConversionOverflow(t reflect.Kind, a any) bool {
	f, fErr := ToFloat(a)
	tLim, tErr := TypeOverflowLimit(t)
	if fErr != nil || tErr != nil {
		return true
	}
	return f > tLim
}

// Abstract type assertions validate whether val is an abstract type
// Abstract types include: String, Int, Float, Uint, Number, Basic, Value
// IsString:	validates whether String
// IsBool:		validates whether Bool
// IsInt:		validates whether Int (any type of int)
// IsFloat:		validates whether Float (any type of float)
// IsUint:		validates whether Uint (any type of uint)
// IsNumber: 	validates whether Number (int, float or uint)
// IsBasic: 	validates whether Basic (string, bool or Number)
// IsTime:		validates whether Time (time.Time)
// IsArray:		validates whether Array (array or slice)
// IsMap: 		validates whether Map (any type of map)
// IsStruct:	validates whether Struct (any typ of struct)
// IsValue 		validates whether Value (Basic, time.Time, array or map)

// IsString evaluates whether 'a' is a string
func IsString(a any) bool {
	if _, ok := a.(string); ok {
		return true
	}
	return false
}

// IsBool evaluates whether 'a' is a bool
func IsBool(a any) bool {
	if _, ok := a.(bool); ok {
		return true
	}
	return false
}

// IsInt evaluates whether 'a' is an Int:
//   int, int8, int16, int32 or int64
func IsInt(a any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

// IsFloat evaluates whether 'a' is a Float:
//   float32 or float64
func IsFloat(a any) bool {
	switch a.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// IsUint evaluates whether 'a' is a Uint:
//   uint, uint8, uint16, uint32 or uint64
func IsUint(a any) bool {
	switch a.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsNumber evaluates whether 'a' is a Number:
//   int, int8, int16, int32, int64,
//   float32, float64,
//   uint, uint8, uint16, uint32 or uint64
func IsNumber(a any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64,
		float32, float64,
		uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsBasic evaluates whether 'a' is a basic type:
//   string, bool,
//   int, int8, int16, int32, int64,
//   float32, float64,
//   uint, uint8, uint16, uint32 or uint64
func IsBasic(a any) bool {
	switch a.(type) {
	case string,
		int, int8, int16, int32, int64,
		float32, float64,
		uint, uint8, uint16, uint32, uint64,
		bool:
		return true
	default:
		return false
	}
}

// IsTime evaluates whether 'a' is a time:
//   time.Time
func IsTime(a any) bool {
	if _, ok := a.(time.Time); ok {
		return true
	}
	return false
}

// IsArray evaluates whether 'a' is an array or slice:
//   [#]any or []any
func IsArray(a any) bool {
	t := reflect.TypeOf(a).Kind()
	return t == reflect.Array || t == reflect.Slice
}

// IsMap evaluates whether 'a' is a map:
//   map[any]any
func IsMap(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Map
}

// IsStruct evaluates whether 'a' is a struct:
//   type Struct struct {}
func IsStruct(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Struct
}

// IsValue evaluates whether 'a' is a single value:
//   string, bool,
//   int, int8, int16, int32, int64,
//   float32, float64,
//   uint, uint8, uint16, uint32, uint64
//   time.Time, slice, map, or struct
func IsValue(a any) bool {
	if IsBasic(a) || IsTime(a) || IsArray(a) || IsMap(a) || IsStruct(a) {
		return true
	}
	return false
}

// IsPtr evaluates whether 'a' is a pointer
func IsPtr(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Pointer
}

// IsFunc evaluates whether 'a' is a function
func IsFunc(a any) bool {
	return reflect.TypeOf(a).Kind() == reflect.Func
}

func IsEmpty(a any) bool {
	if a == nil {
		return true
	}
	switch {
	case IsTime(a):
		return a.(time.Time) == time.Time{}
	case IsArray(a) || IsMap(a):
		return reflect.ValueOf(a).Len() == 0
	case IsPtr(a):
		return IsEmpty(&a)
	default:
		return reflect.ValueOf(a).IsZero()
	}
}

// Equal evaluates whether types of 'x' and 'y' are the same
// the types are strict go types, and not abstract Types
func Equal(x any, y any) bool {
	return fmt.Sprintf("%T", x) == fmt.Sprintf("%T", y)
}

// EqualTypeValues evaluates whether types and values of 'x' and 'y' are the same
// the types are strict go types, and not abstract Types
// the values of arrays, maps and structs are evaluated deeply
func EqualTypeValues(x any, y any) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

// EqualValues evaluates whether values of 'x' and 'y' are loosely the same
// types are ignored in the evaluation (ie. "1" == 1)
// the values of arrays, maps and structs are evaluated deeply
func EqualValues(x any, y any) bool {
	return fmt.Sprintf("%v", x) == fmt.Sprintf("%v", y)
}

// STRING CONVERSION FUNCTIONS
// ToString:		converts any basic type to string	ALTERNATIVE: fmt.Sprint()
// StringToString 	converts any string to a string		ALTERNATIVE: string()
// IntToString: 	converts any int type to string 	ALTERNATIVE: fmt.Sprint()
// FloatToString:	converts any float type to string	ALTERNATIVE: fmt.Sprint()
// UintToString: 	converts any uint type to string 	ALTERNATIVE: fmt.Sprint()
// BoolToString: 	converts a bool to string 			ALTERNATIVE: fmt.Sprint()
// TimeToString:	converts a time.Time to string 		ALTERNATIVE: fmt.Sprint()
// ArrayToString:	converts an array type to string 	ALTERNATIVE: fmt.Sprint()
// MapToString:	converts a map type to string 		ALTERNATIVE: fmt.Sprint()
// StructToString:	converts a struct type to string	ALTERNATIVE: fmt.Sprint()

// StringToString converts a string to asserted string
// Equivilant to string(s)
// Returns error if param 's' type is not string
func StringToString(s any) (string, error) {
	if IsString(s) {
		return string(s.(string)), nil
	} else {
		return "", paramTypeError("StringToString", "string", s)
	}
}

// IntToString converts an int to string
// Equivilant to fmt.Sprint(i)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func IntToString(i any) (string, error) {
	if IsInt(i) {
		return fmt.Sprint(i), nil
	}
	return "", paramTypeError("IntToString", "int, int8, int16, int32 or int64", i)
}

// FloatToString converts a float to string
// Equivilant to fmt.Sprint(f)
// Returns error if param 'f' type is not float32 or float64
func FloatToString(f any) (string, error) {
	if IsFloat(f) {
		return fmt.Sprint(f), nil
	}
	return "", paramTypeError("FloatToString", "float32 or float64", f)
}

// UintToString converts a uint to string
// Equivilant to fmt.Sprint(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func UintToString(u any) (string, error) {
	if IsUint(u) {
		return fmt.Sprint(u), nil
	}
	return "", paramTypeError("UintToString", "uint, uint8, uint16, uint32 or uint64", u)
}

// BoolToString converts a bool to string
// Equivilant to fmt.Sprint(b)
// Returns error if param 'b' type is not bool
func BoolToString(b any) (string, error) {
	if IsBool(b) {
		return fmt.Sprint(b), nil
	} else {
		return "", paramTypeError("BoolToString", "bool", b)
	}
}

// TimeToString converts a time to string
// Equivilant to fmt.Sprint(t)
// Returns error if param 't' type is not time.Time
func TimeToString(t any) (string, error) {
	if IsTime(t) {
		return fmt.Sprint(t.(time.Time).UTC()), nil
	} else {
		return "", paramTypeError("TimeToString", "time.Time", t)
	}
}

// ArrayToString converts an array or slice to string
// Equivilant to fmt.Sprint(a)
// Returns error if param 'a' type is not array or slice
func ArrayToString(a any) (string, error) {
	if IsArray(a) {
		return fmt.Sprint(a), nil
	} else {
		return "", paramTypeError("ArrayToString", "array or slice", a)
	}
}

// MapToString converts a map to string
// Equivilant to fmt.Sprint(m)
// Returns error if param 'm' type is not time.Time
func MapToString(m any) (string, error) {
	if IsMap(m) {
		return fmt.Sprint(m), nil
	} else {
		return "", paramTypeError("MapToString", "map", m)
	}
}

// ToString converts param 'a' of a basic type to string
// Equivilant to fmt.Sprint(i)
// Returns error if param 'a' type is not
// string, int, float, uint, bool, time, slice, map or struct
func ToString(a any) (string, error) {
	if IsValue(a) {
		return fmt.Sprint(a), nil
	}
	return "", paramTypeError("ToString", "string, int, float, uint, bool, time, slice, map or struct", a)
}

// INT CONVERSION FUNCTIONS
// ToInt: 		converts any basic type to int 		CAUTION: performance
// StringToInt: 	converts a string to int 			ALTERNATIVE: strconv.Atoi()
// IntToInt: 	converts any int type to int 		ALTERNATIVE: int()
// FloatToInt:	converts any float type to int		ALTERNATIVE: int()
// UintToInt: 	converts any uint type to int 		ALTERNATIVE: int()
// BoolToInt: 	converts a bool to int 				ALTERNATIVE: none
// TimeToInt:	converts a time.Time to int 		ALTERNATIVE: [time].Unix()

// StringToInt converts a numeric string to a rounded int
// Similar to int(math.Round(strconv.ParseFloat(s, 64)))
// Returns error if param 's' type is not string
// or can't be converted to int
func StringToInt(s any) (int, error) {
	f, err := StringToFloat(s)
	if err != nil {
		return 0, paramTypeError("StringToInt", "numeric like string", s)
	}
	if ConversionOverflow(reflect.Int, f) {
		return 0, typeError("StringToInt", " overflow error")
	}
	return int(math.Round(f)), nil
}

// IntToInt converts any int type to int
// Equivilant to int(i)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func IntToInt(i any) (int, error) {
	switch ii := i.(type) {
	case int:
		return int(ii), nil
	case int8:
		return int(ii), nil
	case int16:
		return int(ii), nil
	case int32:
		return int(ii), nil
	case int64:
		return int(ii), nil
	default:
		return 0, paramTypeError("IntToInt", "int", i)
	}
}

// FloatToInt converts any float type to rounded int
// Equivilant to int(math.Round(f))
// Returns error if param 'f' type is not float32, float64
func FloatToInt(f any) (int, error) {
	if !IsFloat(f) {
		return 0, paramTypeError("FloatToInt", "float", f)
	}
	if ConversionOverflow(reflect.Int, f) {
		return 0, typeError("FloatToInt", " overflow error")
	}
	switch ff := f.(type) {
	case float32:
		return int(math.Round(float64(ff))), nil
	case float64:
		return int(math.Round(ff)), nil
	default:
		return 0, paramTypeError("FloatToInt", "float", f)
	}
}

// UintToInt converts any uint type to int
// Equivilant to int(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func UintToInt(u any) (int, error) {
	if !IsUint(u) {
		return 0, paramTypeError("UintToInt", "uint", u)
	}
	if ConversionOverflow(reflect.Int, u) {
		return 0, typeError("UintToInt", " overflow error")
	}
	switch uu := u.(type) {
	case uint:
		return int(uu), nil
	case uint8:
		return int(uu), nil
	case uint16:
		return int(uu), nil
	case uint32:
		return int(uu), nil
	case uint64:
		return int(uu), nil
	default:
		return 0, paramTypeError("UintToInt", "uint", u)
	}
}

// BoolToInt converts a bool to int
// 1 if true, 0 if false
// Returns error if param 'b' type is not bool
func BoolToInt(b any) (int, error) {
	switch bb := b.(type) {
	case bool:
		if bb {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, paramTypeError("BoolToInt", "bool", b)
	}
}

// TimeToInt converts a time.Time to a unix int
// Equivilant to int(t.Unix())
// Returns error if param 't' type is not time.Time
func TimeToInt(t any) (int, error) {
	switch tt := t.(type) {
	case time.Time:
		return int(tt.Unix()), nil
	default:
		return 0, paramTypeError("TimeToInt", "time.Time", t)
	}
}

// ToInt converts param 'a' of a basic type to int
// Returns error if param 'a' type is not:
//   string, int, float, uint, bool or time
func ToInt(a any) (int, error) {
	switch a.(type) {
	case string:
		return StringToInt(a)
	case int, int8, int16, int32, int64:
		return IntToInt(a)
	case float32, float64:
		return FloatToInt(a)
	case uint, uint8, uint16, uint32, uint64:
		return UintToInt(a)
	case bool:
		return BoolToInt(a)
	case time.Time:
		return TimeToInt(a)
	default:
		return 0, paramTypeError("ToInt", "string, numeric, bool, or time", a)
	}
}

// FLOAT COMVERSION FUNCTIONS
// ToFloat:			converts any basic type to float
// StringToFloat:		converts a string to float			ALTERNATIVE: strconv.ParseFloat(str, 64)
// IntToFloat:		converts any int type to float		ALTERNATIVE: float64()
// FloatToFloat:	converts any float type to float	ALTERNATIVE: float64()
// UnitToFloat:		converts any uint type to float		ALTERNATIVE: float64()
// BoolToFloat:		converts a bool to float			ALTERNATIVE: none
// TimeToFloat:		converts a time to a float			ALTERNATIVE: none
// CurrencyToFloat:	converts a currency to float 		ALTERNATIVE: none

// StringToFloat converts a numeric string to float64
// Similar to strconv.ParseFloat(str, 64)
// Returns error if param 'str' type is not string
// or can't be converted to float64
func StringToFloat(s any) (float64, error) {
	if !IsString(s) {
		return 0, paramTypeError("StringToFloat", "string", s)
	}
	str := strings.ReplaceAll(string(s.(string)), ",", "")
	if str[0] == 40 && str[len(str)-1] == 41 {
		// if first char is '(' and last ')'
		str = "-" + string(str[1:len(str)-1])
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, paramTypeError("StringToFloat", "numeric string", s)
	}
	return f, nil
}

// IntToFloat converts any int type to float64
// Equivilant to float64(i)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func IntToFloat(i any) (float64, error) {
	switch ii := i.(type) {
	case int:
		return float64(ii), nil
	case int8:
		return float64(ii), nil
	case int16:
		return float64(ii), nil
	case int32:
		return float64(ii), nil
	case int64:
		return float64(ii), nil
	default:
		return 0, paramTypeError("IntToFloat", "int", i)
	}
}

// FloatToInt converts any float type to asserted float64
// Equivilant to float64(f)
// Returns error if param 'f' type is not float32, float64
func FloatToFloat(f any) (float64, error) {
	switch ff := f.(type) {
	case float32:
		return float64(ff), nil
	case float64:
		return float64(ff), nil
	default:
		return 0, paramTypeError("FloatToFloat", "float", f)
	}
}

// UintToFloat converts any uint type to float64
// Equivilant to float64(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func UintToFloat(u any) (float64, error) {
	switch uu := u.(type) {
	case uint:
		return float64(uu), nil
	case uint8:
		return float64(uu), nil
	case uint16:
		return float64(uu), nil
	case uint32:
		return float64(uu), nil
	case uint64:
		return float64(uu), nil
	default:
		return 0, paramTypeError("UintToFloat", "uint", u)
	}
}

// BoolToFloat converts a bool to float64
// 1 if true, 0 if false
// Returns error if param 'b' type is not bool
func BoolToFloat(b any) (float64, error) {
	switch bb := b.(type) {
	case bool:
		if bb {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, paramTypeError("BoolToFloat", "bool", b)
	}
}

// TimeToFloat converts a time.Time to a unix float64
// Equivilant to float64(t.Unix())
// Returns error if param 't' type is not time.Time
func TimeToFloat(t any) (float64, error) {
	switch tt := t.(type) {
	case time.Time:
		return float64(tt.Unix()), nil
	default:
		return 0, paramTypeError("TimeToFloat", "time.Time", t)
	}
}

// ToFloat converts param 'a' of a basic type to float64
// Returns error if param 'a' type is not:
//   string, int, float, uint, bool or time
func ToFloat(a any) (float64, error) {
	switch a.(type) {
	case string:
		return StringToFloat(a)
	case int, int8, int16, int32, int64:
		return IntToFloat(a)
	case float32, float64:
		return FloatToFloat(a)
	case uint, uint8, uint16, uint32, uint64:
		return UintToFloat(a)
	case bool:
		return BoolToFloat(a)
	case time.Time:
		return TimeToFloat(a)
	default:
		return 0, paramTypeError("ToFloat", "string, numeric, bool, or time", a)
	}
}

// UINT INVERSION FUNCTIONS
// ToUint:			converts any basic type to unit
// StringToUint:		converts a string to unit			ALTERNATIVE: uint(strconv.ParseFloat())
// IntToUint:		converts any int type to unit		ALTERNATIVE: uint()
// FloatToUint:		converts any float type to unit		ALTERNATIVE: uint()
// UnitToUint:		converts any uint type to unit		ALTERNATIVE: uint()
// BoolToUint:		converts a bool to unit				ALTERNATIVE: none
// TimeToUint:		converts a time to a unit			ALTERNATIVE: none
// CurrencyToUint:	converts a currency to unit 		ALTERNATIVE: none

// StringToUint converts a numeric string to rounded uint
// Similar to uint(strconv.ParseFloat(str, 64))
// Returns error if param 's' type is not string
// or can't be converted to unit
func StringToUint(s any) (uint, error) {
	f, err := StringToFloat(s)
	if err != nil || f < 0 {
		return 0, paramTypeError("StringToUint", "unsigned numeric string", s)
	}
	if ConversionOverflow(reflect.Uint, f) {
		return 0, typeError("StringToUint", " overflow error")
	}
	return uint(math.Round(f)), nil
}

// IntToUint converts any int type to uint
// Equivilant to uint(i)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
// or if 'i' is signed
func IntToUint(i any) (uint, error) {
	switch ii := i.(type) {
	case int:
		return uintIt(ii)
	case int8:
		return uintIt(int(ii))
	case int16:
		return uintIt(int(ii))
	case int32:
		return uintIt(int(ii))
	case int64:
		return uintIt(int(ii))
	default:
		return 0, paramTypeError("IntToUint", "int", i)
	}
}
func uintIt(i int) (uint, error) {
	if ConversionOverflow(reflect.Uint, i) {
		return 0, typeError("IntToUint", " overflow error")
	}
	if i >= 0 {
		return uint(i), nil
	}
	return 0, paramTypeError("ToUint", "unsigned number", i)
}

// FloatToUint converts any float type to asserted uint
// Equivilant to uint(f)
// Returns error if param 'f' type is not float32, float64
// or if 'i' is signed
func FloatToUint(f any) (uint, error) {
	switch ff := f.(type) {
	case float32:
		return uintIt(int(ff))
	case float64:
		return uintIt(int(ff))
	default:
		return 0, paramTypeError("FloatToUint", "unsigned float", f)
	}
}

// UintToUint converts any uint type to uint
// Equivilant to uint(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func UintToUint(u any) (uint, error) {
	switch uu := u.(type) {
	case uint:
		return uint(uu), nil
	case uint8:
		return uint(uu), nil
	case uint16:
		return uint(uu), nil
	case uint32:
		return uint(uu), nil
	case uint64:
		return uint(uu), nil
	default:
		return 0, paramTypeError("UintToUint", "uint", u)
	}
}

// BoolToUint converts a bool to uint
// 1 if true, 0 if false
// Returns error if param 'b' type is not bool
func BoolToUint(b any) (uint, error) {
	switch bb := b.(type) {
	case bool:
		if bb {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, paramTypeError("BoolToUint", "bool", b)
	}
}

// TimeToUint converts a time.Time to a unix uint
// Equivilant to uint(t.Unix())
// Returns error if param 't' type is not time.Time
func TimeToUint(t any) (uint, error) {
	switch tt := t.(type) {
	case time.Time:
		return uint(tt.Unix()), nil
	default:
		return 0, paramTypeError("TimeToUint", "time.Time", t)
	}
}

// ToUint converts param 'a' of a basic type to uint
// Returns error if param 'a' type is not:
//   string, int, float, uint, bool or time
func ToUint(a any) (uint, error) {
	switch a.(type) {
	case string:
		return StringToUint(a)
	case int, int8, int16, int32, int64:
		return IntToUint(a)
	case float32, float64:
		return FloatToUint(a)
	case uint, uint8, uint16, uint32, uint64:
		return UintToUint(a)
	case bool:
		return BoolToUint(a)
	case time.Time:
		return TimeToUint(a)
	default:
		return 0, paramTypeError("ToUint", "string, numeric, bool, or time", a)
	}
}

// BOOL CONVERSION FUNCTIONS
// ToBool:			converts any basic type to bool
// StringToBool:		converts a string to bool			ALTERNATIVE: none
// IntToBool:		converts any int type to bool		ALTERNATIVE: none
// FloatToBool:		converts any float type to bool		ALTERNATIVE: none
// UnitToBool:		converts any uint type to bool		ALTERNATIVE: none
// BoolToBool:		converts a bool to bool				ALTERNATIVE: none

// StringToBool converts a string to bool
// Returns error if param 's' type is not string
// or can't be converted to unit
func StringToBool(s any) (bool, error) {
	if IsString(s) {
		ss := strings.ToLower(s.(string))
		if ss == "t" || ss == "true" || ss == "1" {
			return true, nil
		}
		if ss == "f" || ss == "false" || ss == "0" {
			return false, nil
		}
	}
	return false, paramTypeError("StringToBool", "string of bool", s)
}

// IntToBool converts any int type to bool
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
// or if 'i' is not 0 or 1
func IntToBool(i any) (bool, error) {
	switch ii := i.(type) {
	case int:
		return intToBool(ii)
	case int8:
		return intToBool(int(ii))
	case int16:
		return intToBool(int(ii))
	case int32:
		return intToBool(int(ii))
	case int64:
		return intToBool(int(ii))
	default:
		return false, paramTypeError("IntToBool", "int", i)
	}
}
func intToBool(i int) (bool, error) {
	if i == 1 {
		return true, nil
	}
	if i == 0 {
		return false, nil
	}
	return false, paramTypeError("IntToBool", "0 or 1 int", i)
}

// FloatToBool converts any float type to bool
// Returns error if param 'f' type is not float32, float64
// or if 'i' is not 0 or 1
func FloatToBool(f any) (bool, error) {
	switch ff := f.(type) {
	case float32:
		return intToBool(int(ff))
	case float64:
		return intToBool(int(ff))
	default:
		return false, paramTypeError("FloatToBool", "0 or 1 float", f)
	}
}

// UintToBool converts any uint type to bool
// Equivilant to uint(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
// or if 'i' is not 0 or 1
func UintToBool(u any) (bool, error) {
	switch uu := u.(type) {
	case uint:
		return intToBool(int(uu))
	case uint8:
		return intToBool(int(uu))
	case uint16:
		return intToBool(int(uu))
	case uint32:
		return intToBool(int(uu))
	case uint64:
		return intToBool(int(uu))
	default:
		return false, paramTypeError("UintToBool", "0 or 1 uint", u)
	}
}

// BoolToUint converts a bool to asserted bool
// Returns error if param 'b' type is not bool
func BoolToBool(b any) (bool, error) {
	switch bb := b.(type) {
	case bool:
		if bb {
			return true, nil
		}
		return false, nil
	default:
		return false, paramTypeError("BoolToBool", "bool", b)
	}
}

// ToBool converts param 'a' of a basic type to bool
// Returns error if param 'a' type is not:
//   string, int, float, uint, bool
func ToBool(a any) (bool, error) {
	switch a.(type) {
	case string:
		return StringToBool(a)
	case int, int8, int16, int32, int64:
		return IntToBool(a)
	case float32, float64:
		return FloatToBool(a)
	case uint, uint8, uint16, uint32, uint64:
		return UintToBool(a)
	case bool:
		return BoolToBool(a)
	default:
		return false, paramTypeError("ToBool", "string, numeric or bool", a)
	}
}

// TIME CONVERSION FUNCTIONS
// ToTime:			converts any basic type to time.Time
// StringToTime:	converts a string to time.Time			ALTERNATIVE: time.Parse()
// IntToTime:		converts any int type to time.Time		ALTERNATIVE: time.Parse()
// FloatToTime:		converts any float type to time.Time	ALTERNATIVE: time.Parse()
// UnitToTime:		converts any uint type to time.Time		ALTERNATIVE: time.Parse()
// TimeToTime:		converts a time to a time.Time			ALTERNATIVE: t.(time.Time)
// CurrencyToTime:	converts a currency to time.Time 		ALTERNATIVE: none

// StringToTime converts a numeric string to time.Time
// Similar to time.Parse(format, s)
// Returns error if param 's' type is not string
// or can't be converted to time
func StringToTime(s any) (time.Time, error) {
	if _, ok := s.(string); !ok {
		return time.Time{}, fmt.Errorf("not string")
	}
	f, err := timeStrFormat(s.(string))
	if err != nil {
		return time.Time{}, paramTypeError("StringToTime", "'2006-01-02 15:04:05.000' like date string", s)
	}
	t, err := time.Parse(f, s.(string))
	if err != nil {
		return time.Time{}, paramTypeError("StringToTime", "'2006-01-02 15:04:05.000' like date string", s)
	}
	return t, nil
}

func timeStrFormat(s string) (string, error) {
	// "2006-01-02 15:04:05.000000 +0000 UTC"
	td := `2006-01-02`
	tm := ` 15:04`
	ts := ` 15:04:05`
	tl := ` 15:04:05.000`
	tc := ` 15:04:05.000000`
	tz := ` +0000`
	tn := ` UTC`
	fm := ``
	if regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`).MatchString(string(s)) {
		fm += td
	}
	if regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9].[0-9]{6}`).MatchString(string(s)) {
		fm += tc
	} else if regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9].[0-9]{3}`).MatchString(string(s)) {
		fm += tl
	} else if regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]`).MatchString(string(s)) {
		fm += ts
	} else if regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9]`).MatchString(string(s)) {
		fm += tm
	}
	if fm != td && fm != `` {
		if tz = regexp.MustCompile(` (-|\+)[0-9]{4}[^-]?`).FindString(string(s)); tz != `` {
			fm += tz
		}
		if tn = regexp.MustCompile(` [A-Z]{3}`).FindString(string(s)); tn != `` {
			fm += tn
		}
	}
	if fm == `` {
		return ``, fmt.Errorf(`could not parse format from date string: %s`, s)
	}
	return fm, nil
}

// IntToTime converts any int type representing unix time to time.Time
// Equivilant to time.Unix(i, 0)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func IntToTime(i any) (time.Time, error) {
	switch ii := i.(type) {
	case int:
		return time.Unix(int64(ii), 0), nil
	case int8:
		return time.Unix(int64(ii), 0), nil
	case int16:
		return time.Unix(int64(ii), 0), nil
	case int32:
		return time.Unix(int64(ii), 0), nil
	case int64:
		return time.Unix(int64(ii), 0), nil
	default:
		return time.Time{}, paramTypeError("IntToTime", "unix time int", i)
	}
}

// FloatToTime converts any float type representing unix time to time.Time
// Sumilar to time.Unix(int64(i), n)
// Returns error if param 'f' type is not float32, float64
func FloatToTime(f any) (time.Time, error) {
	switch ff := f.(type) {
	case float32:
		n := int64(math.Round(float64(ff * 1e9)))
		return time.Unix(0, n), nil
	case float64:
		n := int64(math.Round(float64(ff * 1e9)))
		return time.Unix(0, n), nil
	default:
		return time.Time{}, paramTypeError("FloatToTime", "unix time float", f)
	}
}

// UintToTime converts any uint type representing unix time to time.Time
// Equivilant to time.Unix(int64(i), 0)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func UintToTime(u any) (time.Time, error) {
	switch uu := u.(type) {
	case uint:
		return time.Unix(int64(uu), 0), nil
	case uint8:
		return time.Unix(int64(uu), 0), nil
	case uint16:
		return time.Unix(int64(uu), 0), nil
	case uint32:
		return time.Unix(int64(uu), 0), nil
	case uint64:
		return time.Unix(int64(uu), 0), nil
	default:
		return time.Time{}, paramTypeError("UintToTime", "unix time uint", u)
	}
}

// TimeToTime converts a time.Time to asserted time.Time
// Returns error if param 't' type is not time.Time
func TimeToTime(t any) (time.Time, error) {
	switch t.(type) {
	case time.Time:
		return t.(time.Time), nil
	default:
		return time.Time{}, paramTypeError("TimeToTime", "time.Time", t)
	}
}

// ToTime converts param 'a' of a basic type to time.Time
// Returns error if param 'a' type is not:
//   string, int, float, uint or time
func ToTime(a any) (time.Time, error) {
	switch a.(type) {
	case string:
		return StringToTime(a)
	case int, int8, int16, int32, int64:
		return IntToTime(a)
	case float32, float64:
		return FloatToTime(a)
	case uint, uint8, uint16, uint32, uint64:
		return UintToTime(a)
	case time.Time:
		return TimeToTime(a)
	default:
		return time.Time{}, paramTypeError("ToTime", "string, numeric unix time or time", a)
	}
}

// MAP CONVERSION FUNCTIONS
// KeyValArraysToMap	converts two arrays to map 					ALTERNATIVE: none
// KeyValPairsToMap		converts array of key value pairs to map 	ALTERNATIVE: none
// StructToMap			converts struct and substructs to map		ALTERNATIVE: none
// JsonToMap			converts json []byte to map 				ALTERNATIVE: encoding.json.Unmarshal()
// MapKeyType			returns the type of the map keys			ALTERNATIVE: reflect.TypeOf().Key()
// MapValType 			returns the type of the map values			ALTERNATIVE: reflect.TypeOf().Elem()
// DeepTypeOf			returns an array of types at each dememsion	ALTERNATIVE: none

// MapToMap converts a map to map[any]any
func MapToMap(a any) (map[any]any, error) {
	m := map[any]any{}
	if !IsMap(a) {
		return m, paramTypeError("MapToMap", "map", a)
	}
	iter := reflect.ValueOf(a).MapRange()
	for iter.Next() {
		var v any
		if reflect.TypeOf(iter.Value().Interface()).Kind() == reflect.Map {
			v, _ = MapToMap(iter.Value().Interface())
		} else {
			v = iter.Value().Interface()
		}
		m[iter.Key().Interface()] = v
	}
	return m, nil
}

// ArrayToMap converts an array or slice to a map
func ArrayToMap(a any) (map[any]any, error) {
	m := map[any]any{}
	if !IsArray(a) {
		return m, paramTypeError("ArrayToMap", "array or slice", a)
	}
	s := reflect.ValueOf(a)
	for i := 0; i < s.Len(); i++ {
		m[i] = s.Index(i).Interface()
	}
	return m, nil
}

// KeyValArraysToMap converts two arrays to a map
// first array contains the keys
// second array contains the values
func KeyValArraysToMap(k any, v any) (map[any]any, error) {
	m := map[any]any{}
	if !IsArray(k) || !IsArray(v) {
		return m, typeError("KeyValArraysToMap", " expected an array or slice for k and v")
	}
	kv := reflect.ValueOf(k)
	vv := reflect.ValueOf(v)
	if kv.Len() != vv.Len() {
		return m, typeError("KeyValArraysToMap", "  number of elements in keys (%d) and values (%d) do not match", kv.Len(), vv.Len())
	}
	for i := 0; i < kv.Len(); i++ {
		key := kv.Index(i).Interface()
		kt := TypeOf(key)
		if !validMapKeyType(kt) {
			return map[any]any{}, typeError("KeyValArraysToMap", " invalid map key type of %v", kt)
		}
		m[key] = vv.Index(i).Interface()
	}
	return m, nil
}

// KeyValPairsToMap converts an array of key value pairs to a map
// first value at each index is a key
// second value at each index is the associated value
func KeyValPairsToMap(a any) (map[any]any, error) {
	m := map[any]any{}
	if !IsArray(a) {
		return m, paramTypeError("KeyValPairsToMap", "array or slice", a)
	}
	p := reflect.ValueOf(a)
	if p.Len() != 0 {
		kt := TypeOf(p.Index(0).Index(0).Interface())
		if !validMapKeyType(kt) {
			return m, typeError("KeyValPairsToMap", "  invalid key type of %s", kt)
		}
		for i := 0; i < p.Len(); i++ {
			kv := p.Index(i).Interface().([]any)
			if len(kv) != 2 {
				return map[any]any{}, typeError("KeyValPairsToMap", "expected 2 elements at index '%d' of 'a'", i)
			}
			ckt := TypeOf(kv[0])
			if ckt != kt {
				return map[any]any{}, typeError("KeyValPairsToMap", "  inconsistent key type\n  expected %s and received %s", kt, ckt)
			}
			m[kv[0]] = kv[1]
		}
	}
	return m, nil
}

// StructToMap converts a struct to a map[string]any
// also converts embedded structs to maps
// uses struct tag 'json' as an override to key names
func StructToMap(s any) (map[any]any, error) {
	m := map[any]any{}
	sRef := reflect.ValueOf(s)
	if sRef.Kind() == reflect.Pointer {
		sRef = sRef.Elem()
	}
	if sRef.Kind() != reflect.Struct {
		return m, paramTypeError("StructToMap", "a struct", s)
	}
	t := sRef.Type()
	for i := 0; i < sRef.NumField(); i++ {
		f := t.Field(i)
		n := f.Tag.Get("json")
		if n == "" {
			n = f.Name
		}
		if f.Type.Kind() == reflect.Struct {
			m[n], _ = StructToMap(sRef.Field(i).Interface())
		} else {
			m[n] = sRef.Field(i).Interface()
		}
	}
	return m, nil
}

// JsonToMap converts a Json []byte to a map
// Equivilant to encoding.json.Unmarshal(j, map[string]any)
// returns error if j is not []byte type or unable to unmarshal
func JsonToMap(j any) (map[any]any, error) {
	jsn, ok := j.([]byte)
	if !ok {
		return map[any]any{}, paramTypeError("JsonToMap", "json bytes", j)
	}
	m := map[string]any{}
	err := json.Unmarshal(jsn, &m)
	if err != nil {
		paramTypeError("JsonToMap", "json bytes", j)
	}
	ma, _ := MapToMap(m)
	return ma, nil
}

// valueMapKeyType determins if Type 't' can be a key in a map
func validMapKeyType(t Type) bool {
	if t != String && t != Int && t != Uint && t != Float && t != Time {
		return false
	}
	return true
}

// MapKeyType returns the Type of the key in map 'a'
func MapKeyType(a any) (Type, error) {
	if !IsMap(a) {
		return Invalid, paramTypeError("MapKeyType", "map", a)
	}
	typStr := fmt.Sprintf("%T", a)
	typStr = regexp.MustCompile(`\[(.*?)\]`).FindAllString(typStr, -1)[0]
	typStr = strings.Replace(
		regexp.MustCompile(`[^a-zA-Z]`).ReplaceAllString(typStr, ""),
		"interface", "any", 1,
	)
	typ, err := TypeByName(typStr)
	if err != nil {
		return typ, paramTypeError("MapKeyType", "valid key", a)
	} else {
		return typ, nil
	}
}

// MapValType returns the Type of the values in map 'a'
func MapValType(a any) (Type, error) {
	if !IsMap(a) {
		return Invalid, paramTypeError("MapValType", "map", a)
	}
	typStr := fmt.Sprintf("%T", a)
	typStr = regexp.MustCompile(`(\[|\])`).Split(typStr, -1)[2]
	typStr = strings.Replace(
		regexp.MustCompile(`[^a-z]`).ReplaceAllString(typStr, ""),
		"interface", "any", 1,
	)
	typ, err := TypeByName(typStr)
	if err != nil {
		return typ, paramTypeError("MapValType", "valid value", a)
	} else {
		return typ, nil
	}
}

// MapType returns an array of types for each map layer:
// [KeyType, ValueType] or [KeyType, KeyType, ..., ValueType]
func DeepTypeOf(a any) ([]Type, error) {
	types := []Type{}
	typ := TypeOf(a)
	if typ != Map && typ != Array && typ != Struct {
		return types, paramTypeError("DeepTypeof", "map, array, slice or struct", typeNames[typ])
	}
	s := fmt.Sprintf("%T", a)
	s = strings.ReplaceAll(s, `[]`, `array `)
	s = strings.ReplaceAll(s, `interface {}`, `any`)
	s = regexp.MustCompile(`(\[|\])`).ReplaceAllString(s, ` `)
	ts := strings.Split(s, ` `)
	for _, typ := range ts {
		if strings.Contains(typ, ".") {
			typ = "struct"
		}
		t, err := TypeByName(typ)
		if err != nil {
			return []Type{}, paramTypeError("DeepTypeOf", "string, bool, numeric, time, array, map or struct", a)
		}
		types = append(types, t)
	}
	return types, nil
}

// ARRAY CONVERSION FUNCTIONS
// MapKeys			converts the keys of a map to array 		ALTERNATIVE: reflect.MapKeys()
// MapVals 			converts the values of a map to array		ALTERNATIVE: none
// StructKeys		converts the fields of a struct to array	ALTERNATIVE: none
// StructVals		converts the values of a struct to array 	ALTERNATIVE: none
// ArrayValType		returns the type of the vals in array		ALTERNATIVE: none

// MapKeys converts the keys of map 'a' to an array
// Equivilant to reflect.MapKeys()
func MapKeys(a any) ([]any, error) {
	s := []any{}
	if !IsMap(a) {
		return s, paramTypeError("MapKeys", "map", a)
	}
	iter := reflect.ValueOf(a).MapRange()
	for iter.Next() {
		s = append(s, iter.Key().Interface())
	}
	return s, nil
}

// MapVals converts the values of map 'a' to an array
func MapVals(a any) ([]any, error) {
	s := []any{}
	if !IsMap(a) {
		return s, paramTypeError("MapVals", "map", a)
	}
	iter := reflect.ValueOf(a).MapRange()
	for iter.Next() {
		s = append(s, iter.Value().Interface())
	}
	return s, nil
}

// MapToArray converts the values of map 'a' to an array
// Equivilant to MapVals
func MapToArray(a any) ([]any, error) {
	return MapVals(a)
}

// StructFields returns a []string of struct 'a' field names
// uses struct tag 'json' as an override to key names
func StructFields(a any) ([]any, error) {
	s := []any{}
	sRef := reflect.ValueOf(a)
	if sRef.Kind() == reflect.Pointer {
		sRef = sRef.Elem()
	}
	if sRef.Kind() != reflect.Struct {
		return s, paramTypeError("StructFields", "a struct", a)
	}
	t := sRef.Type()
	for i := 0; i < sRef.NumField(); i++ {
		n := t.Field(i).Tag.Get("json")
		if n == "" {
			n = t.Field(i).Name
		}
		s = append(s, n)
	}
	return s, nil
}

// StructValues returns a []any of struct 'a' values
func StructValues(a any) ([]any, error) {
	s := []any{}
	sRef := reflect.ValueOf(a)
	if sRef.Kind() == reflect.Pointer {
		sRef = sRef.Elem()
	}
	if sRef.Kind() != reflect.Struct {
		return s, paramTypeError("StructFields", "a struct", a)
	}
	for i := 0; i < sRef.NumField(); i++ {
		s = append(s, sRef.Field(i).Interface())
	}
	return s, nil
}

// ArrayValType returns the Type of the key in map 'a'
func ArrayValType(a any) (Type, error) {
	if !IsArray(a) {
		return Invalid, paramTypeError("ArrayValType", "array or slice", a)
	}
	typStr := fmt.Sprintf("%T", a)
	typStr = regexp.MustCompile(`[a-z]+`).FindString(typStr)
	typStr = strings.Replace(typStr, "interface", "any", 1)
	typ, err := TypeByName(typStr)
	if err != nil {
		return typ, paramTypeError("ArrayValType", "valid value", a)
	} else {
		return typ, nil
	}
}

// STRUCT CONVERSION FUNCTIONS
// KeyValArraysToStruct		converts
// KeyValPairsToStruct
// MaptoStruct
// JsonToStruct

type Val struct {
	String string
	Bool   bool
	Int    int
	Uint   uint
	Float  float64
	Time   time.Time
	Array  []any
	Map    map[any]any
	Struct any
}

// KeyValArraysToStruct converts two arrays to struct
// first array contains field names (must be strings)
// second array contains the associated value
// returns error if any element in first array is not string
func KeyValArraysToStruct(k any, v any, s ...any) (any, error) {
	m, err := KeyValArraysToMap(k, v)
	if err != nil {
		return nil, typeError("KeyValArraysToStruct", fmt.Sprint(":", err))
	}
	if len(s) > 0 {
		return mapToStructExisting(m, s[0])
	} else {
		return mapToStructNew(m)
	}
}

// KeyValPairsToStruct converts array of key val pairs to struct
// first element in each pair becomes a field name
// second element in each pair becomes the associated value
// returns error if first element in any pair is not string
func KeyValPairsToStruct(a any, s ...any) (any, error) {
	m, err := KeyValPairsToMap(a)
	if err != nil {
		return nil, typeError("KeyValPairsToStruct", fmt.Sprint(":", err))
	}
	if len(s) > 0 {
		return mapToStructExisting(m, s[0])
	} else {
		return mapToStructNew(m)
	}
}

// MapToStruct converts map to struct
// keys become the field name
// values become the associated value
// returns error if keys are not strings
func MapToStruct(m any, s ...any) (any, error) {
	if len(s) > 0 {
		return mapToStructExisting(m, s[0])
	} else {
		return mapToStructNew(m)
	}
}

func mapToStructExisting(m any, s any) (any, error) {
	if !IsMap(m) {
		return nil, paramTypeError("MapToStruct", "map", m)
	}
	// create struct from provided template
	st := reflect.New(reflect.TypeOf(s)).Elem()
	for i := 0; i < st.NumField(); i++ {
		// get field name from struct
		n := st.Type().Field(i).Tag.Get("json")
		if n == "" {
			n = st.Type().Field(i).Name
		}
		// get field value from provided map
		mv, ok := m.(map[any]any)[n]
		sv := st.Field(i).Interface()
		if ok {
			switch TypeOf(sv) {
			case Map:
				if TypeOf(mv) == Map {
					st.Field(i).Set(reflect.ValueOf(mv))
				} else {
					return nil, paramTypeError("MapToStruct", "map", mv)
				}
				break
			case Struct:
				if TypeOf(mv) == Map {
					ns, err := mapToStructExisting(mv, sv)
					if err != nil {
						return nil, err
					}
					st.Field(i).Set(reflect.ValueOf(ns))
				} else {
					return nil, paramTypeError("MapToStruct", "map", mv)
				}
				break
			case String, Bool, Int, Float, Uint:
				iv, err := StrictlyTo(sv, mv)
				if err != nil {
					return nil, paramTypeError("MapToStruct", fmt.Sprint(TypeOf(sv)), mv)
				}
				st.Field(i).Set(reflect.ValueOf(iv[st.Field(i).Kind()]))
				break
			case TypeOf(mv):
				st.Field(i).Set(reflect.ValueOf(mv))
			default:
				return nil, paramTypeError("MapToStruct", fmt.Sprint(TypeOf(sv)), mv)
			}
		}
	}
	return st.Interface(), nil
}

func mapToStructNew(m any) (any, error) {
	if !IsMap(m) {
		return nil, paramTypeError("MapToStruct", "map", m)
	}
	// build a list of struct fields from the map
	fs := []reflect.StructField{}
	i := reflect.ValueOf(m).MapRange()
	vals := map[any]any{}
	for i.Next() {
		// capture key and deep value at i
		k, str := i.Key().Interface().(string)
		if !str {
			return nil, paramTypeError("MapToStruct", "map key to be a string", k)
		}
		mk := strings.ToUpper(k[:1]) + k[1:]
		v := i.Value().Interface()
		if IsMap(v) {
			v, _ = MapToStruct(v)
		}
		vals[mk] = v
		// create struct field
		fs = append(fs, reflect.StructField{
			Name: mk,
			Type: reflect.TypeOf(v),
			Tag:  reflect.StructTag(`json:"` + k + `"`),
		})
	}
	// build an empty struct from struct fields
	s := reflect.New(reflect.StructOf(fs)).Elem()
	//populate the empty struct with values from the map
	for k, v := range vals {
		s.FieldByName(k.(string)).Set(reflect.ValueOf(v))
	}
	return s.Interface(), nil
}

// JsonToStruct converts a json object to struct
// keys become the field name
// values become the associated values
// returns error if keys are not strings
func JsonToStruct(j any, s ...any) (any, error) {
	m, err := JsonToMap(j)
	if err != nil {
		return nil, paramTypeError("JsonToStruct", "json formatted []byte", j)
	}
	if len(s) > 0 {
		return mapToStructExisting(m, s[0])
	} else {
		return mapToStructNew(m)
	}
}

// FieldByTag returns the reflect instance of a Field in struct 's'
// allowing search for a field by tag 't' and its value 'v'
func StructFieldByTag(s any, t string, v string) (any, error) {
	r := reflect.ValueOf(s)
	if r.Kind() != reflect.Struct {
		return nil, fmt.Errorf("reflect.Value.FieldByTag: value not a struct")
	}
	typ := r.Type()
	for i := 0; i < r.NumField(); i++ {
		if v == typ.Field(i).Tag.Get(t) {
			return typ.Field(i), nil
		}
	}
	return nil, fmt.Errorf("reflect.Value.FieldByTag: could not find value '%s' for tag '%s'", v, t)
}

// Class is a container struct that mimics the
// same functionality of a normal struct, allowing
// runtime dynamic us of structs
// Class is comprised of Class times (Fields and Methods)
// Methods are identified with a ClassItem.Type of Func
// Class is instantiated with the func CreateClass(name, []ClassItems)
type Class struct {
	Name  string      // the class name
	List  []string    // a list of field and method names
	Items []ClassItem // the class fields and methods
	inst  bool        // if true, the class has been instantiated
}

// CreateClass instantiates a class by providing it name 'n' and ClassItems 'i'
// returns the instantiated class
// returns error if there is no name or ClassItems
func CreateClass(n string, i []ClassItem) (*Class, error) {
	c := Class{}
	if n == "" {
		return &c, typeError("CreateClass", "missing class name")
	}
	if len(i) == 0 {
		return &c, typeError("CreateClass", "missing class items")
	}
	for _, ci := range i {
		err := c.Append(ci)
		if err != nil {
			return &Class{}, typeError("CreateClass", "class item missing values:\n%v", ci)
		}
	}
	c.Name = n
	c.inst = true
	return &c, nil
}

// IsInst returns true is the Class has been instantiated,
// meaning the Name, and ClassItems have been set
func (c *Class) IsInst() bool {
	return c.inst
}

// Get returns the ClassItem with index or name 'i',
// which must be and Int or String and exist in Class, and
// returns second expression of false if 'i' does not exist
func (c *Class) Get(i any) (*ClassItem, bool) {
	if c.inst {
		t := TypeOf(i)
		if t == String {
			for _, ci := range c.Items {
				if ci.Name == i {
					return &ci, true
				}
			}
		} else if t == Int && i.(int) < len(c.List) {
			return &c.Items[i.(int)], true
		}
	}
	return &ClassItem{}, false
}

// Index returns the ClassItem index of name 'n' and
// returns second expression of false if 'i' does not exist
func (c *Class) IndexOf(n string) (int, bool) {
	if c.inst {
		for i, li := range c.List {
			if n == li {
				return i, true
			}
		}
	}
	return 0, false
}

// Append adds ClassItem 'i' to the Class and
// returns an error if unable to add the item,
// for example, if the class isn't instantiated or
// the item name already exists
func (c *Class) Append(i ClassItem) error {
	if !c.inst {
		return typeError("Class.Append", "target class has not been instantiated")
	}
	if i.Name != "" && i.Type != Invalid {
		i.Class = c
		if i.Value != nil {
			i.set = true
		}
		c.Items = append(c.Items, i)
		c.List = append(c.List, i.Name)
		return nil
	}
	return typeError("Class.Append", "ClassItem missing Name and/or Type:\n%v", i)
}

// Remove removes ClassItem with index or name 'i',
// which must be and Int or String and exist in Class, and
// returns an error if 'i' does not exist
func (c *Class) Remove(i any) error {
	typ := TypeOf(i)
	if typ == String {
		if ii, ok := c.IndexOf(i.(string)); ok {
			i = ii
			typ = Int
		}
	}
	if typ == Int {
		c.List = append(c.List[:i.(int)], c.List[i.(int)+1:]...)
		c.Items = append(c.Items[:i.(int)], c.Items[i.(int)+1:]...)
		return nil
	}
	return paramTypeError("Class.Remove", " existing index int or name string", i)
}

// ClassItem is a field or method of a Class, and
// must have a Name, Type and Class to be instantiated.
// Class may also be instantiated with a Value (matching Type)
type ClassItem struct {
	Name  string // the name of the item
	Type  Type   // the Type of the item
	Tag   string // the item tag eg. `json: tag`
	Value any    // the value of the item
	Class *Class // pointer to the parent Class
	set   bool   // stores whether the Value has been set
}

// IsInst returns true if ClassItem has a Class, Name and Type
func (i *ClassItem) IsInst() bool {
	if i.Class == nil || i.Name == "" || i.Type == Invalid {
		return false
	}
	return true
}

// IsSet returns true if ClassItem is instantiated and has a Value
func (i *ClassItem) IsSet() bool {
	if i.IsInst() && (i.Type == Int) {
		return true
	}
	return false
}

// Set allocates Value 'v' to the ClassItem, and
// returns and error if the Type of 'v' doesn't match
// the ClassItem Type
func (i *ClassItem) Set(v any) error {
	if TypeOf(v) != i.Type {
		paramTypeError("ClassItem.Set", fmt.Sprintf("for item '%v' a %v", i.Name, typeNames[i.Type]), v)
	}
	i.Value = v
	i.set = true
	return nil
}

// Remove removes the ClassItem from it's Class, and
// returns an error if the ClassItem hasn't been instantiated
func (i *ClassItem) Remove() error {
	if i.IsInst() {
		i.Class.Remove(i.Name)
		return nil
	}
	return paramTypeError("ClassItem.Remove", "populated ClassItem", i)
}

// Reset sets the ClassItem Value to the nil equivilant
// for the ClassItem Type (eg. 0, "", false)
//func reset

/*
	ADDITIONAL TASKS

	ToClass should convert to an abstract of reflect.Struct
	with methods provied above

	custom structs
	conversion to from []byte
*/
