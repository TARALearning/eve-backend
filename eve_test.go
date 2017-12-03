package eve

import "testing"

func Test_EveGenEvSchedule(t *testing.T) {
	filePath := "tests/tmp/evschedule_main.go"
	fpath, err := GenEvSchedule("", Uses{}, filePath)
	if err != nil {
		t.Error(err)
	}
	if fpath != filePath {
		t.Error("GenEvSchedule does not work as expected")
	}
}

func Test_EveGenEvBolt(t *testing.T) {
	filePath := "tests/tmp/evbolt_main.go"
	fpath, err := GenEvBolt("", Uses{}, filePath)
	if err != nil {
		t.Error(err)
	}
	if fpath != filePath {
		t.Error("GenEvBolt does not work as expected")
	}
}

func Test_EveGenEvLog(t *testing.T) {
	filePath := "tests/tmp/evlog_main.go"
	fpath, err := GenEvLog("", Uses{}, filePath)
	if err != nil {
		t.Error(err)
	}
	if fpath != filePath {
		t.Error("GenEvLog does not work as expected")
	}
}

func Test_EveGenEvAuth(t *testing.T) {
	filePath := "tests/tmp/evauth_main.go"
	fpath, err := GenEvAuth("", Uses{}, filePath)
	if err != nil {
		t.Error(err)
	}
	if fpath != filePath {
		t.Error("GenEvAuth does not work as expected")
	}
}
