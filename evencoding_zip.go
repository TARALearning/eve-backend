package eve

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	uint32max = (1 << 32) - 1
)

// FileInfo appends the Path to the os.FileInfo struct
type FileInfo struct {
	fName    string
	fSize    int64
	fMode    os.FileMode
	fModTime time.Time
	fisDir   bool
	fSys     interface{}
	fPath    string
}

// Name implements the os.FileInfo interface
func (fi *FileInfo) Name() string {
	return fi.fName
}

// Size implements the os.FileInfo interface
func (fi *FileInfo) Size() int64 {
	return fi.fSize
}

// Mode implements the os.FileInfo interface
func (fi *FileInfo) Mode() os.FileMode {
	return fi.fMode
}

// ModTime implements the os.FileInfo interface
func (fi *FileInfo) ModTime() time.Time {
	return fi.fModTime
}

// Sys implements the os.FileInfo interface
func (fi *FileInfo) Sys() interface{} {
	return fi.fSys
}

// Path appends the path to the the os.FileInfo interface
func (fi *FileInfo) Path() string {
	return fi.fPath
}

// IsDir implements the os.FileInfo interface
func (fi *FileInfo) IsDir() bool {
	return fi.fisDir
}

// NewFileInfo appends the Path to the os.FileInfo struct
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

// saveFiles appends a file to the zipFiles
func saveFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fi := NewFileInfo(info)
	fi.fPath = path
	zipFiles = append(zipFiles, fi)
	return nil
}

// readFiles reads all files from the given src recursive into the zipFiles
func readFiles(src string) error {
	return filepath.Walk(src, saveFiles)
}

// FileInfoHeader returns a zip.FileHeader from the custom FileInfo
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

// Zip compresses the given source into the given target zip file
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
		header, herr := FileInfoHeader(file)
		if herr != nil {
			return herr
		}
		header.Name = file.Path()
		if file.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, werr := w.CreateHeader(header)
		if werr != nil {
			return werr
		}
		if file.IsDir() {
			continue
		}
		fil, oerr := os.Open(file.Path())
		if oerr != nil {
			return oerr
		}
		_, cerr := io.Copy(writer, fil)
		if cerr != nil {
			return cerr
		}
		fil.Close()
	}
	// make sure to check the error on close
	err = w.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(target+".zip", buf.Bytes(), 0777)
}

// UnZip extracts the given archive into the given path
func UnZip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
