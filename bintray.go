package eve

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	binTrayURL      = ""
	binTrayUser     = ""
	binTrayPassword = ""
)

func SetBinTrayURL(URL string) {
	binTrayURL = URL
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
		return errors.New("return status for deleting file <" + url + "> was not 200!")
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
		return errors.New("return status for publishing file <" + url + "> was not 200!")
	}
	return nil
}
