package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (e *Edge) printWall() string {
	color := "#FFF"
	if e.Included {
		color = "#000"
	}
	// color = "#000"

	border := ""

	if e.Direction == "top" {
		border = "border-top"
	} else if e.Direction == "bottom" {
		border = "border-bottom"
	} else if e.Direction == "left" {
		border = "border-left"
	} else if e.Direction == "right" {
		border = "border-right"
	}

	return border + `: 2px solid ` + color + `; `

}

var html string
var grid [][]Node

func (g *Graph) insert(n Node) {
	g.Nodes = append(g.Nodes, n)
}

func main() {
	g := Graph{}

	a := Node{
		Visited: false,
	}

	g.insert(a)

	a.Visited = false

	fmt.Println(a)

}

func mains() {

	html = `
	<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Maze</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>

    <style>
        table {
            border-collapse: collapse;
            border: 0px solid black;
        }

        table td {
            min-width: 20px;
            height: 20px;
        }
    </style>
</head>
<body>

    <table style="border-collapse: collapse;">`

	// Do the maze

	MST()

	// Print the maze
	for i := 0; i < 20; i++ {
		html += "<tr>"
		for j := 0; j < 20; j++ {
			walls := ""
			alt := ""
			for k := 0; k < len(grid[i][j].Edges); k++ {
				walls += grid[i][j].Edges[k].printWall()
				alt += grid[i][j].Edges[k].Label
			}

			html += `<td alt="` + alt + `" style="` + walls + `">` + strconv.Itoa(i) + "," + strconv.Itoa(j) + `</td>`

		}
		html += "</tr>"
	}

	html = html + `    </table>
</body>
</html>
`
	// fmt.Println(maze)

	fmt.Println("Hello World")

	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func Index(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

type Node struct {
	Visited bool
}

type Edge struct {
	Label     string
	Weight    int8
	Included  bool
	Direction string
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

func MakeGrid() [][]Node {
	grid = make([][]Node, 20)
	for i := range grid {
		grid[i] = make([]Node, 20)
	}

	// Connect the grid
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			grid[i][j] = Node{
				Visited: false,
			}

			// Connect to left
			if i > 0 {

				grid[i][j].Edges = append(grid[i][j].Edges, Edge{
					Label:     "E-L-" + strconv.Itoa(i) + "-" + strconv.Itoa(j),
					Weight:    int8(rand.Intn(100)) + 1,
					Node:      grid[i-1][j],
					Direction: "left",
					Included:  false,
				})
			}

			// Connect to right
			if i < 19 {
				grid[i][j].Edges = append(grid[i][j].Edges, Edge{
					Label:     "E-R-" + strconv.Itoa(i) + "-" + strconv.Itoa(j),
					Weight:    int8(rand.Intn(100)) + 1,
					Node:      grid[i+1][j],
					Direction: "right",
					Included:  false,
				})
			}

			// Connect to top
			if j > 0 {
				grid[i][j].Edges = append(grid[i][j].Edges, Edge{
					Label:     "E-T-" + strconv.Itoa(i) + "-" + strconv.Itoa(j),
					Weight:    int8(rand.Intn(100)) + 1,
					Node:      grid[i][j-1],
					Direction: "top",
					Included:  false,
				})
			}

			// Connect to bottom
			if j < 19 {
				grid[i][j].Edges = append(grid[i][j].Edges, Edge{
					Label:     "E-B-" + strconv.Itoa(i) + "-" + strconv.Itoa(j),
					Weight:    int8(rand.Intn(100)) + 1,
					Node:      grid[i][j+1],
					Direction: "bottom",
					Included:  false,
				})
			}

		}
	}

	return grid
}

var totalNodes = 19 * 19

func MST() []*Node {
	grid = MakeGrid()

	// Arbitrary starting node
	rand.Seed(time.Now().UnixNano())

	startX, startY := rand.Intn(19), rand.Intn(19)
	s1 := grid[startX][startY]
	s1.Visited = true

	fmt.Println("Starting at: ", startX, startY)

	// Visited nodes
	t := make([]*Node, 0)
	t = append(t, s1)

	for len(t) < totalNodes {

		var smallestEdge *Edge = nil

		for _, v := range t {

			for _, e := range v.Edges {
				if (*e).Node.Visited == false {
					if smallestEdge == nil || (smallestEdge.Included == false && smallestEdge.Weight > e.Weight) {
						smallestEdge = e
					}
				}
			}

		}

		smallestEdge.Included = true // Print me
		smallestEdge.Label = smallestEdge.Label + "-!"
		fmt.Println("Adding edge: ", smallestEdge.Label, " to the MST")

		smallestEdge.Node.Visited = true

		t = append(t, smallestEdge.Node)
	}

	return t

}
