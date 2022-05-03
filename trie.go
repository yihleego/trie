package trie

import (
	"container/list"
	"fmt"
	"sort"
	"unicode"
	"unicode/utf8"
)

type Emit struct {
	Begin   int
	End     int
	Keyword string
}

func (e *Emit) Length() int {
	return e.End - e.Begin
}

func (e *Emit) Equals(o *Emit) bool {
	return e.Begin == o.Begin && e.End == o.End && e.Keyword == o.Keyword
}

func (e *Emit) Overlaps(o *Emit) bool {
	return e.Begin < o.End && e.End > o.Begin
}

func (e *Emit) Contains(o *Emit) bool {
	return e.Begin <= o.Begin && e.End >= o.End
}

func (e *Emit) String() string {
	return fmt.Sprintf("%d:%d=%s", e.Begin, e.End, e.Keyword)
}

type Token struct {
	Fragment string
	Emit     *Emit
}

func (t *Token) IsMatch() bool {
	return t.Emit != nil
}

func (t *Token) String() string {
	if t.Emit == nil {
		return t.Fragment
	} else {
		return fmt.Sprintf("%s(%v)", t.Fragment, t.Emit)
	}
}

type Keyword struct {
	Value  string
	Length int
}

type State struct {
	depth    int
	success  map[rune]*State
	failure  *State
	keywords []*Keyword
}

func NewState(depth int) *State {
	return &State{depth: depth}
}

func (s *State) NextState(c rune, ignoreCase bool) *State {
	next := s.GetState(c, ignoreCase)
	if next != nil {
		return next
	} else if s.depth == 0 {
		return s
	} else {
		return nil
	}
}

func (s *State) GetState(c rune, ignoreCase bool) *State {
	if s.success == nil {
		return nil
	}
	state, exists := s.success[c]
	if exists {
		return state
	}
	if ignoreCase {
		cc := c
		if unicode.IsLower(c) {
			cc = unicode.ToUpper(c)
		} else if unicode.IsUpper(c) {
			cc = unicode.ToLower(c)
		}
		if c != cc {
			next := s.success[cc]
			return next
		}
	}
	return nil
}

func (s *State) AddState(str string) *State {
	state := s
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		state = state.addState(runes[i])
	}
	return state
}

func (s *State) addState(c rune) *State {
	if s.success == nil {
		s.success = make(map[rune]*State)
	}
	state, exists := s.success[c]
	if exists {
		return state
	}
	ns := NewState(s.depth + 1)
	s.success[c] = ns
	return ns
}

func (s *State) AddKeyword(keyword string) {
	s.keywords = append(s.keywords, &Keyword{keyword, utf8.RuneCountInString(keyword)})
}

func (s *State) AddKeywords(keywords []*Keyword) {
	if len(keywords) > 0 {
		s.keywords = append(s.keywords, keywords...)
	}
}

type Trie struct {
	root *State
}

func NewTrie(keywords ...string) *Trie {
	t := Trie{root: NewState(0)}
	t.Load(keywords...)
	return &t
}

func (t *Trie) Load(keywords ...string) *Trie {
	for _, keyword := range keywords {
		if len(keyword) > 0 {
			t.root.AddState(keyword).AddKeyword(keyword)
		}
	}
	states := list.New()
	for _, state := range t.root.success {
		state.failure = t.root
		states.PushBack(state)
	}
	for states.Len() > 0 {
		state := states.Remove(states.Front()).(*State)
		if state.success == nil {
			continue
		}
		for c, next := range state.success {
			f := state.failure
			fn := f.NextState(c, false)
			for fn == nil {
				f = f.failure
				fn = f.NextState(c, false)
			}
			next.failure = fn
			next.AddKeywords(fn.keywords)
			states.PushBack(next)
		}
	}
	return t
}

func (t *Trie) FindAll(text string, ignoreCase bool) []*Emit {
	emits := make([]*Emit, 0, 10)
	state := t.root
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		state = t.nextState(state, r, ignoreCase)
		for j := 0; j < len(state.keywords); j++ {
			kw := state.keywords[j]
			emits = append(emits, &Emit{i + 1 - kw.Length, i + 1, kw.Value})
		}
	}
	return emits
}

func (t *Trie) FindFirst(text string, ignoreCase bool) *Emit {
	state := t.root
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		state = t.nextState(state, r, ignoreCase)
		if len(state.keywords) > 0 {
			kw := state.keywords[0]
			return &Emit{i + 1 - kw.Length, i + 1, kw.Value}
		}
	}
	return nil
}

func (t *Trie) nextState(state *State, c rune, ignoreCase bool) *State {
	next := state.NextState(c, ignoreCase)
	for next == nil {
		state = state.failure
		next = state.NextState(c, ignoreCase)
	}
	return next
}

func Tokenize(emits []*Emit, source string) []*Token {
	emits = RemoveContains(emits)
	el := len(emits)
	if el == 0 {
		return []*Token{{source, nil}}
	}
	count := 0
	index := 0
	runes := []rune(source)
	tokens := make([]*Token, el*2+1)
	for i := 0; i < el; i++ {
		emit := emits[i]
		if index < emit.Begin {
			tokens[count] = &Token{string(runes[index:emit.Begin]), nil}
			count++
		}
		tokens[count] = &Token{string(runes[emit.Begin:emit.End]), emit}
		count++
		index = emit.End
	}
	last := emits[el-1]
	if last.End < utf8.RuneCountInString(source) {
		tokens[count] = &Token{string(runes[last.End:]), nil}
		count++
	}
	return tokens[:count]
}

func Replace(emits []*Emit, source string, replacement string) string {
	emits = RemoveContains(emits)
	el := len(emits)
	if el == 0 {
		return source
	}
	index := 0
	runes := []rune(source)
	masks := []rune(replacement)
	ml := len(masks)
	for i := 0; i < el; i++ {
		emit := emits[i]
		if index < emit.Begin {
			index = emit.Begin
		}
		for j := emit.Begin; j < emit.End; j++ {
			runes[j] = masks[j%ml]
		}
		index = emit.End
	}
	return string(runes)
}

func RemoveOverlaps(emits []*Emit) []*Emit {
	return removeEmits(emits, func(a, b *Emit) bool {
		return a.Overlaps(b)
	})
}

func RemoveContains(emits []*Emit) []*Emit {
	return removeEmits(emits, func(a, b *Emit) bool {
		return a.Contains(b)
	})
}

func removeEmits(emits []*Emit, predicate func(a, b *Emit) bool) []*Emit {
	el := len(emits)
	if el <= 1 {
		return emits
	}
	sortEmits(emits)
	index := 1
	sorted := make([]*Emit, el)
	sorted[0] = emits[0]
	emit := emits[0]
	for i := 1; i < el; i++ {
		next := emits[i]
		if !predicate(emit, next) {
			sorted[index] = next
			index++
			emit = next
		}
	}
	return sorted[:index]
}

func sortEmits(emits []*Emit) {
	sort.Slice(emits, func(i, j int) bool {
		a := emits[i]
		b := emits[j]
		if a.Begin != b.Begin {
			return a.Begin < b.Begin
		} else {
			return a.End > b.End
		}
	})
}
