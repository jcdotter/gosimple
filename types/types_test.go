// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: James Dotter

package types

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	aType Type = iota + 100
	vType
)

// TEST CASE DATA:
// contains variables that can be converted
// between data types to validate conversion

var (
	str     = "1"
	strb    = "true"
	strt    = "1969-12-31 17:00:01 -0700 MST"
	strf    = "(1,234,567.89)"
	intn    = 1
	floatn  = 1.0
	floats  = -1234567.89
	uintn   = uint(1)
	boolv   = true
	timev   = time.Unix(1, 0)
	array   = []any{1.0, 2.0, 3.0, 4.0}
	arrayk  = []any{"one", "two", "three", "four"}
	arrayv  = []any{1.0, 2.0, map[any]any{"one": "one", "two": "two"}, array}
	arrayvs = []any{1.0, 2.0, sts{"one", "two"}, array}
	arraykv = [][]any{{"one", 1.0}, {"two", 2.0}, {"three", 3.0}, {"four", 4.0}}
	hmap    = map[any]any{"one": 1.0, "two": 2.0, "three": map[any]any{"one": "one", "two": "two"}, "four": array}
	hmapk   = map[any]any{0: "one", 1: "two", 2: "three", 3: "four"}
	hmapkv  = map[any]any{"one": 1.0, "two": 2.0, "three": 3.0, "four": 4.0}
	strct   = st{1, 2, sts{"one", "two"}, []any{1.0, 2.0, 3.0, 4.0}}
	strctkv = stkv{1.0, 2.0, 3.0, 4.0}
	jsonv   = []byte(`{"four":[1.0,2.0,3.0,4.0],"one":1.0,"three":{"one":"one","two":"two"},"two":2.0}`)
)

type st struct {
	One   float64 `json:"one"`
	Two   float64 `json:"two"`
	Three sts     `json:"three"`
	Four  []any   `json:"four"`
}

type sts struct {
	One string `json:"one"`
	Two string `json:"two"`
}

type stkv struct {
	One   float64 `json:"one"`
	Two   float64 `json:"two"`
	Three float64 `json:"three"`
	Four  float64 `json:"four"`
}

// TEST CASES:
// contains cases for testing the conversion
// of variables above between data types

type test struct {
	Type      Type
	Name      string
	Func      any
	Validator any
	Params    []any
}

var convTests = []test{
	// STRING CONVERSION FUNCTIONS
	{String, "StringToString", StringToString, str, []any{str}},
	{String, "IntToString", IntToString, str, []any{intn}},
	{String, "FloatToString", FloatToString, str, []any{floatn}},
	{String, "UintToString", UintToString, str, []any{uintn}},
	{String, "BoolToString", BoolToString, strb, []any{boolv}},
	{String, "TimeToString", TimeToString, strt, []any{timev}},
	// INT CONVERSION FUNCTIONS
	{Int, "StringToInt", StringToInt, intn, []any{str}},
	{Int, "IntToInt", IntToInt, intn, []any{intn}},
	{Int, "FloatToInt", FloatToInt, intn, []any{floatn}},
	{Int, "UintToInt", UintToInt, intn, []any{uintn}},
	{Int, "BoolToInt", BoolToInt, intn, []any{boolv}},
	{Int, "TimeToInt", TimeToInt, intn, []any{timev}},
	// FLOAT CONVERSION FUNCTIONS
	{Float, "StringToFloat", StringToFloat, floatn, []any{str}},
	{Float, "StringToFloat", StringToFloat, floats, []any{strf}},
	{Float, "IntToFloat", IntToFloat, floatn, []any{intn}},
	{Float, "FloatToFloat", FloatToFloat, floatn, []any{floatn}},
	{Float, "UintToFloat", UintToFloat, floatn, []any{uintn}},
	{Float, "BoolToFloat", BoolToFloat, floatn, []any{boolv}},
	{Float, "TimeToFloat", TimeToFloat, floatn, []any{timev}},
	// UINT CONVERSION FUNCTIONS
	{Uint, "StringToUint", StringToUint, uintn, []any{str}},
	{Uint, "IntToUint", IntToUint, uintn, []any{intn}},
	{Uint, "FloatToUint", FloatToUint, uintn, []any{floatn}},
	{Uint, "UintToUint", UintToUint, uintn, []any{uintn}},
	{Uint, "BoolToUint", BoolToUint, uintn, []any{boolv}},
	{Uint, "TimeToUint", TimeToUint, uintn, []any{timev}},
	// BOOL CONVERSION FUNCTIONS
	{Bool, "StringToBool", StringToBool, boolv, []any{str}},
	{Bool, "IntToBool", IntToBool, boolv, []any{intn}},
	{Bool, "FloatToBool", FloatToBool, boolv, []any{floatn}},
	{Bool, "UintToBool", UintToBool, boolv, []any{uintn}},
	{Bool, "BoolToBool", BoolToBool, boolv, []any{boolv}},
	// TIME CONVERSION FUNCTIONS
	{Time, "StringToTime", StringToTime, timev, []any{strt}},
	{Time, "IntToTime", IntToTime, timev, []any{intn}},
	{Time, "FloatToTime", FloatToTime, timev, []any{floatn}},
	{Time, "UintToTime", UintToTime, timev, []any{uintn}},
	{Time, "TimeToTime", TimeToTime, timev, []any{timev}},
	// MAP CONVERSION FUNCTIONS
	{Map, "MapToMap", MapToMap, hmap, []any{hmap}},
	{Map, "ArrayToMap", ArrayToMap, hmapk, []any{arrayk}},
	{Map, "KeyValArraysToMap", KeyValArraysToMap, hmap, []any{arrayk, arrayv}},
	{Map, "KeyValPairsToMap", KeyValPairsToMap, hmapkv, []any{arraykv}},
	{Map, "StructToMap", StructToMap, hmap, []any{strct}},
	{Map, "JsonToMap", JsonToMap, hmap, []any{jsonv}},
	{vType, "validMapKeyType", validMapKeyType, true, []any{String}},
	{vType, "validMapKeyType", validMapKeyType, false, []any{Bool}},
	{aType, "MapKeyType", MapKeyType, Any, []any{hmap}},
	{aType, "MapValType", MapValType, Any, []any{hmap}},
	// ARRAY CONVERSION FUNCTIONS
	{Array, "MapKeys", MapKeys, arrayk, []any{hmap}},
	{Array, "MapVals", MapVals, arrayv, []any{hmap}},
	{Array, "StructFields", StructFields, arrayk, []any{strct}},
	{Array, "StructValues", StructValues, arrayvs, []any{strct}},
	{aType, "ArrayValType", ArrayValType, Any, []any{arraykv}},
	// STRUCT CONVERSION FUNCTIONS
	{Struct, "KeyValArraysToStruct", KeyValArraysToStruct, strctkv, []any{arrayk, array, stkv{}}},
	{Struct, "KeyValPairsToStruct", KeyValPairsToStruct, strctkv, []any{arraykv, stkv{}}},
	{Struct, "MapToStruct", MapToStruct, strct, []any{hmap, st{}}},
	{Struct, "JsonToStruct", JsonToStruct, strct, []any{jsonv, st{}}},
}

func TestConversions(t *testing.T) {
	tm := time.Now()
	for _, c := range convTests {
		fmt.Printf("Testing %s...", c.Name)
		r, e := runConvTest(c.Type, c.Func, c.Params)
		if e != nil {
			fmt.Print("  falied\n")
			t.Fatalf("%s unable handle conversion:\n%s", c.Name, e)
		}
		if c.Type == Array {
			r = sortArr(r)
			c.Validator = sortArr(c.Validator)
		}
		if !reflect.DeepEqual(r, c.Validator) {
			fmt.Print("  failed\n")
			t.Fatalf("%s result error:\nexpected: %#v\nreturned: %#v", c.Name, c.Validator, r)
		}
		fmt.Print("  completed\n")
	}
	d := time.Now().Sub(tm)
	fmt.Println(d)
}

func runConvTest(t Type, f any, p []any) (any, error) {
	switch t {
	case vType:
		return f.(func(Type) bool)(p[0].(Type)), nil
	case aType:
		return f.(func(any) (Type, error))(p[0])
	case String:
		return f.(func(any) (string, error))(p[0])
	case Int:
		return f.(func(any) (int, error))(p[0])
	case Float:
		return f.(func(any) (float64, error))(p[0])
	case Uint:
		return f.(func(any) (uint, error))(p[0])
	case Bool:
		return f.(func(any) (bool, error))(p[0])
	case Time:
		return f.(func(any) (time.Time, error))(p[0])
	case Map:
		switch len(p) {
		case 2:
			return f.(func(any, any) (map[any]any, error))(p[0], p[1])
		default:
			return f.(func(any) (map[any]any, error))(p[0])
		}
	case Array:
		return f.(func(any) ([]any, error))(p[0])
	case Struct:
		switch len(p) {
		case 3:
			return f.(func(any, any, ...any) (any, error))(p[0], p[1], p[2])
		default:
			return f.(func(any, ...any) (any, error))(p[0], p[1])
		}
	default:
		return nil, nil
	}
}

func sortArr(a any) []any {
	if !IsArray(a) {
		panic("trying to sort non array")
	}
	s := []string{}
	r := reflect.ValueOf(a)
	for i := 0; i < r.Len(); i++ {
		s = append(s, r.Index(i).String())
	}
	sort.Strings(s)
	a = []any{}
	for _, v := range s {
		a = append(a.([]any), v)
	}
	return a.([]any)
}

/*
func check(r any, t any) {
	fmt.Println("\nStarting checker...")
	rr, rok := r.(map[any]any)
	tt, tok := t.(map[any]any)
	if !rok || !tok {
		fmt.Println("type mismatch")
		fmt.Printf("r: %T\n", r)
		fmt.Printf("t: %T\n", t)
		return
	}
	if len(tt) != len(rr) {
		fmt.Println("len mismatch")
		fmt.Printf("r: %d\n", len(rr))
		fmt.Printf("t: %d\n", len(tt))
		return
	}
	for k, v := range rr {
		if _, ok := tt[k]; !ok {
			fmt.Println("missing key")
			fmt.Printf("r: %v\n", k)
			return
		}
		if _, vok := v.(map[any]any); vok {
			fmt.Println("found submap")
			check(v, tt[k])
		} else if va, ia := v.([]any); ia {
			for i, av := range va {
				tta, _ := tt[k].([]any)
				if tta[i] != av {
					fmt.Println("sub array value mismatch")
					fmt.Printf("r: %v: %v\n", i, av)
					fmt.Printf("t:  %v: %v\n", i, tta[i])
					return
				}
			}
		} else if tt[k] != v {
			fmt.Println("value mismatch")
			fmt.Printf("r: %v: %v type: %T\n", k, v, v)
			fmt.Printf("t:  %v: %v type: %T\n", k, tt[k], tt[k])
			return
		}
	}
	fmt.Println("validated")
}
*/
