// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package hmap

import (
	"reflect"

	"github.com/jcdotter/gosimple/types"
)

// Is evaluates whether 'a' is a map
func Is(a any) bool {
	return types.IsMap(a)
}

// Keys converts the keys of map 'm' to an array
// Equivilant to reflect.MapKeys()
func Keys(m any) []any {
	k, _ := types.MapKeys(m)
	return k
}

// KeyType returns the Type of the key in map 'a'
func KeyType(a any) (types.Type, error) {
	return types.MapKeyType(a)
}

// Vals converts the values of map 'm' to an array
func Vals(m any) []any {
	v, _ := types.MapVals(m)
	return v
}

// ValType returns the Type of the values in map 'a'
func ValType(a any) (types.Type, error) {
	return types.MapValType(a)
}

// ContainsKeys evaluates whether map 'm'
// contains the keys 'k' provided
func ContainsKeys(m any, k ...any) bool {
	if !Is(m) {
		return false
	}
	mr := reflect.ValueOf(m)
	for _, key := range k {
		if mr.MapIndex(reflect.ValueOf(key)) == (reflect.Value{}) {
			return false
		}
	}
	return true
}

// ContainsKeys evaluates whether map 'm'
// contains the values 'v' provided
func ContainsVals(m any, v ...any) bool {
	if !Is(m) {
		return false
	}
	found := false
	mr := reflect.ValueOf(m)
	for _, val := range v {
		found = false
		vk := reflect.TypeOf(val).Kind()
		i := mr.MapRange()
		for i.Next() {
			if i.Value().Kind() == vk && i.Value().Interface() == val {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return found
}

// FromArray converts an array or slice to a map
func FromArray(a any) (map[any]any, error) {
	return types.ArrayToMap(a)
}

// FromKeyValArrays converts two arrays to a map
// first array contains the keys
// second array contains the values
func FromKeyValArrays(k any, v any) (map[any]any, error) {
	return types.KeyValArraysToMap(k, v)
}

// FromKeyValPairs converts an array of key value pairs to a map
// first value at each index is a key
// second value at each index is the associated value
func FromKeyValPairs(a any) (map[any]any, error) {
	return types.KeyValPairsToMap(a)
}

// FromStruct converts a struct to a map[string]any
// also converts embedded structs to maps
// uses struct tag 'json' as an override to key names
func FromStruct(s any) (map[any]any, error) {
	return types.StructToMap(s)
}

// FromJson converts a Json []byte to a map
// Equivilant to encoding.json.Unmarshal(j, map[string]any)
// returns error if j is not []byte type or unable to unmarshal
func FromJson(j any) (map[any]any, error) {
	return types.JsonToMap(j)
}

// Struct converts map to struct
// keys become the field name
// values become the associated value
// returns error if keys are not strings
func Struct(m any, s ...any) (any, error) {
	return types.MapToStruct(m, s, types.None, "")
}
