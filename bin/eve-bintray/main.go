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
	subject    = ""
	repo       = ""
	rpackage   = ""
	version    = ""
	bintrayURL = ""
	username   = ""
	password   = ""
)

func init() {
	if len(os.Args) < 7 {
		fmt.Println("please specify all the required arguments to run the bintray cleanup script")
		fmt.Println("")
		fmt.Println("eve-bintray {subject} {repo} {password} {package} {version} {api_url} {username} {token}")
		fmt.Println("")
		fmt.Println("example:")
		fmt.Println("")
		fmt.Println("eve-bintray \\")
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

func main() {
	subject = os.Args[1]
	repo = os.Args[2]
	rpackage = os.Args[3]
	version = os.Args[4]
	bintrayURL = os.Args[5]
	username = os.Args[6]
	password = os.Args[7]
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
		if bTFile["repo"].(string) == "eve-backend" {
			err = deleteFile(bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
