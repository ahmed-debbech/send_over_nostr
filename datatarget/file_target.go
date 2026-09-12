package datatarget

import (
	"bufio"
	"log"
	"os"
)

type FileTarget struct {
	DataListenerCh chan []byte
	FilePtr        *os.File
}

func (f *FileTarget) Init() {
	file, err := os.OpenFile(
		"t1.jpg",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	f.FilePtr = file
}

func (f *FileTarget) ListenForData() chan []byte {

	f.DataListenerCh = make(chan []byte, 8192)

	writer := bufio.NewWriterSize(f.FilePtr, 8192)

	go func() {
		for {
			data := <-f.DataListenerCh
			if _, err := writer.Write(data); err != nil {
				log.Fatal(err)
			}
			writer.Flush()
		}
	}()

	return f.DataListenerCh
}

func (f *FileTarget) Destroy() {
	f.FilePtr.Close()
}
