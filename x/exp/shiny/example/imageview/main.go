// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build example
//
// This build tag means that "go install github.com/Andyfoo/golang/x/exp/shiny/..." doesn't
// install this example program. Use "go run main.go" to run it or "go install
// -tags=example" to install it.

// Imageview is a basic image viewer. Supported image formats include BMP, GIF,
// JPEG, PNG, TIFF and WEBP.
package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/Andyfoo/golang/x/exp/shiny/driver"
	"github.com/Andyfoo/golang/x/exp/shiny/screen"
	"github.com/Andyfoo/golang/x/exp/shiny/widget"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/Andyfoo/golang/x/image/bmp"
	_ "github.com/Andyfoo/golang/x/image/tiff"
	_ "github.com/Andyfoo/golang/x/image/webp"
)

// TODO: scrolling, such as when images are larger than the window.

func decode(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode %s: %v", filename, err)
	}
	return m, nil
}

func main() {
	log.SetFlags(0)
	driver.Main(func(s screen.Screen) {
		if len(os.Args) < 2 {
			log.Fatal("no image file specified")
		}
		// TODO: view multiple images.
		src, err := decode(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		w := widget.NewSheet(widget.NewImage(src, src.Bounds()))
		if err := widget.RunWindow(s, w, &widget.RunWindowOptions{
			NewWindowOptions: screen.NewWindowOptions{
				Title: "ImageView Shiny Example",
			},
		}); err != nil {
			log.Fatal(err)
		}
	})
}
