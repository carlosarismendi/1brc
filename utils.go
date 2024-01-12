package main

import (
	"fmt"
	"math"
	"sort"
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

func printResults(stationsMap map[string]*Station) {
	stations := make([]*Station, 0, len(stationsMap))
	for k := range stationsMap {
		stations = append(stations, stationsMap[k])
		delete(stationsMap, k)
	}

	sort.Slice(stations, func(i, j int) bool {
		return stations[i].Name < stations[j].Name
	})

	var sb strings.Builder
	sb.WriteString("{")
	sep := ""
	for _, v := range stations {
		sb.WriteString(fmt.Sprintf("%s%s", sep, v.String()))
		sep = ", "
	}
	sb.WriteString("}")
	fmt.Println(sb.String())
}
