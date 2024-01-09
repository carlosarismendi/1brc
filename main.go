package main

import (
	"bufio"
	"fmt"
	"log"
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

func worker(wg *sync.WaitGroup, input <-chan string, output chan<- *Station) {
	defer wg.Done()

	for line := range input {
		name, temperature := processLine(line)
		output <- &Station{
			Name:  name,
			Min:   temperature,
			Max:   temperature,
			Sum:   temperature,
			Count: 1,
		}
	}
}

func reducer(wg *sync.WaitGroup, stationsMap map[string]*Station, input <-chan *Station) {
	defer wg.Done()

	for station := range input {
		name := station.Name
		temperature := station.Min

		s, ok := stationsMap[name]
		if ok {
			s.Min = min(s.Min, temperature)
			s.Max = max(s.Max, temperature)
			s.Sum += temperature
			s.Count++
			continue
		}
		stationsMap[name] = station
	}
}

func main() {
	measurementsFile := os.Getenv("MEASUREMENTS_FILE")
	if measurementsFile == "" {
		panic("MEASUREMENTS_FILE environment variable not set")
	}

	stationsMap := make(map[string]*Station)

	linesCh := make(chan string, 10000)
	stationsCh := make(chan *Station, 1)

	wgWorkers := &sync.WaitGroup{}
	wgWorkers.Add(Workers)
	for i := 0; i < Workers; i++ {
		go worker(wgWorkers, linesCh, stationsCh)
	}
	log.Println("Workers started.")

	wgReducers := &sync.WaitGroup{}
	wgReducers.Add(Reducers)
	for i := 0; i < Reducers; i++ {
		go reducer(wgReducers, stationsMap, stationsCh)
	}
	log.Println("Reducers started.")

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

	log.Println("Scanner finished.")

	close(linesCh)
	wgWorkers.Wait()
	log.Println("Workers finished.")

	close(stationsCh)
	wgReducers.Wait()
	log.Println("Reducers finished.")

	// Check for errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	printResults(stationsMap)
	fmt.Println("")
}
