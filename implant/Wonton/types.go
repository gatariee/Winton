package main

import (
	"time"
)

type File struct {
	Filename string
	Size     int64
	IsDir    bool
	ModTime  time.Time
}
