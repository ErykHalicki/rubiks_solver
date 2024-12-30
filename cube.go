package main

import (
    "fmt"
)


type Face [3][3]int8

// Flattened rubik's cube representation
// Faces:
//   3
// 2 0 1 5
//   4
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

func (c *Cube) rotate (m int){ // apply the movement in place
    var temp Cube
    copy(temp.faces[:], c.faces[:])

    var face int
    face = m
    if face > 5 {face -= 6}
    c.faces[face].rotate(m < 6) // rotate the selected face (0-5 = clockwise, 6-11 - ccw)
    
    // after selected face is rotated, move the required squares on the other faces
    
    var mapping = [6][4]int8{{2,3,1,4},{0,3,5,4},{5,3,0,4},{2,5,1,0},{2,0,1,5},{1,3,2,4}} 
    // mappings for each rotation (face rotation order)

    for i:=0; i < 3; i++ {
        switch(face){
            case 0:
                c.faces[mapping[face][0]][i][2] = temp.faces[mapping[face][3]][0][i]
                c.faces[mapping[face][1]][2][i] = temp.faces[mapping[face][0]][i][2]
                c.faces[mapping[face][2]][i][0] = temp.faces[mapping[face][1]][2][i]
                c.faces[mapping[face][3]][0][i] = temp.faces[mapping[face][2]][i][0]
            case 1:
                c.faces[mapping[face][0]][i][2] = temp.faces[mapping[face][3]][i][2]
                c.faces[mapping[face][1]][i][2] = temp.faces[mapping[face][0]][i][2]
                c.faces[mapping[face][2]][i][0] = temp.faces[mapping[face][1]][i][2]
                c.faces[mapping[face][3]][i][2] = temp.faces[mapping[face][2]][i][0]
            case 2:
                c.faces[mapping[face][0]][2-i][2] = temp.faces[mapping[face][3]][i][0]
                c.faces[mapping[face][1]][i][0] = temp.faces[mapping[face][0]][2-i][2]
                c.faces[mapping[face][2]][i][0] = temp.faces[mapping[face][1]][i][0]
                c.faces[mapping[face][3]][i][0] = temp.faces[mapping[face][2]][i][0]
            case 3:
                c.faces[mapping[face][0]][0][i] = temp.faces[mapping[face][3]][0][i]
                c.faces[mapping[face][1]][0][i] = temp.faces[mapping[face][0]][0][i]
                c.faces[mapping[face][2]][0][i] = temp.faces[mapping[face][1]][0][i]
                c.faces[mapping[face][3]][0][i] = temp.faces[mapping[face][2]][0][i]
            case 4:
                c.faces[mapping[face][0]][2][i] = temp.faces[mapping[face][3]][2][i]
                c.faces[mapping[face][1]][2][i] = temp.faces[mapping[face][0]][2][i]
                c.faces[mapping[face][2]][2][i] = temp.faces[mapping[face][1]][2][i]
                c.faces[mapping[face][3]][2][i] = temp.faces[mapping[face][2]][2][i]
            case 5:
                c.faces[mapping[face][0]][i][2] = temp.faces[mapping[face][3]][2][2-i]
                c.faces[mapping[face][1]][0][i] = temp.faces[mapping[face][0]][i][2]
                c.faces[mapping[face][2]][2-i][0] = temp.faces[mapping[face][1]][0][i]
                c.faces[mapping[face][3]][2][2-i] = temp.faces[mapping[face][2]][2-i][0]
        }
    }
}

func (c Cube) asString () string{ // returns a string representation of the cube, used for the trie
    return "rrrrrrggggggbbbbbbwwwwwwooooooyyyyyy" // temp value
}

func main() {
    c := initCube()
    fmt.Println(c.faces)
    c.faces[1].rotate(true)
}
