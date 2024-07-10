package sensitive_word_matching

type MatchingStrategy interface {
	IsMatching(keyword string) bool
}
