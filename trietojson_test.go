// Package morsetrie_test provides black-box tests, benchmarks, and
// testable examples for the public API of the morsetrie package.
package morsetrie_test

import (
	"testing"

	"github.com/pierow2k/morsetrie"
)

// TestTrie_ToJSON provides unit tests for the Trie_ToJSON method.
func TestTrie_ToJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func() *morsetrie.Trie
		want    string
		wantErr bool
	}{
		{
			name: "nil trie",
			setup: func() *morsetrie.Trie {
				return nil
			},
			want: "null",
		},
		{
			name: "empty trie",
			setup: func() *morsetrie.Trie {
				return morsetrie.NewTrie()
			},
			want: `{
  "Nodes": [
    {
      "val": "",
      "child": [
        -1,
        -1
      ]
    }
  ]
}`,
		},
		{
			name: "single element E",
			setup: func() *morsetrie.Trie {
				tr, _ := morsetrie.BuildTrie([]morsetrie.MorsePair{
					{Code: ".", R: 'E'},
				})
				return tr
			},
			want: `{
  "Nodes": [
    {
      "val": "",
      "child": [
        1,
        -1
      ]
    },
    {
      "val": "E",
      "child": [
        -1,
        -1
      ]
    }
  ]
}`,
		},
		{
			name: "multiple elements E and T",
			setup: func() *morsetrie.Trie {
				tr, _ := morsetrie.BuildTrie([]morsetrie.MorsePair{
					{Code: ".", R: 'E'},
					{Code: "-", R: 'T'},
				})
				return tr
			},
			want: `{
  "Nodes": [
    {
      "val": "",
      "child": [
        1,
        2
      ]
    },
    {
      "val": "E",
      "child": [
        -1,
        -1
      ]
    },
    {
      "val": "T",
      "child": [
        -1,
        -1
      ]
    }
  ]
}`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			tr := testCase.setup()
			got, gotErr := tr.ToJSON()

			if (gotErr != nil) != testCase.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", gotErr, testCase.wantErr)
				return
			}
			if got != testCase.want {
				t.Errorf("ToJSON() got:\n%v\nwant:\n%v", got, testCase.want)
			}
		})
	}
}
