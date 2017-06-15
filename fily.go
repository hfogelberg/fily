package fily

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func New(r *http.Request, size uint) (string, error) {
	fileName := time.Now().String() + ".jpg"

	// 1. Save file to temp directory
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	out, err := os.Create("./public/tmp/" + fileName)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 2. Reduce image size
	path := "./public/tmp/" + fileName
	fResize, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(fResize)
	if err != nil {
		log.Fatal(err)
	}
	defer fResize.Close()

	name, err := Resize(img, size)
	if err != nil {
		return "", err
	}

	// 3. Remove uncompressed filed
	err = os.Remove(path)
	if err != nil {
		log.Println(err)
	}

	return name, nil
}

func Resize(img image.Image, width uint) (string, error) {
	fileOut := time.Now().String() + ".jpg"

	m := resize.Resize(1000, 0, img, resize.NearestNeighbor)
	out, err := os.Create("./public/tmp/" + fileOut)
	if err != nil {
		return "", err
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)

	return fileOut, nil
}
