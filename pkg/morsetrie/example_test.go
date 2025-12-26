// Package morsetrie_test provides black-box tests and runnable examples
// for the public API of the morsetrie package.
package morsetrie_test

import (
	"fmt"

	"github.com/pierow2k/morsetrie/pkg/morsetrie"
)

// Example_run demonstrates the use of the morsetrie package.
func Example_run() {
	// TODO: variable name 't' is too short for the scope of its usage.
	t, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	in := ".... . .-.. .-.. --- / .-- --- .-. .-.. -.."

	out, err := t.Decode(in)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
	// Output: HELLO WORLD
}
