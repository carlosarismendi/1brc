package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	Workers = 1
)

func processLine(line string) (station string, temperature float64) {
	i := strings.LastIndex(line, ";")
	temp, err := strconv.ParseFloat(line[i+1:], 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing temperature for line %q. Error=%q.", line, err.Error()))
	}
	return line[:i], temp
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
		sb.WriteString(fmt.Sprintf("%s%s=%.1f/%.1f/%.1f", sep, v.Name, v.Min, ((10*v.Sum)/float64(v.Count))/10, v.Max))
		sep = ", "
	}
	sb.WriteString("}")
	fmt.Println(sb.String())
}

func workerOneBrc(o *ChunkedFileReader, output chan<- *StationMap) {
	stationsMap := NewStationMap()

	for {
		line, ok, err := o.GetLine()
		if ok {
			name, temperature := processLine(line)
			stationsMap.Add(name, temperature)
			continue
		}

		if err != nil {
			panic(err)
		}
		break
	}

	output <- stationsMap
}

func oneBrc(measurementsFile string) map[string]*Station {
	var fileSize int64
	func() {
		cfr := NewChunkedFileReader(measurementsFile, 0, 10)
		defer cfr.Close()

		fileStat, err := cfr.file.Stat()
		if err != nil {
			panic(err)
		}

		fileSize = fileStat.Size()
	}()

	chunkSize := fileSize / int64(Workers)

	workerMaps := make(chan *StationMap, Workers)
	offset := int64(0)

	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(Workers)
	for i := 0; i < Workers; i++ {
		go func(offsetWorker int64) {
			defer wgWorkers.Done()

			if i == Workers-1 {
				chunkSize = fileSize
			}

			cfr := NewChunkedFileReader(measurementsFile, uint64(offsetWorker), uint64(chunkSize))
			defer cfr.Close()

			// Set the offset at the beginning of the next line
			if offsetWorker != 0 {
				_, err := cfr.reader.ReadBytes('\n')
				if err != nil && !errors.Is(err, io.EOF) {
					panic(err)
				}
			}

			workerOneBrc(cfr, workerMaps)
		}(offset)

		offset += chunkSize
	}

	var stationsMap map[string]*Station

	wgReducer := sync.WaitGroup{}
	wgReducer.Add(1)
	go func() {
		defer wgReducer.Done()

		sm := <-workerMaps
		for wsm := range workerMaps {
			sm.Merge(wsm)
		}

		stationsMap = sm.m
	}()

	wgWorkers.Wait()
	close(workerMaps)

	wgReducer.Wait()
	return stationsMap
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
