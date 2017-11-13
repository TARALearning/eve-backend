package eve

import "testing"

func Test_TokenEncryptDecrypt(t *testing.T) {
	tokenProps := map[string]string{"u": "1234567890", "c": "2009-11-10 23:00:00 +0000 UTC", "e": "2009-11-10 23:00:00 +0000 UTC", "t": "1234567890", "s": "1234567890"}
	message, err := PlainToken(tokenProps)
	if err != nil {
		t.Error(err)
	}
	key := "123456789012345678901234567890ab"
	secret := "secret_will_be_md5_hashed"
	base64Msg, err := TokenCreate(message, key, secret)
	if err != nil {
		t.Error(err)
	}
	ok, err := TokenCheck(base64Msg, key, secret)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("validation of the created token failed!")
	}
}

func Test_TokenIsExpired(t *testing.T) {
	tokenProps := map[string]string{"u": "1234567890", "c": "2009-11-10 23:00:00 +0000 UTC", "e": "2009-11-10 23:00:00 +0000 UTC", "t": "1234567890", "s": "1234567890"}
	message, err := PlainToken(tokenProps)
	if err != nil {
		t.Error(err)
	}
	key := "123456789012345678901234567890ab"
	secret := "secret_will_be_md5_hashed"
	base64Msg, err := TokenCreate(message, key, secret)
	if err != nil {
		t.Error(err)
	}
	ok, err := TokenIsExpired(base64Msg, key, secret)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("TokenIsExpired can not invalidate invalid token!")
	}
}
