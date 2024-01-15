package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

func processLine(lineBytes []byte) (station string, nameBytes []byte, temperature float64) {
	// line := unsafe.String(unsafe.SliceData(lineBytes), len(lineBytes))
	sep := bytes.LastIndexByte(lineBytes, ';')

	tempBytes := lineBytes[sep+1:]
	tempBytesStr := unsafe.String(unsafe.SliceData(tempBytes), len(tempBytes))
	temp, err := strconv.ParseFloat(tempBytesStr, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing temperature for line %q. Error=%q.", string(lineBytes), err.Error()))
	}

	nameBytes = lineBytes[:sep]
	name := unsafe.String(unsafe.SliceData(nameBytes), len(nameBytes))
	return name, nameBytes, temp
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
