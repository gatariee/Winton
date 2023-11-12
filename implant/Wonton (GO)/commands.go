package main

import (
	"os"
	"time"
	"path/filepath"
)

type File struct {
	Filename string
	Size     int64
	IsDir    bool
	ModTime  time.Time
}

func get_folder_size(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

func ls(path string) ([]File, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var files []File

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			size, err := get_folder_size(path + "/" + fileInfo.Name())
			if err != nil {
				return nil, err
			}

			file := File{
				Filename: fileInfo.Name(),
				Size:     size,
				IsDir:    fileInfo.IsDir(),
				ModTime:  fileInfo.ModTime(),
			}

			files = append(files, file)

		} else {
			file := File{
				Filename: fileInfo.Name(),
				Size:     fileInfo.Size(),
				IsDir:    fileInfo.IsDir(),
				ModTime:  fileInfo.ModTime(),
			}

			files = append(files, file)
		}
	}

	return files, nil
}
