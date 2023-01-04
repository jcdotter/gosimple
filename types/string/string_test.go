// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package string

import (
	"testing"
)

var (
	snakeStr  = "example_string_value"
	camelStr  = "exampleStringValue"
	pascalStr = "ExampleStringValue"
)

func TestCasing(t *testing.T) {
	r := SnakeToCamel(snakeStr)
	if r != camelStr {
		t.Fatalf("string.SnakeToCamel unable to convert \nsnake case string '%s' to camel case string '%s'\n returned: '%s'", snakeStr, camelStr, r)
	}
	r = SnakeToPascal(snakeStr)
	if r != pascalStr {
		t.Fatalf("string.SnakeToPascal unable to convert \nsnake case string '%s' to pascal case string '%s'\n returned: '%s'", snakeStr, pascalStr, r)
	}
	r = PascalToCamel(pascalStr)
	if r != camelStr {
		t.Fatalf("string.PascalToCamel unable to convert \npascal case string '%s' to camel case string '%s'\n returned: '%s'", pascalStr, camelStr, r)
	}
	r = CamelToPascal(camelStr)
	if r != pascalStr {
		t.Fatalf("string.CamelToPascal unable to convert \ncamel case string '%s' to pascal case string '%s'\n returned: '%s'", camelStr, pascalStr, r)
	}
	r = CamelToSnake(camelStr)
	if r != snakeStr {
		t.Fatalf("string.CamelToSnake unable to convert \ncamel case string '%s' to snake case string '%s'\n returned: '%s'", camelStr, snakeStr, r)
	}
	r = PascalToSnake(pascalStr)
	if r != snakeStr {
		t.Fatalf("string.PascalToSnake unable to convert \npascal case string '%s' to snake case string '%s'\n returned: '%s'", pascalStr, snakeStr, r)
	}
}
