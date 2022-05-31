package service

import (
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SaveImage(data, path string) {
	path = "./images/" + path
	data = data[23:]
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Print(err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		log.Print(err)
		return
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()
	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Print(err)
		return
	}
}
