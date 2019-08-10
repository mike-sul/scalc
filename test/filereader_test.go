package scalc

import (
	"github.com/mike-sul/scalc/pkg/scalc"
	"io"
	"testing"
)

func TestFileReader(t *testing.T) {
	// TODO: generate file dynamically and remove it on a test exit
	testMap := map[string][]int{
		"a.txt":  {1, 2, 3},
		"a0.txt": {},
		"aN.txt": {1, 1024, 2048},
	}

	for file, expSet := range testMap {
		fr, err := scalc.NewFileReader(file)
		if err != nil {
			t.Error(err.Error())
		}

		for _, expVal := range expSet {
			val, err := fr.Next()
			if err != nil {
				t.Fatalf("Got an error %s while expected %d", err.Error(), expVal)
			}
			if expVal != val {
				t.Fatalf("Got %d, expected %d", val, expVal)
			}

		}

		val, err := fr.Next()

		if err == nil {
			t.Fatalf("Expected an end of file, got value %d", val)
		}
		if err != io.EOF {
			t.Fatalf("Expected an end of file, got %s", err.Error())
		}
	}
}

func TestFileReaderNegative(t *testing.T) {
	fr, err := scalc.NewFileReader("non-existing-file")
	if err == nil {
		t.Fatal("Expected nil, instead got a file reader")
	}

	fr, err = scalc.NewFileReader("aInv.txt")
	if err != nil {
		t.Fatalf("Expected a new file reader for %s, got nill: %s",
			"aInv.txt", err.Error())
	}

	_, err = fr.Next()
	if err == nil {
		t.Fatal("Expected a type conversion error")
	}

}
