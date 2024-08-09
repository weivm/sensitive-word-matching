package sensitive_word_matching

import (
	"testing"
)

func TestNewAcMatching(t *testing.T) {
	ac := NewAcMatching()
	ac.BatchInsert([]string{"哼xx", "哈xxx", "二", "将", "abcd"})
	ac.buildFailNode()
	results := ac.Search("abcd")
	t.Log(results)
	t.Log(ac.IsMatching("abcd"))
}
