package trie

import (
	"fmt"
	"sort"
	"unicode"
	"unicode/utf8"
)

type List struct {
	elements []any
	size     int
}

func NewList(capacity int) *List {
	if capacity > 0 {
		return &List{elements: make([]any, capacity), size: 0}
	} else {
		return &List{elements: nil, size: 0}
	}
}

func (list *List) Get(i int) any {
	if list.elements == nil || i >= list.size {
		return nil
	}
	return list.elements[i]
}

func (list *List) Add(e any) {
	list.ensureCapacity(list.size + 1)
	list.elements[list.size] = e
	list.size++
}

func (list *List) AddAll(c *List) {
	if c == nil || c.size == 0 {
		return
	}
	list.ensureCapacity(list.size + c.size)
	for i := 0; i < c.size; i++ {
		list.elements[list.size] = c.elements[i]
		list.size++
	}
}

func (list *List) Size() int {
	return list.size
}

func (list *List) IsEmpty() bool {
	return list.size == 0
}

func (list *List) ToArray() []any {
	if list.size == 0 {
		return []any{}
	}
	elements := make([]any, list.size)
	copy(elements, list.elements)
	return elements
}

func (list *List) ensureCapacity(minCapacity int) {
	if list.elements == nil || minCapacity >= len(list.elements) {
		list.grow(max(minCapacity, 10))
	}
}

func (list *List) grow(minCapacity int) {
	oldLength := len(list.elements)
	newLength := oldLength + max(minCapacity-oldLength, oldLength>>1)
	elements := make([]any, newLength)
	copy(elements, list.elements)
	list.elements = elements
}

type Node struct {
	next, prev *Node
	element    any
}

type Queue struct {
	first, last *Node
	size        int
}

func NewQueue() *Queue {
	return &Queue{nil, nil, 0}
}

func (queue *Queue) Add(e any) {
	node := &Node{nil, nil, e}
	last := queue.last
	queue.last = node
	if queue.first == nil {
		queue.first = node
	} else {
		last.next = node
	}
	queue.size++
}

func (queue *Queue) Poll() any {
	first := queue.first
	if first == nil {
		return nil
	}
	element := first.element
	next := first.next
	first.element = nil
	first.next = nil // GC
	queue.first = next
	if next == nil {
		queue.last = nil
	} else {
		next.prev = nil
	}
	queue.size--
	return element
}

func (queue *Queue) Peek() any {
	if queue.first == nil {
		return nil
	} else {
		return queue.first.element
	}
}

func (queue *Queue) Size() int {
	return queue.size
}

func (queue *Queue) IsEmpty() bool {
	return queue.size == 0
}

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

type State struct {
	depth    int
	success  map[rune]*State
	failure  *State
	keywords *List
}

func NewState(depth int) *State {
	return &State{
		depth:    depth,
		keywords: NewList(0),
	}
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
	s.keywords.Add(keyword)
}

func (s *State) AddKeywords(keywords []string) {
	if keywords == nil {
		return
	}
	for _, keyword := range keywords {
		s.keywords.Add(keyword)
	}
}

func (s *State) AddKeywordList(keywords *List) {
	if keywords == nil {
		return
	}
	s.keywords.AddAll(keywords)
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
	states := NewQueue()
	for _, state := range t.root.success {
		state.failure = t.root
		states.Add(state)
	}
	for !states.IsEmpty() {
		state := states.Poll().(*State)
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
			next.AddKeywordList(fn.keywords)
			states.Add(next)
		}
	}
	return t
}

func (t *Trie) FindAll(text string, ignoreCase bool) []*Emit {
	size := 0
	emits := make([]*Emit, 10)
	state := t.root
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		state = t.nextState(state, r, ignoreCase)
		for j := 0; j < state.keywords.Size(); j++ {
			kw := state.keywords.Get(j).(string)
			if size == len(emits) {
				old := emits
				emits = make([]*Emit, size+size>>1)
				copy(emits, old)
			}
			emits[size] = &Emit{i - strlen(kw) + 1, i + 1, kw}
			size++
		}
	}
	return emits[:size]
}

func (t *Trie) FindFirst(text string, ignoreCase bool) *Emit {
	state := t.root
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		state = t.nextState(state, r, ignoreCase)
		if state.keywords.Size() > 0 {
			kw := state.keywords.Get(0).(string)
			return &Emit{i - strlen(kw) + 1, i + 1, kw}
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
	if last.End < strlen(source) {
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

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func strlen(s string) int {
	return utf8.RuneCountInString(s)
}
