package structure

import (
	"log"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

var counter int = 0

type operator func(a, b bool, inputs ...bool) bool

var nand operator = func(a, b bool, inputs ...bool) bool {
	return !and(a, b, inputs...)
}
var xor operator = func(a, b bool, inputs ...bool) bool {
	inputs = append(inputs, a, b)
	i := 0
	for _, Eval := range inputs {
		if Eval {
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
	for _, Eval := range inputs {
		if Eval {
			return true
		}
	}
	return false
}
var and operator = func(a, b bool, inputs ...bool) bool {
	inputs = append(inputs, a, b)
	for _, Eval := range inputs {
		if !Eval {
			return false
		}
	}
	return true
}

type node interface {
	CanonicalNAND()
	Eval(input_map *map[string]bool) bool
	generateDiagram(graph *cgraph.Graph) *cgraph.Node
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
func (n *op_node) CanonicalNAND() {

}

func (n *op_node) Eval(input_map *map[string]bool) bool {
	var inputs []bool
	for _, child := range n.Children {
		result := (child).Eval(input_map)
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

func (n *op_node) generateDiagram(graph *cgraph.Graph) *cgraph.Node {
	node, err := graph.CreateNode(getFuncName(n.Op) + " - " + strconv.FormatInt(int64(counter), 10))
	counter++
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range n.Children {
		graph.CreateEdge("", node, val.generateDiagram(graph))
	}
	return node
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
func (n *input_node) CanonicalNAND() {

}

func (n *input_node) Eval(input_map *map[string]bool) bool {
	return (*input_map)[n.Name]
}

func (n *input_node) generateDiagram(graph *cgraph.Graph) *cgraph.Node {
	node, err := graph.CreateNode(n.Name)
	if err != nil {
		log.Fatal(err)
	}
	return node
}

type not struct {
	Child node
}

func (n *not) CanonicalNAND() {

}

func (n *not) Eval(input_map *map[string]bool) bool {
	return !(n.Child.Eval(input_map))
}

func (n *not) generateDiagram(graph *cgraph.Graph) *cgraph.Node {
	node, err := graph.CreateNode("not - " + strconv.FormatInt(int64(counter), 10))
	counter++
	if err != nil {
		log.Fatal(err)
	}
	if n.Child != nil {
		graph.CreateEdge("", node, n.Child.generateDiagram(graph))
	}
	return node
}

func GenerateDiagram(n node) (*cgraph.Graph, *graphviz.Graphviz) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	n.generateDiagram(graph)
	return graph, g
}

func getFuncName(op *operator) string {
	switch op {
	case &and:
		return "AND"
	case &or:
		return "OR"
	case &xor:
		return "XOR"
	case &nand:
		return "NAND"
	default:
		return "Operator Type Unkown"
	}
}
