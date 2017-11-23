package eve

import (
	"net/http"
	"testing"
)

func Test_EVHttpNewClient(t *testing.T) {
	c := EVHttpNewClient()
	if !c.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("EVHttpNewClient does not work as expected")
	}
}

func Test_EVHttpNewClientCrt(t *testing.T) {
	c, err := EVHttpNewClientCrt("tests/test.client.crt", "tests/test.client.key")
	if err != nil {
		t.Error(err)
	}
	if !c.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("EVHttpNewClient does not work as expected")
	}
}
