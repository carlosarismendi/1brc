package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unsafe"
)

func processLine(lineBytes []byte) (station string, nameBytes []byte, temperature float64) {
	// line := unsafe.String(unsafe.SliceData(lineBytes), len(lineBytes))
	sep := bytes.LastIndexByte(lineBytes, ';')

	tempBytes := lineBytes[sep+1:]
	temp := parseFloat(tempBytes)

	nameBytes = lineBytes[:sep]
	name := unsafe.String(unsafe.SliceData(nameBytes), len(nameBytes))
	return name, nameBytes, temp
}

func parseFloat(s []byte) float64 {
	idx := 0
	sign := int64(1)
	if s[0] == '-' {
		sign = -1
		idx = 1
	}

	num := int64(0)
	for ; idx < len(s); idx++ {
		if s[idx] == '.' {
			idx++
			break
		}

		num = num*10 + int64(s[idx]-'0')
	}

	for ; idx < len(s); idx++ {
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
