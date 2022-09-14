# Trie

[![GoDoc](https://godoc.org/github.com/yihleego/trie?status.svg)](https://godoc.org/github.com/yihleego/trie)
[![Go Report Card](https://goreportcard.com/badge/github.com/yihleego/trie)](https://goreportcard.com/report/github.com/yihleego/trie)

An Aho-Corasick algorithm based string-searching utility for Go. It supports tokenization, ignoring case, replacing text. So you can use it to find keywords in an article, filter sensitive words, etc.

Implemented in Java：[Trie4j](https://github.com/yihleego/trie4j)

基于AC自动机（Aho-Corasick algorithm）实现的关键词、敏感词、非法词、停用词等匹配替换工具，支持结果分词，忽略大小写，替换文本等功能。适用于在文章中查找关键词；过滤聊天、评论、留言中的敏感词等。

Java实现版本：[Trie4j](https://github.com/yihleego/trie4j)

## Introduction

判断一个字符串是否包含另一个字符串，我们通常使用`strings.Index()`或`strings.Contains()`进行判断，其底层实现基于RK、KMP、BM和Sunday等算法。如果要判断一个字符串是否包含多个字符串，比如在一篇文章找几个敏感词，继续使用上述的字符串搜索算法显然是不合适，这种场景就需要用到多模式匹配算法。

[Aho–Corasick](http://cr.yp.to/bib/1975/aho.pdf) 算法是由贝尔实验室的 Alfred V. Aho 和 Margaret J. Corasick 在 1975 年发明的一种字符串搜索算法。它是一种字典匹配算法，可在输入文本中定位有限字符串集（字典）的元素。它同时匹配所有字符串。该算法的复杂性与字符串长度加上搜索文本的长度加上输出匹配的数量成线性关系。

该算法主要依靠构造一个有限状态机来实现，然后通过失配指针在查找字符串失败时进行回退，转向某前缀的其他分支，免于重复匹配前缀，提高算法效率。

## Usage

### 匹配所有关键词

```go
t := trie.New("雨疏", "风骤", "残酒", "卷帘人", "知否")
emits := t.FindAll("昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。", false)
```

```text
[2:4=雨疏, 4:6=风骤, 11:13=残酒, 16:19=卷帘人, 27:29=知否, 30:32=知否]
```

### 匹配首个关键词

```go
t := trie.New("雨疏", "风骤", "残酒", "卷帘人", "知否")
emit := t.FindFirst("昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。", false)
```

```text
2:4=雨疏
```

### 匹配所有关键词 忽略大小写

```go
t := trie.New("poetry", "TRANSLATION")
emits := t.FindAll("Poetry is what gets lost in translation.", true)
```

```text
[0:6=poetry, 28:39=TRANSLATION]
```

### 匹配首个关键词 忽略大小写

```go
t := trie.New("poetry", "TRANSLATION")
emit := t.FindFirst("Poetry is what gets lost in translation.", true)
```

```text
0:6=poetry
```

### 切分词

```go
s := "常记溪亭日暮，沉醉不知归路。兴尽晚回舟，误入藕花深处。争渡，争渡，惊起一滩鸥鹭。"
t := trie.New("溪亭", "归路", "藕花", "争渡")
emits := t.FindAll(s, false)
tokens := trie.Tokenize(emits, s)
```

```text
["常记", "溪亭(2:4=溪亭)", "日暮，沉醉不知", "归路(11:13=归路)", "。兴尽晚回舟，误入", "藕花(22:24=藕花)", "深处。", "争渡(27:29=争渡)", "，", "争渡(30:32=争渡)", "，惊起一滩鸥鹭。"]
```

### 替换关键词

```go
s := "我正在参加砍价，砍到0元就可以免费拿啦。亲~帮我砍一刀呗，咱们一起免费领好货。"
t := trie.New("0元", "砍一刀", "免费拿", "免费领")
emits := t.FindAll(s, false)
r1 := trie.Replace(emits, s, "*")
r2 := trie.Replace(emits, s, "@#$%^&*")
```

```text
我正在参加砍价，砍到**就可以***啦。亲~帮我***呗，咱们一起***好货。
我正在参加砍价，砍到%^就可以#$%啦。亲~帮我%^&呗，咱们一起&*@好货。
```

## Contact

- [提交问题](https://github.com/yihleego/trie/issues)

## License

This project is under the MIT license. See the [LICENSE](LICENSE) file for details.
