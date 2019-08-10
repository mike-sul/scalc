package scalc

import (
	"github.com/mike-sul/scalc/pkg/scalc"
	"os/exec"
	"strings"
	"testing"
)

func TestScalcExeSmoke(t *testing.T) {
	exp := "[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]"
	expectedOut := "1\n3\n4\n"

	out, err := exec.Command("scalc", strings.Fields(exp)...).Output()
	if err != nil {
		t.Error(err)
	}

	if string(out) != expectedOut {
		t.Errorf("Got %s expected %s", string(out), expectedOut)
	}
}

func TestScalcSmoke(t *testing.T) {
	exp := "[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]"

	ds, _, err := scalc.ParseExpression(strings.Fields(exp), 0)
	if err != nil {
		t.Error(err)
	}

	expected := []int{1, 3, 4}
	ii := 0
	for val, err := ds.Next(); err != nil; val, err = ds.Next() {
		if expected[ii] != val {
			t.Errorf("Got %d expected %d at position %d", val, expected[ii], ii)
		}
		ii++
	}
}

func TestNegativeEdgeCases(t *testing.T) {
	invExp := []string{
		"",
		"[",
		"[ ]",
		"B B",
		"[ SUM ]",
		"[ SUM [",
		"[ SUM [ ]",
		"[ SUM [ ] ]",
		"[ NOTSUPPORTED a.txt ]",
		"[ NOTSUPPORTED [ ] ]",
		"[ DIF [ SUM [ INT [ DIF [ ] ] ] ] ]",
	}

	for _, exp := range invExp {
		_, _, err := scalc.ParseExpression(strings.Fields(exp), 0)
		if err == nil {
			t.Errorf("Failed for %s", exp)
		}
	}
}

func TestPositiveEdgeCases(t *testing.T) {
	exps := []string{
		"[ SUM a0.txt ]",
		"[ INT a0.txt ]",
		"[ DIF a0.txt ]",
		"[ INT [ SUM a0.txt ] ]",
		"[ INT [ DIF a0.txt ] ]",
		"[ INT [ INT a0.txt ] ]",
		"[ DIF [ INT a0.txt ] ]",
		"[ DIF [ SUM a0.txt ] ]",
		"[ DIF [ DIF a0.txt ] ]",
		"[ SUM [ INT a0.txt ] ]",
		"[ SUM [ DIF a0.txt ] ]",
		"[ SUM [ DIF a0.txt ] ]",
		"[ SUM [ DIF a0.txt ] [ INT a0.txt ] ]",
		"[ DIF a0.txt a.txt b.txt c.txt d.txt ]",
		"[ INT a0.txt a.txt b.txt c.txt d.txt ]",
		"[ INT [ DIF a1.txt c.txt ] [ SUM b.txt c.txt ] ]",
	}

	for _, exp := range exps {
		ds, _, err := scalc.ParseExpression(strings.Fields(exp), 0)
		if err != nil {
			t.Errorf("Failed for %s", exp)
		}
		_, err1 := ds.Peek()

		if err1 == nil {
			t.Errorf("Failed for %s", exp)
		}
	}
}

func TestScalcSanity(t *testing.T) {
	testMap := map[string][]int{
		"[ SUM a.txt ]":                                    {1, 2, 3},
		"[ INT a.txt ]":                                    {1, 2, 3},
		"[ DIF a.txt ]":                                    {1, 2, 3},
		"[ SUM [ SUM [ SUM a1.txt ] ] ]":                   {1},
		"[ INT [ INT [ INT a1.txt ] ] ]":                   {1},
		"[ DIF [ DIF [ DIF a1.txt ] ] ]":                   {1},
		"[ DIF aN.txt c.txt ]":                             {1, 1024, 2048},
		"[ SUM [ INT a.txt c.txt ] [ DIF aN.txt c.txt ] ]": {1, 3, 1024, 2048},
		"[ SUM [ INT a1.txt c.txt ] [ DIF a.txt b.txt ] ]": {1},
		"[ SUM [ DIF a.txt b.txt ] [ INT a1.txt c.txt ] ]": {1},
	}

	for exp, expDataStream := range testMap {
		ds, _, err := scalc.ParseExpression(strings.Fields(exp), 0)
		if err != nil {
			t.Errorf("Failed for %s", exp)
		}

		for _, expVal := range expDataStream {
			actVal, err := ds.Next()
			if err != nil {
				t.Errorf("Failed for %s", exp)
			}
			if expVal != actVal {
				t.Errorf("Failed for %s. Expected: %d, got: %d", exp, expVal, actVal)
			}
		}
	}
}
