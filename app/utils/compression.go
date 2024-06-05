package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
)

type compressorFunc func(b []byte) ([]byte, error)

type Compressor struct {
	encodingTypes map[string]compressorFunc
}

func NewCompressor() *Compressor {
	et := make(map[string]compressorFunc)
	et["gzip"] = compressGzip

	return &Compressor{
		encodingTypes: et,
	}
}

func (c *Compressor) Compress(encodingTypes []string, input []byte) ([]byte, error, string) {
	//compresses the data using the first valid encoding type
	for _, et := range encodingTypes {
		et = strings.TrimSpace(et)
		f, found := c.encodingTypes[et]
		if found {
			out, err := f(input)
			return out, err, et
		}
	}

	return input, fmt.Errorf("No supported encoding type"), ""
}

func compressGzip(input []byte) ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	gzipWriter := gzip.NewWriter(buffer)

	_, err := gzipWriter.Write(input)
	if err != nil {
		return b, err
	}
	gzipWriter.Close()

	return buffer.Bytes(), nil
}
