package sensitive_word_matching

import (
	"sync"
)

type TrieMatching struct {
	Trie
}

func NewTiredMatching() *TrieMatching {
	return &TrieMatching{
		Trie: Trie{
			root: NewTrieNode(),
		},
	}
}

type Trie struct {
	root *TrieNode
	mu   sync.RWMutex
}

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

func (t *Trie) BatchInsert(words []string) {
	for _, word := range words {
		t.Insert(word)
	}
}

func (t *Trie) Insert(word string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, char := range []rune(word) {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

func (t *Trie) IsMatching(word string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, char := range []rune(word) {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

func (t *Trie) SearchPrefix(prefix string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, char := range []rune(prefix) {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}
	return true
}

func (t *Trie) Delete(word string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.delete(t.root, []rune(word), 0)
}

func (t *Trie) delete(node *TrieNode, word []rune, depth int) bool {
	if node == nil {
		return false
	}

	if depth == len(word) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return len(node.children) == 0
	}

	char := word[depth]
	childNode, exists := node.children[char]
	if !exists {
		return false
	}

	isDelete := t.delete(childNode, word, depth+1)

	if isDelete {
		delete(node.children, char)
		return len(node.children) == 0 && !node.isEnd
	}

	return false
}
