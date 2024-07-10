package sensitive_word_matching

import (
	"testing"
)

func TestNewDfaMatching(t *testing.T) {

	dfa := NewDfaMatching()
	var keywords = []string{"哼xx", "哈", "二", "将"}
	dfa.GenerateSensitiveKeyWords(keywords)

	var words = []string{"哼xx", "哈", "猪八戒", "不如"}

	for _, word := range words {
		if dfa.IsMatching(word) {
			t.Logf("TestDfaMatching success, %s should be matched", word)
		}
	}
}
