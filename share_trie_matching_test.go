package sensitive_word_matching

import (
	"testing"
)

func TestNewShareTiredMatching(t *testing.T) {
	shareTiredMatching := NewShareTiredMatching(32)
	shareTiredMatching.BatchInsert([]string{"哼xx", "哈xxx", "二", "将"})
	var words = []string{"哼xxxxx", "哈", "二", "不如"}
	is := shareTiredMatching.Delete("二")
	t.Log(is)
	for _, v := range words {
		if shareTiredMatching.IsMatching(v) {
			t.Logf("TestShareTiredMatching IsMatching success, %s should be matched", v)
		}

		if shareTiredMatching.SearchPrefix(v) {
			t.Logf("TestShareTiredMatching SearchPrefix success, %s should be matched", v)
		}

	}
}
