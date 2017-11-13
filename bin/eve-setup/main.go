package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"evalgo.org/eve"
)

var (
	EVUserStorage       = ""
	EVUserStorageDB     = "users.db"
	EVUserStorageBucket = "users"

	EVSecretStorage       = ""
	EVSecretStorageDB     = "secrets.db"
	EVSecretStorageBucket = "secrets"

	EVSecretEncKey      = "TokenKeyEnc"
	EVSecretEncKeyValue = ""
	EVSecretSigKey      = "TokenKeySig"
	EVSecretSigKeyValue = ""
)

func main() {

	if len(os.Args) < 6 {
		fmt.Println("please specify all the required arguments to initialize the eve environment")
		fmt.Println("")
		fmt.Println("eve-setup {evbolt_url} {email} {password} {encryption_key} {signature_key}")
		fmt.Println("")
		fmt.Println("example:")
		fmt.Println("")
		fmt.Println("eve-setup \\")
		fmt.Println("    http://localhost:9092/0.0.1/eve/evbolt \\")
		fmt.Println("    francisc.simon@evalgo.org \\")
		fmt.Println("    secret \\")
		fmt.Println("    123456789012345678901234567890ab \\")
		fmt.Println("    secret.sig.key")
		fmt.Println("")
		os.Exit(1)
	}

	EVUserStorage = os.Args[1]
	EVSecretStorage = os.Args[1]
	id := eve.Sha1(os.Args[2])
	message := eve.Sha1(os.Args[3])
	EVSecretEncKeyValue = os.Args[4]
	EVSecretSigKeyValue = os.Args[5]

	resp, err := eve.EVHttpSendForm(http.MethodPost, EVUserStorage, url.Values{"database": {EVUserStorageDB}, "bucket": {EVUserStorageBucket}, "key": {id}, "message": {message}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the user id is :: ", string(rB))

	resp, err = eve.EVHttpSendForm(http.MethodPost, EVSecretStorage, url.Values{"database": {EVSecretStorageDB}, "bucket": {EVSecretStorageBucket}, "key": {EVSecretEncKey}, "message": {EVSecretEncKeyValue}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the id for the encryption secret is :: ", string(rB))

	resp, err = eve.EVHttpSendForm(http.MethodPost, EVSecretStorage, url.Values{"database": {EVSecretStorageDB}, "bucket": {EVSecretStorageBucket}, "key": {EVSecretSigKey}, "message": {EVSecretSigKeyValue}, "evbolt.msgtype": {"string"}})
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
