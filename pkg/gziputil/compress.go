package gziputil

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

var magicGzip = []byte{0x1f, 0x8b, 0x08}

func IsCompressed(path string) (bool, error) {
	fs, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer fs.Close()

	reader := bufio.NewReader(fs)
	lr := io.LimitReader(reader, 8)
	b, err := ioutil.ReadAll(lr)
	if err != nil {
		return false, err
	}

	if bytes.Equal(b[0:3], magicGzip) {
		return true, nil
	}

	return false, nil
}

func Decompress(fromPath, toPath string) error {
	fromFS, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFS.Close()

	gr, err := gzip.NewReader(bufio.NewReader(fromFS))
	if err != nil {
		return err
	}

	toFS, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFS.Close()

	_, err = io.Copy(toFS, gr)
	if err != nil {
		return err
	}

	return nil
}

func CompressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return []byte{}, err
	}
	if _, err = w.Write(data); err != nil {
		return []byte{}, err
	}
	w.Close()

	return buf.Bytes(), nil
}
