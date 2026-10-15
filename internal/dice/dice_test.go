package dice

import (
	"reflect"
	"strings"
	"testing"
)

// seq returns an intn that yields the given die faces (1-based) in order.
func seq(faces ...int) func(int) int {
	i := 0
	return func(int) int {
		v := faces[i] - 1
		i++
		return v
	}
}

func TestRollFixed(t *testing.T) {
	tests := []struct {
		notation string
		faces    []int
		rolls    []int
		total    int
	}{
		{"1d20+3", []int{14}, []int{14}, 17},
		{"2d20kh1+5", []int{3, 18}, []int{3, 18}, 23},
		{"2d20kl1", []int{3, 18}, []int{3, 18}, 3},
		{"4d6kl3", []int{6, 1, 4, 2}, []int{6, 1, 4, 2}, 7},
		{"4d6kh3", []int{6, 1, 4, 2}, []int{6, 1, 4, 2}, 12},
		{"d100", []int{42}, []int{42}, 42},
		{"3d8-1", []int{1, 2, 3}, []int{1, 2, 3}, 5},
		{" 1D6 + 2d4 - 1d4 ", []int{5, 1, 2, 3}, []int{5, 1, 2, 3}, 5},
		{"7", nil, []int{}, 7},
		{"-1d4+10", []int{3}, []int{3}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			e, err := Parse(tt.notation)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			got := e.roll(seq(tt.faces...))
			if !reflect.DeepEqual(got.Rolls, tt.rolls) || got.Total != tt.total {
				t.Fatalf("got rolls %v total %d, want %v %d", got.Rolls, got.Total, tt.rolls, tt.total)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	bad := []string{
		"", "   ", "abc", "d", "1d", "0d6", "101d6", "1d1", "1d1001", "1d20+", "+", "1d20++3",
		"4d6kh", "4d6kh0", "4d6kh5", "4d6kx1", "1d20*2", "1d20+1234567", "1d20+100001",
		"60d6+60d6", strings.Repeat("1+", 60) + "1", "1d6;", "１d6",
		"1d20" + strings.Repeat(" ", MaxRaw) + "+3",
	}
	for _, s := range bad {
		if _, err := Parse(s); err == nil {
			t.Errorf("Parse(%q) = nil error, want error", s)
		}
	}
}

func TestParseLimitsOK(t *testing.T) {
	for _, s := range []string{"100d1000", "50d6+50d6", "1d2", "100d6kl100", "100000"} {
		if _, err := Parse(s); err != nil {
			t.Errorf("Parse(%q): %v", s, err)
		}
	}
}

func TestRollCryptoRange(t *testing.T) {
	for i := 0; i < 200; i++ {
		r, err := Roll("3d6+1")
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Rolls) != 3 || r.Total < 4 || r.Total > 19 {
			t.Fatalf("out of range: %+v", r)
		}
		for _, v := range r.Rolls {
			if v < 1 || v > 6 {
				t.Fatalf("die out of range: %d", v)
			}
		}
	}
}
