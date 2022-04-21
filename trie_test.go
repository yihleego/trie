package trie

import (
	"fmt"
	"testing"
)

func TestFindAll(t *testing.T) {
	text := "昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。"
	trie := NewTrie("雨疏", "风骤", "残酒", "卷帘人", "知否")
	emits := trie.FindAll(text, false)
	fmt.Println(emits)
	EqualEmit(t, emits[0], 2, 4, "雨疏")
	EqualEmit(t, emits[1], 4, 6, "风骤")
	EqualEmit(t, emits[2], 11, 13, "残酒")
	EqualEmit(t, emits[3], 16, 19, "卷帘人")
	EqualEmit(t, emits[4], 27, 29, "知否")
	EqualEmit(t, emits[5], 30, 32, "知否")
}

func TestFindFirst(t *testing.T) {
	text := "昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。"
	trie := NewTrie("雨疏", "风骤", "残酒", "卷帘人", "知否")
	emit := trie.FindFirst(text, false)
	fmt.Println(emit)
	EqualEmit(t, emit, 2, 4, "雨疏")
}

func TestFindAllIgnoreCase(t *testing.T) {
	text := "Poetry is what gets lost in translation."
	trie := NewTrie("poetry", "TRANSLATION")
	emits := trie.FindAll(text, true)
	fmt.Println(emits)
	EqualEmit(t, emits[0], 0, 6, "poetry")
	EqualEmit(t, emits[1], 28, 39, "TRANSLATION")
}

func TestFindFirstIgnoreCase(t *testing.T) {
	text := "Poetry is what gets lost in translation."
	trie := NewTrie("poetry", "TRANSLATION")
	emit := trie.FindFirst(text, true)
	fmt.Println(emit)
	EqualEmit(t, emit, 0, 6, "poetry")
}

func TestTokenize(t *testing.T) {
	text := "常记溪亭日暮，沉醉不知归路。兴尽晚回舟，误入藕花深处。争渡，争渡，惊起一滩鸥鹭。"
	trie := NewTrie("溪亭", "归路", "藕花", "争渡")
	emits := trie.FindAll(text, false)
	tokens := Tokenize(emits, text)
	fmt.Println(len(tokens), tokens)
	EqualToken(t, tokens[0], -1, -1, "常记")
	EqualToken(t, tokens[1], 2, 4, "溪亭")
	EqualToken(t, tokens[2], -1, -1, "日暮，沉醉不知")
	EqualToken(t, tokens[3], 11, 13, "归路")
	EqualToken(t, tokens[4], -1, -1, "。兴尽晚回舟，误入")
	EqualToken(t, tokens[5], 22, 24, "藕花")
	EqualToken(t, tokens[6], -1, -1, "深处。")
	EqualToken(t, tokens[7], 27, 29, "争渡")
	EqualToken(t, tokens[8], -1, -1, "，")
	EqualToken(t, tokens[9], 30, 32, "争渡")
	EqualToken(t, tokens[10], -1, -1, "，惊起一滩鸥鹭。")
}

func TestReplace(t *testing.T) {
	text := "我正在参加砍价，砍到0元就可以免费拿啦。亲~帮我砍一刀呗，咱们一起免费领好货。"
	trie := NewTrie("0元", "砍一刀", "免费拿", "免费领")
	emits := trie.FindAll(text, false)
	r1 := Replace(emits, text, "*")
	r2 := Replace(emits, text, "@#$%^&*")
	fmt.Println(r1)
	fmt.Println(r2)
	EqualString(t, "我正在参加砍价，砍到**就可以***啦。亲~帮我***呗，咱们一起***好货。", r1)
	EqualString(t, "我正在参加砍价，砍到%^就可以#$%啦。亲~帮我%^&呗，咱们一起&*@好货。", r2)
}

func TestOverlaps(t *testing.T) {
	text := "12345"
	trie := NewTrie("1", "2", "12", "23", "34", "45", "123")
	emits := trie.FindAll(text, false)
	removed := RemoveOverlaps(emits)
	fmt.Println(emits)
	fmt.Println(removed)
	EqualEmit(t, removed[0], 0, 3, "123")
	EqualEmit(t, removed[1], 3, 5, "45")

}

func TestContains(t *testing.T) {
	text := "12345"
	trie := NewTrie("1", "2", "12", "23", "34", "45", "123")
	emits := trie.FindAll(text, false)
	removed := RemoveContains(emits)
	fmt.Println(emits)
	fmt.Println(removed)
	EqualEmit(t, removed[0], 0, 3, "123")
	EqualEmit(t, removed[1], 2, 4, "34")
	EqualEmit(t, removed[2], 3, 5, "45")
}

func TestLoad(t *testing.T) {
	text := "Hello, World!"
	trie := NewTrie()
	trie.Load("hello", "world")
	emits := trie.FindAll(text, true)
	fmt.Println(emits)
	EqualEmit(t, emits[0], 0, 5, "hello")
	EqualEmit(t, emits[1], 7, 12, "world")
}

func TestList(t *testing.T) {
	list := NewList(0)
	for i := 0; i < 1000; i++ {
		list.Add(i)
	}
	size := list.Size()
	for i := 0; i < size; i++ {
		fmt.Print(list.Get(i), ",")
	}
	fmt.Println()
	EqualInt(t, 1000, size)
}

func TestQueue(t *testing.T) {
	queue := NewQueue()
	for i := 0; i < 1000; i++ {
		queue.Add(i)
	}
	size := queue.Size()
	for !queue.IsEmpty() {
		fmt.Print(queue.Poll(), ",")
	}
	fmt.Println()
	EqualInt(t, 1000, size)
}

func EqualInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Error(expected, actual)
	}
}

func EqualString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Error(expected, actual)
	}
}

func EqualEmit(t *testing.T, emit *Emit, begin int, end int, kw string) {
	if emit.Begin != begin || emit.End != end || emit.Keyword != kw {
		t.Error(emit)
	}
}

func EqualToken(t *testing.T, token *Token, begin int, end int, kw string) {
	if token.Fragment != kw {
		t.Error(token)
	}
	if token.IsMatch() {
		EqualEmit(t, token.Emit, begin, end, kw)
	}
}
