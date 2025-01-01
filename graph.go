package main

import(
    "fmt"
    "strconv"
    "container/heap"
)

type CubeNode struct{
    cube Cube
    score int
    visited bool
    lastMove int8
    distance int
}

func (c CubeNode) generateMove(m int8) CubeNode{
    c.cube.rotate(m)
    c.score = c.cube.calculateScore()
    c.visited = false
    c.lastMove = m
    c.distance = 10000000
    return c 
}

func generateMoves (c CubeNode) []CubeNode{
    if c.visited {return nil}
    var result []CubeNode

    for move := 0; move < 12; move ++ {
        //if(int8(move) == c.lastMove) {continue}
        temp := c.generateMove(int8(move))
        result = append(result, temp)
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

type CubeHeap []CubeNode

var mult float32

func (h CubeHeap) Len() int           { return len(h) }
func (h CubeHeap) Less(i, j int) bool { return float32(h[i].distance) - float32(h[i].score)*mult <  float32(h[j].distance) - float32(h[j].score)*mult }
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
    mult = 0.5
    current := nodeFromString(root)

    nodeIndex := make(map[string]CubeNode)
    openSet := &CubeHeap{}

    // min heap
    heap.Init(openSet)
    heap.Push(openSet, current)
    nodeIndex[current.cube.asString()] = current

    topScore := 0
    
    for openSet.Len() > 0{
        current = heap.Pop(openSet).(CubeNode)
        current = nodeIndex[current.cube.asString()] //getting the "true version"
        
        if current.score == 294 {break}

        children := generateMoves(current)

        current.visited = true
        if(current.score > topScore) {
            topScore = current.score
            mult = 0.4 * float32(topScore) / 260.0
            if(topScore > 255) {mult = 5}
            fmt.Println(topScore, current.cube.asString())
        }
        for _, child := range children {
            node, ok := nodeIndex[child.cube.asString()] // get the "true version" of the child
            if !ok { // if not in the map
                child.distance = current.distance + 1
                heap.Push(openSet, child)
                nodeIndex[child.cube.asString()] = child // if there isnt one already, create one
            } else { // if the node has already been added to the node index
                if current.distance + 1 < node.distance && !node.visited{
                    node.distance = current.distance + 1
                    heap.Push(openSet, node)
                }
                
            }
        }
        
    }
    counter := 0
    for current.cube.asString() != root{ 
        result = append(result, current.lastMove)
        current.cube.draw("cubes/" + strconv.Itoa(counter) + ".png")
        node, _ := nodeIndex[current.cube.asString()]
        if node.lastMove < 6 {
            //println(current.lastMove+6)
            current.cube.rotate(node.lastMove+6)
        } else {
            //println(current.lastMove-6)
            current.cube.rotate(node.lastMove-6)
        }
        current = nodeIndex[current.cube.asString()]    
    }
    temp := make([]int8, len(result))
    copy(temp[:], result[:])
    for i := len(result)-1; i>=0; i-- {
        result[len(result)-i-1] = temp[i] // reverse
    }
    return result
}

func main() {
    var root CubeNode
    root.cube = initCube()
    // green = 0, yellow = 1, white = 2, red = 3, orange = 4, blue = 5 
    root.cube = cubeFromString("005500444343215154131024350300232431215341232210254515")
    //scramble := root.cube.scramble(40)
    root.cube.draw("cubes/start.png")
 
    solve := AStarSearch(root.cube.asString())

     
    for i, move := range(solve) {
        var moveString string
        switch move {
            case 0:
        }
        fmt.Println("move " + strconv.Itoa(i) + ": " + moveString)
    }
    
    //fmt.Println("scramble moves: " , scramble)
}
