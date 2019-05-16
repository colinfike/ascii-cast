package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/nfnt/resize"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

const asciistr = "MND8OZ$7I?+=~:,.."

// 70 Characters
// const asciistr = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
// 64 characters
// const asciistr = "$@B%8&WM#*ohkbdpqwmZO0QCJYXzcvunxrjt/\\|()1{[]?-_+~<>!lI;:,\"^`'. "

func getImage(fpath string) image.Image {
	f, err := os.Open("frames/" + fpath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	return img
}

func scaleImage(img image.Image, terminalWidth, terminalHeight float64) (image.Image, int, int) {
	sz := img.Bounds()
	maxX := float64(sz.Max.X)
	maxY := float64(sz.Max.Y)

	scaledHeight := maxY
	scaledWidth := maxX

	// fmt.Println("MaxX: ", maxX)
	// fmt.Println("MaxY: ", maxY)
	// fmt.Println("terminalWidth: ", terminalWidth)
	// fmt.Println("terminalHeight: ", terminalHeight)

	// You have image of dimensions X and Y
	// You have terminal of dimension XT and YT

	// if X > Y

	if maxY >= maxX {
		shrinkFactor := terminalHeight / maxY
		// fmt.Println(shrinkFactor)
		scaledHeight = terminalHeight
		scaledWidth = maxX * shrinkFactor
	} else {
		shrinkFactor := terminalWidth / maxX
		// fmt.Println(shrinkFactor)
		scaledWidth = terminalWidth
		scaledHeight = maxY * shrinkFactor
	}

	scaledWidth = scaledWidth * 3

	img = resize.Resize(uint(scaledWidth), uint(scaledHeight), img, resize.Lanczos3)
	return img, int(scaledWidth), int(scaledHeight)
}

func convert2Ascii(img image.Image, w, h int) []byte {
	table := []byte(asciistr)
	buf := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			// pos := int(y / 4)
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// 0-255 Divide by 70

func main() {
	width, _ := terminal.Width()
	height, _ := terminal.Height()

	files, err := ioutil.ReadDir("./frames")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		image := getImage(file.Name())
		p := convert2Ascii(scaleImage(image, float64(width), float64(height-1)))
		print("\033[0;0H")
		os.Stdout.Write(p)
		time.Sleep(26 * time.Millisecond)
	}
}
