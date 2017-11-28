package eve

import (
	"os"
	"strings"
	"testing"
)

func Test_EVBoltItob(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	testID := 5
	rBytes := EVBoltItob(testID)
	if testID != EVBoltBoti(rBytes) {
		t.Error("EVBoltItob and EVBoltBoti does not work as expected")
	}
	errorBytes := make([]byte, 10)
	for a := 0; a < 10; a++ {
		errorBytes = append(errorBytes, byte(a))
	}
	if EVBoltBoti(errorBytes) != -1 {
		t.Error("EVBoltBoti does not work as expected")
	}
}

func Test_EVBoltPut(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	res, err := EVBoltPut("key", "value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if res != "key" {
		t.Error("EVBoltPut does not wor as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltCustomPut(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	res, err := EVBoltCustomPut("key", "value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if res != "key" {
		t.Error("EVBoltCustomPut does not wor as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltAutoPut(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	res, err := EVBoltAutoPut("value", "testnr.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if res != "1" {
		t.Error("EVBoltAutoPut does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "testnr.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltCustomUpdate(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	value := "value"
	key, err := EVBoltCustomPut("key", value, "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	key, err = EVBoltCustomUpdate("value2", "test.db", "testbucket", key)
	if err != nil {
		t.Error(err)
	}
	value2, err := EVBoltCustomGet("test.db", "testbucket", key)
	if err != nil {
		t.Error(err)
	}
	if value2 != "value2" {
		t.Error("EVBoltCustomUpdate does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltAutoUpdate(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	key, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	key, err = EVBoltAutoUpdate("value2", "test.db", "testbucket", key)
	if err != nil {
		t.Error(err)
	}
	value2, err := EVBoltAutoGet("test.db", "testbucket", key)
	if err != nil {
		t.Error(err)
	}
	if value2 != "value2" {
		t.Error("EVBoltAutoUpdate does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltLast(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	value3, err := EVBoltLast("test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if value3 != "value3" {
		t.Error("EVBoltLast does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltFirst(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	value1, err := EVBoltFirst("test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if value1 != "value" {
		t.Error("EVBoltFirst does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

// todo fix this test
func Test_EVBoltAllHTML(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	htmlString, err := EVBoltAllHTML("test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	// if htmlString != `<div itemscope="" itemtype="https://schema.org/ItemList"><div itemprop="itemListElement" itemscope="" itemtype="http://schema.org/ListItem"><a itemprop="item" evboltkey="" href="http://localhost:9092/`+VERSION+`eve/evbolt.html?database=test.db&bucket=testbucket&key="><div itemprop="name">value2</div><meta itemprop="position" content="1"/></a></div><div itemprop="itemListElement" itemscope="" itemtype="http://schema.org/ListItem"><a itemprop="item" evboltkey="" href="http://localhost:9092/`+VERSION+`/eve/evbolt.html?database=test.db&bucket=testbucket&key="><div itemprop="name">value3</div><meta itemprop="position" content="2"/></a></div><div itemprop="itemListElement" itemscope="" itemtype="http://schema.org/ListItem"><a itemprop="item" evboltkey="" href="http://localhost:9092/`+VERSION+`/eve/evbolt.html?database=test.db&bucket=testbucket&key="><div itemprop="name">value</div><meta itemprop="position" content="3"/></a></div></div>` {
	// 	t.Log(htmlString)
	// 	t.Error("EVBoltAllHTML does not work as expected")
	// }
	t.Log(htmlString)
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

// todo fix thist test auto put keys need to be converted to int and then to string
func Test_EVBoltAllJSON(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	htmlJSON, err := EVBoltAllJSON("test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	// if htmlJson != `{"\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0001":"value","\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0002":"value2","\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0003":"value3"}` {
	// 	t.Log(htmlJson)
	// 	t.Error("EVBoltAllJSON does not work as expected")
	// }
	t.Log(htmlJSON)
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

// todo fix this test because it returns the key also but in the string it is not displayed
// check the bytes format for more information
func Test_EVBoltAllString(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	htmlString, err := EVBoltAllString("test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(htmlString, "\n ") != `value` {
		for _, ch := range []byte(htmlString) {
			t.Log(ch)
		}
		// t.Error("EVBoltAllString does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltAutoDelete(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltAutoPut("value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltAutoPut("value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	res, err := EVBoltAutoDelete("test.db", "testbucket", "2")
	if err != nil {
		t.Error(err)
	}
	if res != "2" {
		t.Error("EVBoltAutoDelete does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_EVBoltCustomDelete(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	_, err := EVBoltCustomPut("key1", "value", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltCustomPut("key2", "value2", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	_, err = EVBoltCustomPut("key3", "value3", "test.db", "testbucket")
	if err != nil {
		t.Error(err)
	}
	res, err := EVBoltCustomDelete("test.db", "testbucket", "key2")
	if err != nil {
		t.Error(err)
	}
	if res != "key2" {
		t.Error("EVBoltCustomDelete does not work as expected")
	}
	os.Remove(evBoltRoot + string(os.PathSeparator) + "test.db")
	os.Remove(evBoltRoot + string(os.PathSeparator) + "evbolt.meta.db")
}

func Test_SetEvBoltRoot(t *testing.T) {
	DEBUG = true
	defer func() { DEBUG = false }()
	old := evBoltRoot
	defer SetEvBoltRoot(old)
	SetEvBoltRoot("test")
	if evBoltRoot != "test" {
		t.Error("SetEvBoltRoot does not work as expected")
	}
}
