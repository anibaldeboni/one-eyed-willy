package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"sync"
)

func readFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readFiles(files []*multipart.FileHeader) (filesBytes [][]byte, err error) {
	var wg sync.WaitGroup
	ch := make(chan []byte, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			f, err := readFile(file)
			if err != nil {
				return
			}
			ch <- f
		}(file)
	}
	wg.Wait()
	close(ch)

	var fb [][]byte
	for file := range ch {
		fb = append(fb, file)
	}
	return fb, err
}
