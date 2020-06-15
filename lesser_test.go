// Copyright 2020 Brad Fitzpatrick. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lesser

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

type TStringInt struct {
	S string
	I int
}

func TestOf(t *testing.T) {
	tests := []struct {
		name     string
		in, want interface{}
	}{
		{
			name: "int",
			in:   []int{2, 4, 1, 3, 0, -1, 5},
			want: []int{-1, 0, 1, 2, 3, 4, 5},
		},
		{
			name: "string",
			in:   []string{"foo", "quux", "baz", "bar"},
			want: []string{"bar", "baz", "foo", "quux"},
		},
		{
			name: "string_int",
			in:   []TStringInt{{"a", 2}, {"b", 2}, {"b", 1}, {"a", 1}},
			want: []TStringInt{{"a", 1}, {"a", 2}, {"b", 1}, {"b", 2}},
		},
		{
			name: "bool",
			in:   []bool{false, true, false, false, true},
			want: []bool{false, false, false, true, true},
		},
		{
			name: "complex64",
			in:   []complex64{complex(1, 2), complex(2, 1), complex(1, 1), complex(2, 2)},
			want: []complex64{complex(1, 1), complex(1, 2), complex(2, 1), complex(2, 2)},
		},
		{
			name: "complex128",
			in:   []complex128{complex(1, 2), complex(2, 1), complex(1, 1), complex(2, 2)},
			want: []complex128{complex(1, 1), complex(1, 2), complex(2, 1), complex(2, 2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lesser := Of(tt.in)
			sort.Slice(tt.in, lesser)
			if !reflect.DeepEqual(tt.in, tt.want) {
				t.Errorf("wrong:\n got: %v\nwant: %v\n", tt.in, tt.want)
			}
		})
	}
}

func BenchmarkStructSort_native(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(123)
	unsorted := make([]TStringInt, 10000)
	for i := range unsorted {
		unsorted[i].S = fmt.Sprint(rand.Intn(1e9))
		unsorted[i].I = rand.Intn(1e9)
	}
	buf := make([]TStringInt, len(unsorted))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, unsorted)
		sort.Slice(buf, func(i, j int) bool {
			va, vb := &unsorted[i], &unsorted[j]
			if va.S == vb.S {
				return va.I < vb.I
			}
			return va.S < vb.S
		})
	}
}

func BenchmarkStructSort_lesser(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(123)
	unsorted := make([]TStringInt, 10000)
	for i := range unsorted {
		unsorted[i].S = fmt.Sprint(rand.Intn(1e9))
		unsorted[i].I = rand.Intn(1e9)
	}
	buf := make([]TStringInt, len(unsorted))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(buf, unsorted)
		sort.Slice(buf, Of(buf))
	}
}

func BenchmarkStructSort_lesser_reuse(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(123)
	unsorted := make([]TStringInt, 10000)
	for i := range unsorted {
		unsorted[i].S = fmt.Sprint(rand.Intn(1e9))
		unsorted[i].I = rand.Intn(1e9)
	}
	buf := make([]TStringInt, len(unsorted))
	b.ResetTimer()
	lesser := Of(buf)
	for i := 0; i < b.N; i++ {
		copy(buf, unsorted)
		sort.Slice(buf, lesser)
	}
}
