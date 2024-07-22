package sensitive_word_matching

import (
	"hash/fnv"
	"sync"
)

type ShareTiredMatching[T int] struct {
	*ShareTried[T]
}

func NewShareTiredMatching[T int](size T) *ShareTiredMatching[T] {
	return &ShareTiredMatching[T]{
		ShareTried: NewShareTrie[T](size),
	}
}

type ShareTried[T int] struct {
	roots []*ShareTrieNode
	mus   []sync.RWMutex
	size  int
}

type ShareTrieNode struct {
	children map[rune]*ShareTrieNode
	isEnd    bool
}

func NewShareTrieNode() *ShareTrieNode {
	return &ShareTrieNode{
		children: make(map[rune]*ShareTrieNode),
		isEnd:    false,
	}
}

func NewShareTrie[T int](size T) *ShareTried[T] {
	st := &ShareTried[T]{
		roots: make([]*ShareTrieNode, size),
		mus:   make([]sync.RWMutex, size),
		size:  int(size),
	}

	for i := 0; i < int(size); i++ {
		st.roots[i] = &ShareTrieNode{children: make(map[rune]*ShareTrieNode)}
	}
	return st
}

func getShareHash(word string, size uint32) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(word))
	return int(h.Sum32() % size)
}

func (s *ShareTried[T]) BatchInsert(words []string) {
	for _, word := range words {
		s.Insert(word)
	}
}

func (s *ShareTried[T]) Insert(word string) {
	shard := getShareHash(word, uint32(s.size))
	s.mus[shard].Lock()
	defer s.mus[shard].Unlock()

	node := s.roots[shard]
	for _, char := range []rune(word) {
		if _, exists := node.children[char]; !exists {
			node.children[char] = NewShareTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

func (s *ShareTried[T]) IsMatching(word string) bool {
	shard := getShareHash(word, uint32(s.size))
	s.mus[shard].RLock()
	defer s.mus[shard].RUnlock()

	node := s.roots[shard]
	for _, char := range []rune(word) {
		if _, exists := node.children[char]; !exists {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

func (s *ShareTried[T]) SearchPrefix(prefix string) bool {
	for shard := 0; shard < s.size; shard++ {
		s.mus[shard].RLock()
		node := s.roots[shard]
		isExist := true

		for _, char := range prefix {
			if _, exists := node.children[char]; !exists {
				isExist = false
				break
			}
			node = node.children[char]
		}
		s.mus[shard].RUnlock()
		if isExist {
			return true
		}
	}
	return false
}

func (s *ShareTried[T]) Delete(word string) bool {
	shard := getShareHash(word, uint32(s.size))
	s.mus[shard].Lock()
	defer s.mus[shard].Unlock()
	return s.delete(s.roots[shard], []rune(word), 0)
}

func (s *ShareTried[T]) delete(node *ShareTrieNode, word []rune, depth int) bool {
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

	shouldDeleteChild := s.delete(childNode, word, depth+1)

	if shouldDeleteChild {
		delete(node.children, char)
		return len(node.children) == 0 && !node.isEnd
	}

	return false
}
