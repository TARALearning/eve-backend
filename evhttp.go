package eve

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func EVHttpSendForm(method string, url string, data url.Values) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.Do(req)
}

func EVHttpSendText(method string, url string, text string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "plain/text")
	return client.Do(req)
}

func EVHttpReceive(method string, url string, values url.Values) (*http.Response, error) {
	urlValues := ""
	if values != nil {
		urlValues = "?" + values.Encode()
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url+urlValues, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func ResponseBodyAll(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

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

// returnErrorMessage generates a error message and write it into the response object
func ReturnErrorMessage(w http.ResponseWriter, statusCode int, err error, format string) {
	w.WriteHeader(statusCode)
	switch format {
	case "json":
		w.Write([]byte(`{"Response":{"Status":"failed","StatusCode":` + strconv.Itoa(statusCode) + `,"Failed":true,"Message":"` + err.Error() + `"}}`))
	case "html":
		w.Write([]byte(`<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">failed</div><div itemprop="StatusCode">` + strconv.Itoa(statusCode) + `</div><div itemprop="Failed">true</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">` + err.Error() + `</div></div></div>`))
	default:
		w.Write([]byte(err.Error()))
	}
}

// returnResult generates response message depending on given format and writes it into the response object
func ReturnResult(w http.ResponseWriter, statusCode int, response, format string) {
	w.WriteHeader(statusCode)
	resultString := ""
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		resultString = `{"Response":{"Status":"success","StatusCode":` + strconv.Itoa(statusCode) + `,"Failed":false,"Message":"` + response + `"}}`
	case "html":
		w.Header().Set("Content-Type", "text/html")
		resultString = `<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">success</div><div itemprop="StatusCode">` + strconv.Itoa(statusCode) + `</div><div itemprop="Failed">false</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">` + response + `</div></div></div>`
	default:
		w.Header().Set("Content-Type", "text/plain")
		resultString = response
	}
	w.Write([]byte(resultString))
}

// defineKeyType sets key type to predefined values depending on user input
func DefineKeyType(keyType string) string {
	switch keyType {
	case "c", "custom":
		keyType = EVBOLT_KEY_TYPE_CUSTOM
	default:
		keyType = EVBOLT_KEY_TYPE_AUTO
	}
	return keyType
}

// decodeMessage decodes message depending on message type
func DecodeMessage(msg, msgType string) (string, error) {
	message := ""
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
