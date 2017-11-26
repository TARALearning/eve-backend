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
	evUserStorage       = ""
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

func init() {
	if len(os.Args) < 6 {
		fmt.Println("please specify all the required arguments to initialize the eve environment")
		fmt.Println("")
		fmt.Println("eve-setup {evbolt_url} {email} {password} {encryption_key} {signature_key}")
		fmt.Println("")
		fmt.Println("example:")
		fmt.Println("")
		fmt.Println("eve-setup \\")
		fmt.Println("    http://localhost:9092/" + eve.VERSION + "/eve/evbolt \\")
		fmt.Println("    francisc.simon@evalgo.org \\")
		fmt.Println("    secret \\")
		fmt.Println("    123456789012345678901234567890ab \\")
		fmt.Println("    secret.sig.key")
		fmt.Println("")
		os.Exit(1)
	}
}

func main() {
	evUserStorage = os.Args[1]
	evSecretStorage = os.Args[1]
	id := eve.Sha1(os.Args[2])
	message := eve.Sha1(os.Args[3])
	evSecretEncKeyValue = os.Args[4]
	evSecretSigKeyValue = os.Args[5]

	resp, err := eve.EVHttpSendForm(http.MethodPost, evUserStorage, url.Values{"database": {evUserStorageDB}, "bucket": {evUserStorageBucket}, "key": {id}, "message": {message}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the user id is :: ", string(rB))

	resp, err = eve.EVHttpSendForm(http.MethodPost, evSecretStorage, url.Values{"database": {evSecretStorageDB}, "bucket": {evSecretStorageBucket}, "key": {evSecretEncKey}, "message": {evSecretEncKeyValue}, "evbolt.msgtype": {"string"}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	rB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the id for the encryption secret is :: ", string(rB))

	resp, err = eve.EVHttpSendForm(http.MethodPost, evSecretStorage, url.Values{"database": {evSecretStorageDB}, "bucket": {evSecretStorageBucket}, "key": {evSecretSigKey}, "message": {evSecretSigKeyValue}, "evbolt.msgtype": {"string"}})
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
