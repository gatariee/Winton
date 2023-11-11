package main

import (
	"os"
)

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
		file := File{
			Filename: fileInfo.Name(),
			Size:     fileInfo.Size(),
			IsDir:    fileInfo.IsDir(),
			ModTime:  fileInfo.ModTime(),
		}

		files = append(files, file)

	}

	return files, nil
}
