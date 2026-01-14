// Package morsetrie implements trie-based decoding for morse code.
package morsetrie

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	// defaultTrieCapacity is the initial allocation size for the Trie node slice.
	defaultTrieCapacity = 64
	// decodeAllocDivisor estimates decoded text size (approx 1/3 of morse input).
	decodeAllocDivisor = 3
)

const (
	rootIdx     = int16(0)
	invalidIdx  = int16(-2) // Represents traversing off the trie.
	missingNode = int16(-1) // Represents a child definition that doesn't exist.
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
	// ErrTrieFull is returned when the Trie exceeds the maximum number of nodes
	// addressable by int16.
	ErrTrieFull = errors.New("trie capacity exceeded")
)

// MorsePair is a single Morse-code mapping entry.
type MorsePair struct {
	Code string
	R    rune
}

// MorseTable is the ITU M.1677 standard International Morse code mapping for.
var MorseTable = []MorsePair{
	{".-", 'A'},
	{"-...", 'B'},
	{"-.-.", 'C'},
	{"-..", 'D'},
	{".", 'E'},
	{"..-..", 'É'},
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

	{".-.-.-", '.'},
	{"--..--", ','},
	{"---...", ':'},
	{"..--..", '?'},
	{".----.", '’'},
	{"-....-", '–'},
	{"-..-.", '/'},
	{"-.--.", '('},
	{"-.--.-", ')'},
	{".-..-.", '"'},
	{"-...-", '='},
	{".-.-.", '+'},
	{".--.-.", '@'},
}

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

// MarshalJSON customizes the JSON representation of a Node.
// It converts the Val rune into a human-readable string.
func (n Node) MarshalJSON() ([]byte, error) {
	var valStr string
	if n.Val != 0 {
		valStr = string(n.Val)
	}

	// We use an anonymous struct to define the JSON shape.
	return json.Marshal(struct {
		Val   string   `json:"val"`
		Child [2]int16 `json:"child"`
	}{
		Val:   valStr,
		Child: n.Child,
	})
}

// ToJSON returns the Trie and all its nodes as a formatted JSON string.
func (t *Trie) ToJSON() (string, error) {
	if t == nil {
		return "null", nil
	}
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// NewTrie creates a new, empty Morse decode Trie.
func NewTrie() *Trie {
	trie := &Trie{Nodes: make([]Node, 1, defaultTrieCapacity)}

	// Initialize root children to indicate they are missing.
	trie.Nodes[rootIdx].Child[0] = missingNode
	trie.Nodes[rootIdx].Child[1] = missingNode

	return trie
}

// add inserts a Morse code sequence and its corresponding rune into the Trie.
// It traverses the Trie based on the input code ('.' for child 0, '-' for child 1),
// creating new nodes as necessary.
//
// It returns ErrInvalidElement if the code contains characters other than '.' or '-'.
// It returns ErrDuplicate if the code sequence is already registered.
// It returns ErrTrieFull if the Trie capacity is exceeded.
func (t *Trie) add(code string, symbol rune) error {
	idx := rootIdx

	for charIdx := range len(code) {
		var bit int

		switch code[charIdx] {
		case '.':
			bit = 0
		case '-':
			bit = 1
		default:
			return fmt.Errorf("%w: %q in %q", ErrInvalidElement, code[charIdx], code)
		}

		next := t.Nodes[idx].Child[bit]
		if next == missingNode {
			// Check for integer overflow before casting to int16.
			if len(t.Nodes) > math.MaxInt16 {
				return ErrTrieFull
			}

			// Disabled gosec linter warning. Bounds are checked prior to conversion.
			next = int16(len(t.Nodes)) //nolint:gosec

			// Create new node with missing children.
			t.Nodes = append(t.Nodes, Node{Child: [2]int16{missingNode, missingNode}})
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

// BuildTrie constructs a new Trie from a slice of MorsePair. It initializes a
// new Trie and populates it by adding each pair from the input slice.
//
// An error is returned if the input pairs contain invalid data, such as
// duplicate codes (ErrDuplicate) or codes with invalid characters
// (ErrInvalidElement).
func BuildTrie(pairs []MorsePair) (*Trie, error) {
	trie := NewTrie()
	for _, p := range pairs {
		if err := trie.add(p.Code, p.R); err != nil {
			return nil, err
		}
	}

	return trie, nil
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
