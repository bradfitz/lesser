// Copyright 2020 Brad Fitzpatrick. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lesser generates less functions for sort.Slice.
package lesser

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// Of returns a less function suitable to passing to sort.Slice.
//
// The slice argument must be a slice.
//
// The ordering rules are more general than with Go's < operator:
//
//  - bool compares false before true
//  - ints, floats, and strings order by <
//  - NaN compares less than non-NaN floats
//  - complex compares real, then imag
//  - pointers, chan, func and map compare by
//    machine address
//  - structs compare each field in turn
//  - arrays compare each non-blank element in turn
//
// Performance should be comparable to writing a native sort.Slice
// function.
func Of(slice interface{}) (less func(i, j int) bool) {
	rv := reflect.ValueOf(slice)
	t := rv.Type()
	if t.Kind() != reflect.Slice {
		panic("slice argument is not a slice")
	}
	if rv.Len() == 0 {
		return nil // won't be called
	}
	et := t.Elem()
	addr0 := unsafe.Pointer(rv.Index(0).UnsafeAddr())
	return forAddr(addr0, et.Size(), 0, et, nil)
}

func forAddr(addr0 unsafe.Pointer, size, off uintptr, t reflect.Type, optEq less) less {
	var makeLess func(addr0 unsafe.Pointer, size, off uintptr, optEq less) less
	switch t.Kind() {
	case reflect.Bool:
		makeLess = lessBool
	case reflect.Int:
		makeLess = lessInt
	case reflect.Int8:
		makeLess = lessInt8
	case reflect.Int16:
		makeLess = lessInt16
	case reflect.Int32:
		makeLess = lessInt32
	case reflect.Int64:
		makeLess = lessInt64
	case reflect.Uint:
		makeLess = lessUint
	case reflect.Uint8:
		makeLess = lessUint8
	case reflect.Uint16:
		makeLess = lessUint16
	case reflect.Uint32:
		makeLess = lessUint32
	case reflect.Uint64:
		makeLess = lessUint64
	case reflect.Uintptr:
		makeLess = lessUintptr
	case reflect.Float32:
		makeLess = lessFloat32
	case reflect.Float64:
		makeLess = lessFloat64
	case reflect.Complex64:
		makeLess = lessComplex64
	case reflect.Complex128:
		makeLess = lessComplex128
	case reflect.Array:
		ret := optEq
		et := t.Elem()
		for i := t.Len() - 1; i >= 0; i-- {
			ret = forAddr(addr0, size, et.Size()*uintptr(i), et, ret)
		}
		return ret
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		makeLess = lessUintptr
	case reflect.String:
		makeLess = lessString
	case reflect.Struct:
		// Walk fields from the back, building up the
		// tie-breaker chain in reverse.
		ret := optEq
		for i := t.NumField() - 1; i >= 0; i-- {
			sf := t.Field(i)
			if sf.Name == "_" {
				continue
			}
			ret = forAddr(addr0, size, sf.Offset, sf.Type, ret)
		}
		return ret
	case reflect.Interface:
		// TODO
	case reflect.Slice:
		// TODO
	}
	if makeLess == nil {
		panic(fmt.Sprintf("un-sortable type %v (kind %v)", t, t.Kind()))
	}
	return makeLess(addr0, size, off, optEq)
}

type less func(i, j int) bool

func addr(addr0 unsafe.Pointer, size, off uintptr, i int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(addr0) + size*uintptr(i) + off)
}

func lessBool(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*bool)(addr(addr0, size, off, i)), *(*bool)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va == false
	}
}

func lessString(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*string)(addr(addr0, size, off, i)), *(*string)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessInt(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*int)(addr(addr0, size, off, i)), *(*int)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessInt8(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*int8)(addr(addr0, size, off, i)), *(*int8)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessInt16(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*int16)(addr(addr0, size, off, i)), *(*int16)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessInt32(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*int32)(addr(addr0, size, off, i)), *(*int32)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessInt64(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*int64)(addr(addr0, size, off, i)), *(*int64)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUint(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uint)(addr(addr0, size, off, i)), *(*uint)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUint8(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uint8)(addr(addr0, size, off, i)), *(*uint8)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUint16(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uint16)(addr(addr0, size, off, i)), *(*uint16)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUint32(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uint32)(addr(addr0, size, off, i)), *(*uint32)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUint64(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uint64)(addr(addr0, size, off, i)), *(*uint64)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessUintptr(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*uintptr)(addr(addr0, size, off, i)), *(*uintptr)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb
	}
}

func lessFloat32(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*float32)(addr(addr0, size, off, i)), *(*float32)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb || isNaN32(va) && !isNaN32(vb)
	}
}

func lessFloat64(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return func(i, j int) bool {
		va, vb := *(*float64)(addr(addr0, size, off, i)), *(*float64)(addr(addr0, size, off, j))
		if va == vb {
			if optEq != nil {
				return optEq(i, j)
			}
			return false
		}
		return va < vb || math.IsNaN(va) && !math.IsNaN(vb)
	}
}

func lessComplex64(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return lessFloat32(addr0, size, off, lessFloat32(addr0, size, off+4, optEq))
}

func lessComplex128(addr0 unsafe.Pointer, size, off uintptr, optEq less) less {
	return lessFloat64(addr0, size, off, lessFloat64(addr0, size, off+8, optEq))
}

func isNaN32(f float32) bool { return f != f }
