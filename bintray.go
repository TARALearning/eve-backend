package eve

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	binTrayURL         = ""
	binTrayDownloadURL = ""
	binTrayUser        = ""
	binTrayPassword    = ""
)

func SetBinTrayURL(URL string) {
	binTrayURL = URL
}

func SetBinTrayDownloadURL(URL string) {
	binTrayDownloadURL = URL
}

func SetBinTrayUser(username string) {
	binTrayUser = username
}

func SetBinTrayPassword(password string) {
	binTrayPassword = password
}

func BinTrayDeleteFile(subject, repo, filepath string) error {
	client := EvHTTPNewClient()
	url := binTrayURL + "/content/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
	fmt.Println("delete file ::", url, "...")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(binTrayUser, binTrayPassword)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("return status for deleting file <" + url + "> was not 200! it was " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func BinTrayPublishFile(subject, repo, filepath string) error {
	client := EvHTTPNewClient()
	fmt.Println("waint 5 second before publishing " + filepath)
	time.Sleep(time.Second * 5)
	url := binTrayURL + "/file_metadata/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
	fmt.Println("publish file ::", url, "...")
	buff := bytes.NewBuffer(nil)
	buff.WriteString(`{"list_in_downloads":true}`)
	req, err := http.NewRequest(http.MethodPut, url, buff)
	if err != nil {
		return err
	}
	req.SetBasicAuth(binTrayUser, binTrayPassword)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("return status for publishing file <" + url + "> was not 200! it was " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func BinTrayDownloadFile(subject, repo, filepath string) ([]byte, error) {
	client := EvHTTPNewClient()
	fmt.Println("waint 3 second before downloading " + filepath)
	time.Sleep(time.Second * 3)
	url := binTrayDownloadURL + "/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
	fmt.Println("downloading file ::", url, "...")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(binTrayUser, binTrayPassword)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("return status for downloading file <" + url + "> was not 200! it was " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
