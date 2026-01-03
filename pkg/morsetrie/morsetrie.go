// Package morsetrie implements trie based decoding for morse code.
package morsetrie

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidElement is returned when a code string contains characters
	// other than '.' or '-'.
	ErrInvalidElement = errors.New("invalid morse element")
	// ErrDuplicate is returned when a morse code sequence is registered twice.
	ErrDuplicate = errors.New("duplicate morse code")
	// ErrUnexpectedChar is returned when the input string contains
	// unsupported characters.
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
// Optimized Node: 8 bytes total.
// Child indices are int16 (max 32767 nodes, plenty for Morse).
type Node struct {
	Val   rune
	Child [2]int16
}

// Trie is a compact, array-backed Morse decode Trie.
type Trie struct {
	Nodes []Node
}

// NewTrie creates a new, empty Morse decode Trie.
func NewTrie() *Trie {
	t := &Trie{Nodes: make([]Node, 1, 64)} // Pre-alloc cap to avoid early appends
	// Use -1 to indicate no child.
	t.Nodes[0].Child[0] = -1
	t.Nodes[0].Child[1] = -1
	return t
}

func (t *Trie) add(code string, symbol rune) error {
	idx := int16(0) // start at root

	for i := 0; i < len(code); i++ {
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
			next = int16(len(t.Nodes))
			// Create new node with empty children
			t.Nodes = append(t.Nodes, Node{Child: [2]int16{-1, -1}})
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

// Decode decodes the input string.
// Optimization: flattened logic, stack variables, buffer pre-allocation.
func (t *Trie) Decode(input string) (string, error) {
	var b strings.Builder
	// Heuristic: Decoded text is roughly 1/3 the size of Morse input.
	b.Grow(len(input) / 3)

	const (
		rootIdx     = int16(0)
		invalidIdx  = int16(-2) // Represents we walked off the trie
		missingNode = int16(-1) // Represents a child definition that doesn't exist
	)

	curr := rootIdx
	lastWasSpace := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch char {
		case '.':
			if curr != invalidIdx {
				curr = t.Nodes[curr].Child[0]
				// If child is -1, we went off the path
				if curr == missingNode {
					curr = invalidIdx
				}
			}
		case '-':
			if curr != invalidIdx {
				curr = t.Nodes[curr].Child[1]
				if curr == missingNode {
					curr = invalidIdx
				}
			}
		case ' ', '\t', '\n', '\r':
			// Commit letter
			if curr != rootIdx {
				if curr == invalidIdx {
					b.WriteByte('?')
				} else {
					val := t.Nodes[curr].Val
					if val == 0 {
						b.WriteByte('?')
					} else {
						b.WriteRune(val)
					}
				}
				curr = rootIdx
				lastWasSpace = false
			}
		case '/':
			// Commit letter (if any)
			if curr != rootIdx {
				if curr == invalidIdx {
					b.WriteByte('?')
				} else {
					val := t.Nodes[curr].Val
					if val == 0 {
						b.WriteByte('?')
					} else {
						b.WriteRune(val)
					}
				}
				curr = rootIdx
			}
			// Commit word break
			if b.Len() > 0 && !lastWasSpace {
				b.WriteByte(' ')
				lastWasSpace = true
			}
		default:
			return "", fmt.Errorf("%w: %c", ErrUnexpectedChar, char)
		}
	}

	// Final commit (end of string)
	if curr != rootIdx {
		if curr == invalidIdx {
			b.WriteByte('?')
		} else {
			val := t.Nodes[curr].Val
			if val == 0 {
				b.WriteByte('?')
			} else {
				b.WriteRune(val)
			}
		}
	}

	return b.String(), nil
}
