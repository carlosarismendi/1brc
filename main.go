package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"runtime/trace"
	"sync"
)

func workerOneBrc(o *ChunkedFileReader, output chan<- *StationMap) {
	stationsMap := NewStationMap()

	for {
		nameBytes, tempBytes, ok, err := o.GetLine()
		if ok {
			name, temperature := processLine(nameBytes, tempBytes)
			stationsMap.Add(name, nameBytes, temperature)
			continue
		}

		if err != nil {
			panic(err)
		}
		break
	}

	output <- stationsMap
}

func oneBrc(measurementsFile string, maxWorkers, maxRam int) *StationMap {
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

	// log.Printf("workers=%v, chunk=%v, ram=%v, fileSize=%v", maxWorkers, chunkSize, maxRam, fileSize)
	workerMaps := make(chan *StationMap, maxWorkers)
	var stationsMap *StationMap
	wgReducer := sync.WaitGroup{}
	wgReducer.Add(1)
	go func() {
		defer wgReducer.Done()

		sm := <-workerMaps
		for wsm := range workerMaps {
			sm.Merge(wsm)
		}

		stationsMap = sm
	}()

	wgWorkers := sync.WaitGroup{}

	left := fileSize - chunkSize
	right := fileSize
	quit := false
	for !quit {
		// log.Printf("left: %v, right: %v\n", left, right)
		if left < 0 {
			left = 0
			quit = true
		}
		cfr := NewChunkedFileReader(measurementsFile, left, right)

		// Set the offset at the beginning of the next line
		var diff int64
		if left > 0 {
			n, err := cfr.MoveReaderToStartOfNextLine()
			if err != nil {
				panic(err)
			}
			diff = n
		}

		wgWorkers.Add(1)
		go func(workerCfr *ChunkedFileReader) {
			defer workerCfr.Close()
			defer wgWorkers.Done()

			err := workerCfr.MMap()
			if err != nil {
				panic(err)
			}
			workerOneBrc(workerCfr, workerMaps)
		}(cfr)

		right = left + diff
		left = left - chunkSize
	}

	wgWorkers.Wait()
	close(workerMaps)

	wgReducer.Wait()
	return stationsMap
}

const GB = 1024 * 1024 * 1024

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var traceprofile = flag.String("traceprofile", "", "write trace profile to file")
var measurementsFile = flag.String("measurements-file", "measurements.txt", "measurements file")
var maxRam = flag.Int("max-ram", 2, "max ram to use (GB)")
var maxWorkers = flag.Int("max-workers", 1, "max workers to use")

func main() {
	flag.Parse()
	// log.Printf("file: %v, ram: %v, workers: %v\n", *measurementsFile, *maxRam, *maxWorkers)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *traceprofile != "" {
		f, err := os.Create(*traceprofile)
		if err != nil {
			log.Fatal(err)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	debug.SetMemoryLimit(int64(*maxRam * GB))
	results := oneBrc(*measurementsFile, *maxWorkers, *maxRam*GB)

	printResults(results)
}
