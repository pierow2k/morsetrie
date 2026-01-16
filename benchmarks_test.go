// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

// input is the morse code of the text:
// THE QUICK BROWN FOOX JUMPS OVER THE LAZY DOG 0123456789
const input = "- .... . / --.- ..- .. -.-. -.- / -... .-. --- .-- -. / " +
	"..-. --- -..-. --- -..- / .--- ..- -- .--. ... / --- ...- . .-. / " +
	"- .... . / .-.. .- --.. -.-- / -.. --- --. / ----- .---- ..--- " +
	"...-- ....- ..... -.... --... ---.. ----."

var (
	trieSink   *morsetrie.Trie
	decodeSink string
)

// BenchmarkDecode measures the performance of the Decode method.
// It does not include the time to build the trie, measuring only
// the process of decoding.
func BenchmarkDecode(b *testing.B) {
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		out, _ := trie.Decode(input)
		decodeSink = out
	}
}

// BenchmarkDecode_LoadingTrie measures the performance of the Decode
// method and includes the process of building the trie using the data
// from MorseTable.
func BenchmarkDecode_LoadingTrie(b *testing.B) {
	for b.Loop() {
		trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
		if err != nil {
			b.Fatal(err)
		}
		out, _ := trie.Decode(input)
		decodeSink = out
		trieSink = trie
	}
}

// BenchmarkDecode_Using_StaticTrie measures the performance of the Decode
// method using the static trie.
func BenchmarkDecode_Using_StaticTrie(b *testing.B) {
	trie := morsetrie.StaticTrie
	for b.Loop() {
		out, _ := trie.Decode(input)
		decodeSink = out
	}
}

// BenchmarkBuildTrie measures the performance of the BuildTrie function.
// It builds the trie using data from the MorseTable.
func BenchmarkBuildTrie(b *testing.B) {
	for b.Loop() {
		trie, _ := morsetrie.BuildTrie(morsetrie.MorseTable)
		trieSink = trie
	}
}

// BenchmarkLoadStaticTrie measures the performance of loading the static
// trie.
func BenchmarkLoadStaticTrie(b *testing.B) {
	for b.Loop() {
		trieSink = morsetrie.StaticTrie
	}
}
