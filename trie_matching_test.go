package sensitive_word_matching

import (
	"testing"
)

func TestNewTiredMatching(t *testing.T) {
	tiredMatching := NewTiredMatching()
	tiredMatching.BatchInsert([]string{"哼xx", "哈xxx", "二", "将"})
	var words = []string{"哼", "哈", "二", "不如"}
	is := tiredMatching.Delete("二")
	t.Log(is)
	for _, v := range words {
		if tiredMatching.IsMatching(v) {
			t.Logf("TestTiredMatching IsMatching success, %s should be matched", v)
		}

		if tiredMatching.SearchPrefix(v) {
			t.Logf("TestTiredMatching SearchPrefix success, %s should be matched", v)
		}

	}
}
