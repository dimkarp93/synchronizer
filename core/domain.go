package core

import (
	"bufio"
	"os"
	"time"
)

const BUF_SIZE int = 4096

type Source struct {
	Items []SourceItem
}

type SourceItem struct {
	History History
	Data    map[string]Batch
}

type History struct {
	Time     time.Time
	Revision string
}

type Batch interface {
	GetBytes() []byte
	Length() int
}

type FileBatch struct {
	file   *os.File
	reader *bufio.Reader
}
