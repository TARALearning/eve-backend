package main

import (
	"bytes"
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
	command    = ""
	subject    = ""
	repo       = ""
	rpackage   = ""
	version    = ""
	bintrayURL = ""
	username   = ""
	password   = ""
)

func init() {
	if len(os.Args) < 8 {
		fmt.Println("please specify all the required arguments to run the bintray cleanup script")
		fmt.Println("")
		fmt.Println("eve-bintray {command} {subject} {repo} {password} {package} {version} {api_url} {username} {token}")
		fmt.Println("")
		fmt.Println("example:")
		fmt.Println("")
		fmt.Println("eve-bintray \\")
		fmt.Println("    delete | list | publish \\")
		fmt.Println("    evalgo \\")
		fmt.Println("    eve-backend \\")
		fmt.Println("    core \\")
		fmt.Println("    0.0.1 \\")
		fmt.Println("    https://api.bintray.com \\")
		fmt.Println("    {username} \\")
		fmt.Println("    {secret}")
		fmt.Println("")
		os.Exit(1)
	}
}

func deleteFile(filepath string) error {
	url := bintrayURL + "/content/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
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

func publishFile(filepath string) error {
	url := bintrayURL + "/file_metadata/" + subject + "/" + repo + "/" + url.QueryEscape(filepath)
	fmt.Println("publish file ::", url, "...")
	buff := bytes.NewBuffer(nil)
	buff.WriteString(`{"list_in_downloads":true}`)
	req, err := http.NewRequest(http.MethodPut, url, buff)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
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
	command = os.Args[1]
	subject = os.Args[2]
	repo = os.Args[3]
	rpackage = os.Args[4]
	version = os.Args[5]
	bintrayURL = os.Args[6]
	username = os.Args[7]
	password = os.Args[8]
	fmt.Println("getting all files for ", subject, repo, rpackage, version, "...")
	client := &http.Client{}
	url := bintrayURL + "/packages/" + subject + "/" + repo + "/" + rpackage + "/versions/" + version + "/files?include_unpublished=1"
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
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	respObj := []interface{}{}
	err = json.Unmarshal(respJSON, &respObj)
	if err != nil {
		log.Fatal(err)
	}
	for _, binTrayFile := range respObj {
		bTFile := binTrayFile.(map[string]interface{})
		switch command {
		case "delete":
			err = deleteFile(bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		case "list":
			fmt.Println(bTFile["path"].(string))
		case "publish":
			err = publishFile(bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(errors.New("the given command " + command + " is not supported!"))
		}
	}
}
