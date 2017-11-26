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

func Test_TokenMakeSig(t *testing.T) {
	res := tokenMakeSig("message", "secret")
	if res != "89134b490722b4cc16ed03a8793a8287e0991244f79f30a041e812a589945245" {
		t.Error("tokenMakeSig does not work as expected")
	}
}

func Test_TokenToString(t *testing.T) {
	token, err := TokenCreate("message", "123456789012345678901234567890ab", "secret")
	if err != nil {
		t.Error(err)
	}
	str, err := TokenToString(token, "123456789012345678901234567890ab")
	if err != nil {
		t.Error(err)
	}
	if str != "message#89134b490722b4cc16ed03a8793a8287e0991244f79f30a041e812a589945245" {
		t.Error("TokenToString does not work as expected")
	}
}

func Test_TokenCheckFailBase64(t *testing.T) {
	errToken := "errorToken"
	_, err := TokenCheck(errToken, "wrongSecret", "wrongSig")
	if err == nil {
		t.Error("TokenCheck does not work as expected")
	}
}

func Test_TokenCheckFailSignature(t *testing.T) {
	token, err := TokenCreate("message", "123456789012345678901234567890ab", "secret")
	if err != nil {
		t.Error(err)
	}
	_, err = TokenCheck(token, "123456789012345678901234567890ab", "wrongSig")
	if err == nil {
		t.Error("TokenCheck does not work as expected")
	}
}

func Test_TokenCreateFailSignature(t *testing.T) {
	_, err := TokenCreate("message", "wrongEncSecret", "secret")
	if err == nil {
		t.Error("TokenCreate does not work as expected")
	}
}
