package main

import (
	"bytes"
	"fmt"
	str "jjlehner/logicAnalyser/structure"
	"log"

	"github.com/goccy/go-graphviz"
)

func main() {
	inputLines := "!(I&&J||K)"
	root, _ := str.DescendExpression([]rune(inputLines), 0)
	a := root.Eval(&map[string]bool{
		"abc": true,
		"d":   true,
	})
	fmt.Println(a)

	graph, g := str.GenerateDiagram(root)
	// 1. write encoded PNG data to buffer
	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}

	// 2. get as image.Image instance
	_, err := g.RenderImage(graph)
	if err != nil {
		log.Fatal(err)
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}
}
