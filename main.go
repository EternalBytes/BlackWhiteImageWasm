package main

//
// This code was written by @EternalBytes github.com/EternalBytes
//
import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"strings"
	"syscall/js"
)

func main() {
	wait := make(chan error)
	fmt.Println("Wasm running.")
	jsFunc := js.FuncOf(Transform)
	js.Global().Set("Transform", jsFunc)
	<-wait // Just to avoid the program to end execution
}

func Transform(this js.Value, args []js.Value) any {
	// Decode the image
	img, err := jpeg.Decode(decodeBase64(args[0].String()[22:]))
	if err != nil {
		panic(errors.New(err.Error() + " // LINHA 32"))
	}

	// Create a new grayscale image with the same dimensions as the original image
	grayImg := image.NewGray(img.Bounds())

	// Convert each pixel to grayscale
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// Get the original color of the pixel
			originalColor := img.At(x, y)
			// Convert the color to grayscale
			grayColor := color.GrayModel.Convert(originalColor).(color.Gray)
			// Set the grayscale color for the corresponding pixel in the new image
			grayImg.Set(x, y, grayColor)
		}
	}

	// Create a new file to save the black and white image
	var outputImage bytes.Buffer
	// Encode the grayscale image as a JPEG
	err = jpeg.Encode(&outputImage, grayImg, nil)
	if err != nil {
		panic(errors.New(err.Error() + " // LINHA 56"))
	}

	var outputBase64 strings.Builder

	var wrtCloser = encodeBase64(&outputBase64)
	wrtCloser.Write(outputImage.Bytes())
	wrtCloser.Close()

	return outputBase64.String()
}

// Receives a base64 string without the type specification ex.: data:image/png;base64,
// It receives only the base64 payload.
func decodeBase64(base string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(base))
}

// Receives a pointer to a value of type io.Writer to write to.
// It returns a io.WriteCloser which is used to write data to the argument.
func encodeBase64(base io.Writer) io.WriteCloser {
	return base64.NewEncoder(base64.StdEncoding, base)
}
