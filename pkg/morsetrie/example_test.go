// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"fmt"

	"github.com/pierow2k/morsetrie/pkg/morsetrie"
)

// ExampleTrie_Decode provides an example to demonstrate the use of the
// Decode method.
func ExampleTrie_Decode() {
	// Build a trie using the built-in MorseTable data.
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	morseCode := ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."

	text, err := trie.Decode(morseCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)
	// Output:
	// HELLO WORLD
}
