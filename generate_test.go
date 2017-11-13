package eve

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func Test_GenMain(t *testing.T) {
	res, err := EVGenMain(ConfigForTest("default"))
	if string(res) != "TestValue" {
		t.Error("the expected result <TestValue> does not match <" + string(res) + ">")
	}
	if err != nil {
		t.Error(err)
	}
}

func Test_GenRestMain(t *testing.T) {
	srv := ConfigForTest("rest")
	srvJson, err := json.Marshal(srv)
	if err != nil {
		t.Error(err)
	}
	t.Log(srv)
	err = ioutil.WriteFile("tests/tmp/main.json", srvJson, 0777)
	if err != nil {
		t.Error(err)
	}
	res, err := EVGenMain(srv)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("tests/tmp/main.go", res, 0777)
	if err != nil {
		t.Error(err)
	}
	os.Remove("tests/tmp/main.json")
	os.Remove("tests/tmp/main.go")
}

func Test_GenScheduleMain(t *testing.T) {
	srv := ConfigForTest("schedule")
	srvJson, err := json.Marshal(srv)
	if err != nil {
		t.Error(err)
	}
	t.Log(srv)
	err = ioutil.WriteFile("tests/tmp/main.json", srvJson, 0777)
	if err != nil {
		t.Error(err)
	}
	res, err := EVGenMain(srv)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("tests/tmp/main.go", res, 0777)
	if err != nil {
		t.Error(err)
	}
	// os.Remove("tests/tmp/main.json")
	// os.Remove("tests/tmp/main.go")
}
