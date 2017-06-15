package fily

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func SaveFile(r *http.Request) (string, error) {
	var name = ""

	// 1. Save file to temp directory
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	fmt.Printf("Type of file %T\n", file)

	out, err := os.Create("./public/tmp/" + header.Filename)
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

	log.Println(header.Filename)

	// 2. Reduce image size and remove original file

	// 3. Save to Db

	return name, nil

}
