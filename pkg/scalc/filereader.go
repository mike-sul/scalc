package scalc

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

type FileReader struct {
	file   *os.File
	reader *bufio.Reader

	buffer  *int
	lastErr error
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

	return &FileReader{file: file, reader: reader, buffer: nil, lastErr: nil}, nil
}

func (fr *FileReader) Peek() (int, error) {
	if fr.buffer == nil {
		fr.buffer = new(int)
		fr.readValue()
	}
	return *fr.buffer, fr.lastErr
}

func (fr *FileReader) Next() (int, error) {
	if fr.buffer != nil {
		fr.readValue()
	}
	return fr.Peek()
}

func (fr *FileReader) Close() {
	fr.buffer = nil
	fr.reader = nil
	fr.file.Close()
}

func (fr *FileReader) readValue() {
	line, err := fr.reader.ReadString('\n')
	if err != nil {
		*fr.buffer = 0
		fr.lastErr = err
		return
	}

	line = strings.TrimSuffix(line, "\n")
	if len(line) == 0 {
		*fr.buffer = 0
		fr.lastErr = io.EOF
		return
	}

	val, err := strconv.Atoi(line)
	if err != nil {
		*fr.buffer = 0
		fr.lastErr = strconv.ErrSyntax
		return
	}

	*fr.buffer = val
	fr.lastErr = nil
}
