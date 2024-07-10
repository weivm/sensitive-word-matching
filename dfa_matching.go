package sensitive_word_matching

import (
	"fmt"
	"unicode"
)

type DfaMatching struct {
	startState       string                     //状态初始状态
	states           map[string]struct{}        //存储所有节点状态
	acceptStates     map[string]struct{}        //每个关键词终态
	alphabet         map[rune]struct{}          //存储所有字符集
	transitionStates map[string]map[rune]string //存储状态转移
}

// 动态生成支持的字符集
func generateAlphabet() map[rune]struct{} {
	alphabet := make(map[rune]struct{})

	// 添加 ASCII 字符（包括字母和数字）
	for r := rune(0); r < 128; r++ {
		if unicode.IsPrint(r) {
			alphabet[r] = struct{}{}
		}
	}

	// 添加汉字字符范围（简体字和繁体字）
	for r := rune(0x4E00); r <= rune(0x9FFF); r++ {
		alphabet[r] = struct{}{}
	}

	// 添加常见标点符号
	for _, r := range []rune{' ', '.', ',', '!', '?', ';', ':', '-', '_'} {
		alphabet[r] = struct{}{}
	}

	return alphabet
}

func NewDfaMatching() *DfaMatching {
	return &DfaMatching{
		states:           make(map[string]struct{}),
		alphabet:         generateAlphabet(),
		transitionStates: make(map[string]map[rune]string),
		startState:       StartState,
		acceptStates:     make(map[string]struct{}),
	}
}

func (d *DfaMatching) GenerateSensitiveKeyWords(keyWords []string) {
	for _, keyword := range keyWords {
		state := StartState
		for _, char := range keyword {
			if _, exists := d.transitionStates[state]; !exists {
				d.transitionStates[state] = make(map[rune]string)
			}
			if nextState, exists := d.transitionStates[state][char]; !exists {
				nextState = fmt.Sprintf("state_%d", len(d.states)+1)
				d.states[nextState] = struct{}{}
				d.transitionStates[state][char] = nextState
			}
			state = d.transitionStates[state][char]
		}
		d.acceptStates[state] = struct{}{}
	}
	d.generateFailAlphabet()
}

func (d *DfaMatching) generateFailAlphabet() {
	// 添加失败状态，如果已经存在则不添加
	if _, exists := d.states["error"]; !exists {
		d.states["error"] = struct{}{}
		d.transitionStates["error"] = make(map[rune]string)
		for char := range d.alphabet {
			d.transitionStates["error"][char] = "error"
		}
	}

	// 为每个状态添加对失败状态的转移
	for _, transitions := range d.transitionStates {
		for char := range d.alphabet {
			if _, exists := transitions[char]; !exists {
				transitions[char] = "error"
			}
		}
	}
}

func (d *DfaMatching) IsMatching(keyword string) bool {
	currentState := d.startState
	for _, char := range keyword {
		if _, ok := d.alphabet[char]; !ok {
			continue // 忽略不在字母表中的字符
		}

		nextState, ok := d.transitionStates[currentState][char]
		if !ok {
			return false // 字符转移不存在，返回不匹配
		}
		currentState = nextState
	}
	_, accepted := d.acceptStates[currentState]
	return accepted
}
