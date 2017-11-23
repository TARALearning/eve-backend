package eve

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// EVHttpNewClient creates a http client with the default values
func EVHttpNewClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{
		Transport: tr,
	}
}

// EVHttpNewClientCrt creates a http client with certificate authentification enabled with the given cert and key
func EVHttpNewClientCrt(crt, key string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{
		Transport: tr,
	}, nil
}

// EVHttpSendForm sends form data with a given method and values to the given url
func EVHttpSendForm(method string, url string, data url.Values) (*http.Response, error) {
	client := EVHttpNewClient()
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.Do(req)
}

// EVHttpSendText sends a text with the given method url and values
func EVHttpSendText(method string, url string, text string) (*http.Response, error) {
	client := EVHttpNewClient()
	req, err := http.NewRequest(method, url, strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "plain/text")
	return client.Do(req)
}

// EVHttpReceive is mostly used for all types of GET requests like DELETE etc.
func EVHttpReceive(method string, url string, values url.Values) (*http.Response, error) {
	urlValues := ""
	if values != nil {
		urlValues = "?" + values.Encode()
	}
	client := EVHttpNewClient()
	req, err := http.NewRequest(method, url+urlValues, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// ResponseBodyAll returns []byte and error from response body
func ResponseBodyAll(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

// RequestBodyAll returns []byte and error from the request body
func RequestBodyAll(request *http.Request) ([]byte, error) {
	defer request.Body.Close()
	return ioutil.ReadAll(request.Body)
}

// CheckRequiredFormValues checks if the required form values are available in the request object
func CheckRequiredFormValues(r *http.Request, values map[string]bool) error {
	for val, required := range values {
		if required {
			if r.FormValue(val) == "" {
				return errors.New("CheckRequiredFormValues the required value" + val + " was not provided!")
			}
		}
	}
	return nil
}

// ReturnErrorMessage generates a error message and write it into the response object
func ReturnErrorMessage(w http.ResponseWriter, statusCode int, err error, format string) {
	w.WriteHeader(statusCode)
	switch format {
	case ".json":
		w.Write([]byte(`{"Response":{"Status":"failed","StatusCode":` + strconv.Itoa(statusCode) + `,"Failed":true,"Message":"` + err.Error() + `"}}`))
	case ".html":
		w.Write([]byte(`<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">failed</div><div itemprop="StatusCode">` + strconv.Itoa(statusCode) + `</div><div itemprop="Failed">true</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">` + err.Error() + `</div></div></div>`))
	default:
		w.Write([]byte(err.Error()))
	}
}

// ReturnResult generates response message depending on given format and writes it into the response object
func ReturnResult(w http.ResponseWriter, statusCode int, response, format string) {
	w.WriteHeader(statusCode)
	var resultString string
	switch format {
	case ".json":
		w.Header().Set("Content-Type", "application/json")
		resultString = `{"Response":{"Status":"success","StatusCode":` + strconv.Itoa(statusCode) + `,"Failed":false,"Message":"` + response + `"}}`
	case ".html":
		w.Header().Set("Content-Type", "text/html")
		resultString = `<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">success</div><div itemprop="StatusCode">` + strconv.Itoa(statusCode) + `</div><div itemprop="Failed">false</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">` + response + `</div></div></div>`
	default:
		w.Header().Set("Content-Type", "text/plain")
		resultString = response
	}
	w.Write([]byte(resultString))
}

// DefineKeyType sets key type to predefined values depending on user input
func DefineKeyType(keyType string) string {
	switch keyType {
	case "c", "custom":
		keyType = evBoltKeyTypeCustom
	default:
		keyType = evBoltKeyTypeAuto
	}
	return keyType
}

// DecodeMessage decodes message depending on message type
func DecodeMessage(msg, msgType string) (string, error) {
	var message string
	switch msgType {
	case "base64":
		bMsg, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			return "", err
		}
		message = string(bMsg)
	default:
		message = msg
	}
	return message, nil
}

// EVHttpSendFile sends a file or with other words it is a file upload
func EVHttpSendFile(uri, filename, filepath string) (*http.Response, error) {
	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)
	fileWriter, err := bodyWriter.CreateFormFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	fh, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	client := EVHttpNewClient()
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return client.Do(req)
}
