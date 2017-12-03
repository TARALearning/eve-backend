package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

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
	frepo      = ""
	rpackage   = ""
	fpackage   = ""
	version    = ""
	fversion   = ""
	URL        = ""
	DURL       = ""
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
	package
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
    frepo
	rpackage
    fpackage
	version
	fversion
	url
    durl
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
	flags.StringVar(&frepo, "frepo", "", "")
	flags.StringVar(&rpackage, "rpackage", "", "")
	flags.StringVar(&fpackage, "fpackage", "", "")
	flags.StringVar(&version, "version", "", "")
	flags.StringVar(&fversion, "fversion", "", "")
	flags.StringVar(&URL, "url", "", "")
	flags.StringVar(&DURL, "durl", "", "")
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
			eve.SetBinTrayURL(URL)
			eve.SetBinTrayUser(username)
			eve.SetBinTrayPassword(password)
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "list":
		switch subcommand {
		case "bintray":
			eve.SetBinTrayURL(URL)
			eve.SetBinTrayUser(username)
			eve.SetBinTrayPassword(password)
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "package":
		switch subcommand {
		case "bintray":
			eve.SetBinTrayURL(URL)
			eve.SetBinTrayDownloadURL(DURL)
			eve.SetBinTrayUser(username)
			eve.SetBinTrayPassword(password)
			bintray()
		default:
			fmt.Println("error: the given " + command + "/" + subcommand + " is not supported yet")
			GenUsage()
			os.Exit(2)
		}
	case "publish":
		switch subcommand {
		case "bintray":
			eve.SetBinTrayURL(URL)
			eve.SetBinTrayUser(username)
			eve.SetBinTrayPassword(password)
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
	case "evweb":
		mainFile, err = eve.GenEvWeb(config, use, target)
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

func generateVersionFile(filepath string) error {
	content := bytes.NewBuffer([]byte("VERSIONS"))
	content.WriteString("\n\n")
	content.WriteString(subject + "::" + repo + "::" + rpackage + "::" + version + "\n")
	content.WriteString(subject + "::" + frepo + "::" + fpackage + "::" + fversion + "\n")
	content.WriteString("\n\n")
	fmt.Println("writing VERSIONS file to", filepath)
	return ioutil.WriteFile(filepath, content.Bytes(), 0777)
}

func generateReadmeFile(filepath string) error {
	content := bytes.NewBuffer([]byte("# eve"))
	content.WriteString("\n\n")
	content.WriteString("## start services")
	content.WriteString("\n\n")
	content.WriteString("./evlog http\n")
	content.WriteString("./evbolt http\n")
	content.WriteString("./evweb http -webroot webroot\n")
	content.WriteString("./evauth http\n\n")
	content.WriteString("## initialize services\n\n")
	content.WriteString("./eve setup evauth \\\n")
	content.WriteString("    -url http://127.0.0.1:9092/" + version + "/eve/evbolt\\\n")
	content.WriteString("    -username francisc.simon@evalgo.org \\\n")
	content.WriteString("    -password secret \\\n")
	content.WriteString("    -encrypt 123456789012345678901234567890ab\\\n")
	content.WriteString("    -sign signatureSecret")
	content.WriteString("\n\n")
	fmt.Println("writing README.md file to", filepath)
	return ioutil.WriteFile(filepath, content.Bytes(), 0777)
}

func getAllBintrayFiles(subject, repo, rpackage, version string, unpublished int) ([]string, error) {
	allFiles := []string{}
	fmt.Println("getting all files for ", subject, repo, rpackage, version, "...")
	client := eve.EvHTTPNewClient()
	url := URL + "/packages/" + subject + "/" + repo + "/" + rpackage + "/versions/" + version + "/files?include_unpublished=" + strconv.Itoa(unpublished)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("return status code is not 200 it is " + strconv.Itoa(resp.StatusCode))
	}
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	respObj := []interface{}{}
	err = json.Unmarshal(respJSON, &respObj)
	if err != nil {
		return nil, err
	}
	for _, binTrayFile := range respObj {
		bTFile := binTrayFile.(map[string]interface{})
		allFiles = append(allFiles, bTFile["path"].(string))
	}
	return allFiles, nil
}

func bintray() {
	var packageDarwin = []string{}
	var packageLinux = []string{}
	var packageWindows = []string{}
	unpublished := 1
	if command == "package" {
		unpublished = 0
	}

	// download the backend files
	allFiles, err := getAllBintrayFiles(subject, repo, rpackage, version, unpublished)
	if err != nil {
		log.Fatal(err)
	}
	for _, binTrayFile := range allFiles {
		switch command {
		case "delete":
			err = eve.BinTrayDeleteFile(subject, repo, binTrayFile)
			if err != nil {
				log.Fatal(err)
			}
		case "list":
			fmt.Println(binTrayFile)
		case "package":
			fname := path.Base(binTrayFile)
			switch true {
			case strings.Contains(fname, "darwin"):
				packageDarwin = append(packageDarwin, binTrayFile)
			case strings.Contains(fname, "linux"):
				packageLinux = append(packageLinux, binTrayFile)
			case strings.Contains(fname, "windows"):
				packageWindows = append(packageWindows, binTrayFile)
			default:
				fmt.Println("skip not needed file for packaging", fname)
			}
		case "publish":
			err = eve.BinTrayPublishFile(subject, repo, binTrayFile)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(errors.New("the given command " + command + " is not supported!"))
		}
	}
	if command == "package" {
		// donwload the frontend files
		webrootFiles, err := getAllBintrayFiles(subject, frepo, fpackage, fversion, unpublished)
		if err != nil {
			log.Fatal(err)
		}
		if len(packageDarwin) > 0 && len(packageLinux) > 0 && len(packageWindows) > 0 {
			fmt.Println("packaging requested files for darwin to given target " + target + "...")
			for _, pFile := range packageDarwin {
				content, err := eve.BinTrayDownloadFile(subject, repo, pFile)
				if err != nil {
					log.Fatal(err)
				}
				targetFile := target + string(os.PathSeparator) + "darwin" + string(os.PathSeparator) + strings.Replace(filepath.Base(pFile), "darwin-amd64-"+version+"_", "", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			err := generateVersionFile(target + string(os.PathSeparator) + "darwin" + string(os.PathSeparator) + "VERSIONS")
			if err != nil {
				log.Fatal(err)
			}
			err = generateReadmeFile(target + string(os.PathSeparator) + "darwin" + string(os.PathSeparator) + "README.md")
			if err != nil {
				log.Fatal(err)
			}
			for _, wFile := range webrootFiles {
				content, err := eve.BinTrayDownloadFile(subject, frepo, wFile)
				if err != nil {
					log.Fatal(err)
				}
				rootFolder := strings.Split(wFile, "/")
				targetFile := target + string(os.PathSeparator) + "darwin" + string(os.PathSeparator) + strings.Replace(wFile, rootFolder[0], "webroot", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println("packaging requested files for linux to given target " + target + "...")
			for _, pFile := range packageLinux {
				content, err := eve.BinTrayDownloadFile(subject, repo, pFile)
				if err != nil {
					log.Fatal(err)
				}
				targetFile := target + string(os.PathSeparator) + "linux" + string(os.PathSeparator) + strings.Replace(filepath.Base(pFile), "linux-amd64-"+version+"_", "", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = generateVersionFile(target + string(os.PathSeparator) + "linux" + string(os.PathSeparator) + "VERSIONS")
			if err != nil {
				log.Fatal(err)
			}
			err = generateReadmeFile(target + string(os.PathSeparator) + "linux" + string(os.PathSeparator) + "README.md")
			if err != nil {
				log.Fatal(err)
			}
			for _, wFile := range webrootFiles {
				content, err := eve.BinTrayDownloadFile(subject, frepo, wFile)
				if err != nil {
					log.Fatal(err)
				}
				rootFolder := strings.Split(wFile, "/")
				targetFile := target + string(os.PathSeparator) + "linux" + string(os.PathSeparator) + strings.Replace(wFile, rootFolder[0], "webroot", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println("zipping the linux version...")
			err = eve.Zip("linux", "linux")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("packaging requested files for windows to given target " + target + "...")
			for _, pFile := range packageWindows {
				content, err := eve.BinTrayDownloadFile(subject, repo, pFile)
				if err != nil {
					log.Fatal(err)
				}
				targetFile := target + string(os.PathSeparator) + "windows" + string(os.PathSeparator) + strings.Replace(filepath.Base(pFile), "windows-amd64-"+version+"_", "", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			err = generateVersionFile(target + string(os.PathSeparator) + "windows" + string(os.PathSeparator) + "VERSIONS")
			if err != nil {
				log.Fatal(err)
			}
			err = generateReadmeFile(target + string(os.PathSeparator) + "windows" + string(os.PathSeparator) + "README.md")
			if err != nil {
				log.Fatal(err)
			}
			for _, wFile := range webrootFiles {
				content, err := eve.BinTrayDownloadFile(subject, frepo, wFile)
				if err != nil {
					log.Fatal(err)
				}
				rootFolder := strings.Split(wFile, "/")
				targetFile := target + string(os.PathSeparator) + "windows" + string(os.PathSeparator) + strings.Replace(wFile, rootFolder[0], "webroot", 1)
				root := filepath.Dir(targetFile)
				err = os.MkdirAll(root, 0777)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("writing file", targetFile, "...")
				err = ioutil.WriteFile(targetFile, content, 0777)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println("zipping the linux version...")
			err = eve.Zip("windows", "windows")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
