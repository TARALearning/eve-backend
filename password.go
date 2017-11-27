package eve

import (
	"github.com/kless/osutil/user/crypt/sha512_crypt"
)

// PasswordLinuxCreate generates a linux password to be used for the user login
func PasswordLinuxCreate(password string) (string, error) {
	c := sha512_crypt.New()
	hash, err := c.Generate([]byte(password), []byte("$6$SomeSaltSomeSalt$"))
	if err != nil {
		return "", err
	}
	return hash, nil
}
