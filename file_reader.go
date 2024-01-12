package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
	"unsafe"
)

type ChunkedFileReader struct {
	mutex  sync.Mutex
	file   *os.File
	reader *bufio.Reader

	maxReadBytes int64
	bytesRead    int64
}

func NewChunkedFileReader(fileName string, offset, maxReadBytes uint64) *ChunkedFileReader {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	if offset > 0 {
		_, err := file.Seek(int64(offset), io.SeekStart)
		if err != nil {
			panic(err)
		}
	}

	return &ChunkedFileReader{
		file:         file,
		reader:       bufio.NewReader(file),
		maxReadBytes: int64(maxReadBytes),
		bytesRead:    0,
	}
}

func (o *ChunkedFileReader) Close() {
	o.file.Close()
}

func (o *ChunkedFileReader) GetLine() (string, bool, error) {
	if o.bytesRead >= o.maxReadBytes {
		return "", false, nil
	}

	lineBytes, err := o.reader.ReadBytes('\n')

	if err != nil {
		if !errors.Is(err, io.EOF) {
			return "", false, err
		}

		if len(lineBytes) == 0 {
			return "", false, nil
		}
	}

	o.bytesRead += int64(len(lineBytes))
	s := unsafe.String(unsafe.SliceData(lineBytes), len(lineBytes)-1)
	return s, true, nil
}
