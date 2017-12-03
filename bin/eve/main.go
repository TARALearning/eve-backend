package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"evalgo.org/eve"
)

var (
	flags      *flag.FlagSet
	eveVersion = eve.VERSION
	command    = ""
	subcommand = ""
	service    = ""
	use        eve.Uses
	debug      = false
	config     = ""
	target     = ""
	subject    = ""
	repo       = ""
	rpackage   = ""
	version    = ""
	URL        = ""
	username   = ""
	password   = ""
	encSecret  = ""
	sigSecret  = ""

	evUserStorageDB     = "users.db"
	evUserStorageBucket = "users"

	evSecretStorage       = ""
	evSecretStorageDB     = "secrets.db"
	evSecretStorageBucket = "secrets"

	evSecretEncKey      = "TokenKeyEnc"
	evSecretEncKeyValue = ""
	evSecretSigKey      = "TokenKeySig"
	evSecretSigKeyValue = ""
)

// GenUsage displays the help/usage instructions
func GenUsage() {
	fmt.Println(`
usage: 
    eve \
      {command} \
      {subcommand} \
      -{argument} {value} \
      -{argument} {value} ...

bintray example:
    eve \
      list \
      bintray \
      -subject evalgo \
      -repo eve-backend \
      -rpackage core /
      -version 0.0.1 \
      -url https://api.bintray.com \
      -username {username} \
      -password {secret}

generate example:
    eve \
      generate \
      golang \
      -service evauth \
      -use debug \
      -target evauth_main.go

setup example:
    eve \
      setup \
      evauth \
      -url evauth \
      -username debug \
      -password evauth_main.go \
      -ecrypt 123456789012345678901234567890ab \
      -sign secretSignature

commands:
    delete
    help
    generate
    list
    publish
    setup
    version

subcommand:
    bintray
    golang
    evauth

types:
    service

arguments:
    debug
    config
    target
    subject
    repo
    rpackage
    version
    url
    username
    password

services:
    evauth
    evbolt
    evlog
    evschedule

use flags:
    debug
	`)
}

func init() {
	if len(os.Args) <= 1 {
		GenUsage()
		os.Exit(2)
	}
	command = os.Args[1]
	flags = flag.NewFlagSet(command, flag.ExitOnError)
	flags.BoolVar(&debug, "debug", false, "turns debug mode on or off")
	flags.StringVar(&service, "service", "", "service to be generated")
	flags.StringVar(&config, "config", "", "file with a given configuration to be used to build the service")
	flags.StringVar(&target, "target", ".", "path to the directory where the fiel should be generated")
	flags.Var(&use, "use", "use a specific module/feature for the given service to be created")
	flags.StringVar(&subject, "subject", "", "")
	flags.StringVar(&repo, "repo", "", "")
	flags.StringVar(&rpackage, "rpackage", "", "")
	flags.StringVar(&version, "version", "", "")
	flags.StringVar(&URL, "url", "", "")
	flags.StringVar(&username, "username", "", "")
	flags.StringVar(&password, "password", "", "")
	flags.StringVar(&encSecret, "encrypt", "", "")
	flags.StringVar(&sigSecret, "sign", "", "")
	flags.Usage = GenUsage
}

func main() {
	// check it the user wants only info
	if len(os.Args) == 2 {
		flags.Parse(os.Args[2:])
		switch command {
		case "help":
			GenUsage()
			os.Exit(0)
		case "version":
			fmt.Println("eve-gen version:", eveVersion)
		default:
			fmt.Println("error: the given " + command + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
		return
	}
	// run the commands
	flags.Parse(os.Args[3:])
	subcommand = os.Args[2]
	if debug {
		eve.SetDebug(true)
	}
	switch command {
	case "generate":
		switch subcommand {
		case "golang":
			generate()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "delete":
		switch subcommand {
		case "bintray":
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "list":
		switch subcommand {
		case "bintray":
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "publish":
		switch subcommand {
		case "bintray":
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "setup":
		switch subcommand {
		case "evauth":
			setup()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	default:
		fmt.Println("error: the given command <" + command + "> is not supported yet")
		GenUsage()
		os.Exit(2)
	}
}

func generate() {
	fmt.Println("eve-gen :: check if service name was provided...")
	if service == "" {
		fmt.Println("error: plese specify a service name")
		GenUsage()
		os.Exit(2)
	}
	fmt.Println("eve-gen :: found service name <" + service + ">")
	var mainFile string
	var err error
	switch service {
	case "evauth":
		mainFile, err = eve.GenEvAuth(config, use, target)
	case "evbolt":
		mainFile, err = eve.GenEvBolt(config, use, target)
	case "evlog":
		mainFile, err = eve.GenEvLog(config, use, target)
	case "evschedule":
		mainFile, err = eve.GenEvSchedule(config, use, target)
	default:
		fmt.Println("error: the given service name <" + service + "> is not supported yet")
		GenUsage()
		os.Exit(2)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("eve-gen :: generated file " + mainFile + "...")
}

func setup() {
	resp, err := eve.EvHTTPSendForm(http.MethodPost, URL, url.Values{"database": {evUserStorageDB}, "bucket": {evUserStorageBucket}, "key": {username}, "message": {password}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the user id is :: ", string(rB))

	resp, err = eve.EvHTTPSendForm(http.MethodPost, URL, url.Values{"database": {evSecretStorageDB}, "bucket": {evSecretStorageBucket}, "key": {evSecretEncKey}, "message": {encSecret}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the id for the encryption secret is :: ", string(rB))

	resp, err = eve.EvHTTPSendForm(http.MethodPost, URL, url.Values{"database": {evSecretStorageDB}, "bucket": {evSecretStorageBucket}, "key": {evSecretSigKey}, "message": {sigSecret}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the id for the signature key is :: ", string(rB))

}

/*
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
*/

func bintray() {
	fmt.Println("getting all files for ", subject, repo, rpackage, version, "...")
	client := eve.EvHTTPNewClient()
	url := URL + "/packages/" + subject + "/" + repo + "/" + rpackage + "/versions/" + version + "/files?include_unpublished=1"
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
			err = eve.BinTrayDeleteFile(subject, repo, bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		case "list":
			fmt.Println(bTFile["path"].(string))
		case "publish":
			err = eve.BinTrayPublishFile(subject, repo, bTFile["path"].(string))
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(errors.New("the given command " + command + " is not supported!"))
		}
	}
}
