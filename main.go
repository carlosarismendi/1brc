package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	Workers = 20
)

type OneBRC struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewOneBRC(fileName string) *OneBRC {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	return &OneBRC{
		file:    file,
		scanner: scanner,
	}
}

func (o *OneBRC) Close() {
	o.file.Close()
}

func (o *OneBRC) GetLine() (string, bool, error) {
	ok := o.scanner.Scan()
	if ok {
		return o.scanner.Text(), true, nil
	}

	return "", false, o.scanner.Err()
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
	for k := range stationsMap {
		s := stationsMap[k]
		stations = append(stations, s)
		delete(stationsMap, s.Name)
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

func worker(wg *sync.WaitGroup, input <-chan string, stations chan<- Station) {
	defer wg.Done()

	for line := range input {
		name, temperature := processLine(line)
		stations <- Station{
			Name:  name,
			Min:   temperature,
			Max:   temperature,
			Sum:   temperature,
			Count: 1,
		}
	}
}

func oneBrc(measurementsFile string) map[string]*Station {
	o := NewOneBRC(measurementsFile)
	defer o.Close()

	stationsMap := NewStationMap()
	linesCh := make(chan string, 1000)

	stationsCh := make(chan Station, 1000)
	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(Workers)
	for i := 0; i < Workers; i++ {
		go worker(&wgWorkers, linesCh, stationsCh)
	}

	wgReducer := sync.WaitGroup{}
	wgReducer.Add(1)
	go func() {
		defer wgReducer.Done()
		for station := range stationsCh {
			stationsMap.Update(station)
		}
	}()

	// Read the file line by line
	for {
		line, ok, err := o.GetLine()
		if !ok {
			if err != nil {
				panic(err)
			}
			break
		}

		linesCh <- line
	}

	close(linesCh)
	wgWorkers.Wait()

	close(stationsCh)
	wgReducer.Wait()

	return stationsMap.m
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	measurementsFile := os.Getenv("MEASUREMENTS_FILE")
	if measurementsFile == "" {
		panic("MEASUREMENTS_FILE environment variable not set")
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	results := oneBrc(measurementsFile)

	printResults(results)
	fmt.Println("")
}
