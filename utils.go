package main

import (
	"fmt"
	"math"
	"strings"
	"unsafe"
)

func processLine(nameBytes, tempBytes []byte) (station string, temperature float64) {
	temp := parseFloat(tempBytes)
	name := unsafe.String(unsafe.SliceData(nameBytes), len(nameBytes))
	return name, temp
}

func parseFloat(s []byte) float64 {
	var sign int64
	var idx int
	if s[0] == '-' {
		sign = -1
		idx = 1
	} else {
		sign = 1
	}

	num := int64(0)
	for ; idx < len(s); idx++ {
		if s[idx] == '.' {
			continue
		}

		num = num*10 + int64(s[idx]-'0')
	}

	return float64(sign*num) / 10.0
}

func round(f float64) float64 {
	const ratio = 10.0
	return math.Round(f*ratio) / ratio
}

func printResults(stationsMap *StationMap) {
	var sb strings.Builder
	sb.WriteString("{")

	sb.WriteString(stationsMap.String())

	sb.WriteString("}")
	fmt.Println(sb.String())
}
