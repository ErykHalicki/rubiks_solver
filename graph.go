package main

import(
    "fmt"
    "container/heap"
    "github.com/dominikbraun/graph"
    "github.com/dominikbraun/graph/draw"
    "os"
)

type CubeNode struct{
    cube Cube
    score int
    visited bool
    lastMove int8
    distance int
}

func cubeHash (c CubeNode) string{
    return c.cube.asString()
}

func (c CubeNode) generateMove(m int8) CubeNode{
    c.cube.rotate(m)
    c.score = c.cube.calculateScore()
    c.visited = false
    c.lastMove = m
    c.distance = 10000000
    return c 
}

func generateMoves (g graph.Graph[string, CubeNode], c CubeNode) []CubeNode{
    if c.visited {return nil}
    var result []CubeNode
    g.AddVertex(c)

    for move := 0; move < 12; move ++ {
        // if(int8(move) == c.lastMove) {continue}
        temp := c.generateMove(int8(move))
        node, err := g.Vertex(temp.cube.asString()) 
        if err != nil { // if not in the heap
            g.AddVertex(temp)
            g.AddEdge(c.cube.asString(), temp.cube.asString())
            result = append(result, temp)
        } else {
            node.lastMove = int8(move)
            result = append(result, node)
        }
    }
    c.visited = true
    return result
} 

func nodeFromString (data string) CubeNode {
    var result CubeNode
    result.cube = cubeFromString(data)
    result.score = result.cube.calculateScore()
    result.visited = false
    result.lastMove = -1 // no previous move
    result.distance = 0
    return result
}

func drawGraph (g graph.Graph[string, CubeNode]){
    file, _ := os.Create("graphs/graph.gv")
    draw.DOT(g, file)
}

type CubeHeap []CubeNode

func (h CubeHeap) Len() int           { return len(h) }
func (h CubeHeap) Less(i, j int) bool { return int(h[i].score) >  int(h[j].score) }
func (h CubeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CubeHeap) Push(x any) {
    *h = append(*h, x.(CubeNode))
}

func (h *CubeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func AStarSearch (root string) []int8 {
    // start at root, generate moves and select the move with the best score
    var result []int8
    current := nodeFromString(root)
    g := graph.New(cubeHash)

    h := &CubeHeap{}

    // min heap
    heap.Init(h)
    heap.Push(h, current)
    
    for current.score < 294 && h.Len() > 0{
        current = heap.Pop(h).(CubeNode)
        children := generateMoves(g, current)
        current.visited = true
        for _, child := range children {
            currScore := current.distance + 1
            if currScore < child.distance {
                child.distance = currScore
            }
            if(!child.visited){
                heap.Push(h, child)
            }
        }
        fmt.Println(current.score, current.cube.asString())
    }
    // result = append(result, current.lastMove)
    drawGraph(g)
    return result
}


func BestFirstSearch (root string) []int8 {
    // start at root, generate moves and select the move with the best score
    var result []int8
    current := nodeFromString(root)
    g := graph.New(cubeHash, graph.PreventCycles())
    
    for current.score != 54{
        // generate moves from current node
        children := generateMoves(g, current)
        max_ := 0
        best := children[0]
        for _, child := range children {
            if(child.score > max_) {
                max_ = child.score
                best = child
            }
        }
        current = best
        result = append(result, current.lastMove)
    }


    drawGraph(g)
    return result
}

func main() {
    var root CubeNode
    root.cube = initCube()
    root.score = root.cube.calculateScore()
    println(root.score)

    root.cube.rotate(0)
    root.cube.rotate(3)
    root.cube.rotate(5)
    root.cube.rotate(1)
    root.cube.rotate(2)
    // root.cube.rotate(3)
    root.cube.draw()
 
    //fmt.Println(BestFirstSearch(root.cube.asString()))
    fmt.Println(AStarSearch(root.cube.asString()))
}
