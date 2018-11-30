package fcompl

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	pinyin "github.com/mozillazg/go-pinyin"
)

const (
	zh  string = "chs"
	eng string = "en"
)

var (
	stopWords = []string{
		"I",
		"a",
		"about",
		"an",
		"are",
		"as",
		"at",
		"be",
		"by",
		"com",
		"for",
		"from",
		"how",
		"in",
		"is",
		"it",
		"of",
		"on",
		"or",
		"that",
		"the",
		"this",
		"to",
		"was",
		"what",
		"when",
		"where",
		"who",
		"will",
		"with",
		"the",
		"www",
	}
	stopDict  map[string]bool
	pinyinCfg = pinyin.NewArgs()
)

type TrieNode struct {
	Children  map[string]*TrieNode
	ContainID []int
}

func newTrieNode() *TrieNode {
	return &TrieNode{make(map[string]*TrieNode), nil}
}

func (tn *TrieNode) insertPhrase(phrase string, index int) {
	n := tn
	lang := detectLang(phrase)
	var keys []string
	switch lang {
	case eng:
		keys = strings.Split(phrase, " ")
	case zh:
		keys = pinyin.LazyConvert(phrase, &pinyinCfg)
		//log.Println(keys)
	}

	//log.Printf("Insert %s with index %d", phrse, index)
	for _, k := range keys {
		if stopDict != nil {
			if _, ok := stopDict[k]; ok {
				continue
			}
		}
		if _, has := n.Children[k]; !has {
			n.Children[k] = &TrieNode{make(map[string]*TrieNode), []int{index}}
		} else {
			//log.Printf("%s has key %s", phrse, k)
			n.Children[k].ContainID = append(n.Children[k].ContainID, index)
		}
		n = n.Children[k]
	}
}

func (root *TrieNode) Find(text string) []int {
	text = strings.ToLower(text)
	n := root
	foundPhrases := []int{}
	candis := strings.Split(text, " ")

	for i := 0; i < len(candis); {
		if stopDict != nil {
			// ignore stop words
			if _, ok := stopDict[candis[i]]; ok {
				i++
				continue
			}
		}
		// deeper Children found
		if node, ok := n.Children[candis[i]]; ok {
			if len(node.ContainID) == 1 {
				return node.ContainID
			}
			n = node
			i++
			continue
		}
		// no found one yet
		if n == root {
			i++
			continue
		}
		// go back match from root
		n = root
	}

	if n != root && len(n.ContainID) != 0 {
		log.Printf("append %v", n.ContainID)
		foundPhrases = append(foundPhrases, n.ContainID...)
	}
	return foundPhrases
}

// Compress cut the node below that no more branches to reduce the capacity of the whole tree
func Compress(root *TrieNode) {
	if len(root.ContainID) == 1 {
		for k := range root.Children {
			delete(root.Children, k)
		}
		return
	}
	for _, n := range root.Children {
		Compress(n)
	}
}

func (root *TrieNode) print() {
	nodes := []struct {
		Val   *TrieNode
		Depth int
	}{{root, 0}}
	curDepth := 0
	for len(nodes) > 0 {
		node := nodes[0]
		if node.Depth > curDepth {
			curDepth = node.Depth
			fmt.Println("")
		}
		nodes = nodes[1:]

		for k, child := range node.Val.Children {
			fmt.Printf(" %s ", k)
			nodes = append(nodes, struct {
				Val   *TrieNode
				Depth int
			}{child, node.Depth + 1})
		}
	}
}

// Save use gob to serialize TrieNode struct into file, remember to compress whole tree before save it
func (root *TrieNode) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(root)
	}
	file.Close()
	return err
}

// Load unserialize root TrieNode from file and decode using gob
func Load(filePath string) (*TrieNode, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	root := new(TrieNode)
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(root)
	if err != nil {
		return nil, err
	}
	return root, err
}

// New return the Trie root
func New(rd *bufio.Reader, stopFilter bool) *TrieNode {
	if stopFilter {
		stopDict = make(map[string]bool)
		for _, w := range stopWords {
			stopDict[w] = true
		}
	}
	root := newTrieNode()
	id := 0
	for {
		phrase, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return nil
		}
		phrase = strings.TrimSuffix(phrase, "\n")
		root.insertPhrase(strings.ToLower(phrase), id)
		id++
	}
	return root
}

func detectLang(phrase string) string {
	c := []rune(phrase)
	if unicode.Is(unicode.Scripts["Han"], c[0]) {
		return zh
	}
	return eng
}
