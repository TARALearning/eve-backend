package eve

import "testing"

func Test_NewEVServiceFlagDebug(t *testing.T) {
	debug := NewEVServiceFlagDebug()
	if debug.FName != "debug" {
		t.Error("NewEVServiceFlagDebug does not work as expected")
	}
}

func Test_NewEVServiceFlagHelpHTTP(t *testing.T) {
	hhttp := NewEVServiceFlagHelpHTTP()
	if hhttp.FName != "hhttp" {
		t.Error("NewEVServiceFlagHelpHTTP does not work as expected")
	}
}

func Test_NewEVServiceFlagHTTPAddress(t *testing.T) {
	address := NewEVServiceFlagHTTPAddress()
	if address.FName != "address" {
		t.Error("NewEVServiceFlagHTTPAddress does not work as expected")
	}
}

func Test_NewEVServiceFlagVersion(t *testing.T) {
	version := NewEVServiceFlagVersion()
	if version.FName != "version" {
		t.Error("NewEVServiceFlagVersion does not work as expected")
	}
}

func Test_NewEVServiceFlagSslCrt(t *testing.T) {
	crt := NewEVServiceFlagSslCrt()
	if crt.FName != "crt" {
		t.Error("NewEVServiceFlagSslCrt does not work as expected")
	}
}

func Test_NewEVServiceFlagSslKey(t *testing.T) {
	key := NewEVServiceFlagSslKey()
	if key.FName != "key" {
		t.Error("NewEVServiceFlagSslKey does not work as expected")
	}
}

func Test_NewEVServiceFlagWebRoot(t *testing.T) {
	webroot := NewEVServiceFlagWebRoot()
	if webroot.FName != "webroot" {
		t.Error("NewEVServiceFlagWebRoot does not work as expected")
	}
}
