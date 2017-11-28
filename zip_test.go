package eve

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_Zip(t *testing.T) {
	testFolder := "test-zip-folder"
	testFolder2 := "test-zip-folder-2"
	testFolder3 := "test-zip-folder-3"
	err := os.MkdirAll(testFolder, 0777)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(testFolder+string(os.PathSeparator)+"testfile1.txt", []byte("test1"), 0777)
	if err != nil {
		t.Error(err)
	}
	err = os.MkdirAll(testFolder+string(os.PathSeparator)+testFolder2, 0777)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(testFolder+string(os.PathSeparator)+testFolder2+string(os.PathSeparator)+"testfile2.txt", []byte("test2"), 0777)
	if err != nil {
		t.Error(err)
	}
	err = os.MkdirAll(testFolder+string(os.PathSeparator)+testFolder3, 0777)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(testFolder+string(os.PathSeparator)+testFolder3+string(os.PathSeparator)+"testfile3.txt", []byte("test3"), 0777)
	if err != nil {
		t.Error(err)
	}
	err = Zip(testFolder, testFolder)
	if err != nil {
		t.Error(err)
	}
	err = os.RemoveAll(testFolder)
	if err != nil {
		t.Error(err)
	}
	err = UnZip(testFolder+".zip", testFolder)
	if err != nil {
		t.Error(err)
	}
	err = os.RemoveAll(testFolder)
	if err != nil {
		t.Error(err)
	}
	err = os.Remove(testFolder + ".zip")
	if err != nil {
		t.Error(err)
	}
}

func Test_FileInfo(t *testing.T) {
	testfile := "tests/tmp/test.txt"
	defer os.Remove(testfile)
	err := ioutil.WriteFile(testfile, []byte("content"), 0777)
	if err != nil {
		t.Error(err)
	}
	f, err := os.Open(testfile)
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	fInfo, err := f.Stat()
	if err != nil {
		t.Error(err)
	}
	fi := NewFileInfo(fInfo)
	t.Log(fi.Name())
	t.Log(fi.IsDir())
	t.Log(fi.ModTime())
	t.Log(fi.Mode())
	t.Log(fi.Path())
	t.Log(fi.Size())
	t.Log(fi.Sys())
}
