package utils

import (
	"archive/zip"
	"bytes"
)

func ReadZipData(data []byte) (*zip.Reader, error) {
	return zip.NewReader(bytes.NewReader(data), int64(len(data)))
} 