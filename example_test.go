// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"fmt"

	"github.com/pierow2k/morsetrie"
)

// ExampleDecode provides an example to demonstrate the package-level
// DecodeStatic function.
func ExampleDecode() {
	morseCode := "- .... .. ... / .. ... / -- --- .-. ... . - .-. .. ."

	// Use the static trie to decode the morse code.
	text, err := morsetrie.Decode(morseCode)
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
	trie := morsetrie.StaticTrie

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
	trie := morsetrie.StaticTrie

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
	trie := morsetrie.StaticTrie

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

// ExampleTrie_FindCandidates provides an example to demonstrate the use of
// the FindCandidates method to find potential candidates for a Morse code
// sequence.
func ExampleTrie_FindCandidates() {
	trie := morsetrie.StaticTrie

	sequence := "....--."

	// FindCandidates returns every possible decoding for the Morse code
	// sequence.
	candidates := trie.FindCandidates(sequence)

	// For demonstration purposes we set expectedCandidates. Typically,
	// the resulting candidate would be unknown and would be passed
	// downstream for further processing.
	expectedCandidates := []string{
		"4N", "4TE", "HG", "HME", "HTN", "HTTE", "SP", "SWE", "SAN", "SATE",
		"SEG", "SEME", "SETN", "SETTE", "IUN", "IUTE", "IIG", "IIME", "IITN",
		"IITTE", "IEP", "IEWE", "IEAN", "IEATE", "IEEG", "IEEME", "IEETN",
		"IEETTE", "E3E", "EVN", "EVTE", "ESG", "ESME", "ESTN", "ESTTE",
		"EIP", "EIWE", "EIAN", "EIATE", "EIEG", "EIEME", "EIETN", "EIETTE",
		"EEUN", "EEUTE", "EEIG", "EEIME", "EEITN", "EEITTE", "EEEP", "EEEWE",
		"EEEAN", "EEEATE", "EEEEG", "EEEEME", "EEEETN", "EEEETTE",
	}

	for i := range candidates {
		if expectedCandidates[i] != candidates[i] {
			fmt.Printf("Candidate %d does not match. Expected %s, got %s.",
				i, expectedCandidates[i], candidates[i])
		}
	}

	fmt.Println("Candidates match!")

	// Output:
	// Candidates match!
}
