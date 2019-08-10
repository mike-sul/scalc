package scalc

import (
	"github.com/mike-sul/scalc/pkg/scalc"
	"io"
	"testing"
)

type InputSetsToExpected struct {
	operator   scalc.OperatorID
	inputfiles []string
	expected   []int
}

func testOperator(t *testing.T, testMap []InputSetsToExpected) {
	for indx, val := range testMap {

		readers := make([]scalc.SetReader, len(val.inputfiles))

		for ii, filename := range val.inputfiles {
			var err error
			readers[ii], err = scalc.NewFileReader(filename)
			if err != nil {
				t.Fatal(err.Error())
			}
		}

		uo, err := scalc.GetOperatorRegistry().Create(val.operator, readers)
		if err != nil {
			t.Fatalf("Failed to create an operator: %s", err.Error())
		}

		for _, expVal := range val.expected {
			val, err := uo.Next()
			if err != nil {
				t.Fatalf("Got an error %s while expected %d", err.Error(), expVal)
			}
			if expVal != val {
				t.Fatalf("Got %d, expected %d", val, expVal)
			}

		}
		_, err = uo.Next()
		if err == nil {
			t.Fatalf("Expected EOF, got value for test %d", indx)
		}
		if err != io.EOF {
			t.Fatalf("Expected EOF, got %s", err.Error())
		}
	}
}

// TODO: Consider dynamic generation of an input sets/files instead of usage of predefined set files
func TestUnionOperator(t *testing.T) {

	testMap := []InputSetsToExpected{
		{scalc.UnionOperatorId, []string{"a.txt"}, []int{1, 2, 3}},
		{scalc.UnionOperatorId, []string{"a.txt", "b.txt"}, []int{1, 2, 3, 4}},
		{scalc.UnionOperatorId, []string{"a.txt", "c.txt"}, []int{1, 2, 3, 4, 5}},
	}

	testOperator(t, testMap)
}

func TestInterOperator(t *testing.T) {

	testMap := []InputSetsToExpected{
		{scalc.InterOperatorId, []string{"a.txt"}, []int{1, 2, 3}},
		{scalc.InterOperatorId, []string{"a.txt", "b.txt"}, []int{2, 3}},
		{scalc.InterOperatorId, []string{"a.txt", "c.txt"}, []int{3}},
	}

	testOperator(t, testMap)
}

func TestDiffOperator(t *testing.T) {
	testMap := []InputSetsToExpected{
		{scalc.DifOperatorId, []string{"a.txt", "b.txt"}, []int{1}},
		{scalc.DifOperatorId, []string{"a.txt", "c.txt"}, []int{1, 2}},
		{scalc.DifOperatorId, []string{"a.txt", "b.txt", "c.txt"}, []int{1}},
		{scalc.DifOperatorId, []string{"a.txt"}, []int{1,2,3}},
		{scalc.DifOperatorId, []string{"a1.txt"}, []int{1}},
		{scalc.DifOperatorId, []string{"aN.txt", "b.txt", "c.txt"}, []int{1, 1024, 2048}},
		{scalc.DifOperatorId, []string{"c.txt", "a.txt"}, []int{4, 5}},
	}

	testOperator(t, testMap)
}
