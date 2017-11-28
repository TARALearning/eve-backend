package eve

import "testing"

func Test_PasswordLinuxCreate(t *testing.T) {
	pass, err := PasswordLinuxCreate("secret")
	if err != nil {
		t.Error(err)
	}
	if pass != "$6$SomeSaltSomeSalt$qXgrIF758PGlzo9woIBQWixNNdflHyVazP6pQiuLrNpG/afRWsGcJF7QH.Btz7ct9hdLl8.K9/ZWH4X45JYH1." {
		t.Log(pass)
		t.Error("PasswordLinuxCreate does not work as expected")
	}
}
