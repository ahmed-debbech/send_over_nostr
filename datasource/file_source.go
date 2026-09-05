package datasource

import (
	"bufio"
	"os"
)

type FileSource struct {
	Scanner *bufio.Scanner
	FilePtr *os.File
	FileCh  chan []byte
}

func (f *FileSource) Init() {
	fileLocal, err := os.Open("test.jpg")
	if err != nil {
		panic(err)
	}
	f.FilePtr = fileLocal
	f.Scanner = bufio.NewScanner(f.FilePtr)
	f.FileCh = make(chan []byte)
}

func (f *FileSource) StartFetch() chan []byte {
	go func() {
		for f.Scanner.Scan() {
			dat := f.Scanner.Bytes()
			f.FileCh <- dat
		}
		close(f.FileCh)
	}()
	return f.FileCh
}

func (f *FileSource) Destroy() {
	f.FilePtr.Close()
}
