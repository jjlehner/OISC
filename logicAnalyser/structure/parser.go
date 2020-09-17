package structure

import (
	"fmt"
	"reflect"
	"unicode"
)

type stack struct {
	size  int
	store []interface{}
}

func (b *stack) pop() (r interface{}) {
	b.size--
	r = b.store[b.size]
	b.store = b.store[:b.size]
	return
}

func (b *stack) push(entry interface{}) {
	b.size++
	b.store = append(b.store, entry)
	return
}
func (b stack) getSize() int {
	return b.size
}

var operatorNodeTokens = map[string]*operator{
	"&&": &and,
	"+":  &or,
	"||": &or,
	"⊕":  &xor,
}

type item struct {
	Subject interface{}
	Prev    *item
	Next    *item
}
type linkedList struct {
	start *item
	end   *item
}

func (list *linkedList) attach(i interface{}) {
	toAttach := item{
		i,
		list.end,
		nil,
	}
	if list.end != nil {
		list.end.Next = &toAttach
	}
	list.end = &toAttach
	if list.start == nil {
		list.start = &toAttach
	}
}
func (list *linkedList) removeNeighbours(i *item) {
	if i == list.end || i == list.start {
		panic("Can't remove neighbours for first or last item in list")
	}
	if i.Next == list.end {
		list.end = i
	}
	if i.Prev == list.start {
		list.start = i
	}
	if i.Prev.Prev != nil {
		i.Prev.Prev.Next = i
	}
	if i.Next.Next != nil {
		i.Next.Next.Prev = i
	}
	i.Prev = i.Prev.Prev
	i.Next = i.Next.Next
}
func (list *linkedList) removeAbove(i *item) {
	if i == list.end {
		panic("Can't remove above for last item in list")
	}
	if i.Next == list.end {
		list.end = i
	}
	if i.Next.Next != nil {
		i.Next.Next.Prev = i
	}
	i.Next = i.Next.Next
}

func createRootNode(list *linkedList) (root node) {
	for i := list.start; ; i = i.Next {
		if i == nil {
			break
		}
		not, ok := i.Subject.(*not)
		if ok {
			next, _ := i.Next.Subject.(node)
			not.Child = next
			list.removeAbove(i)
		}
		if i == list.end {
			break
		}
	}
	priority := [...]*operator{&and, &xor, &or}
	for _, operation := range priority {
		for i := list.start; ; i = i.Next {
			if i == nil {
				break
			}
			op, ok := i.Subject.(*op_node)
			if ok && operation == op.Op {
				if op.hasChildren() == true {
					continue
				}
				next, _ := i.Next.Subject.(node)
				prev, _ := i.Prev.Subject.(node)
				op.Children = append(op.Children, next, prev)
				list.removeNeighbours(i)
			}
			if i == list.end {
				break
			}
		}
	}
	val, _ := list.start.Subject.(node)
	list.testOutputList()
	return val
}
func DescendExpression(expr []rune, cursorStart int) (root node, endPosition int) {
	strLen := len(expr)
	longVariable := false
	var list linkedList
	for i := cursorStart; i < strLen; i++ {
		var c rune = expr[i]
		if !longVariable {
			switch c {
			case '(':
				var child node
				child, i = DescendExpression(expr, cursorStart+1)
				cursorStart = i + 1
				list.attach(child)
			case ')':
				return createRootNode(&list), i
			case '"':
				cursorStart = i + 1
				longVariable = true
			case ' ':
				continue
			case '!':
				notNode := node(&not{})
				list.attach(notNode)
				cursorStart = i + 1
			default:
				if op, present := operatorNodeTokens[string(expr[cursorStart:i+1])]; present {
					opNode := node(Op_node_constructor(op))
					list.attach(opNode)
				} else if unicode.IsLetter(c) {
					inputNode := node(Input_node_constructor(string(c)))
					list.attach(inputNode)
					cursorStart = i + 1
				}
			}
		} else {
			if c == '"' {
				longVariableNode := node(Input_node_constructor(string(expr[cursorStart:i])))
				list.attach(longVariableNode)
				cursorStart = i + 1
				longVariable = false
			}
		}
	}
	return createRootNode(&list), strLen
}

func (list *linkedList) testOutputList() {
	for i := list.start; ; i = i.Next {
		if i == nil {
			break
		}
		fmt.Println(reflect.TypeOf(i.Subject).String())
		if i == list.end {
			break
		}
	}
	fmt.Println()
}
