package scalc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type FileReader struct {
	file   *os.File
	reader *bufio.Reader

	buffer  *int
	lastErr error
	pipe chan SetVal
}

func NewFileReader(filename string) (SetReader, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("failed to open file: " + filename)
	}

	reader := bufio.NewReader(file)
	if reader == nil {
		file.Close()
		return nil, errors.New("failed to create a file reader: " + filename)
	}

	fr := FileReader{file: file, reader: reader, buffer: nil, lastErr: nil, pipe:make(chan SetVal)}
	go fr.pushValueToChannel()
	return &fr, nil
}

func (fr *FileReader) Peek() (int, error) {
	if fr.buffer == nil {
		fr.buffer = new(int)
		fr.readValueFromChannel()
		//fr.readValueFromFile()
	}
	return *fr.buffer, fr.lastErr
}

func (fr *FileReader) Next() (int, error) {
	if fr.buffer != nil {
		fr.readValueFromChannel()
		//fr.readValueFromFile()
	}
	return fr.Peek()
}

func (fr *FileReader) Close() {
	fr.buffer = nil
	fr.reader = nil
	fr.file.Close()
}

func (fr *FileReader) readValue() (int, error) {
	line, err := fr.reader.ReadString('\n')
	if err != nil {
		return math.MaxInt64, err
	}

	line = strings.TrimSuffix(line, "\n")
	if len(line) == 0 {
		return math.MaxInt64, io.EOF
	}

	val, err := strconv.Atoi(line)
	if err != nil {
		return math.MaxInt64, fmt.Errorf("invalid value `%s` in %s", line, fr.file.Name())
	}

	return val, nil
}

func (fr *FileReader) pushValueToChannel() {
	var val int
	var err error
	for val, err = fr.readValue(); err == nil; val, err = fr.readValue() {
		fr.pipe <- SetVal{val: val, err: err}
	}
	fr.pipe <- SetVal{val: val, err: err}
	close(fr.pipe)
	fr.file.Close()
}

func (fr *FileReader) readValueFromChannel() {
	val, ok := <-fr.pipe
	if ok {
		fr.lastErr = val.err
		*fr.buffer = val.val
	} else {
		fr.lastErr = io.EOF
		*fr.buffer = math.MaxInt64
	}
}

func (fr *FileReader) readValueFromFile() {
	*fr.buffer, fr.lastErr = fr.readValue()
}