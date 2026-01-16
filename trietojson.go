// Package morsetrie implements trie-based decoding for morse code.
package morsetrie

import "encoding/json"

// MarshalJSON customizes the JSON representation of a Node.
// It converts the Val rune into a human-readable string.
func (n Node) MarshalJSON() ([]byte, error) {
	var valStr string
	if n.Val != 0 {
		valStr = string(n.Val)
	}

	// We use an anonymous struct to define the JSON shape.
	return json.Marshal(struct {
		Val   string   `json:"val"`
		Child [2]int16 `json:"child"`
	}{
		Val:   valStr,
		Child: n.Child,
	})
}

// ToJSON returns the Trie and all its nodes as a formatted JSON string.
func (t *Trie) ToJSON() (string, error) {
	if t == nil {
		return "null", nil
	}
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
