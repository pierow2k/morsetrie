// Package morsetrie provides white-box tests for unexported functions in
// the morsetrie package.
//

package morsetrie

import (
	"strings"
	"testing"
)

// TestTrie_commit provides unit tests for commit, including specifically
// testing for the `if val == 0` condition.
func TestTrie_commit(t *testing.T) {
	t.Parallel()

	type fields struct {
		Nodes []Node
	}

	type args struct {
		curr int16
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "root_node_no_output",
			fields: fields{Nodes: []Node{{Val: 0}}},
			args:   args{curr: rootIdx},
			want:   "",
		},
		{
			name:   "invalid_node_question_mark",
			fields: fields{Nodes: []Node{{Val: 0}}},
			args:   args{curr: invalidIdx},
			want:   "?",
		},
		{
			name: "valid_node_value",
			fields: fields{Nodes: []Node{
				{Val: 0},   // root
				{Val: 'A'}, // index 1
			}},
			args: args{curr: 1},
			want: "A",
		},
		{
			name: "node_with_zero_value_question_mark",
			fields: fields{Nodes: []Node{
				{Val: 0}, // root
				{Val: 0}, // index 1, no value
			}},
			args: args{curr: 1},
			want: "?",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			trie := &Trie{
				Nodes: testCase.fields.Nodes,
			}

			var builder strings.Builder
			trie.commit(&builder, testCase.args.curr)

			if got := builder.String(); got != testCase.want {
				t.Errorf("commit() = %q, want %q", got, testCase.want)
			}
		})
	}
}

//nolint:gosmopolitan
func Test_reverseString(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty_string",
			args: args{s: ""},
			want: "",
		},
		{
			name: "single_character",
			args: args{s: "a"},
			want: "a",
		},
		{
			name: "alphabet",
			args: args{s: "abcdefghijklmnopqrstuvwxyz"},
			want: "zyxwvutsrqponmlkjihgfedcba",
		},
		{
			name: "non_ascii",
			args: args{s: "你好世界"},
			want: "界世好你",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := reverseString(tt.args.s); got != tt.want {
				t.Errorf("reverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}
