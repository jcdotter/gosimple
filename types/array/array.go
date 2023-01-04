// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package array

import (
	"reflect"

	"github.com/jcdotter/gosimple/types"
)

// Is evaluates whether 'a' is an array or slice
func Is(a any) bool {
	return types.IsArray(a)
}

// Contains evaluates whether slice (or array) 's'
// contains the elems provided
func Contains(s any, elems ...any) bool {
	if !Is(s) {
		return false
	}
	array := reflect.ValueOf(s)
	arrayType := reflect.TypeOf(s).Elem().Kind()
	found := false
	for _, el := range elems {
		found = false
		elType := reflect.TypeOf(el).Kind()
		if arrayType == reflect.Interface || elType == arrayType {
			for i := 0; i < array.Len(); i++ {
				if el == array.Index(i).Interface() {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return found
}

// ValType returns the Type of the key in map 'a'
func ValType(a any) (types.Type, error) {
	return types.ArrayValType(a)
}

// FromMapKeys converts the keys of map 'a' to an array
// Equivilant to reflect.MapKeys()
func FromMapKeys(a any) ([]any, error) {
	return types.MapKeys(a)
}

// FromMapVals converts the values of map 'a' to an array
func FromMapVals(a any) ([]any, error) {
	return types.MapVals(a)
}

// FromStructFields returns a []string of struct 'a' field names
// uses struct tag 'json' as an override to key names
func FromStructFields(a any) ([]any, error) {
	return types.StructFields(a)
}

// FromStructValues returns a []any of struct 'a' values
func FromStructValues(a any) ([]any, error) {
	return types.StructValues(a)
}
