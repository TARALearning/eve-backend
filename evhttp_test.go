package eve

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_EvHTTPNewClient(t *testing.T) {
	c := EvHTTPNewClient()
	if !c.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("EvHTTPNewClient does not work as expected")
	}
}

func Test_EvHTTPNewClientCrt(t *testing.T) {
	c, err := EvHTTPNewClientCrt("tests/test.client.crt", "tests/test.client.key")
	if err != nil {
		t.Error(err)
	}
	if !c.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("EvHTTPNewClient does not work as expected")
	}
}

func Test_EvHTTPNewClientCrtFail(t *testing.T) {
	_, err := EvHTTPNewClientCrt("crt", "key")
	if err == nil {
		t.Error("EvHTTPNewClientCrt does not work as expected")
	}
}

func Test_EvHTTPSendForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.FormValue("testkey"))
	}))
	defer ts.Close()
	res, err := EvHTTPSendForm(http.MethodPost, ts.URL, url.Values{"testkey": []string{"testvalue"}})
	if err != nil {
		t.Error(err)
	}
	testvalue, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(testvalue), "\n ") != "testvalue" {
		t.Error("EvHTTPSendForm does not work as expected")
	}
}

func Test_EvHTTPSendText(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testvalue, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(testvalue)
	}))
	defer ts.Close()
	res, err := EvHTTPSendText(http.MethodPost, ts.URL, "testvalue")
	if err != nil {
		t.Error(err)
	}
	testvalue, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(testvalue), "\n ") != "testvalue" {
		t.Error("EvHTTPSendText does not work as expected")
	}
}

func Test_EvHTTPReceive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("testvalue"))
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	testvalue, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(testvalue), "\n ") != "testvalue" {
		t.Error("EvHTTPReceive does not work as expected")
	}
}

func Test_ResponseBodyAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("testvalue1 testvalue2 testvalue3"))
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	testvalue, err := ResponseBodyAll(res)
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(testvalue), "\n ") != "testvalue1 testvalue2 testvalue3" {
		t.Error("ResponseBodyAll does not work as expected")
	}
}

func Test_RequestBodyAll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := RequestBodyAll(r)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()
	res, err := EvHTTPSendText(http.MethodPost, ts.URL, "testvalue")
	if err != nil {
		t.Error(err)
	}
	ok, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(ok), "\n ") != "OK" {
		t.Error("RequestBodyAll does not work as expected")
	}
}

func Test_CheckRequiredFormValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := CheckRequiredFormValues(r, map[string]bool{"testvalue": true})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()
	res, err := EvHTTPSendForm(http.MethodPost, ts.URL, url.Values{"testvalue": []string{"testvalue"}})
	if err != nil {
		t.Error(err)
	}
	ok, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(ok), "\n ") != "OK" {
		t.Error("CheckRequiredFormValues does not work as expected")
	}
}

func Test_ReturnErrorMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnErrorMessage(w, 500, errors.New("testerror"), ".txt")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	okErr, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(okErr), "\n ") != "testerror" {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_ReturnErrorMessageJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnErrorMessage(w, 500, errors.New("testerror"), ".json")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	okErr, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(okErr))
	if strings.Trim(string(okErr), "\n ") != `{"Response":{"Status":"failed","StatusCode":500,"Failed":true,"Message":"testerror"}}` {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_ReturnErrorMessageHtml(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnErrorMessage(w, 500, errors.New("testerror"), ".html")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	okErr, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(okErr), "\n") != `<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">failed</div><div itemprop="StatusCode">500</div><div itemprop="Failed">true</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">testerror</div></div></div>` {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_ReturnResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnResult(w, 200, "testvalue", ".txt")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(resp), "\n") != `testvalue` {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_ReturnResultJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnResult(w, 200, "testvalue", ".json")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(resp), "\n") != `{"Response":{"Status":"success","StatusCode":200,"Failed":false,"Message":"testvalue"}}` {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_ReturnResultHtml(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnResult(w, 200, "testvalue", ".html")
	}))
	defer ts.Close()
	res, err := EvHTTPReceive(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(resp), "\n") != `<div itemscope="" itemtype="https://evalgo.org/schema/Response"><div itemprop="Status">success</div><div itemprop="StatusCode">200</div><div itemprop="Failed">false</div><div itemprop="Message" itemscope="" itemtype="https://evalgo.org/schema/Message"><div itemprop="Content">testvalue</div></div></div>` {
		t.Error("ReturnErrorMessage does not work as expected")
	}
}

func Test_DefineKeyType(t *testing.T) {
	eyType := DefineKeyType("custom")
	if eyType != EvBoltKeyTypeCustom {
		t.Error("DefineKeyType does not work as expected")
	}
	eyType = DefineKeyType("default")
	if eyType != EvBoltKeyTypeAuto {
		t.Error("DefineKeyType does not work as expected")
	}
}

func Test_DecodeMessage(t *testing.T) {
	txtMessage := "message"
	dMessage, err := DecodeMessage(txtMessage, "text")
	if err != nil {
		t.Error(err)
	}
	if txtMessage != dMessage {
		t.Error("DecodeMessage does not work as expected")
	}
	b64Msg := base64.StdEncoding.EncodeToString([]byte(txtMessage))
	eMessage, err := DecodeMessage(b64Msg, "base64")
	if err != nil {
		t.Error(err)
	}
	if txtMessage != eMessage {
		t.Error("DecodeMessage does not work as expected")
	}
}

func Test_EvHTTPSendFile(t *testing.T) {
	testfile := "tests/tmp/upload.file.txt"
	uploadkey := "fileupload"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, header, err := r.FormFile(uploadkey)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(header.Filename))
	}))
	defer ts.Close()
	err := ioutil.WriteFile(testfile, []byte("content"), 0777)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(testfile)
	res, err := EvHTTPSendFile(ts.URL, uploadkey, testfile)
	if err != nil {
		t.Error(err)
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(string(resp), "\n") != testfile {
		t.Log(string(resp))
		t.Error("EvHTTPSendFile does not work as expected")
	}
}
