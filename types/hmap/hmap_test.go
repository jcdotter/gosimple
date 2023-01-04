// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package hmap

import (
	"testing"
)

var (
	hmap  = map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}
	tKeys = []any{"two", "three"}
	tVals = []any{2, 3}
)

func TestContainsKeys(t *testing.T) {
	if !ContainsKeys(hmap, tKeys...) {
		t.Fatalf("hmap.ContainsKeys unable to match \nkeys: %v\n in map: %v", tKeys, hmap)
	}
}

func TestContainsVals(t *testing.T) {
	if !ContainsVals(hmap, tVals...) {
		t.Fatalf("hmap.ContainsVals unable to match \nvals: %v\n in map: %v", tVals, hmap)
	}
}
