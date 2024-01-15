package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func processLine(line string) (station string, temperature float64) {
	i := strings.LastIndex(line, ";")
	temp, err := strconv.ParseFloat(line[i+1:], 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing temperature for line %q. Error=%q.", line, err.Error()))
	}
	return line[:i], temp
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
