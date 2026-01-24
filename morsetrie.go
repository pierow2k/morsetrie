// Package morsetrie implements trie-based decoding for morse code.
package morsetrie

import (
	"bytes"
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

// ErrUnexpectedChar is returned when the morse code string contains
// unsupported characters.
var ErrUnexpectedChar = errors.New("unexpected character in morse input")

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
		builder.WriteByte(byte(val)) //nolint:gosec
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

// traverse recursively explores the trie to build all valid letter combinations
// from a Morse code sequence without separators.
//
// It uses backtracking with a shared buffer to efficiently explore all possible
// segmentations. At each step, it pursues two simultaneous strategies:
//
//  1. Continue: Treat current symbols as a prefix for a longer letter.
//  2. Branch: If current path forms a valid letter, commit it and restart at root.
//
// The mark/restore pattern avoids allocations at branch points: when a valid
// letter is found, the buffer position is saved, the letter is appended,
// recursion proceeds, and the buffer is truncated back to the saved position.
//
// Parameters:
//   - sequence: The full Morse code string being decoded.
//   - idx: Current position in the sequence.
//   - buf: Shared buffer for building decoded strings (reused across branches).
//   - nodeIdx: Current position in the trie (index into t.Nodes).
//   - results: Pointer to slice accumulating valid decodings.
func (t *Trie) traverse(sequence string, idx int, buf *bytes.Buffer, nodeIdx int16, results *[]string) {
	// Base case: All symbols have been consumed.
	if idx == len(sequence) {
		// Valid decodings must end on a letter boundary (back at root).
		// Incomplete sequences end mid-trie and are discarded.
		if nodeIdx == rootIdx {
			*results = append(*results, buf.String())
		}

		return
	}

	symbol := sequence[idx]

	var bit int

	switch symbol {
	case '.':
		bit = 0
	case '-':
		bit = 1
	default:
		// Invalid characters terminate this path.
		return
	}

	childIdx := t.Nodes[nodeIdx].Child[bit]
	if childIdx == missingNode {
		// Dead end: sequence follows a path that doesn't exist in the trie.
		return
	}

	nextIdx := idx + 1

	// STRATEGY 1: CONTINUE
	// Extend the current prefix without committing to a letter boundary.
	t.traverse(sequence, nextIdx, buf, childIdx, results)

	// STRATEGY 2: BRANCH
	// If the current path forms a valid letter, commit to it and restart.
	if t.Nodes[childIdx].Val != 0 {
		mark := buf.Len()
		buf.WriteRune(t.Nodes[childIdx].Val)
		t.traverse(sequence, nextIdx, buf, rootIdx, results)
		buf.Truncate(mark) // Restore buffer for sibling branches.
	}
}

// FindCandidates returns all valid decodings for a Morse code sequence
// without separators. It attempts to segment the raw symbol stream into
// every possible letter combination.
//
// This is useful for decoding ambiguous inputs, such as Morse code bracelets,
// where letter boundaries are not marked.
//
// Parameters:
//   - sequence: A string containing only Morse symbols ('.' and '-').
//
// Returns:
//   - A slice of all valid decodings. Returns an empty slice if the sequence
//     cannot be fully segmented into valid letters.
func (t *Trie) FindCandidates(sequence string) []string {
	var (
		candidates []string
		buf        bytes.Buffer
	)

	buf.Grow(len(sequence) / decodeAllocDivisor)
	t.traverse(sequence, 0, &buf, rootIdx, &candidates)

	return candidates
}

// FindCandidates provides a package-level function to return all valid
// decodings for a Morse code sequence using the default trie.
func FindCandidates(morseCode string) []string {
	return StaticTrie.FindCandidates(morseCode)
}
