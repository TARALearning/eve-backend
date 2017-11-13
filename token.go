package eve

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"sort"
	"strings"
	"time"
)

type TokenSecretHandler interface {
	SecretGet() (encSecret string, sigSecret string, err error)
}

func tokenMakeSig(message, secret string) string {
	secretHash := md5.New()
	secretHash.Write([]byte(secret))
	key := secretHash.Sum(nil)
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

func tokenCheckSig(message, secret string) bool {
	mkey := strings.Split(message, "#")
	newToken := tokenMakeSig(mkey[0], secret)
	return hmac.Equal([]byte(newToken), []byte(mkey[1]))
}

func TokenIsExpired(base64Msg, encSecret, sigSecret string) (bool, error) {
	check := time.Now()
	// decrypt the base64 message into the encrypted message
	base64EncMsg, err := base64.StdEncoding.DecodeString(base64Msg)
	if err != nil {
		return false, err
	}
	// decrypt the encrypted message with the given key
	decMessage, err := tokenDecrypt([]byte(encSecret), base64EncMsg)
	if err != nil {
		return false, err
	}
	log.Println(string(decMessage))
	mkey := strings.Split(string(decMessage), "#")
	info := strings.Split(mkey[0], "\n")
	created := ""
	expire := ""
	for k := range info {
		infoChars := []byte(info[k])
		if infoChars[0] == 'c' {
			created = string(infoChars[2:])
		}
		if infoChars[0] == 'e' {
			expire = string(infoChars[2:])
		}
	}
	log.Println(created, expire)
	log.Println(len(created))
	start, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", created)
	if err != nil {
		return true, err
	}
	end, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", expire)
	if err != nil {
		return true, err
	}
	inbetween := check.After(start) && check.Before(end)
	if inbetween {
		return false, nil
	}
	return true, nil
}

func tokenEncrypt(key, text []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(string(text)))
	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return ciphertext, nil
}

func tokenDecrypt(key, ciphertext []byte) (plaintext []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	plaintext = ciphertext
	return plaintext, nil
}

func TokenCreate(message, encSecret, sigSecret string) (string, error) {
	message = strings.Replace(message, "&#43;", "+", -1)
	// create hmac message signature
	sHmacSignature := tokenMakeSig(message, sigSecret)
	// encrypt the message and the hmac signature with the given key
	encMessage, err := tokenEncrypt([]byte(encSecret), []byte(message+"#"+sHmacSignature))
	if err != nil {
		return "", err
	}
	// encrypt the encrypted message with base64
	return base64.StdEncoding.EncodeToString(encMessage), nil
}

func TokenToString(base64Msg, encSecret string) (string, error) {
	// decrypt the base64 message into the encrypted message
	base64EncMsg, err := base64.StdEncoding.DecodeString(base64Msg)
	if err != nil {
		return "", err
	}
	// decrypt the encrypted message with the given key
	decMessage, err := tokenDecrypt([]byte(encSecret), base64EncMsg)
	if err != nil {
		return "", err
	}
	return strings.Replace(string(decMessage), "&#43;", "+", -1), nil
}

func TokenCheck(base64Msg, encSecret, sigSecret string) (bool, error) {
	nDecMessage, err := TokenToString(base64Msg, encSecret)
	if err != nil {
		return false, err
	}
	// check if the decrypted message shmac signature matches
	if !tokenCheckSig(nDecMessage, sigSecret) {
		return false, errors.New("decrypted message was not signed with the given signature in the message")
	}
	return true, nil
}

func PlainToken(tokenProps map[string]string) (string, error) {
	mapKeys := make([]string, 0)
	for k, _ := range tokenProps {
		mapKeys = append(mapKeys, k)
	}
	sort.Strings(mapKeys)
	plainToken := ""
	for id, key := range mapKeys {
		if (len(mapKeys) - 1) == id {
			plainToken += key + " " + tokenProps[key]
		} else {
			plainToken += key + " " + tokenProps[key] + "\n"
		}
	}
	return plainToken, nil
}
