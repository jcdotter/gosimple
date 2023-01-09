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
	phraseStr = "Example string value"
)

func TestCasing(t *testing.T) {
	r := ToCamelCase(snakeStr)
	if r != camelStr {
		t.Fatalf("string.ToCamelCase unable to convert \nsnake case string '%s' to camel case string '%s'\n returned: '%s'", snakeStr, camelStr, r)
	}
	r = ToCamelCase(pascalStr)
	if r != camelStr {
		t.Fatalf("string.ToCamelCase unable to convert \npascal case string '%s' to camel case string '%s'\n returned: '%s'", pascalStr, camelStr, r)
	}
	r = ToPascalCase(snakeStr)
	if r != pascalStr {
		t.Fatalf("string.ToPascalCase unable to convert \nsnake case string '%s' to pascal case string '%s'\n returned: '%s'", snakeStr, pascalStr, r)
	}
	r = ToPascalCase(camelStr)
	if r != pascalStr {
		t.Fatalf("string.ToPascalCase unable to convert \ncamel case string '%s' to pascal case string '%s'\n returned: '%s'", camelStr, pascalStr, r)
	}
	r = ToSnakeCase(camelStr)
	if r != snakeStr {
		t.Fatalf("string.CamelToSnake unable to convert \ncamel case string '%s' to snake case string '%s'\n returned: '%s'", camelStr, snakeStr, r)
	}
	r = ToSnakeCase(pascalStr)
	if r != snakeStr {
		t.Fatalf("string.PascalToSnake unable to convert \npascal case string '%s' to snake case string '%s'\n returned: '%s'", pascalStr, snakeStr, r)
	}
	r = ToPhraseCase(camelStr, true)
	if r != phraseStr {
		t.Fatalf("string.ToPhraseCase unable to convert \ncamel case string '%s' to phrase case string '%s'\n returned: '%s'", camelStr, phraseStr, r)
	}
	r = ToPhraseCase(pascalStr, true)
	if r != phraseStr {
		t.Fatalf("string.ToPhraseCase unable to convert \npascal case string '%s' to phrase case string '%s'\n returned: '%s'", pascalStr, phraseStr, r)
	}
	r = ToPhraseCase(snakeStr, true)
	if r != phraseStr {
		t.Fatalf("string.ToPhraseCase unable to convert \nsnake case string '%s' to phrase case string '%s'\n returned: '%s'", snakeStr, phraseStr, r)
	}
}
