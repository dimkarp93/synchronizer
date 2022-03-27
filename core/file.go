package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func MakeFileBatch(name string) (*FileBatch, error) {
	file, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return &FileBatch{file: file, reader: bufio.NewReaderSize(file, BUF_SIZE)}, nil
}

func (f FileBatch) GetBytes() ([]byte, error) {
	var result []byte
	for {
		data, err := f.reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}

		result = append(result, data...)
	}

	return result, nil
}

func (f FileBatch) Length() int {
	info, err := f.file.Stat()
	if err != nil {
		panic(fmt.Errorf("cannot get stat of file: '%v'", err))
	}

	return int(info.Size())
}
