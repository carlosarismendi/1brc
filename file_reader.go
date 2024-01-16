package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
)

type ChunkedFileReader struct {
	mutex  sync.Mutex
	file   *os.File
	reader *bufio.Reader

	offset   int64
	maxBytes int64

	text []byte
}

func NewChunkedFileReader(fileName string, offset, maxReadBytes int64) *ChunkedFileReader {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	if offset > 0 {
		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			panic(err)
		}
	}

	return &ChunkedFileReader{
		file:     file,
		reader:   bufio.NewReader(file),
		offset:   offset,
		maxBytes: maxReadBytes,
		text:     make([]byte, maxReadBytes-offset),
	}
}

func (o *ChunkedFileReader) Close() {
	o.file.Close()
	// runtime.GC()
	// debug.FreeOSMemory()
}

func (o *ChunkedFileReader) MoveReaderToStartOfNextLine() (bytesJumped int64, rErr error) {
	b, err := o.reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	readBytes := int64(len(b))
	o.offset += readBytes
	o.text = o.text[readBytes:]
	return readBytes, nil
}

func (o *ChunkedFileReader) MMap() error {
	n, err := io.ReadFull(o.reader, o.text)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}

	o.text = o.text[:n]
	return nil
}

func (o *ChunkedFileReader) GetLine() (name []byte, temp []byte, ok bool, rErr error) {
	if len(o.text) == 0 {
		return nil, nil, false, nil
	}

	idx := bytes.IndexByte(o.text, ';')
	name = o.text[:idx]

	o.text = o.text[idx+1:]
	idx = bytes.IndexByte(o.text, '\n')
	if idx < 0 {
		return name, o.text, true, nil
	}

	temp = o.text[:idx]
	o.text = o.text[idx+1:]
	return name, temp, true, nil
}
