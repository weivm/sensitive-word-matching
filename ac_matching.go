package sensitive_word_matching

import (
	"container/list"
	"sync"
)

type AcMatching struct {
	root *AcNode
	mu   sync.RWMutex
}

func NewAcMatching() *AcMatching {
	return &AcMatching{
		root: newNode()}
}

type AcNode struct {
	children map[rune]*AcNode //每个节点的子节点
	fail     *AcNode          //失败节点指针
	isEnd    bool             //是否是结束节点
	output   []string         //节点内容
}

func newNode() *AcNode {
	return &AcNode{
		children: make(map[rune]*AcNode),
		fail:     nil,
		output:   []string{},
		isEnd:    false,
	}
}

func (ac *AcMatching) BatchInsert(words []string) {
	for _, word := range words {
		ac.Insert(word)
	}
}

func (ac *AcMatching) Insert(word string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	node := ac.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = newNode()
		}
		node = node.children[ch]
	}
	node.output = append(node.output, word)
	node.isEnd = true
}

func (ac *AcMatching) buildFailNode() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	queue := list.New()
	for _, node := range ac.root.children {
		node.fail = ac.root  //将所有根节点的字节点里失败节点全出指向初始节点
		queue.PushBack(node) //将根节点的子节点加入队列
	}

	for queue.Len() > 0 { //广度优先遍历
		front := queue.Front()
		if front == nil {
			continue
		}

		current, ok := queue.Remove(front).(*AcNode)
		if !ok {
			continue
		}

		for ch, child := range current.children {
			queue.PushBack(child)

			failNode := current.fail
			for failNode != nil { //不断向上找符合失败节点中子节点是否符合当前字符数据
				if failChild, ok := failNode.children[ch]; ok {
					child.fail = failChild
					child.output = append(child.output, failChild.output...)
					break
				}
				failNode = failNode.fail
			}
			if failNode == nil { //如果找不到，指向根节点
				child.fail = ac.root
			}
		}
	}
}

func (ac *AcMatching) Search(word string) map[string][]int {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	node := ac.root
	result := make(map[string][]int)

	for i, w := range word {
		for node != ac.root && node.children[w] == nil { //非根节点不断向上遍历失败节点,找到节点的字节点中有包含该字符的节点
			node = node.fail
		}

		if nextNode, ok := node.children[w]; ok {
			node = nextNode
		} else {
			node = ac.root
		}

		tempNode := node
		for tempNode != nil {
			if tempNode.isEnd {
				for _, pattern := range tempNode.output {
					if _, ok := result[pattern]; !ok {
						result[pattern] = []int{}
					}

					start := i - len(pattern) + 1
					if start > 0 {
						result[pattern] = append(result[pattern], start)
					}
				}
			}
			tempNode = tempNode.fail
		}
	}

	return result
}

func (ac *AcMatching) IsMatching(word string) bool {
	results := ac.Search(word)
	return len(results) > 0
}
