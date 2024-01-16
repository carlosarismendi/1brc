package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMapAccess(b *testing.B) {
	arr := []string{"1.1", "22.0", "33.8", "-4.0", "-85.7"}
	b.Run("strconv.ParseFloat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range arr {
				n, _ := strconv.ParseFloat(string(v), 64)
				n = n
			}
		}
	})

	arr2 := [][]byte{[]byte("1.1"), []byte("22.0"), []byte("33.8"), []byte("-4.0"), []byte("-85.7")}
	b.Run("parseFloat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range arr2 {
				n := parseFloat(v)
				n = n
			}
		}
	})
}

func TestParseFloat(t *testing.T) {
	tests := [][]byte{[]byte("1.1"), []byte("22.0"), []byte("33.8"), []byte("-4.0"), []byte("-85.7")}
	expected := []float64{1.1, 22.0, 33.8, -4.0, -85.7}

	for i := range tests {
		v := tests[i]
		t.Run(string(v), func(t *testing.T) {
			n := parseFloat(v)
			assert.Equal(t, expected[i], n)
		})
	}
}
