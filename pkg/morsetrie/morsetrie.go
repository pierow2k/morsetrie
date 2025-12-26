// Package morsetrie implements trie based decoding for morse code.
package morsetrie

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidElement is returned when a code string contains characters other than '.' or '-'.
	ErrInvalidElement = errors.New("invalid morse element")
	// ErrDuplicate is returned when a morse code sequence is registered twice.
	ErrDuplicate = errors.New("duplicate morse code")
	// ErrUnexpectedChar is returned when the input string contains unsupported characters.
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

// NewTrie creates a new, empty Morse decode Trie.
func NewTrie() *Trie {
	t := &Trie{Nodes: make([]Node, 1)} // node 0 is the root
	t.Nodes[0].Child[0] = -1
	t.Nodes[0].Child[1] = -1

	return t
}

func (t *Trie) add(code string, symbol rune) error {
	idx := 0 // start at root

	for charIndex := range len(code) {
		var bit int

		switch code[charIndex] {
		case '.':
			bit = 0
		case '-':
			bit = 1
		default:
			return fmt.Errorf("%w: %q in %q", ErrInvalidElement, code[charIndex], code)
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

	t.Nodes[idx].Val = symbol

	return nil
}

// BuildTrie constructs a Trie from the provided list of Morse code pairs.
func BuildTrie(pairs []MorsePair) (*Trie, error) {
	trie := NewTrie()
	for _, p := range pairs {
		if err := trie.add(p.Code, p.R); err != nil {
			return nil, err
		}
	}

	return trie, nil
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

	var builder strings.Builder

	lastWasSpace := false

	commitLetter := func() {
		switch state {
		case inLetter:
			if t.Nodes[idx].Val == 0 {
				builder.WriteRune('?')
			} else {
				builder.WriteRune(t.Nodes[idx].Val)
			}

			lastWasSpace = false
		case invalidLetter:
			builder.WriteRune('?')

			lastWasSpace = false
		}
		// Reset to root.
		state = atRoot
		idx = 0
	}

	commitWordBreak := func() {
		// Avoid leading or repeated spaces.
		if builder.Len() > 0 && !lastWasSpace {
			builder.WriteByte(' ')

			lastWasSpace = true
		}
	}

	for inputIndex := range len(input) {
		char := input[inputIndex]
		switch char {
		case '.', '-':
			if state == invalidLetter {
				// Already invalid; keep consuming until a separator.
				continue
			}

			bit := 0
			if char == '-' {
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
			return "", fmt.Errorf("%w: "+string([]byte{char}), ErrUnexpectedChar) //nolint:err113
		}
	}

	// End-of-input flush.
	if state != atRoot {
		commitLetter()
	}

	return builder.String(), nil
}
