package main

import (
    "fmt"
    "math"
)

type Face [3][3]int8

// Flattened rubik's cube representation
// Faces:
//   1
// 4 0 2 5
//   3
// 0 = front face, 5 = back face, 4 = bottom face

type Cube struct {
    faces [6]Face
    score int8 // # of correct squares
}

func initFace (color int8) Face{
    var f Face
    for x:=0; x<3; x++ {
        for y:=0; y<3; y++ {
            f[x][y] = color
        }
    }
    return f
}

func initCube () Cube{
    var c Cube
    
    var i int8
    for ; i<6; i++{
        c.faces[i] = initFace(i)
    }
    
    c.score = c.calculateScore() // all squares are in correct position
    return c
}

func (f Face) calculateScore () int8{
    var counts [6]int8 // # of squares of each color
    var score int8
    for i := 0; i < 3; i++{
        for j := 0; j < 3; j++{
            counts[f[i][j]]++;
            if counts[f[i][j]] > score{
                score = counts[f[i][j]]
            }
        }
    }
    return score // returns the how many matching squares there are
}

func (c Cube) calculateScore () int8{
    var score int8
    for i := 0; i<6; i++ {
        score += c.faces[i].calculateScore()
    }
    return score // returns how many matching squares there are on the whole cube
}

func (f *Face) rotate (clockwise bool){ // rotates face in place
    var temp Face
    copy(temp[:], f[:])
    if clockwise {
        for i := 0; i<3; i++{
            f[i][2] = temp[0][i] // top -> right
            f[2][i] = temp[i][2] // right -> bottom
            f[i][0] = temp[2][i] // bottom -> left
            f[0][i] = temp[i][0] // left -> top
        }
    } else {
        for i := 0; i<3; i++{
            f[0][i] = temp[i][2] // top <- right
            f[i][2] = temp[2][i] // right <- bottom
            f[2][i] = temp[i][0] // bottom <- left
            f[i][0] = temp[0][i] // left <- top
        }
    }
}

func (c *Cube) move (m int) Cube{ // apply the movement in place
    // movement switch
    var abs int
    abs = m
    if m < 0 {abs = -m}
    c.faces[abs](m < 0)

    return c
}

func (c Cube) asString () string{ // returns a string representation of the cube, used for the trie
    return "rrrrrrggggggbbbbbbwwwwwwooooooyyyyyy" // temp value
}

func main() {
    c := initCube()
    fmt.Println(c.score)
    c.faces[1].rotate(true)
}
