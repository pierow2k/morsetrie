// Package morsetrie implements trie-based decoding for morse code.
package morsetrie

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	// decodeAllocDivisor estimates decoded text size (approx 1/3 of morse input).
	decodeAllocDivisor = 3

	rootIdx     = int16(0)
	invalidIdx  = int16(-2) // Represents traversing off the trie.
	missingNode = int16(-1) // Represents a child definition that doesn't exist.
)

var (
	// ErrUnexpectedChar is returned when the morse code string contains
	// unsupported characters.
	ErrUnexpectedChar = errors.New("unexpected character in morse input")
)

// Node is a node in the decoding Trie.
type Node struct {
	// Val is the decoded rune. If 0, this node is not a valid symbol end.
	Val rune
	// Child stores indices for the next node.
	// Child[0] is the '.' edge; Child[1] is the '-' edge.
	// We use int16 to reduce memory footprint.
	Child [2]int16
}

// Trie is a compact, array-backed Morse decode Trie.
type Trie struct {
	Nodes []Node
}

// advance returns the index of the child node for the given bit (0 for
// dot, 1 for dash).
// If the current path is already invalid, it stays invalid.
func (t *Trie) advance(curr int16, bit int) int16 {
	if curr == invalidIdx {
		return invalidIdx
	}

	next := t.Nodes[curr].Child[bit]
	if next == missingNode {
		return invalidIdx
	}

	return next
}

// commit writes the character corresponding to the current node index to
// the builder.
// If the path was invalid or the node has no value, it writes '?'.
// If the current node is the root (meaning consecutive separators),
// nothing is written.
func (t *Trie) commit(builder *strings.Builder, curr int16) {
	if curr == rootIdx {
		return
	}

	if curr == invalidIdx {
		builder.WriteByte('?')

		return
	}

	val := t.Nodes[curr].Val
	if val == 0 {
		builder.WriteByte('?')
		return
	}

	if val < utf8.RuneSelf {
		// Avoid WriteRune for ASCII values
		builder.WriteByte(byte(val))
	} else {
		builder.WriteRune(val)
	}
}

// Decode converts a string of Morse code into its corresponding text
// representation.
// It interprets '.' and '-' as Morse signals and uses whitespace (space,
// tab, newline, carriage return) to delimit encoded characters. The
// forward slash ('/') is treated as a word separator and is converted to a
// space in the output.
//
// If the input contains characters other than '.', '-', '/', or
// whitespace, Decode returns an empty string and ErrUnexpectedChar.
// Unknown Morse sequences are represented by '?' in the output.
func (t *Trie) Decode(morseCode string) (string, error) {
	var builder strings.Builder
	builder.Grow(len(morseCode) / decodeAllocDivisor)

	curr := rootIdx
	lastWasSpace := false

	for i := range len(morseCode) {
		char := morseCode[i]

		switch char {
		case '.':
			lastWasSpace = false
			curr = t.advance(curr, 0)
		case '-':
			lastWasSpace = false
			curr = t.advance(curr, 1)
		case ' ', '\t', '\n', '\r':
			t.commit(&builder, curr)
			curr = rootIdx
		case '/':
			t.commit(&builder, curr)
			curr = rootIdx
			// Commit word break if we aren't already spacing.
			if builder.Len() > 0 && !lastWasSpace {
				builder.WriteByte(' ')

				lastWasSpace = true
			}
		default:
			return "", fmt.Errorf("%w: %c", ErrUnexpectedChar, char)
		}
	}

	t.commit(&builder, curr)

	return builder.String(), nil
}

// Decode provides a package-level decode function that uses the
// static trie to decode a string.
func Decode(morseCode string) (string, error) {

	return StaticTrie.Decode(morseCode)
}
