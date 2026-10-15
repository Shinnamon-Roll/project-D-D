// Package dice parses and rolls dice notation: terms like NdM, NdMkhK, NdMklK and integer
// constants joined by + or -. Grammar: docs/design/protocol.md "Dice notation".
package dice

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	MaxDice     = 100 // per term, and per whole expression (bounds the shared log size)
	MaxSides    = 1000
	MaxConst    = 100000 // keeps totals far from int overflow
	MaxNotation = 100    // characters, after whitespace removal
	MaxRaw      = 200    // bytes before whitespace removal; callers store the raw string
	keepNone    = 0
	keepHighest = 1
	keepLowest  = 2
)

type term struct {
	sign     int // +1 or -1
	n, sides int // n == 0 means a constant term
	keep     int // keepNone, keepHighest, keepLowest
	k        int
	constant int
}

// Expr is a parsed dice expression.
type Expr struct{ terms []term }

// Result of rolling an Expr: every die in roll order, and the total of kept dice plus constants.
type Result struct {
	Rolls []int
	Total int
}

// Parse checks notation against the grammar and limits. Whitespace is ignored, case-insensitive.
func Parse(notation string) (Expr, error) {
	if len(notation) > MaxRaw {
		return Expr{}, fmt.Errorf("dice notation longer than %d bytes", MaxRaw)
	}
	s := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, notation)
	if s == "" {
		return Expr{}, errors.New("empty dice notation")
	}
	if len(s) > MaxNotation {
		return Expr{}, fmt.Errorf("dice notation longer than %d characters", MaxNotation)
	}
	p := parser{s: s}
	var e Expr
	totalDice := 0
	for i := 0; p.i < len(s); i++ {
		sign := 1
		if c := p.peek(); c == '+' || c == '-' {
			if c == '-' {
				sign = -1
			}
			p.i++
		} else if i > 0 {
			return Expr{}, fmt.Errorf("expected + or - at position %d", p.i+1)
		}
		t, err := p.term()
		if err != nil {
			return Expr{}, err
		}
		t.sign = sign
		totalDice += t.n
		if totalDice > MaxDice {
			return Expr{}, fmt.Errorf("at most %d dice per roll", MaxDice)
		}
		e.terms = append(e.terms, t)
	}
	return e, nil
}

type parser struct {
	s string
	i int
}

func (p *parser) peek() byte {
	if p.i < len(p.s) {
		return p.s[p.i]
	}
	return 0
}

// number reads digits; ok is false when there are none.
func (p *parser) number() (int, bool, error) {
	start := p.i
	for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
	}
	if p.i == start {
		return 0, false, nil
	}
	if p.i-start > 6 {
		return 0, true, fmt.Errorf("number %q too large", p.s[start:p.i])
	}
	n, _ := strconv.Atoi(p.s[start:p.i])
	return n, true, nil
}

func (p *parser) term() (term, error) {
	n, hasN, err := p.number()
	if err != nil {
		return term{}, err
	}
	if p.peek() != 'd' {
		if !hasN {
			return term{}, fmt.Errorf("expected a number or dice at position %d", p.i+1)
		}
		if n > MaxConst {
			return term{}, fmt.Errorf("constant %d larger than %d", n, MaxConst)
		}
		return term{constant: n}, nil
	}
	p.i++ // 'd'
	if !hasN {
		n = 1
	}
	if n < 1 || n > MaxDice {
		return term{}, fmt.Errorf("dice count must be 1-%d", MaxDice)
	}
	sides, ok, err := p.number()
	if err != nil {
		return term{}, err
	}
	if !ok {
		return term{}, fmt.Errorf("missing die size at position %d", p.i+1)
	}
	if sides < 2 || sides > MaxSides {
		return term{}, fmt.Errorf("die size must be 2-%d", MaxSides)
	}
	t := term{n: n, sides: sides}
	if strings.HasPrefix(p.s[p.i:], "kh") || strings.HasPrefix(p.s[p.i:], "kl") {
		t.keep = keepHighest
		if p.s[p.i+1] == 'l' {
			t.keep = keepLowest
		}
		p.i += 2
		k, ok, err := p.number()
		if err != nil {
			return term{}, err
		}
		if !ok || k < 1 || k > n {
			return term{}, fmt.Errorf("keep count must be 1-%d", n)
		}
		t.k = k
	}
	return t, nil
}

// Roll parses and rolls notation with crypto/rand.
func Roll(notation string) (Result, error) {
	e, err := Parse(notation)
	if err != nil {
		return Result{}, err
	}
	return e.roll(cryptoIntn), nil
}

func cryptoIntn(n int) int {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err) // crypto/rand never fails on supported platforms
	}
	return int(v.Int64())
}

// roll uses intn(n) in [0, n) so tests can inject a fixed sequence.
func (e Expr) roll(intn func(int) int) Result {
	res := Result{Rolls: []int{}}
	for _, t := range e.terms {
		if t.n == 0 {
			res.Total += t.sign * t.constant
			continue
		}
		rolls := make([]int, t.n)
		for i := range rolls {
			rolls[i] = intn(t.sides) + 1
		}
		res.Rolls = append(res.Rolls, rolls...)
		kept := rolls
		if t.keep != keepNone {
			kept = append([]int(nil), rolls...)
			sort.Ints(kept)
			if t.keep == keepHighest {
				kept = kept[len(kept)-t.k:]
			} else {
				kept = kept[:t.k]
			}
		}
		for _, v := range kept {
			res.Total += t.sign * v
		}
	}
	return res
}
