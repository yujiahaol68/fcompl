# fcompl

![build](https://img.shields.io/travis/yujiahaol68/fcompl/master.svg)

A fast phrases completion library using trie tree. Also provide path compression for lower memory usage

Common stopwords in English will also be removed

Support language: English, Chinese

## Usage

Insert Phrases

```go
s := "A ball in the ground\nThe bat in the sky\nThe ball hit his head\n"
rd := bufio.NewReader(strings.NewReader(s))
f := fcompl.New(rd, true)   // true if you want to get rid of stopwords
```

Find completion IDs

```go
f.Find("Ball") // [0, 2]
f.Find("bat") // [1]
```

Save and Load

```go
// path compress and save into a file
fcompl.Compress(f)
f.Save("./foo.bin") // default using gob to serialize struct into file

// Load
f, err := fcompl.Load("./foo.bin")
// ...
```

When complete Chinese phrases, you should find in its pinyin format like

```go
// eg.Want to complete '我是'
f.find("wo shi")
```

## Test coverage

```
PASS
coverage: 94.4% of statements
ok  	github.com/yujiahaol68/fcompl	1.049s	coverage: 94.4% of statements
```

Feel free to star if this small package helps. Issue is also welcome