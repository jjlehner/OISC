package main

type operator func(a, b bool, inputs ...bool) bool

var nand operator = func(a, b bool, inputs ...bool) bool {
	return !and(a, b, inputs...)
}
var xor operator = func(a, b bool, inputs ...bool) bool {
	inputs = append(inputs, a, b)
	i := 0
	for _, eval := range inputs {
		if eval {
			i += 1
		}
		if i > 1 {
			return false
		}
	}
	if i == 0 {
		return false
	}
	return true
}
var or operator = func(a, b bool, inputs ...bool) bool {
	inputs = append(inputs, a, b)
	for _, eval := range inputs {
		if eval {
			return true
		}
	}
	return false
}
var and operator = func(a, b bool, inputs ...bool) bool {
	inputs = append(inputs, a, b)
	for _, eval := range inputs {
		if !eval {
			return false
		}
	}
	return true
}

type node interface {
	canonicalNAND()
	eval(input_map *map[string]bool) bool
}

type op_node struct {
	Children []node
	Op       *operator
}

func (n *op_node) hasChildren() bool {
	if len(n.Children) > 0 {
		return true
	}
	return false
}
func (n *op_node) canonicalNAND() {

}

func (n *op_node) eval(input_map *map[string]bool) bool {
	var inputs []bool
	for _, child := range n.Children {
		result := (child).eval(input_map)
		inputs = append(inputs, result)
	}
	if len(inputs) > 2 {
		return (*n.Op)(inputs[0], inputs[1], inputs[2:]...)
	}
	if len(inputs) == 2 {
		return (*n.Op)(inputs[0], inputs[1])
	}
	if len(inputs) == 1 {
		return (*n.Op)(inputs[0], inputs[0])
	}
	panic("Evaluation error")
}

func Op_node_constructor(opr *operator, nodes ...node) *op_node {
	nodes = append(nodes, nodes...)
	root := op_node{
		Children: nodes,
		Op:       opr,
	}
	return &root
}

type input_node struct {
	Name string
}

func Input_node_constructor(name string) *input_node {
	root := input_node{
		Name: name,
	}
	return &root
}
func (n *input_node) canonicalNAND() {

}

func (n *input_node) eval(input_map *map[string]bool) bool {
	return (*input_map)[n.Name]
}

type not struct {
	Child node
}

func (n *not) canonicalNAND() {

}

func (n *not) eval(input_map *map[string]bool) bool {
	return !(n.Child.eval(input_map))
}
