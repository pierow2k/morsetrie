// Package morsetrie_test provides black-box tests, benchmarks, and runnable
// examples for the public API of the morsetrie package.
//

package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

// TestDecode tests the package-level decode function for each
// character that is supported by the static trie.
func TestDecode(t *testing.T) {
	t.Parallel()

	morseCode := ".- / -... / -.-. / -.. / . / ..-.. / ..-. / --. / .... / " +
		".. / .--- / -.- / .-.. / -- / -. / --- / .--. / --.- / .-. / ... / - / " +
		"..- / ...- / .-- / -..- / -.-- / --.. / ----- / .---- / ..--- / ...-- / " +
		"....- / ..... / -.... / --... / ---.. / ----. / .-.-.- / --..-- / " +
		"---... / ..--.. / .----. / -....- / -..-. / -.--. / -.--.- / .-..-. / " +
		"-...- / .-.-. / .--.-."

	want := "A B C D E É F G H I J K L M N O P Q R S T U V W X Y Z " +
		"0 1 2 3 4 5 6 7 8 9 . , : ? ’ – / ( ) \" = + @"

	result, _ := morsetrie.Decode(morseCode)
	if result != want {
		t.Errorf("TestDecode() = %v, want %v", result, want)
	}
}

// TestTrie_Decode provides unit tests for the Decode method.
func TestTrie_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		morseCode string
		want      string
		wantErr   bool
	}{
		{
			name:      "SOS",
			morseCode: "... --- ...",
			want:      "SOS",
			wantErr:   false,
		},
		{
			name:      "Hello World",
			morseCode: ".... . .-.. .-.. --- / .-- --- .-. .-.. -..",
			want:      "HELLO WORLD",
			wantErr:   false,
		},
		{
			name:      "Unknown Morse Code Sequence",
			morseCode: "........",
			want:      "?",
			wantErr:   false,
		},
		{
			name:      "Invalid Input",
			morseCode: "abcd",
			want:      "",
			wantErr:   true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			trie := morsetrie.StaticTrie

			got, gotErr := trie.Decode(testCase.morseCode)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("Decode() failed: %v", gotErr)
				}

				return
			}

			if testCase.wantErr {
				t.Fatal("Decode() succeeded unexpectedly")
			}

			if got != testCase.want {
				t.Errorf("Decode() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestFindCandidates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		morseCode string
		want      []string
	}{
		{
			name:      "Single Dot",
			morseCode: ".",
			want:      []string{"E"},
		},
		{
			name:      "Single Dash",
			morseCode: "-",
			want:      []string{"T"},
		},
		{
			name:      "Two Dots",
			morseCode: "..",
			want:      []string{"I", "EE"},
		},
		{
			name:      "Dot Dash",
			morseCode: ".-",
			want:      []string{"A", "ET"},
		},
		{
			name:      "Three Dots",
			morseCode: "...",
			want:      []string{"S", "IE", "EI", "EEE"},
		},
		{
			name:      "Invalid Characters",
			morseCode: "ABC",
			want:      nil,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := morsetrie.FindCandidates(testCase.morseCode)

			if len(got) != len(testCase.want) {
				t.Errorf("FindCandidates() got %d candidates, want %d", len(got), len(testCase.want))

				return
			}

			for i := range got {
				if got[i] != testCase.want[i] {
					t.Errorf("FindCandidates() candidate[%d] = %v, want %v", i, got[i], testCase.want[i])
				}
			}
		})
	}
}

func TestGenerateCyclicRotations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sequence string
		reverse  bool
		want     []string
	}{
		{
			name:     "Empty Sequence",
			sequence: "",
			reverse:  false,
			want:     nil,
		},
		{
			name:     "Single Character",
			sequence: "A",
			reverse:  false,
			want:     []string{"A"},
		},
		{
			name:     "Two Characters",
			sequence: "AB",
			reverse:  false,
			want:     []string{"AB", "BA"},
		},
		{
			name:     "Three Characters with Reverse",
			sequence: "ABC",
			reverse:  true,
			want:     []string{"ABC", "BCA", "CAB", "CBA", "BAC", "ACB"},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := morsetrie.GenerateCyclicRotations(testCase.sequence, testCase.reverse)

			if len(got) != len(testCase.want) {
				t.Errorf("FindCandidates() got %d candidates, want %d", len(got), len(testCase.want))

				return
			}

			for i := range got {
				if got[i] != testCase.want[i] {
					t.Errorf("GenerateCyclicRotations() sequence[%d] = %v, want %v", i, got[i], testCase.want[i])
				}
			}
		})
	}
}
