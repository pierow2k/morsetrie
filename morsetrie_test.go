// Black-box tests for exported functions in the morsetrie package.
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
