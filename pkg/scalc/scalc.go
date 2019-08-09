package scalc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type DataStream interface {
	Peek() (int, bool)
	Next() (int, bool)
	Close()
}

type FileDataStream struct {
	file   *os.File
	reader *bufio.Reader
	cur    *int
	ok     bool
}

func NewFileDataStream(filename string) DataStream {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open file: " + filename)
		return nil
	}

	reader := bufio.NewReader(file)
	if reader == nil {
		log.Fatal("Failed to create a file reader: " + filename)
		return nil
	}

	return &FileDataStream{file: file, reader: reader, cur: nil, ok: false}
}

func (fds *FileDataStream) Close() {
	fds.cur = nil
	fds.reader = nil
	fds.file.Close()
}

func (fds *FileDataStream) Peek() (int, bool) {
	if fds.cur == nil {
		fds.cur = new(int)
		*fds.cur, fds.ok = fds.readValue()
	}
	return *fds.cur, fds.ok
}

func (fds *FileDataStream) Next() (int, bool) {
	if fds.cur == nil {
		fds.cur = new(int)
	}
	*fds.cur, fds.ok = fds.readValue()

	return *fds.cur, fds.ok
}

func (fds *FileDataStream) readValue() (int, bool) {
	line, err := fds.reader.ReadString('\n')
	if err != nil {
		return 0, false
	}
	if len(line) == 0 {
		return 0, false
	}

	line = strings.TrimSuffix(line, "\n")

	if len(line) == 0 {
		return 0, false
	}

	res, err := strconv.Atoi(line)
	if err != nil {
		return 0, false
	}
	return res, true
}

type SumDataStream struct {
	inputStreams []DataStream
	cur          *int
	ok           bool
}

func NewSumDataStream(streams []DataStream) DataStream {
	return &SumDataStream{inputStreams: streams, cur: nil, ok: false}
}

func (sds *SumDataStream) Peek() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
		*sds.cur, sds.ok = sds.nextValue()
	}
	return *sds.cur, sds.ok
}

func (sds *SumDataStream) Next() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
	} else {
		for ii := 0; ii < len(sds.inputStreams); ii++ {
			val, ok := sds.inputStreams[ii].Peek()
			if ok && val == *sds.cur {
				sds.inputStreams[ii].Next()
			}
		}
	}

	*sds.cur, sds.ok = sds.nextValue()

	return *sds.cur, sds.ok
}

func (sds *SumDataStream) nextValue() (int, bool) {
	min, ok := getInitialMinValue(sds.inputStreams)
	if !ok {
		return 0, false
	}

	for ii := 0; ii < len(sds.inputStreams); ii++ {
		val, ok := sds.inputStreams[ii].Peek()
		if !ok {
			continue
		}
		if val < min {
			min = val
		}
	}
	return min, ok
}

func (sds *SumDataStream) Close() {
	sds.cur = nil
}

func getInitialMinValue(streams []DataStream) (int, bool) {
	for ii := 0; ii < len(streams); ii++ {
		val, ok := streams[ii].Peek()
		if ok {
			return val, true
		}
	}
	return 0, false
}

type IntDataStream struct {
	inputStreams []DataStream
	cur          *int
	ok           bool
}

func NewIntDataStream(streams []DataStream) DataStream {
	return &IntDataStream{inputStreams: streams, cur: nil, ok: true}
}

func (sds *IntDataStream) Peek() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
		*sds.cur, sds.ok = sds.nextValue()
	}
	return *sds.cur, sds.ok
}

func (sds *IntDataStream) Next() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
	} else {
		if sds.ok {
			for ii := 0; ii < len(sds.inputStreams); ii++ {
				val, ok := sds.inputStreams[ii].Peek()
				if ok && val == *sds.cur {
					sds.inputStreams[ii].Next()
				}
			}
		}
	}

	if sds.ok {
		*sds.cur, sds.ok = sds.nextValue()
	}

	return *sds.cur, sds.ok
}

func (sds *IntDataStream) nextMinValue() (int, bool) {
	min, ok := getInitialMinValue(sds.inputStreams)
	if !ok {
		return 0, false
	}

	for ii := 0; ii < len(sds.inputStreams); ii++ {
		val, ok := sds.inputStreams[ii].Peek()
		if !ok {
			continue
		}
		if val < min {
			min = val
		}
	}
	return min, ok
}

func (sds *IntDataStream) nextValue() (int, bool) {

	res := 0
	foundIntersection := false

	for {
		min, ok := sds.nextMinValue()
		if !ok {
			return 0, false
		}

		foundIntersection = true

		for ii := 0; ii < len(sds.inputStreams); ii++ {
			val, ok := sds.inputStreams[ii].Peek()
			if !ok {
				return 0, false
			}
			if val != min {
				// move next for the streams with minimum value
				for ii := 0; ii < len(sds.inputStreams); ii++ {
					val, ok := sds.inputStreams[ii].Peek()
					if ok && val == min {
						sds.inputStreams[ii].Next()
					}
				}
				foundIntersection = false
				break
			}
		}

		if foundIntersection {
			res = min
			break
		}
	}
	return res, foundIntersection
}

func (sds *IntDataStream) Close() {
	sds.cur = nil
}

func DumpDataStrem(stream DataStream) {
	for val, ok := stream.Next(); ok; val, ok = stream.Next() {
		fmt.Println(val)
	}
}

type DiffDataStream struct {
	inputStreams []DataStream
	cur          *int
	ok           bool
}

func NewDiffDataStream(streams []DataStream) DataStream {
	return &DiffDataStream{inputStreams: streams, cur: nil, ok: true}
}

func (sds *DiffDataStream) Peek() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
		*sds.cur, sds.ok = sds.nextValue()
	}
	return *sds.cur, sds.ok
}

func (sds *DiffDataStream) nextValue() (int, bool) {

	diff_min, ok := sds.inputStreams[0].Peek()
	if !ok {
		return 0, false
	}

	cont := false

	for {
		for ii := 1; ii < len(sds.inputStreams); ii++ {
			val, ok := sds.inputStreams[ii].Peek()
			if !ok {
				continue
			}

			if val < diff_min {
				sds.inputStreams[ii].Next()
				cont = true
			}

			if val == diff_min {
				diff_min, ok = sds.inputStreams[0].Next()
				if !ok {
					return 0, false
				}
				cont = true
				break
			}

		}
		if !cont {
			break
		}
	}

	return diff_min, true
}

func (sds *DiffDataStream) Next() (int, bool) {
	if sds.cur == nil {
		sds.cur = new(int)
	} else {
		if sds.ok {
			sds.inputStreams[0].Next()
		}
	}

	if sds.ok {
		*sds.cur, sds.ok = sds.nextValue()
	}

	return *sds.cur, sds.ok
}

func (sds *DiffDataStream) Close() {
	sds.cur = nil
}

func ParseExpression(exp []string, indx int) (DataStream, int, error) {
	operators := map[string]func([]DataStream) DataStream{
		"SUM": NewSumDataStream,
		"INT": NewIntDataStream,
		"DIF": NewDiffDataStream,
	}

	if exp[indx] != "[" {
		return nil, indx, errors.New("invalid expression format")
	}

	operator := exp[indx+1]
	inputDataStreams := make([]DataStream, 0, 10)

	if exp[indx+2] == "]" {
		return nil, indx + 2, errors.New("invalid expression format")
	}

	ii := indx + 2
	for ; exp[ii] != "]"; ii++ {
		if exp[ii] != "[" {
			inputDataStreams = append(inputDataStreams, NewFileDataStream(exp[ii]))
		} else {
			ds, pos, err := ParseExpression(exp, ii)
			if err != nil {
				return nil, pos, errors.New("invalid expression format")
			}
			inputDataStreams = append(inputDataStreams, ds)
			ii = pos
		}
	}

	return operators[operator](inputDataStreams), ii, nil
}
