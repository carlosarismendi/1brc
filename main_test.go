package main

import (
	"hash/maphash"
	"testing"

	"unsafe"
)

func BenchmarkMapAccess(b *testing.B) {

	seed := maphash.MakeSeed()

	b.Run("map[uint64]*Station", func(b *testing.B) {
		line := "line"
		m := make(map[uint64]*Station)
		for i := 0; i < b.N; i++ {
			b := unsafe.Slice(unsafe.StringData(line), len(line))
			n := maphash.Bytes(seed, b)

			m[n] = &Station{}

			line += "1"
			if i%100 == 0 {
				line = "line"
			}
		}
	})

	b.Run("4*map[uint64]int", func(b *testing.B) {
		line := "line"

		m1 := make(map[uint64]int)
		m2 := make(map[uint64]int)
		m3 := make(map[uint64]int)
		m4 := make(map[uint64]int)
		for i := 0; i < b.N; i++ {
			b := unsafe.Slice(unsafe.StringData(line), len(line))
			n := maphash.Bytes(seed, b)

			m1[n] = min(m1[n], i)
			m2[n] = max(m2[n], i)
			m3[n]++
			m4[n]++

			if i%100 == 0 {
				line = "line"
			}
		}
	})
}
