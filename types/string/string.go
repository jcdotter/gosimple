// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package string

import (
	"time"

	"github.com/jcdotter/gosimple/types"
)

// Is evaluates whether 'a' is a string
func Is(a any) bool {
	return types.IsString(a)
}

// From converts param 'a' of a basic type to string
// Equivilant to fmt.Sprint(i)
// Returns error if param 'a' type is not
// string, int, float, uint, bool, time, slice, map or struct
func From(a any) (string, error) {
	return types.ToString(a)
}

// FromInt converts an int to string
// Equivilant to fmt.Sprint(i)
// Returns error if param 'i' type is not int, int8, int16, int32 or int64
func FromInt(i any) (string, error) {
	return types.IntToString(i)
}

// FromFloat converts a float to string
// Equivilant to fmt.Sprint(f)
// Returns error if param 'f' type is not float32 or float64
func FromFloat(f any) (string, error) {
	return types.FloatToString(f)
}

// FromUint converts a uint to string
// Equivilant to fmt.Sprint(u)
// Returns error if param 'u' type is not uint, uint8, uint16, uint32 or uint64
func FromUint(u any) (string, error) {
	return types.UintToString(u)
}

// FromBool converts a bool to string
// Equivilant to fmt.Sprint(b)
// Returns error if param 'b' type is not bool
func FromBool(b any) (string, error) {
	return types.BoolToString(b)
}

// FromTime converts a time to string
// Equivilant to fmt.Sprint(t)
// Returns error if param 't' type is not time.Time
func FromTime(t any) (string, error) {
	return types.TimeToString(t)
}

// Int converts a numeric string to a rounded int
// Similar to int(math.Round(strconv.ParseFloat(s, 64)))
// Returns error if param 's' type is not string
// or can't be converted to int
func Int(s any) (int, error) {
	return types.StringToInt(s)
}

// Float converts a numeric string to float64
// Similar to strconv.ParseFloat(str, 64)
// Returns error if param 'str' type is not string
// or can't be converted to float64
func Float(s any) (float64, error) {
	return types.StringToFloat(s)
}

// Uint converts a numeric string to rounded uint
// Similar to uint(strconv.ParseFloat(str, 64))
// Returns error if param 's' type is not string
// or can't be converted to unit
func Uint(s any) (uint, error) {
	return types.StringToUint(s)
}

// Bool converts a string to bool
// Returns error if param 's' type is not string
// or can't be converted to unit
func Bool(s any) (bool, error) {
	return types.StringToBool(s)
}

// Time converts a numeric string to time.Time
// Similar to time.Parse(format, s)
// Returns error if param 's' type is not string
// or can't be converted to time
func Time(s any) (time.Time, error) {
	return types.StringToTime(s)
}

// ToCamelCase converts string 's' example_string
// to camel case format exampleString
func ToCamelCase(s string) string {
	return types.ToCamelString(s)
}

// ToPascalCase converts string 's' example_string
// to pascal case format ExampleString
func ToPascalCase(s string) string {
	return types.ToPascalString(s)
}

// ToSnakeCase converts string 's' exampleString
// to snake case format example_string
func ToSnakeCase(s string) string {
	return types.ToSnakeString(s)
}

// ToPhraseCase converts string 's' exampleString
// to phrase case format Example string and
// if case sensative 'c', creating new word at each capital letter
func ToPhraseCase(s string, c bool) string {
	return types.ToPhraseString(s, c)
}
