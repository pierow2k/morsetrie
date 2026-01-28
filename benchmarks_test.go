// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

// input is the morse code of the text:
// THE QUICK BROWN FOOX JUMPS OVER THE LAZY DOG 0123456789.
const input = "- .... . / --.- ..- .. -.-. -.- / -... .-. --- .-- -. / " +
	"..-. --- -..-. --- -..- / .--- ..- -- .--. ... / --- ...- . .-. / " +
	"- .... . / .-.. .- --.. -.-- / -.. --- --. / ----- .---- ..--- " +
	"...-- ....- ..... -.... --... ---.. ----."

const sequence = "....--."

var (
	trieSink   *morsetrie.Trie
	decodeSink string
)

// BenchmarkDecode measures the performance of the Decode function.
func BenchmarkDecode(b *testing.B) {
	trie := morsetrie.StaticTrie
	for b.Loop() {
		out, _ := trie.Decode(input)
		decodeSink = out
	}
}

// BenchmarkLoadStaticTrie measures the performance of loading the static
// trie.
func BenchmarkLoadStaticTrie(b *testing.B) {
	for b.Loop() {
		trieSink = morsetrie.StaticTrie
	}
}

// BenchmarkFindCandidates measures the performance of the FindCandidates
// function.
func BenchmarkFindCandidates(b *testing.B) {
	trie := morsetrie.StaticTrie

	// FindCandidates returns every possible decoding for the Morse code
	// sequence.
	for b.Loop() {
		candidates := trie.FindCandidates(sequence)
		_ = candidates
	}
}
