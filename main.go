package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`usage: gif2pngs <path>

  gif2pngs will write a .png file for every frame of the gif.
  The images are stored at the same path as the .gif but with numbers added to
  the end of the file names and with .png as the extension.`)
		return
	}
	path := os.Args[1]
	if !strings.HasSuffix(path, ".gif") {
		fmt.Println("The input file must be a .gif file")
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading .gif file:", err.Error())
		return
	}
	r := bytes.NewReader(data)
	img, err := gif.DecodeAll(r)
	if err != nil {
		fmt.Println("Error decoding .gif file:", err.Error())
		return
	}
	dir, file := filepath.Split(path)
	file = strings.TrimSuffix(file, ".gif")
	base := filepath.Join(dir, file)
	if len(img.Image) == 0 {
		return
	}
	dest := image.NewRGBA(img.Image[0].Bounds())
	for i, img := range img.Image {
		draw.Draw(dest, img.Bounds(), img, img.Bounds().Min, draw.Over)
		var w bytes.Buffer
		err := png.Encode(&w, dest)
		if err != nil {
			fmt.Println("Error encoding .png file:", err.Error())
			return
		}
		err = ioutil.WriteFile(
			fmt.Sprintf(base+"%05d.png", i),
			w.Bytes(),
			0666,
		)
		if err != nil {
			fmt.Println("Error writing .png file:", err.Error())
			return
		}
	}
}
