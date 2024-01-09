package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	Workers  = 20
	Reducers = 1
)

type Station struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

type StationMap struct {
	mutex sync.Mutex
	m     map[string]*Station
}

func NewStationMap() *StationMap {
	return &StationMap{
		m: make(map[string]*Station),
	}
}

func (sm *StationMap) Update(name string, temperature float64) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	s, ok := sm.m[name]
	if ok {
		if temperature < s.Min {
			s.Min = temperature
		}

		if temperature > s.Max {
			s.Max = temperature
		}

		s.Sum += temperature
		s.Count++
	} else {
		sm.m[name] = &Station{
			Name:  name,
			Min:   temperature,
			Max:   temperature,
			Sum:   temperature,
			Count: 1,
		}
	}
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

	var sb strings.Builder
	sb.WriteString("{")
	sep := ""
	for _, v := range stations {
		sb.WriteString(fmt.Sprintf("%s%s=%.1f/%.1f/%.1f", sep, v.Name, v.Min, ((10*v.Sum)/float64(v.Count))/10, v.Max))
		sep = ", "
	}
	sb.WriteString("}")
	fmt.Println(sb.String())
}

func worker(wg *sync.WaitGroup, input <-chan string, stationsMap *StationMap) {
	defer wg.Done()

	for line := range input {
		name, temperature := processLine(line)
		stationsMap.Update(name, temperature)
	}
}


func main() {
	measurementsFile := os.Getenv("MEASUREMENTS_FILE")
	if measurementsFile == "" {
		panic("MEASUREMENTS_FILE environment variable not set")
	}

	stationsMap := NewStationMap()

	linesCh := make(chan string, 10000)

	wgWorkers := &sync.WaitGroup{}
	wgWorkers.Add(Workers)
	for i := 0; i < Workers; i++ {
		go worker(wgWorkers, linesCh, stationsMap)
	}
	
	// Open the file
	file, err := os.Open(measurementsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		linesCh <- scanner.Text()
	}

	close(linesCh)
	wgWorkers.Wait()


	// Check for errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	printResults(stationsMap.m)
	fmt.Println("")
}
