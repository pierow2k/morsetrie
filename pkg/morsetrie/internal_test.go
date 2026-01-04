// Package morsetrie provides white-box tests for unexported functions in
// the morsetrie package.
//
//nolint:funlen
package morsetrie

import (
	"strings"
	"testing"
)

// TestTrie_add provides unit testing for add, specifically testing for the
// ErrInvalidElement and ErrDuplicate returns.
func TestTrie_add(t *testing.T) {
	t.Parallel()

	type entry struct {
		code   string
		symbol rune
	}

	tests := []struct {
		name     string
		existing []entry
		code     string
		symbol   rune
		wantErr  bool
	}{
		{
			name:    "valid_entry",
			code:    ".-",
			symbol:  'A',
			wantErr: false,
		},
		{
			name:    "invalid_element",
			code:    ".*",
			symbol:  '*',
			wantErr: true,
		},
		{
			name: "duplicate_entry",
			existing: []entry{
				{code: ".", symbol: 'E'},
			},
			code:    ".",
			symbol:  'E',
			wantErr: true,
		},
		{
			name: "duplicate_collision",
			existing: []entry{
				{code: ".", symbol: 'E'},
			},
			code:    ".",
			symbol:  'I',
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			trie := NewTrie()
			for _, e := range testCase.existing {
				if err := trie.add(e.code, e.symbol); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			gotErr := trie.add(testCase.code, testCase.symbol)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("add() failed: %v", gotErr)
				}

				return
			}

			if testCase.wantErr {
				t.Fatal("add() succeeded unexpectedly")
			}
		})
	}
}

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
