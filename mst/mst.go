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

type Connection struct {
	Vertex    *Vertex
	Edge      *Edge
	Direction string
}

type Edge struct {
	Label      string
	Weight     int
	IsIncluded bool
	Printed    bool
}

type Vertex struct {
	Label       string
	Connections []*Connection
	IsVisited   bool
}

func (v *Vertex) AddConnection(e *Edge, v2 *Vertex, d string) {
	if debug {
		fmt.Println("Adding connection from", v.Label, "to", v2.Label)
	}
	v.Connections = append(v.Connections, &Connection{
		Vertex:    v2,
		Edge:      e,
		Direction: d,
	})
}

var size int
var html string
var debug bool

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func main() {
	debug = false
	rand.NewSource(time.Now().UnixNano())
	size = 20 // 2^n

	graph := make([][]*Vertex, size)

	// Fill graph
	fillGraph(&graph, size)

	// Run MST
	prims(&graph)

	// Create the HTML for the MST
	html = drawGraph(&graph)

	// Serve
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

	if debug {
		fmt.Println(html)
	}

}

func Index(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func fillGraph(graph *[][]*Vertex, size int) {

	for i := 0; i < size; i++ {

		(*graph)[i] = make([]*Vertex, size)

		for j := 0; j < size; j++ {

			(*graph)[i][j] = &Vertex{
				Label:       "V" + strconv.Itoa(i) + strconv.Itoa(j),
				Connections: make([]*Connection, 0),
				IsVisited:   false,
			}
			if debug {
				fmt.Println("Adding vertex ", (*graph)[i][j].Label)
			}
		}

	}

	// Connect in grid mode
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {

			// Connect to right
			if i < size-1 {
				w := randInt(0, 100)
				eR := &Edge{
					Label:      (*graph)[i][j].Label + "-" + (*graph)[i+1][j].Label + ":" + strconv.Itoa(w),
					Weight:     w,
					IsIncluded: false,
				}

				if debug {
					fmt.Println("Adding edge from", (*graph)[i][j].Label, "to", (*graph)[i+1][j].Label, "with weight", w)
				}

				(*graph)[i][j].AddConnection(eR, (*graph)[i+1][j], "bottom")
				(*graph)[i+1][j].AddConnection(eR, (*graph)[i][j], "top")
			}

			// Connect to bottom
			if j < size-1 {
				w := randInt(0, 100)
				eB := &Edge{
					Label:      (*graph)[i][j].Label + "-" + (*graph)[i][j+1].Label + ":" + strconv.Itoa(w),
					Weight:     w,
					IsIncluded: false,
				}

				if debug {
					fmt.Println("Adding edge from", (*graph)[i][j].Label, "to", (*graph)[i][j+1].Label, "with weight", w)
				}

				(*graph)[i][j].AddConnection(eB, (*graph)[i][j+1], "right")
				(*graph)[i][j+1].AddConnection(eB, (*graph)[i][j], "left")

			}

		}
	}
}

func prims(graph *[][]*Vertex) {

	// Pick a random vertex
	x, y := randInt(0, len(*graph)), randInt(0, len(*graph))

	if debug {
		fmt.Println("Starting with vertex", (*graph)[x][y].Label)
	}

	// Set it as visited
	(*graph)[x][y].IsVisited = true

	// Spanning tree edges
	var graphMST []*Vertex
	graphMST = make([]*Vertex, 0)
	graphMST = append(graphMST, (*graph)[x][y])

	// Loop until all vertices are visited
	for hasUnvisitedVertices(graph) {

		var minEdge *Edge
		var minVertex *Vertex

		for _, v := range graphMST {

			// Visit all my unvisited neighbors
			for _, c := range v.Connections {
				if debug {
					fmt.Println(c.Edge, c.Vertex)
				}
				if !c.Vertex.IsVisited && !c.Edge.IsIncluded {
					if minEdge == nil || c.Edge.Weight < minEdge.Weight {
						if debug {
							fmt.Println("Changing to vertex ", c.Vertex.Label)
						}
						minEdge = c.Edge
						minVertex = c.Vertex
					}
				}
			}

		}
		if minEdge == nil || minVertex == nil {
			if debug {
				fmt.Println("No more unvisited vertices")
			}
			return
		} else {

			if debug {
				fmt.Println("[Visible] Adding edge: " + minEdge.Label)
			}

			minEdge.IsIncluded = true

			minVertex.IsVisited = true

			graphMST = append(graphMST, minVertex)
		}

	}
}

func hasUnvisitedVertices(graph *[][]*Vertex) bool {

	for i := 0; i < len(*graph); i++ {
		for j := 0; j < len(*graph); j++ {
			if !(*graph)[i][j].IsVisited {
				return true
			}
		}
	}

	return false
}

func drawGraph(graph *[][]*Vertex) string {

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
	<script>
        function draw() {
            var canvas = document.getElementById('tutorial');
            if (canvas.getContext) {
                var ctx = canvas.getContext('2d');

				`

	// Print all vertex
	edgeOffset := 2
	for i := 0; i < size; i++ {
		offsetY := i * 20
		offsetX := 0
		for j := 0; j < size; j++ {
			offsetX = j * 20
			html += `
						ctx.beginPath();
						ctx.arc(` + strconv.Itoa(2+offsetX) + `, ` + strconv.Itoa(2+offsetY) + `, 2, 0, 2 * Math.PI, true);
						ctx.fillStyle = '#0095DD';
						ctx.fill();
						`

			// Draw Connections
			for _, c := range (*graph)[i][j].Connections {

				if c.Edge.IsIncluded && !c.Edge.Printed {

					if c.Direction == "right" {
						c.Edge.Printed = true
						html += `
						ctx.beginPath();
						ctx.moveTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
						ctx.lineTo(` + strconv.Itoa(edgeOffset+offsetX+20) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
						ctx.stroke();
						`
					} else if c.Direction == "bottom" {
						c.Edge.Printed = true
						html += `
							ctx.beginPath();
							ctx.moveTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
							ctx.lineTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY+20) + `);
							ctx.stroke();
							`
					} else if c.Direction == "left" {
						c.Edge.Printed = true
						html += `
							ctx.beginPath();
							ctx.moveTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
							ctx.lineTo(` + strconv.Itoa(edgeOffset+offsetX-20) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
							ctx.stroke();
							`
					} else if c.Direction == "top" {
						c.Edge.Printed = true
						html += `
							ctx.beginPath();
							ctx.moveTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY) + `);
							ctx.lineTo(` + strconv.Itoa(edgeOffset+offsetX) + `, ` + strconv.Itoa(edgeOffset+offsetY-20) + `);
							ctx.stroke();
							`
					}

				}

			}

		}

	}

	// Draw vertex
	//ctx.beginPath();
	//ctx.moveTo(0, 2);
	//ctx.arc(2, 2, 2, 0, Math.PI * 2, true); // Outer circle
	//ctx.fill();

	// Draw edge
	//ctx.beginPath();
	//ctx.moveTo(0, 2);
	//ctx.lineTo(42, 2);
	//ctx.stroke();

	html += `}
		}
	</script>

</head>
<body onload="draw();">
    <canvas id="tutorial" ></canvas>`

	html += `<table style="border-collapse: collapse;">`

	for i := 0; i < size; i++ {
		html += `<tr>`

		offsetY := i * 20
		offsetX := 0
		if debug {
			fmt.Println(offsetX, offsetY)
		}
		for j := 0; j < size; j++ {
			offsetX = j * 20

			border := 2
			borderLeft, borderTop, borderRight, borderBottom := border, border, border, border

			for _, c := range (*graph)[i][j].Connections {
				if c.Edge.IsIncluded {
					if c.Direction == "right" {
						borderRight = 0
					} else if c.Direction == "bottom" {
						borderBottom = 0
					} else if c.Direction == "left" {
						borderLeft = 0
					} else if c.Direction == "top" {
						borderTop = 0
					}

				}
			}

			html += fmt.Sprintf(`<td style="border-left: %dpx solid black; border-top: %dpx solid black; border-right: %dpx solid black; border-bottom: %dpx solid black;"></td>`, borderLeft, borderTop, borderRight, borderBottom)
		}

		html += `</tr>`

	}

	html += `</table>`

	for _, l := range *graph {
		for _, v := range l {
			for _, c := range v.Connections {
				if c.Edge.IsIncluded {
					html += `<p> Vertex: ` + v.Label + ` Direction: ` + c.Direction + ` Edge: ` + c.Edge.Label + `</p>`
				}
			}
		}
	}

	html += `</body>
	</html>
	`
	return html
}
