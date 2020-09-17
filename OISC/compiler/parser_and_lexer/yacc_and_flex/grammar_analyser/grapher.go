package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type node struct {
	Name                 string
	Construction_parents []string
	Graph_node           *cgraph.Node
}

var createdTypes map[string]*node

func main() {
	createdTypes = make(map[string]*node)
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scan := bufio.NewScanner(file)
	scan.Scan()
	for a := scan.Text(); a != "%%"; a = scan.Text() {
		scan.Scan()
	}

	scan.Scan()
	reg := regexp.MustCompile("(?:^|\\n)(\\w+)")
	for scan.Scan() != false {
		line := scan.Text()
		extracted_line := reg.FindString(line)
		reg_inner := regexp.MustCompile("(\\w+|','|'{'|'}'|'\\('|'\\)'|'\\['|'\\]')")
		if len(extracted_line) != 0 {
			scan.Scan()
			fmt.Println(extracted_line)
			if _, ok := createdTypes[extracted_line]; !ok {
				child := new(node)
				child.Name = extracted_line
				createdTypes[extracted_line] = child
			}
			for i := reg_inner.FindAllString(scan.Text(), -1); len(i) != 0; i = reg_inner.FindAllString(scan.Text(), -1) {
				for _, token := range i {
					if val, ok := createdTypes[token]; !ok {
						child := new(node)
						child.Name = token
						child.Construction_parents = append(child.Construction_parents, extracted_line)
						createdTypes[token] = child
					} else {
						val.Construction_parents = append(val.Construction_parents, extracted_line)
					}
				}
				scan.Scan()
			}
		}
	}
	GenerateDiagram()
}

func GenerateDiagram() {
	g := graphviz.New()
	graph, err := g.Graph()
	for _, val := range createdTypes {
		node, _ := graph.CreateNode(val.Name)
		val.Graph_node = node
	}
	for _, val := range createdTypes {
		for _, x := range val.Construction_parents {
			graph.CreateEdge("", val.Graph_node, createdTypes[x].Graph_node)
		}
	}
	// 1. write encoded PNG data to buffer
	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}

	// 2. get as image.Image instance
	_, err2 := g.RenderImage(graph)
	if err2 != nil {
		log.Fatal(err)
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}

}
