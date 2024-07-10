package sensitive_word_matching

import "fmt"

type State int

const (
	FirstState State = iota
)

var StartState = fmt.Sprintf("state_%d", FirstState)
