package eve

import "testing"

func Test_Sha1(t *testing.T) {
	SECRETSALT = ""
	hashed := Sha1("test")
	if hashed != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		t.Error("Sha1 does not work as expected")
	}
}
