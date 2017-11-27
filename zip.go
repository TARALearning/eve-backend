package eve

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	uint32max = (1 << 32) - 1
)

type FileInfo struct {
	fName    string
	fSize    int64
	fMode    os.FileMode
	fModTime time.Time
	fisDir   bool
	fSys     interface{}
	fPath    string
}

func (fi *FileInfo) Name() string {
	return fi.fName
}

func (fi *FileInfo) Size() int64 {
	return fi.fSize
}

func (fi *FileInfo) Mode() os.FileMode {
	return fi.fMode
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.fModTime
}

func (fi *FileInfo) Sys() interface{} {
	return fi.fSys
}

func (fi *FileInfo) Path() string {
	return fi.fPath
}

func (fi *FileInfo) IsDir() bool {
	return fi.fisDir
}

func NewFileInfo(info os.FileInfo) FileInfo {
	fi := new(FileInfo)
	fi.fName = info.Name()
	fi.fSize = info.Size()
	fi.fMode = info.Mode()
	fi.fModTime = info.ModTime()
	fi.fisDir = info.IsDir()
	fi.fPath = ""
	return *fi
}

var files = make([]os.FileInfo, 0)
var zipFiles = make([]FileInfo, 0)

func GetFolderName(originalName string) (string, error) {
	splited := strings.Split(originalName, "-")
	moduleNr, err := strconv.Atoi(splited[0])
	if err != nil {
		return "", err
	}
	return "m" + strconv.Itoa(moduleNr), nil
}

func saveFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fi := NewFileInfo(info)
	fi.fPath = path
	zipFiles = append(zipFiles, fi)
	return nil
}

func readFiles(src string) error {
	return filepath.Walk(src, saveFiles)
}

func FileInfoHeader(fi FileInfo) (*zip.FileHeader, error) {
	size := fi.Size()
	fh := &zip.FileHeader{
		Name:               fi.Name(),
		UncompressedSize64: uint64(size),
	}
	fh.SetModTime(fi.ModTime())
	fh.SetMode(fi.Mode())
	if fh.UncompressedSize64 > uint32max {
		fh.UncompressedSize = uint32max
	} else {
		fh.UncompressedSize = uint32(fh.UncompressedSize64)
	}
	return fh, nil
}

func Zip(src, target string) error {
	fmt.Println("zipping", src, "into", target+".zip")
	// reset zipFiles
	zipFiles = make([]FileInfo, 0)
	// read zip files
	err := readFiles(src)
	if err != nil {
		return err
	}
	// create a buffer to write our archive to
	buf := new(bytes.Buffer)
	// create a new zip archive
	w := zip.NewWriter(buf)
	for _, file := range zipFiles {
		header, err := FileInfoHeader(file)
		if err != nil {
			return err
		}
		header.Name = file.Path()
		if file.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		if file.IsDir() {
			continue
		}
		fil, err := os.Open(file.Path())
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fil)
		fil.Close()
	}
	// make sure to check the error on close
	err = w.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(target+".zip", buf.Bytes(), 0777)
}
