package main

import(
    "fmt"
    "strconv"
    "container/heap"
    "os"
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
    var solveCubes []Cube
    for current.cube.asString() != root{
        result = append(result, current.lastMove)
        solveCubes = append(solveCubes, current.cube)
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
        solveCubes[i].draw("cubes/" + strconv.Itoa(len(result)-i-1) + ".png")
        result[len(result)-i-1] = temp[i] // reverse
    }
    return result
}

func main() {
    // Check if cube string argument is provided
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run . <cube_string>")
        fmt.Println("Example: go run . \"000000000111111111222222222333333333444444444555555555\"")
        os.Exit(1)
    }
    
    cubeString := os.Args[1]
    
    // Validate cube string length
    if len(cubeString) != 54 {
        fmt.Printf("Error: Cube string must be exactly 54 characters. Got %d characters.\n", len(cubeString))
        os.Exit(1)
    }
    
    // Validate cube string characters
    for i, char := range cubeString {
        if char < '0' || char > '5' {
            fmt.Printf("Error: Invalid character '%c' at position %d. Only digits 0-5 are allowed.\n", char, i)
            os.Exit(1)
        }
    }
    
    var root CubeNode
    root.cube = cubeFromString(cubeString)
    
    // Validate cube configuration
    if !root.cube.isValid() {
        fmt.Println("Error: Invalid cube configuration. Each color must appear exactly 9 times.")
        os.Exit(1)
    }
    
    root.cube.draw("cubes/start.png")
 
    solve := AStarSearch(root.cube.asString())

    moveStrings := []string{"F","R","L","U","D","B","F'","R'","L'","U'","D'","B'"} 
    
    // Apply each move and output both the move and resulting cube state
    solveCube := root.cube
    fmt.Printf("step 0: START|%s\n", solveCube.asString())
    
    for i, move := range(solve) {
        solveCube.rotate(move)
        fmt.Printf("step %d: %s|%s\n", i+1, moveStrings[move], solveCube.asString())
    }
}
