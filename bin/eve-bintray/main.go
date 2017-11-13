package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	client     = &http.Client{}
	subject    = os.Args[1]
	repo       = os.Args[2]
	rpackage   = os.Args[3]
	version    = os.Args[4]
	bintrayUrl = "https://api.bintray.com"
	username   = os.Args[5]
	password   = os.Args[6]
)

func deleteFile(filepath string) error {
	url := bintrayUrl + "/content/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
	fmt.Println("delete file ::", url, "...")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("return status for deleting file <" + url + "> was not 200!")
	}
	return nil
}

func main() {
	fmt.Println("getting all files for ", subject, repo, rpackage, version, "...")
	client := &http.Client{}
	url := bintrayUrl + "/packages/" + subject + "/" + repo + "/" + rpackage + "/versions/" + version + "/files?include_unpublished=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("return status code is not 200")
	}
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	respObj := []interface{}{}
	err = json.Unmarshal(respJson, &respObj)
	if err != nil {
		log.Fatal(err)
	}
	for _, binTrayFile := range respObj {
		bTFile := binTrayFile.(map[string]interface{})
		if bTFile["repo"].(string) == "eve-backend" {
			err = deleteFile(bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
