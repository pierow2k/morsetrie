// Package morsetrie implements a tried for morse code.
package morsetrie

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// TODO: ErrInvalidElement is returned when...
	ErrInvalidElement = errors.New("invalid morse element")
	// TODO: ErrDuplicate is returned when...
	ErrDuplicate = errors.New("duplicate morse code")
	// TODO: ErrUnexpectedChar is returned when...
	ErrUnexpectedChar = errors.New("unexpected character in morse input")
)

// MorsePair is a single Morse-code mapping entry.
type MorsePair struct {
	Code string
	R    rune
}

// MorseTable is the morse code mapping.
var MorseTable = []MorsePair{
	{".-", 'A'},
	{"-...", 'B'},
	{"-.-.", 'C'},
	{"-..", 'D'},
	{".", 'E'},
	{"..-.", 'F'},
	{"--.", 'G'},
	{"....", 'H'},
	{"..", 'I'},
	{".---", 'J'},
	{"-.-", 'K'},
	{".-..", 'L'},
	{"--", 'M'},
	{"-.", 'N'},
	{"---", 'O'},
	{".--.", 'P'},
	{"--.-", 'Q'},
	{".-.", 'R'},
	{"...", 'S'},
	{"-", 'T'},
	{"..-", 'U'},
	{"...-", 'V'},
	{".--", 'W'},
	{"-..-", 'X'},
	{"-.--", 'Y'},
	{"--..", 'Z'},

	{"-----", '0'},
	{".----", '1'},
	{"..---", '2'},
	{"...--", '3'},
	{"....-", '4'},
	{".....", '5'},
	{"-....", '6'},
	{"--...", '7'},
	{"---..", '8'},
	{"----.", '9'},
}

// Node is a binary trie Node.
// child[0] is the '.' edge; child[1] is the '-' edge.
// val==0 means "no symbol at this Node".
type Node struct {
	Child [2]int
	Val   rune
}

// Trie is a compact, array-backed Morse decode Trie.
type Trie struct {
	Nodes []Node
}

// TODO: NewTrie...
func NewTrie() *Trie {
	t := &Trie{Nodes: make([]Node, 1)} // node 0 is the root
	t.Nodes[0].Child[0] = -1
	t.Nodes[0].Child[1] = -1

	return t
}

// TODO: parameter name 'r' is too short for the scope of its usage.
func (t *Trie) add(code string, r rune) error {
	idx := 0 // start at root

	// TODO: variable name 'i' is too short for the scope of its usage
	for i := range len(code) {
		var bit int

		switch code[i] {
		case '.':
			bit = 0
		case '-':
			bit = 1
		default:
			return fmt.Errorf("%w: %q in %q", ErrInvalidElement, code[i], code)
		}

		next := t.Nodes[idx].Child[bit]
		if next == -1 {
			next = len(t.Nodes)
			// New nodes start with missing children.
			t.Nodes = append(t.Nodes, Node{Child: [2]int{-1, -1}})
			t.Nodes[idx].Child[bit] = next
		}

		idx = next
	}

	if t.Nodes[idx].Val != 0 {
		return fmt.Errorf("%w: %q (already maps to %q)", ErrDuplicate, code, t.Nodes[idx].Val)
	}

	t.Nodes[idx].Val = r

	return nil
}

// TODO: BuildTrie...
func BuildTrie(pairs []MorsePair) (*Trie, error) {
	// TODO: variable name 't' is too short for the scope of its usage
	t := NewTrie()
	for _, p := range pairs {
		if err := t.add(p.Code, p.R); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Decode implements the 3-state streaming state machine:
//   - AtRoot: not currently in a letter
//   - InLetter: traversing '.'/'-' edges for a letter
//   - InvalidLetter: current letter token became invalid
//
// Separators:
//   - ' ' (and common ASCII whitespace) ends a letter
//   - '/' ends a letter and emits a word break
func (t *Trie) Decode(input string) (string, error) {
	const (
		atRoot = iota
		inLetter
		invalidLetter
	)

	state := atRoot
	idx := 0

	// TODO: variable name 'b' is too short for the scope of its usage
	var b strings.Builder

	lastWasSpace := false

	commitLetter := func() {
		switch state {
		case inLetter:
			if t.Nodes[idx].Val == 0 {
				b.WriteRune('?')
			} else {
				b.WriteRune(t.Nodes[idx].Val)
			}

			lastWasSpace = false
		case invalidLetter:
			b.WriteRune('?')

			lastWasSpace = false
		}
		// Reset to root.
		state = atRoot
		idx = 0
	}

	commitWordBreak := func() {
		// Avoid leading or repeated spaces.
		if b.Len() > 0 && !lastWasSpace {
			b.WriteByte(' ')

			lastWasSpace = true
		}
	}

	for i := range len(input) {
		// TODO: variable name 'c' is too short for the scope of its usage
		c := input[i]
		switch c {
		case '.', '-':
			if state == invalidLetter {
				// Already invalid; keep consuming until a separator.
				continue
			}

			bit := 0
			if c == '-' {
				bit = 1
			}

			next := t.Nodes[idx].Child[bit]
			if next == -1 {
				state = invalidLetter

				continue
			}

			idx = next
			state = inLetter

		case ' ', '\t', '\n', '\r':
			// Letter separator.
			if state != atRoot {
				commitLetter()
			}
			// If at root, ignore extra whitespace.

		case '/':
			// Word break.
			if state != atRoot {
				commitLetter()
			}

			commitWordBreak()

		default:
			return "", fmt.Errorf("%w: "+string([]byte{c}), ErrUnexpectedChar) //nolint:err113
		}
	}

	// End-of-input flush.
	if state != atRoot {
		commitLetter()
	}

	return b.String(), nil
}
