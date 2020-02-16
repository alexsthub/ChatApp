package indexes

import (
	"testing"
)

func TestEmptyTrie(t *testing.T) {
	trie := NewTrie()
	if trie.Len() != 0 {
		t.Errorf("Length should be 0. Recieved %d", trie.Len())
	}
}

func TestAddSmallTrie(t *testing.T) {
	trie := NewTrie()
	trie.Add("alex", 1)
	trie.Add("poop", 2)
	trie.Add("pee", 3)
	trie.Add("toilet", 4)
	trie.Add("bathroom", 5)
	if trie.Len() != 5 {
		t.Errorf("Length should be 5. Recieved %d", trie.Len())
	}
}

func TestDuplicates(t *testing.T) {
	trie := NewTrie()
	trie.Add("alex", 10)
	trie.Add("alex", 10)
	if trie.Len() != 1 {
		t.Errorf("Should only be 1 value in the trie. Got %d", trie.Len())
	}
}

func TestSimpleFind(t *testing.T) {
	trie := NewTrie()
	trie.Add("alex", 10)
	values := trie.Find("alex", 1)
	if len(values) != 1 || values[0] != 10 {
		t.Errorf("Should get just value of 10 back but received %d values", len(values))
	}
}

func TestFind(t *testing.T) {
	trie := NewTrie()
	results := trie.Find("test", 1)
	if len(results) != 0 {
		t.Errorf("Searching an empty tree should return a nil slice")
	}
	trie.Add("test", 10)
	emptyPrefixResults := trie.Find("", 1)
	if emptyPrefixResults != nil {
		t.Errorf("Searching with empty prefix should return a nil slice but did not")
	}
	maxResults := trie.Find("test", 0)
	if maxResults != nil {
		t.Errorf("Searching with max of 0 should return a nil slice but did not")
	}
	badPrefixResults := trie.Find("bad", 1)
	if badPrefixResults != nil {
		t.Errorf("Searching with a prefix that isn't in the trie should return a nil slice but did not")
	}
}

func TestFindSorting(t *testing.T) {

}

func TestDelete(t *testing.T) {
	trie := NewTrie()
}
