// Benchmarks for the public API of the morsetrie package.
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

var trieSink *morsetrie.Trie

// BenchmarkDecode measures the performance of the Decode function.
func BenchmarkDecode(b *testing.B) {
	// Input is a Morse encoding of the text:
	// "THE QUICK BROWN FOOX JUMPS OVER THE LAZY DOG 0123456789"
	const input = `- .... . / --.- ..- .. -.-. -.- / -... .-. --- .-- -. / ` +
		`..-. --- -..-. --- -..- / .--- ..- -- .--. ... / --- ...- . .-. / ` +
		`- .... . / .-.. .- --.. -.-- / -.. --- --. / ----- .---- ..--- ` +
		`...-- ....- ..... -.... --... ---.. ----.`

	for b.Loop() {
		out, _ := morsetrie.Decode(input)
		_ = out
	}
}

// BenchmarkLoadStaticTrie measures the performance of loading the static
// trie.
func BenchmarkLoadStaticTrie(b *testing.B) {
	for b.Loop() {
		trieSink = morsetrie.StaticTrie
	}
}
