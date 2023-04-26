package util

import (
	"image"
	"io/fs"
)

func LoadImage(filename string, filesystem fs.FS) (image.Image, string, error) {
	f, err := filesystem.Open(filename)
	if err != nil {
		return nil, "", err
	}
	return image.Decode(f)
}
