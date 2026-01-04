// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie/pkg/morsetrie"
)

// BenchmarkDecode measures the performance of the Decode function.
func BenchmarkDecode(b *testing.B) {
	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		b.Fatal(err)
	}

	input := "- .... . / --.- ..- .. -.-. -.- / -... .-. --- .-- -. / ..-. " +
		"--- -..- / .--- ..- -- .--. ... / --- ...- . .-. / - .... . / .-.. " +
		".- --.. -.-- / -.. --- --. / ----- .---- ..--- ...-- ....- ..... " +
		"-.... --... ---.. ----."

	for b.Loop() {
		_, _ = trie.Decode(input)
	}
}
