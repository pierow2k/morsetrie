// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"fmt"

	"github.com/pierow2k/morsetrie"
)

// ExampleDecode provides an example to demonstrate the package-level
// Decode function.
func ExampleDecode() {
	morseCode := "- .... .. ... / .. ... / -- --- .-. ... . - .-. .. ."

	// Build a trie and decode the morse code.
	text, err := morsetrie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	// Output:
	// THIS IS MORSETRIE
}

// ExampleDecodeStatic provides an example to demonstrate the package-level
// DecodeStatic function.
func ExampleDecodeStatic() {
	morseCode := "- .... .. ... / .. ... / -- --- .-. ... . - .-. .. ."

	// Use the static trie to decode the morse code.
	text, err := morsetrie.DecodeStatic(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	// Output:
	// THIS IS MORSETRIE
}

// ExampleTrie_Decode_alphanumeric provides an example to demonstrate the
// use of the Decode method to decode alphanumeric characters.
func ExampleTrie_Decode_alphanumeric() {
	// Build the trie using the built-in morse code data from morsetrie's
	// `MorseTable`.
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	// Define morse code input using standard alphanumeric characters.
	morseCode := ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."

	// Call trie.Decode to decode the morse code data.
	text, err := trie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println("Alphanumeric:")
	fmt.Println(text)

	// Output:
	// Alphanumeric:
	// HELLO WORLD
}

// ExampleTrie_Decode_extended provides an example to demonstrate the
// use of the Decode method to decode the extended character set.
func ExampleTrie_Decode_extended() {
	// Build the trie using the built-in morse code data from morsetrie's
	// `MorseTable`.
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	// The package also supports the ITU specification for the accented 'E'
	// as well as punctuation symbols.
	extendedMorse := "..-.. .-.-.- --..-- ---... ..--.. .----. " +
		"-....- -..-. -.--. -.--.- .-..-. -...- .-.-. .--.-."

	extended, err := trie.Decode(extendedMorse)
	if err != nil {
		panic(err)
	}

	fmt.Println("Accented 'E' and Punctuation:")
	fmt.Println(extended)

	// Output:
	// Accented 'E' and Punctuation:
	// É.,:?’–/()"=+@
}

func ExampleTrie_Decode_prosign() {
	// Procedural sign (or prosign) shorthand signals are not supported
	// by the package since these can not be directly mapped to a rune.
	// Prosigns are decoded as an unknown character, represented as '?'.
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	prosign := "-.-.-"

	prosignText, err := trie.Decode(prosign)
	if err != nil {
		panic(err)
	}

	fmt.Println("Prosign:")
	fmt.Println(prosignText)

	// Output:
	// Prosign:
	// ?
}

// ExampleBuildTrie provides an example to demonstrate the use of the
// BuildTrie function.
func ExampleBuildTrie() {
	// The BuildTrie function can be used to build the trie with a custom
	// data set.

	// Define custom morse code/rune pairs.
	pairs := []morsetrie.MorsePair{
		{Code: ".-", R: 'A'},
		{Code: "-...", R: 'B'},
	}

	// Build the trie by calling BuildTrie.
	trie, err := morsetrie.BuildTrie(pairs)
	if err != nil {
		panic(err)
	}

	// Define the morse code data.
	morseCode := ".- -... -... .-"

	// Decode uses the trie to decode the morse code data.
	text, err := trie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	// Output:
	// ABBA
}
