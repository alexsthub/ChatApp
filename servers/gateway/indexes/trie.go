package indexes

import (
	"sort"
	"sync"
)

//PRO TIP: if you are having troubles and want to see
//what your trie structure looks like at various points,
//either use the debugger, or try this package:
//https://github.com/davecgh/go-spew

//Trie implements a trie data structure mapping strings to int64s
//that is safe for concurrent use.
type Trie struct {
	root *Node
	size int
	mu   sync.RWMutex
}

// Node is a node in the Trie
type Node struct {
	key      rune
	value    []int64
	children map[rune]*Node
	parent   *Node
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{root: &Node{children: make(map[rune]*Node)}}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.size
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	currentNode := t.root
	runes := []rune(key)
	for i := range runes {
		char := runes[i]
		if _, exists := currentNode.children[char]; !exists {
			newNode := &Node{
				parent:   currentNode,
				key:      char,
				children: make(map[rune]*Node),
			}
			currentNode.children[char] = newNode
		}
		currentNode = currentNode.children[char]
	}
	// If the value does not exist in the set
	if !Contains(currentNode.value, value) {
		currentNode.value = append(currentNode.value, value)
		t.size++
	}
}

// RuneSlice is a helper to sort slice of runes
type RuneSlice []rune

func (p RuneSlice) Len() int           { return len(p) }
func (p RuneSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p RuneSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Contains takes a slice and looks for an element in it. If found it will true
func Contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var results []int64
	currentNode := t.root
	if t.Len() == 0 || prefix == "" || max == 0 {
		return results
	}
	runes := []rune(prefix)
	for i := range runes {
		char := runes[i]
		if nextNode, exists := currentNode.children[char]; exists {
			currentNode = nextNode
		} else {
			return results
		}
	}
	// Made it to the end of the prefix. DFS to get values up to `max`
	queue := []*Node{currentNode}
	for len(queue) > 0 {
		currentNode = queue[0]
		queue = queue[1:]
		// If value exists add to results
		for _, v := range currentNode.value {
			results = append(results, int64(v))
			if len(results) >= max {
				return results
			}
		}
		// Sort keys alphabetically and BFS
		keys := make([]rune, 0)
		for k := range currentNode.children {
			keys = append(keys, k)
		}
		sort.Sort(RuneSlice(keys))
		for _, k := range keys {
			nextNode, _ := currentNode.children[k]
			queue = append(queue, nextNode)
		}
	}
	return results
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	currentNode := t.root
	runes := []rune(key)
	for i := range runes {
		char := runes[i]
		if nextNode, exists := currentNode.children[char]; exists {
			currentNode = nextNode
		}
	}
	// Current node is the node to remove value from
	var index int
	valueList := currentNode.value
	for i, v := range valueList {
		if v == value {
			index = int(i)
		}
	}
	valueList = append(valueList[0:index], valueList[index+1:]...)
	currentNode.value = valueList
	t.size--
	for {
		// If node is leafNode && valueList is empty, trim up until there is a value
		if len(currentNode.children) == 0 && len(currentNode.value) == 0 {
			delNode := currentNode
			currentNode = currentNode.parent
			delete(currentNode.children, delNode.key)
		} else {
			break
		}
	}
}
