package datasource

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
)

type FileSource struct {
	Scanner *bufio.Scanner
	FilePtr *os.File
	FileCh  chan []byte
}

func (f *FileSource) Init() error {
	fileLocal, err := os.Open("t.jpg")
	if err != nil {
		log.Println("Error opening file to start sending it because:", err)
		return errors.New("Error opening file to start sending it")
	}
	f.FilePtr = fileLocal
	f.Scanner = bufio.NewScanner(f.FilePtr)
	f.FileCh = make(chan []byte)
	return nil
}

func (f *FileSource) StartFetch() chan []byte {
	go func() {

		for {
			buf := make([]byte, 8192)

			n, err := f.FilePtr.Read(buf)
			if n > 0 {
				dat := buf[:n]
				f.FileCh <- dat
			}

			if err == io.EOF {
				break
			}
		}
		close(f.FileCh)

	}()

	return f.FileCh
}

func (f *FileSource) Destroy() {
	f.FilePtr.Close()
}
