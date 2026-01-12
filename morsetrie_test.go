// Package morsetrie_test provides black-box tests, benchmarks, and runnable
// examples for the public API of the morsetrie package.
//
//nolint:funlen
package morsetrie_test

import (
	"errors"
	"testing"

	"github.com/pierow2k/morsetrie"
)

// TestTrie_Decode provides unit tests for the Decode method.
func TestTrie_Decode(t *testing.T) {
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

			trie, err := morsetrie.BuildTrie(testCase.fields.pairs)
			if err != nil {
				t.Fatalf("BuildTrie() error = %v", err)
			}

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

// TestBuildTrie provides unit tests for BuildTrie, specifically for the error
// return path.
func TestBuildTrie(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pairs     []morsetrie.MorsePair
		wantErr   bool
		wantErrIs error
	}{
		{
			name:    "valid_build",
			pairs:   morsetrie.MorseTable,
			wantErr: false,
		},
		{
			name: "duplicate_code",
			pairs: []morsetrie.MorsePair{
				{Code: ".-", R: 'A'},
				{Code: ".-", R: 'B'},
			},
			wantErr:   true,
			wantErrIs: morsetrie.ErrDuplicate,
		},
		{
			name: "invalid_element",
			pairs: []morsetrie.MorsePair{
				{Code: ".*", R: 'A'},
			},
			wantErr:   true,
			wantErrIs: morsetrie.ErrInvalidElement,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := morsetrie.BuildTrie(testCase.pairs)
			if (err != nil) != testCase.wantErr {
				t.Errorf("BuildTrie() error = %v, wantErr %v", err, testCase.wantErr)
			}

			if testCase.wantErrIs != nil && !errors.Is(err, testCase.wantErrIs) {
				t.Errorf("BuildTrie() error = %v, wantErrIs %v", err, testCase.wantErrIs)
			}
		})
	}
}
