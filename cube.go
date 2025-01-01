package main

import (
    "strconv"
    "math/rand"
)

type Face [3][3]int8

// Flattened rubik's cube representation
//:
//   3
// 2 0 1 5
//   4
// 0 = front face, 5 = back face, 4 = bottom face

type Cube [6]Face

func initFace (color int8) Face{
    var f Face
    for x:=0; x<3; x++ {
        for y:=0; y<3; y++ {
            f[y][x] = color
        }
    }
    return f
}

func initCube () Cube{
    var c Cube
    
    var i int8
    for ; i<6; i++{
        c[i] = initFace(i)
    }
    
    return c
}

func (f Face) calculateScore () int{
    // var counts [6]int // # of squares of each color
    var score int
    for i := 0; i < 3; i++{
        for j := 0; j < 3; j++{
            /*
            counts[f[i][j]]++;
            if counts[f[i][j]] > score{
                score = counts[f[i][j]]
            }
            */
            for x:=-1; x<=1; x++ {
                if i + x < 0 || i + x > 2 {continue}
                for y:=-1; y<=1; y++ {
                    if j + y < 0 || j + y > 2 {continue}
                    if f[i][j] == f[i+x][j+y] {score ++}
                }
            }
        }
    }
    return score // returns the how many matching squares there are
}

func (c *Cube) calculateScore () int{
    var score int
    for i := 0; i<6; i++ {
        score += c[i].calculateScore()
    }
    return score // returns how many matching squares there are on the whole cube
}

func (f *Face) rotate (clockwise bool){ // rotates face in place
    var temp Face
    copy(temp[:], f[:])
    if clockwise {
        f[0][0] = temp[2][0]
        f[0][1] = temp[1][0]
        f[0][2] = temp[0][0]  
        f[1][0] = temp[2][1]
        f[1][2] = temp[0][1]
        f[2][0] = temp[2][2]
        f[2][1] = temp[1][2]
        f[2][2] = temp[0][2]
    } else {
        f[2][0] = temp[0][0]
        f[1][0] = temp[0][1]
        f[0][0] = temp[0][2]  
        f[2][1] = temp[1][0]
        f[0][1] = temp[1][2]
        f[2][2] = temp[2][0]
        f[1][2] = temp[2][1]
        f[0][2] = temp[2][2]
    }
}

func (c *Cube) rotate (m int8){ // apply the movement in place
    var temp Cube
    copy(temp[:], c[:])

    var face int8
    repeats := 1
    face = m
    if m > 5 {
        face = m-6
        repeats = 3
    }
    c[face].rotate(m < 6) // rotate the selected face (0-5 = clockwise, 6-11 - ccw)
    
    // after selected face is rotated, move the required squares on the other
    
    var mapping = [6][4]int8{{2,3,1,4},{0,3,5,4},{5,3,0,4},{2,5,1,0},{2,0,1,5},{1,3,2,4}} 
    // mappings for each rotation (face rotation order)

    for repeat := 0; repeat < repeats; repeat++ { // rotate 3 times if counter clockwise 
    for i :=0; i < 3; i++ {
        switch(face){
            case 0:
                c[mapping[face][0]][i][2] = temp[mapping[face][3]][0][i]
                c[mapping[face][1]][2][2-i] = temp[mapping[face][0]][i][2]
                c[mapping[face][2]][i][0] = temp[mapping[face][1]][2][i]
                c[mapping[face][3]][0][2-i] = temp[mapping[face][2]][i][0]
            case 1:
                c[mapping[face][0]][i][2] = temp[mapping[face][3]][i][2]
                c[mapping[face][1]][i][2] = temp[mapping[face][0]][i][2]
                c[mapping[face][2]][2-i][0] = temp[mapping[face][1]][i][2]
                c[mapping[face][3]][i][2] = temp[mapping[face][2]][2-i][0]
            case 2:
                c[mapping[face][0]][2-i][2] = temp[mapping[face][3]][i][0]
                c[mapping[face][1]][i][0] = temp[mapping[face][0]][2-i][2]
                c[mapping[face][2]][i][0] = temp[mapping[face][1]][i][0]
                c[mapping[face][3]][i][0] = temp[mapping[face][2]][i][0]
            case 3:
                c[mapping[face][0]][0][i] = temp[mapping[face][3]][0][i]
                c[mapping[face][1]][0][i] = temp[mapping[face][0]][0][i]
                c[mapping[face][2]][0][i] = temp[mapping[face][1]][0][i]
                c[mapping[face][3]][0][i] = temp[mapping[face][2]][0][i]
            case 4:
                c[mapping[face][0]][2][i] = temp[mapping[face][3]][2][i]
                c[mapping[face][1]][2][i] = temp[mapping[face][0]][2][i]
                c[mapping[face][2]][2][i] = temp[mapping[face][1]][2][i]
                c[mapping[face][3]][2][i] = temp[mapping[face][2]][2][i]
            case 5:
                c[mapping[face][0]][i][2] = temp[mapping[face][3]][2][2-i]
                c[mapping[face][1]][0][i] = temp[mapping[face][0]][i][2]
                c[mapping[face][2]][2-i][0] = temp[mapping[face][1]][0][i]
                c[mapping[face][3]][2][2-i] = temp[mapping[face][2]][2-i][0]
        }
    }
    copy(temp[:], c[:])
    }
}

func (c *Cube) asString () string{ // returns a string representation of the cube, used for the trie
    var result string

    for f := 0; f < 6; f++{
        for x := 0; x < 3; x++{
            for y := 0; y < 3; y++{
                result += strconv.Itoa(int(c[f][x][y]))
            }
        }
    }
    return result
}

func  cubeFromString (data string) Cube{ // returns a string representation of the cube, used for the trie
    var counter int
    var c Cube
    for f := 0; f < 6; f++{
        for x := 0; x < 3; x++{
            for y := 0; y < 3; y++{
                c[f][x][y] = int8(data[counter]-48)
                counter ++
            }
        }
    }
    return c
}

func (c *Cube) scramble (moves int) []int8 {
    var result []int8
    for i := 0; i < moves; i++ {
        move := int8(rand.Intn(12))
        c.rotate(move)
        result = append(result, move)
        //println(move)
    }
    return result
}
