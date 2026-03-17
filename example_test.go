// Testable examples for the public API of the morsetrie package.
package morsetrie_test

import (
	"fmt"

	"github.com/pierow2k/morsetrie"
)

// The Decode function decodes a string of Morse code.
func ExampleDecode() {
	morseCode := "- .... .. ... / .. ... / -- --- .-. ... . - .-. .. ."

	text, err := morsetrie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	// Output:
	// THIS IS MORSETRIE
}

// The forward slash '/' is treated as a word separator.
func ExampleDecode_wordSeparator() {
	text, _ := morsetrie.Decode("... --- ... / ... --- ...")
	fmt.Println(text)

	// Output:
	// SOS SOS
}

// The default static trie supports standard alphanumeric characters,
// punctuation symbols, and the accented 'E'.
func ExampleDecode_extended() {
	morseCode := `..-.. .-.-.- --..-- ---... ..--.. .----. ` +
		`-....- -..-. -.--. -.--.- .-..-. -...- .-.-. .--.-.`

	text, err := morsetrie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	// Output:
	// É.,:?’–/()"=+@
}

// Decode returns ErrUnexpectedChar for invalid input characters.
func ExampleDecode_invalidInput() {
	text, err := morsetrie.Decode("... --- ...!")
	fmt.Println(text)
	fmt.Println(err)

	// Output:
	//
	// unexpected character in morse input: !
}

// Unknown Morse sequences are represented by '?' in the output.
func ExampleDecode_unknownSequence() {
	text, err := morsetrie.Decode(".......") // 7 dots — not a valid sequence
	fmt.Println(text)
	fmt.Println(err)

	// Output:
	// ?
	// <nil>
}

// The Decode method can be called on a custom Trie instance,
// allowing alternative Morse code mappings.
func ExampleTrie_Decode() {
	myTrie := &morsetrie.Trie{
		Nodes: []morsetrie.Node{
			{Val: 0, Child: [2]int16{1, 2}},     // root
			{Val: 'E', Child: [2]int16{-1, -1}}, // just "." -> E
			{Val: 'T', Child: [2]int16{-1, -1}}, // just "-" -> T
		},
	}

	result, err := myTrie.Decode(". -")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	// Output:
	// ET
}
