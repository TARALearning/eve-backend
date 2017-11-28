package eve

import (
	"os"
	"os/exec"
	"testing"
)

func Test_NewSCmd(t *testing.T) {
	c := NewSCmd()
	if c.ID != "N/A" {
		t.Error("NewSCmd does not work as expected")
	}
}

func Test_NewCronCmd(t *testing.T) {
	c := NewCronCmd()
	if c.Running {
		t.Error("NewCronCmd does not work as expected")
	}
}

func Test_SchedulerDeleteCmd(t *testing.T) {
	s := NewScheduler()
	err := s.AppendCmd("test", "./test", "testuser", []string{})
	if err != nil {
		t.Error(err)
	}
	err = s.DeleteCmd("test")
	if err != nil {
		t.Error(err)
	}
	if len(s.Cmds) > 0 {
		t.Error("Scheduler DeleteCmd does not work as expected")
	}
}

func Test_SchedulerReplaceCmd(t *testing.T) {
	s := NewScheduler()
	err := s.AppendCmd("test", "./test", "testuser", []string{})
	if err != nil {
		t.Error(err)
	}
	err = s.ReplaceCmd("test", "./test2", "testuser2", []string{})
	if err != nil {
		t.Error(err)
	}
	if s.Cmds[0].ID != "test" {
		t.Error("Scheduler ReplaceCmd does not work as expected")
	}
}

func Test_SchedulerEnableProcess(t *testing.T) {
	s := NewScheduler()
	err := s.AppendCmd("test", "./test", "testuser", []string{})
	if err != nil {
		t.Error(err)
	}
	err = s.EnableProcess("test")
	if err != nil {
		t.Error(err)
	}
	if s.Cmds[0].ServiceType != "enabled" {
		t.Error("Scheduler EnableProcess does not work as expected")
	}
}

func Test_SchedulerDisableProcess(t *testing.T) {
	s := NewScheduler()
	err := s.AppendCmd("test", "./test", "testuser", []string{})
	if err != nil {
		t.Error(err)
	}
	err = s.DisableProcess("test")
	if err != nil {
		t.Error(err)
	}
	if s.Cmds[0].ServiceType != "disabled" {
		t.Error("Scheduler DisableProcess does not work as expected")
	}
}

// Test_SchedulerShutdown is used to test the scheduler shutdown
func Test_SchedulerShutdown(t *testing.T) {
	s := NewScheduler()
	gobin := "go"
	if os.Getenv("GOROOT") != "" {
		gobin = os.Getenv("GOROOT") + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "go"
	}
	out, err := exec.Command(gobin, "build", "-o", "tests/tmp/eve-test.exe", "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	err = s.AppendCmd("eve-test.exe", "./tests/tmp/eve-test.exe", "", []string{"1m"})
	if err != nil {
		t.Error(err)
	}
	// todo check this test here ... from time to time it does not work
	/*
		go func(s *Scheduler) {
			log.Println("wait 10 sec and then quit the scheduler")
			time.Sleep(10 * time.Second)
			log.Println(s.Shutdown())
			return
		}(s)
		err = s.Run()
		if err != nil {
			t.Error(err)
		}
	*/
}

// Test_CronScheduler is used to test the scheduler cron capabilities
func Test_CronScheduler(t *testing.T) {
	s := NewScheduler()
	err := s.AppendCronCmd("echo", []string{"test"}, "*", "*", "*", "*", "*")
	if err != nil {
		t.Error(err)
	}
	err = s.AppendCronCmd("echo", []string{"test 2"}, "*", "*", "*", "*", "*")
	if err != nil {
		t.Error(err)
	}
	err = s.RunCron()
	if err != nil {
		t.Error(err)
	}
}
