// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package array

import "testing"

var (
	array = []int{1, 2, 3, 4}
	tVals = []any{2, 0, 3}
)

func TestContains(t *testing.T) {
	if !Contains(array, tVals...) {
		t.Fatalf("array.Contains unable to match \nvals: %v\n in array: %v", tVals, array)
	}
}
