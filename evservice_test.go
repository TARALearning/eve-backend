package eve

import "testing"

func Test_NewEVServiceConfig(t *testing.T) {
	srv := NewEVServiceConfigObj()
	srv.NewEVServiceConfig("rest")
	if srv.Config.Main != "EVREST" {
		t.Error("NewEVServiceConfig does not work as expected")
	}
	t.Log(srv)
}

func Test_NewEVServiceConfigRestAll(t *testing.T) {
	srv := NewEVServiceConfigObj()
	srv.NewEVServiceConfig("rest_all")
	if !srv.Config.Vars["USE_EVLOG"].(bool) {
		t.Error("NewEVServiceConfigRestAll does not work as expected")
	}
	t.Log(srv)
}

func Test_EVServiceConfiguration(t *testing.T) {
	srv := NewEVServiceConfigObj()
	conf := srv.EVServiceConfiguration()
	if conf.Main != "TestMain" {
		t.Error("EVServiceConfiguration does not work as expected")
	}
}

func Test_SrvConfigXXX(t *testing.T) {
	if SrvConfigMain() != "EVREST" {
		t.Error("SrvConfigMain does not work as expected")
	}
	if SrvConfigCommands()[0].Name != "help" {
		t.Error("SrvConfigCommands does not work as expected")
	}
	if len(SrvConfigTemplates()) == 0 {
		t.Error("SrvConfigTemplates does not work as expected")
	}
}

func Test_SetDefaultCType(t *testing.T) {
	cType := defaultCType
	SetDefaultCType("testvalue")
	if defaultCType != "testvalue" {
		t.Error("SetDefaultCType does not work as expected")
	}
	defaultCType = cType
}
