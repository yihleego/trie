package trie

import (
	"testing"
	"unicode/utf8"
)

func TestFindAll(t *testing.T) {
	text := "昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。"
	trie := NewTrie("雨疏", "风骤", "残酒", "卷帘人", "知否")
	emits := trie.FindAll(text, false)
	t.Log(emits)
	EqualEmit(t, emits[0], 2, 4, "雨疏")
	EqualEmit(t, emits[1], 4, 6, "风骤")
	EqualEmit(t, emits[2], 11, 13, "残酒")
	EqualEmit(t, emits[3], 16, 19, "卷帘人")
	EqualEmit(t, emits[4], 27, 29, "知否")
	EqualEmit(t, emits[5], 30, 32, "知否")
	EqualInt(t, 6, len(emits))
}

func TestFindFirst(t *testing.T) {
	text := "昨夜雨疏风骤，浓睡不消残酒。试问卷帘人，却道海棠依旧。知否，知否？应是绿肥红瘦。"
	trie := NewTrie("雨疏", "风骤", "残酒", "卷帘人", "知否")
	emit := trie.FindFirst(text, false)
	t.Log(emit)
	EqualEmit(t, emit, 2, 4, "雨疏")
}

func TestFindAllIgnoreCase(t *testing.T) {
	text := "Poetry is what gets lost in translation."
	trie := NewTrie("poetry", "TRANSLATION")
	emits := trie.FindAll(text, true)
	t.Log(emits)
	EqualEmit(t, emits[0], 0, 6, "poetry")
	EqualEmit(t, emits[1], 28, 39, "TRANSLATION")
	EqualInt(t, 2, len(emits))
}

func TestFindFirstIgnoreCase(t *testing.T) {
	text := "Poetry is what gets lost in translation."
	trie := NewTrie("poetry", "TRANSLATION")
	emit := trie.FindFirst(text, true)
	t.Log(emit)
	EqualEmit(t, emit, 0, 6, "poetry")
}

func TestIgnoreCase(t *testing.T) {
	text := "TurninG OnCe AgAiN BÖRKÜ"
	trie := NewTrie("turning", "once", "again", "börkü")
	emits := trie.FindAll(text, true)
	t.Log(emits)
	EqualEmit(t, emits[0], 0, 7, "turning")
	EqualEmit(t, emits[1], 8, 12, "once")
	EqualEmit(t, emits[2], 13, 18, "again")
	EqualEmit(t, emits[3], 19, 24, "börkü")
	EqualInt(t, 4, len(emits))
}

func TestTokenize(t *testing.T) {
	text := "常记溪亭日暮，沉醉不知归路。兴尽晚回舟，误入藕花深处。争渡，争渡，惊起一滩鸥鹭。"
	trie := NewTrie("溪亭", "归路", "藕花", "争渡")
	emits := trie.FindAll(text, false)
	tokens := Tokenize(emits, text)
	t.Log(len(emits), emits)
	t.Log(len(tokens), tokens)
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
	EqualInt(t, 5, len(emits))
	EqualInt(t, 11, len(tokens))
}

func TestReplace(t *testing.T) {
	text := "我正在参加砍价，砍到0元就可以免费拿啦。亲~帮我砍一刀呗，咱们一起免费领好货。"
	trie := NewTrie("0元", "砍一刀", "免费拿", "免费领")
	emits := trie.FindAll(text, false)
	r1 := Replace(emits, text, "*")
	r2 := Replace(emits, text, "@#$%^&*")
	t.Log(emits)
	t.Log(r1)
	t.Log(r2)
	EqualString(t, "我正在参加砍价，砍到**就可以***啦。亲~帮我***呗，咱们一起***好货。", r1)
	EqualString(t, "我正在参加砍价，砍到%^就可以#$%啦。亲~帮我%^&呗，咱们一起&*@好货。", r2)
	EqualInt(t, 4, len(emits))
}

func TestOverlaps(t *testing.T) {
	text := "a123,456b"
	trie := NewTrie("123", "12", "23", "45", "56")
	emits := trie.FindAll(text, false)
	t.Log(emits)
	removed := RemoveOverlaps(emits)
	t.Log(emits)
	t.Log(removed)
	EqualEmit(t, removed[0], 1, 4, "123")
	EqualEmit(t, removed[1], 5, 7, "45")
	EqualInt(t, 5, len(emits))
	EqualInt(t, 2, len(removed))
}

func TestContains(t *testing.T) {
	text := "a123,456b"
	trie := NewTrie("123", "12", "23", "45", "56")
	emits := trie.FindAll(text, false)
	t.Log(emits)
	removed := RemoveContains(emits)
	t.Log(emits)
	t.Log(removed)
	EqualEmit(t, removed[0], 1, 4, "123")
	EqualEmit(t, removed[1], 5, 7, "45")
	EqualEmit(t, removed[2], 6, 8, "56")
	EqualInt(t, 5, len(emits))
	EqualInt(t, 3, len(removed))
}

func TestDuplicate(t *testing.T) {
	text := "123456"
	trie := NewTrie("123", "123", "456", "456")
	emits := trie.FindAll(text, false)
	t.Log(emits)
	EqualEmit(t, emits[0], 0, 3, "123")
	EqualEmit(t, emits[1], 3, 6, "456")
	EqualInt(t, 2, len(emits))
}

func TestAddKeywords(t *testing.T) {
	text := "ushers"
	trie1 := NewTrie("he", "she", "his", "hers")
	trie2 := NewTrie().AddKeywords("he", "she", "his", "hers")
	trie3 := NewTrie().AddKeywords("he").AddKeywords("she").AddKeywords("his").AddKeywords("hers")
	emits1 := trie1.FindAll(text, false)
	emits2 := trie2.FindAll(text, false)
	emits3 := trie3.FindAll(text, false)
	t.Log(emits1)
	t.Log(emits2)
	t.Log(emits3)
	EqualEmits(t, emits1, emits2)
	EqualEmits(t, emits1, emits3)
	EqualEmits(t, emits2, emits3)
}

func TestEmoji(t *testing.T) {
	t.Log("utf8.RuneCountInString(\"🐼\") >>", utf8.RuneCountInString("🐼"))
	t.Log("len(\"🐼\") >>", len("🐼"))
	EqualInt(t, 1, utf8.RuneCountInString("🐼"))
	EqualInt(t, 4, len("🐼"))
	text := "I love 🐼 very much."
	trie := NewTrie("🐼", "🐻")
	emits := trie.FindAll(text, false)
	t.Log(emits)
	EqualEmit(t, emits[0], 7, 8, "🐼")
	EqualInt(t, 1, len(emits))
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

func EqualEmits(t *testing.T, emits1 []*Emit, emits2 []*Emit) {
	if len(emits1) != len(emits2) {
		t.Error(emits1, emits2)
		return
	}
	for i := 0; i < len(emits1); i++ {
		emit1, emit2 := emits1[i], emits2[i]
		if !emit1.Equals(emit2) {
			t.Error(emits1, emits2)
			return
		}
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
