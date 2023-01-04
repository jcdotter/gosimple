// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package string

import (
	"gosimple/types"
	"strings"
	"time"
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

// SnakeToCamel converts snake case string 's'
// example_string to camel case format exampleString
func SnakeToCamel(s string) string {
	return snakeTo(s, 1)
}

// SnakeToCamel converts snake case string 's'
// example_string to pascal case format ExampleString
func SnakeToPascal(s string) string {
	return snakeTo(s, 0)
}

func snakeTo(s string, n int) string {
	s = strings.ToLower(s)
	np := strings.Split(s, `_`)
	for i := n; i < len(np); i++ {
		np[i] = strings.ToUpper(np[i][:1]) + np[i][1:]
	}
	return strings.Join(np, ``)
}

// PascalToCamel converts pascal case string 's'
// ExampleString to camel case format exampleString
func PascalToCamel(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

// CamelToPascal converts camel case string 's'
// exampleString to pascal case format ExampleString
func CamelToPascal(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// CamelToSnake converts camel case string 's'
// exampleString to snake case format example_string
func CamelToSnake(s string) string {
	r := string(s[0])
	for i := 1; i < len(s); i++ {
		if 91 > s[i] && s[i] > 67 {
			r += `_`
		}
		r += string(s[i])
	}
	return strings.ToLower(r)
}

// PascalToSnake converts pacsal case string 's'
// ExampleString to snake case format example_string
func PascalToSnake(s string) string {
	return CamelToPascal(s)
}
