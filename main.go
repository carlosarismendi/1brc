package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

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

func oneBrc(measurementsFile string, maxWorkers, maxRam int) map[string]*Station {
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

	chunkSize := fileSize / int64(maxWorkers)
	// log.Printf("File size: %v, chunk size: %v\n", fileSize, chunkSize)

	workerMaps := make(chan *StationMap, maxWorkers)
	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(maxWorkers)

	left := fileSize - chunkSize

	right := fileSize
	for i := maxWorkers; i > 0; i-- {
		if i == 1 {
			left = 0
			// log.Printf("(left, right)=(%v, %v)\n", left, right)
		}
		cfr := NewChunkedFileReader(measurementsFile, uint64(left), uint64(right))

		// Set the offset at the beginning of the next line
		var diff int64
		if left > 0 {
			b, err := cfr.reader.ReadBytes('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}

			diff = int64(len(b))
			// log.Printf("(left, right)=(%v, %v)\n", left+diff, right)
		}

		go func(workerCfr *ChunkedFileReader) {
			defer workerCfr.Close()
			defer wgWorkers.Done()

			workerOneBrc(workerCfr, workerMaps)
		}(cfr)

		right = left + diff
		left = left - chunkSize
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
var measurementsFile = flag.String("measurements-file", "measurements.txt", "measurements file")
var maxRam = flag.Int("max-ram", 2, "max ram to use (GB)")
var maxWorkers = flag.Int("max-workers", 1, "max workers to use")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	results := oneBrc(*measurementsFile, *maxWorkers, *maxRam)

	printResults(results)
}
