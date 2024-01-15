package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"unsafe"
)

type ChunkedFileReader struct {
	mutex  sync.Mutex
	file   *os.File
	reader *bufio.Reader

	offset   uint64
	maxBytes uint64
	text     []byte
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
		file:     file,
		reader:   bufio.NewReader(file),
		offset:   offset,
		maxBytes: maxReadBytes,
		text:     make([]byte, maxReadBytes-offset),
	}
}

func (o *ChunkedFileReader) Close() {
	o.file.Close()
}

func (o *ChunkedFileReader) MMap() error {
	n, err := io.ReadFull(o.reader, o.text)
	// n, err := o.reader.Read(o.text)
	// log.Printf("Offset: %v | MaxBytes: %v, BytesRead: %v, len(o.text): %v, | Text: \n%v\n",
	// 	o.offset, o.maxBytes, n, len(o.text),
	// 	string(o.text),
	// )

	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return err
	}

	o.text = o.text[:n]
	return nil
}

func (o *ChunkedFileReader) GetLine() (string, bool, error) {
	if len(o.text) == 0 {
		return "", false, nil
	}

	idx := bytes.IndexByte(o.text, byte('\n'))
	var lineBytes []byte
	if idx > 0 {
		lineBytes = o.text[:idx]
	} else {
		log.Println(string(o.text))
		lineBytes = o.text
	}

	if len(lineBytes) == 0 {
		return "", false, nil
	}

	o.text = o.text[idx+1:]
	s := unsafe.String(unsafe.SliceData(lineBytes), len(lineBytes))
	return s, true, nil
}
