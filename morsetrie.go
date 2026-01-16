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

// GenerateCyclicRotations returns all cyclic permutations of a sequence.
//
// Because Morse code bracelets are circular, there is no defined starting
// point for the message. This function addresses that ambiguity by generating
// every possible rotation (shifting the start index from 0 to N-1).
//
// If the reverse flag is true, the function also generates rotations for the
// reversed sequence. This accounts for the additional ambiguity of reading
// direction (clockwise vs. counter-clockwise).
//
// The resulting slice contains N forward rotations followed by N reverse
// rotations (if applicable), totaling 2N elements.
func GenerateCyclicRotations(sequence string, reverse bool) []string {
	const doubleSize = 2

	sequenceLength := len(sequence)
	if sequenceLength == 0 {
		return nil
	}

	// Pre-allocate capacity to avoid repeated slice growth.
	// We need space for N forward rotations, plus N reverse if requested.
	capacity := sequenceLength
	if reverse {
		capacity = doubleSize * sequenceLength
	}

	rotations := make([]string, 0, capacity)

	// Generate forward rotations:
	// For each index i, the sequence is split and swapped:
	//   sequence[i:] (suffix) + sequence[:i] (prefix)
	for i := range sequenceLength {
		rotations = append(rotations, sequence[i:]+sequence[:i])
	}

	if reverse {
		// Generate reverse rotations to handle reading-direction ambiguity.
		// A bracelet read backwards produces a different Morse string entirely.
		reversed := reverseString(sequence)
		for i := range sequenceLength {
			rotations = append(rotations, reversed[i:]+reversed[:i])
		}
	}

	return rotations
}

// reverseString returns a new string with the characters in reverse order.
//
// It converts the input to a rune slice to correctly handle multi-byte
// Unicode characters. While Morse code uses only single-byte ASCII,
// this approach ensures the function is robust for any string input.
func reverseString(s string) string {
	// Convert to runes to operate on Unicode code points rather than bytes.
	// This prevents corrupting multi-byte characters during the swap.
	runes := []rune(s)

	// Swap characters from the outside in using two pointers.
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// traverse recursively explores the trie to build all valid letter combinations.
//
// It solves the "ambiguous segmentation" problem inherent in separator-less
// Morse code by exploring two simultaneous strategies at every step:
//
//  1. Continue: Treat the current symbols as a prefix for a longer letter
//     (e.g., interpreting '.' as the start of 'A' ".-" rather than 'E').
//
//  2. Branch: If the current path forms a valid letter, commit to that letter
//     (add it to the result) and restart traversal at the root for the next
//     symbols.
//
// Parameters:
//   - remaining: The Morse symbols yet to be processed.
//   - current: The decoded string accumulated so far.
//   - nodeIdx: The current position in the trie (index into t.Nodes).
//   - results: A pointer to the slice accumulating valid decoding results.
func (t *Trie) traverse(remaining, current string, nodeIdx int16, results *[]string) {
	// Base case: All symbols have been consumed.
	if remaining == "" {
		// Valid decodings must end on a letter boundary.
		// If we are back at rootIdx, it means the previous symbol completed
		// a valid letter (triggering a branch) or the string was empty.
		// If we are deep in the trie (nodeIdx != rootIdx), the sequence ended
		// with an incomplete letter code (e.g., trailing "." without termination).
		if nodeIdx == rootIdx {
			*results = append(*results, current)
		}

		return
	}

	// Determine which branch of the trie to take.
	symbol := remaining[0]

	var bit int

	switch symbol {
	case '.':
		bit = 0
	case '-':
		bit = 1
	default:
		// Invalid characters terminate this path. Given the constraint that
		// FindCandidates expects valid Morse symbols, this acts as a safeguard.
		return
	}

	// Look up the child node.
	childIdx := t.Nodes[nodeIdx].Child[bit]

	// Dead end: The sequence follows a path that doesn't exist in the trie.
	// This path is abandoned.
	if childIdx == missingNode {
		return
	}

	nextRemaining := remaining[1:]

	// STRATEGY 1: CONTINUE
	// We do not commit to a letter yet. We simply move deeper into the trie
	// with the same accumulated `current` string.
	t.traverse(nextRemaining, current, childIdx, results)

	// STRATEGY 2: BRANCH
	// Check if the current path forms a valid letter.
	// Val != 0 indicates a valid letter ends at childIdx.
	if t.Nodes[childIdx].Val != 0 {
		// We commit to this letter by appending it to `current` and restarting
		// the traversal at `rootIdx` for the `nextRemaining` symbols.
		newCurrent := current + string(t.Nodes[childIdx].Val)
		t.traverse(nextRemaining, newCurrent, rootIdx, results)
	}
}

// FindCandidates returns all valid decodings (segmentations) for a Morse code
// sequence. Unlike Decode, which requires separators, this function attempts
// to segment the raw symbol stream into all possible letter combinations.
//
// This is useful for decoding ambiguous inputs, such as Morse code bracelets,
// where word boundaries and letter boundaries are not marked.
//
// Parameters:
//   - sequence: A string containing only Morse symbols ('.' and '-').
//
// Returns:
//   - A slice of strings containing all valid decodings. If the sequence
//     cannot be fully segmented into valid letters, an empty slice is returned.
func (t *Trie) FindCandidates(sequence string) []string {
	// We do not pre-allocate the results slice capacity because the number of
	// valid candidates varies wildly based on sequence length and ambiguity.
	var results []string
	t.traverse(sequence, "", rootIdx, &results)

	return results
}

// FindCandidates provides a package-level decode function to return all
// valid decodings for a Morse code sequence.
func FindCandidates(morseCode string) []string {
	return StaticTrie.FindCandidates(morseCode)
}
