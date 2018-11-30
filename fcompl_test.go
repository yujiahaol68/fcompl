package fcompl

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	pinyin "github.com/mozillazg/go-pinyin"
	"github.com/stretchr/testify/assert"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Test_completion_fromfile(t *testing.T) {
	f, err := os.OpenFile("phrases.txt", os.O_RDONLY, os.ModePerm)
	check(err)

	defer f.Close()

	rd := bufio.NewReader(f)
	root := New(rd, true)

	table := []struct {
		input  string
		expect []int
	}{
		{"the batman", []int{0, 7}},
		{"american", []int{2, 3}},
		{"wonder", []int{4, 5}},
		{"A robin", []int{6}},
		{"wonder anything", []int{}},
		{"我", []int{8, 9}},
	}
	root.print()

	t1 := time.Now()
	for _, ca := range table {
		if detectLang(ca.input) == zh {
			pin := pinyin.LazyConvert(ca.input, &pinyinCfg)
			contentsEqual(root.Find(strings.Join(pin, " ")), ca.expect)
			continue
		}
		contentsEqual(root.Find(ca.input), ca.expect)
	}
	log.Println(time.Since(t1))
}

func Test_completion_fromString(t *testing.T) {
	s := "A ball in the ground\nThe bat in the sky\nThe ball hit his head\n"
	rd := bufio.NewReader(strings.NewReader(s))
	f := New(rd, true)

	contentsEqual(f.Find("Ball"), []int{0, 2})
	contentsEqual(f.Find("bat"), []int{1})
}

func Test_compress(t *testing.T) {
	f, err := os.OpenFile("phrases.txt", os.O_RDONLY, os.ModePerm)
	check(err)

	defer f.Close()

	rd := bufio.NewReader(f)
	root := New(rd, true)
	log.Println("Before compress:")
	root.print()
	log.Println("After compress:")
	Compress(root)
	root.print()
}

func Test_serialize(t *testing.T) {
	f, err := os.OpenFile("phrases.txt", os.O_RDONLY, os.ModePerm)
	check(err)

	defer f.Close()

	rd := bufio.NewReader(f)
	root := New(rd, true)
	Compress(root)
	err = root.Save("./fcompl.bin")
	defer os.Remove("./fcompl.bin")
	check(err)
	loadedRoot, err := Load("./fcompl.bin")
	check(err)
	assert.Equal(t, *root, *loadedRoot)
}

func contentsEqual(input []int, expect []int) {
	if len(input) != len(expect) {
		log.Fatalf("expect %v, but got %v", expect, input)
	}

	for i, v := range input {
		if v != expect[i] {
			log.Fatalf("expect %v, but got %v", expect, input)
		}
	}
}
