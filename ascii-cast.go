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

// Could clean this up a bit
func scaleImage(img image.Image, terminalWidth, terminalHeight float64) (image.Image, int, int) {
	imageRect := img.Bounds()
	maxImageX := float64(imageRect.Max.X)
	maxImageY := float64(imageRect.Max.Y)

	scaledHeight := maxImageY
	scaledWidth := maxImageX

	if maxImageY >= maxImageX {
		shrinkFactor := terminalHeight / maxImageY
		scaledHeight = terminalHeight
		scaledWidth = maxImageX * shrinkFactor
	} else {
		shrinkFactor := terminalWidth / maxImageX
		scaledWidth = terminalWidth
		scaledHeight = maxImageY * shrinkFactor
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
