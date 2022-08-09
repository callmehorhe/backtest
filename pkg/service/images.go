package service

import (
	"encoding/base64"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func SaveImage(data, path string) {
	path = "./images/" + path
	data = data[23:]
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		logrus.Warnf("SaveImage error: Name: %v Decode error: %v", path, err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		logrus.Warnf("SaveImage error: Name: %v Makedirerror: %v", path, err)
		return
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Warnf("SaveImage error: Name: %s Openfile error: %v", path, err)
		return
	}
	defer f.Close()
	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		logrus.Warnf("SaveImage error: Name: %s, Encode error: %v", err)
		return
	}
}
