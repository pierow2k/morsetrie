// Package morsetrie_test provides black-box tests, benchmarks, and runnable
// examples for the public API of the morsetrie package.
//
//nolint:funlen
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

// TestTrie_Decode_StaticTrie provides unit tests for the Decode method
// using the static trie.
func TestTrie_Decode_StaticTrie(t *testing.T) {
	t.Parallel()

	type fields struct {
		pairs []morsetrie.MorsePair
	}

	type args struct {
		input string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic_sos",
			fields: fields{
				pairs: []morsetrie.MorsePair{
					{Code: "...", R: 'S'},
					{Code: "---", R: 'O'},
				},
			},
			args: args{
				input: "... --- ...",
			},
			want:    "SOS",
			wantErr: false,
		},
		{
			name: "standard_hello_world",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: ".... . .-.. .-.. --- / .-- --- .-. .-.. -..",
			},
			want:    "HELLO WORLD",
			wantErr: false,
		},
		{
			name: "resume with accented e",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: ".-. ..-.. ... ..- -- ..-..",
			},
			want:    "RÉSUMÉ",
			wantErr: false,
		},
		{
			name: "punctuation",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: ".-.-.- --..-- ---... ..--.. .----. -....- -..-. " +
					"-.--. -.--.- .-..-. -...- .-.-. .--.-.",
			},
			want:    ".,:?’–/()\"=+@",
			wantErr: false,
		},
		{
			name: "multiple consecutive word separators",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: "- .... . // --.- ..- .. -.-. -.- / / " +
					"-... .-. --- .-- -. / ..-. --- -..-",
			},
			want:    "THE QUICK BROWN FOX",
			wantErr: false,
		},
		{
			name: "unknown_sequence",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: "........",
			},
			want:    "?",
			wantErr: false,
		},
		{
			name: "mixed_valid_and_unknown",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: ".... . .-.. .-.. --- / ........ / .-- --- .-. .-.. -..",
			},
			want:    "HELLO ? WORLD",
			wantErr: false,
		},
		{
			name: "invalid_character",
			fields: fields{
				pairs: morsetrie.MorseTable,
			},
			args: args{
				input: "abc",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Load the StaticTrie instead of building.
			trie := morsetrie.StaticTrie

			got, err := trie.Decode(testCase.args.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("Trie.Decode() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if got != testCase.want {
				t.Errorf("Trie.Decode() = %v, want %v", got, testCase.want)
			}
		})
	}
}
