package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Station struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func processLine(line string) (station string, temperature float64) {
	v := strings.Split(line, ";")
	temperature, err := strconv.ParseFloat(v[1], 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing temperature for line %q. Error=%q.", line, err.Error()))
	}
	return v[0], temperature
}

func printResults(stationsMap map[string]*Station) {
	stations := make([]*Station, 0, len(stationsMap))
	for _, v := range stationsMap {
		stations = append(stations, v)
		delete(stationsMap, v.Name)
	}

	sort.Slice(stations, func(i, j int) bool {
		return stations[i].Name < stations[j].Name
	})

	fmt.Print("{")
	sep := ""
	for _, v := range stations {
		fmt.Printf("%s%s=%.1f/%.1f/%.1f", sep, v.Name, v.Min, ((10*v.Sum)/float64(v.Count))/10, v.Max)
		sep = ", "
	}
	fmt.Print("}\n")
}

func main() {
	measurementsFile := os.Getenv("MEASUREMENTS_FILE")
	if measurementsFile == "" {
		panic("MEASUREMENTS_FILE environment variable not set")
	}

	// Open the file
	file, err := os.Open(measurementsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	stationsMap := make(map[string]*Station)

	// Read the file line by line
	for scanner.Scan() {
		// Print the line
		line := scanner.Text()
		stationName, temperature := processLine(line)
		s, ok := stationsMap[stationName]
		if ok {
			s.Min = min(s.Min, temperature)
			s.Max = max(s.Max, temperature)
			s.Sum += temperature
			s.Count++
			continue
		}
		stationsMap[stationName] = &Station{Name: stationName, Min: temperature, Max: temperature, Sum: temperature, Count: 1}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	printResults(stationsMap)
}
