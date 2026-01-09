// Package morsetrie_test provides black-box tests, benchmarks, and runnable
// examples for the public API of the morsetrie package.
package morsetrie_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pierow2k/morsetrie"
)

// FuzzDecode performs fuzz testing for the Decode method.
func FuzzDecode(f *testing.F) {
	f.Add(".... . .-.. .-.. --- / .-- --- .-. .-.. -..")
	f.Add("-- --- .-. ... . - .-. .. .")
	f.Add("Abracadabra - Plain text should fail.")
	f.Add("")
	f.Add("       ")
	f.Add("///")

	trie, err := morsetrie.BuildTrie(morsetrie.MorseTable)
	if err != nil {
		panic(err)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_, err := trie.Decode(input)

		// Check if input contains invalid characters.
		isValid := true

		for _, r := range input {
			if !strings.ContainsRune(".- / \t\n\r", r) {
				isValid = false

				break
			}
		}

		if isValid && err != nil {
			t.Errorf("Decode(%q) unexpected error: %v", input, err)
		} else if !isValid && !errors.Is(err, morsetrie.ErrUnexpectedChar) {
			t.Errorf("Decode(%q) expected ErrUnexpectedChar, got %v", input, err)
		}
	})
}

// FuzzBuildTrie performs fuzz testing for the BuildTrie function.
func FuzzBuildTrie(f *testing.F) {
	f.Add(".-", 'A')
	f.Add("...", 'S')
	f.Add("invalid", 'X')

	f.Fuzz(func(t *testing.T, code string, r rune) {
		pairs := []morsetrie.MorsePair{
			{Code: code, R: r},
		}

		_, err := morsetrie.BuildTrie(pairs)

		// Check for invalid characters in code.
		isValid := true

		for _, c := range code {
			if c != '.' && c != '-' {
				isValid = false

				break
			}
		}

		if !isValid {
			if err == nil {
				t.Errorf("BuildTrie with invalid code %q should fail", code)
			} else if !errors.Is(err, morsetrie.ErrInvalidElement) {
				t.Errorf("BuildTrie with invalid code %q: expected ErrInvalidElement, got %v", code, err)
			}
		} else if err != nil && !errors.Is(err, morsetrie.ErrTrieFull) {
			t.Errorf("BuildTrie with valid code %q failed: %v", code, err)
		}
	})
}
