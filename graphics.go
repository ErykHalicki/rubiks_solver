package main

import(
    "image"
	"image/color"
	"image/png"
	"os"
)

func drawRect (img *image.RGBA, x int, y int, size int, color_ int8) {
    // Define square color
    var squareColor color.RGBA
    switch(color_){
        case 0:
	        squareColor = color.RGBA{0, 255, 0, 255} // green
        case 1:
	        squareColor = color.RGBA{255, 255, 0, 255} // yellow
        case 2:
	        squareColor = color.RGBA{255, 255, 255, 255} // white
        case 3:
	        squareColor = color.RGBA{255, 0, 0, 255} // red
        case 4:
	        squareColor = color.RGBA{255, 165, 0, 255} // orange
        case 5:
	        squareColor = color.RGBA{0, 0, 255, 255} // blue
    }

	// Draw a square
	for i := x; i < x+size; i++ {
		for j := y; j < y+size; j++ {
			img.Set(i, j, squareColor)
		}
	}
}

func (c *Cube) draw(filePath string) {
	// Image dimensions
	width, height := 600, 500

	// Create a new blank image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    for f := 0; f < 6; f++{
        for x := 0; x < 3; x++{
            for y := 0; y < 3; y++{
                switch(f){
                    case 3:
                        drawRect(img, 200+x*33, 100+y*33, 32, c[f][y][x])
                    case 4:
                        drawRect(img, 200+x*33, 300+y*33, 32, c[f][y][x])
                    case 0:
                        drawRect(img, 200+x*33, 200+y*33, 32, c[f][y][x])
                    case 2:
                        drawRect(img, 100+x*33, 200+y*33, 32, c[f][y][x])
                    case 1:
                        drawRect(img, 300+x*33, 200+y*33, 32, c[f][y][x])
                    case 5:
                        drawRect(img, 400+x*33, 200+y*33, 32, c[f][y][x])
                }
            }
        }
    }

	
	// Save the image to a file
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	// println("cube image saved in " + filePath)
}
