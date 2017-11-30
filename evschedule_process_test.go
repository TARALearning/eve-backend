package eve

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func Test_NewSCmd(t *testing.T) {
	debug = true
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

func Test_NewScheduler(t *testing.T) {
	t.Log(NewScheduler())
}

func Test_AppendServiceCmd(t *testing.T) {
	s := NewScheduler()
	s.AppendServiceCmd("sleep", []string{"10"})
	if len(s.Cmds) < 1 {
		t.Error("Scheduler.AppendServiceCmd does not work as expected")
	}
}

func Test_SchedulerServicesRun(t *testing.T) {
	gobin := "go"
	if os.Getenv("GOROOT") != "" {
		gobin = os.Getenv("GOROOT") + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "go"
	}
	out, err := exec.Command(gobin, "build", "-o", "tests/tmp/eve-test.exe", "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("tests/tmp/eve-test.exe")
	t.Log(out)
	startRoutines := runtime.NumGoroutine()
	s := NewScheduler()
	sleep := NewSCmd()
	sleep.Enabled = false
	sleep.ID = "sleep"
	sleep.Cmd = exec.Command("sleep", "60")
	s.AppendService(sleep)
	s.AppendServiceCmd("./tests/tmp/eve-test.exe", []string{"1m"})
	err = s.ServicesRun(&wg)
	if err != nil {
		t.Error(err)
	}
	// stop the command and all the routines running for the command
	err = s.Shutdown()
	if err != nil {
		t.Error(err)
	}
	s.CmdKillerQuitChannel <- true
	// wait for all routines to finish
	wg.Wait()
	// close all channels
	s.ServicesStopChannels()
	close(s.CmdKillerQuitChannel)
	current := runtime.NumGoroutine()
	if startRoutines != current {
		t.Log(current)
		t.Error("SchedulerServicesRun does not work as expected")
	}
}

func Test_EnableDisableCmd(t *testing.T) {
	cmd := NewSCmd()
	cmd.Enable()
	cmd.Disable()
	if cmd.Enabled {
		t.Error("Disable does not work as expected")
	}
}

func Test_ResetCmd(t *testing.T) {
	cmd := NewSCmd()
	cmd.Enable()
	cmd.Running = true
	cmd.ResetCmd("test", []string{})
	if cmd.Enabled && cmd.Running {
		t.Error("ResetCmd does not work as expected")
	}
}

func Test_SchedulerServicesRestart(t *testing.T) {
	// SetDebug(true)
	startRoutines := runtime.NumGoroutine()
	srv1 := "evtest.exe"
	srv2 := "evtest-2.exe"
	defer os.Remove("tests/tmp/" + srv1)
	defer os.Remove("tests/tmp/" + srv2)
	gobin := "go"
	if os.Getenv("GOROOT") != "" {
		gobin = os.Getenv("GOROOT") + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "go"
	}
	out, err := exec.Command(gobin, "build", "-o", "tests/tmp/"+srv1, "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	out, err = exec.Command(gobin, "build", "-o", "tests/tmp/"+srv2, "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	s := NewScheduler()
	s.AppendServiceCmd("./tests/tmp/"+srv1, []string{"2m"})
	s.AppendServiceCmd("./tests/tmp/"+srv2, []string{"2m"})
	err = s.ServicesRun(&wg)
	if err != nil {
		t.Error(err)
	}
	err = s.ServiceRestart(0, "./tests/tmp/"+srv1, []string{"1m"})
	if err != nil {
		t.Error(err)
	}
	// stop the command and all the routines running for the command
	fmt.Println("shutdown all services...")
	err = s.Shutdown()
	if err != nil {
		t.Error(err)
	}
	s.CmdKillerQuitChannel <- true
	// wait for all routines to finish
	wg.Wait()
	// close all channels
	s.ServicesStopChannels()
	close(s.CmdKillerQuitChannel)
	current := runtime.NumGoroutine()
	if startRoutines != current {
		t.Log(current)
		t.Error("SchedulerServicesRestart does not work as expected")
	}
}

func Test_WaitGroup(t *testing.T) {
	s := NewScheduler()
	if s.WaitGroup() == nil {
		t.Error("WaitGroup does not work as expected")
	}
}

func Test_ServicesStop(t *testing.T) {
	srv1 := "evtest.exe"
	srv2 := "evtest-2.exe"
	defer os.Remove("tests/tmp/" + srv1)
	defer os.Remove("tests/tmp/" + srv2)
	gobin := "go"
	if os.Getenv("GOROOT") != "" {
		gobin = os.Getenv("GOROOT") + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "go"
	}
	out, err := exec.Command(gobin, "build", "-o", "tests/tmp/"+srv1, "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	out, err = exec.Command(gobin, "build", "-o", "tests/tmp/"+srv2, "tests/evschedule/main.go").Output()
	if err != nil {
		t.Error(err)
	}
	t.Log(out)
	s := NewScheduler()
	s.AppendServiceCmd("./tests/tmp/"+srv1, []string{"2m"})
	s.AppendServiceCmd("./tests/tmp/"+srv2, []string{"2m"})
	err = s.ServicesRun(&wg)
	if err != nil {
		t.Error(err)
	}
	err = s.ServicesStop()
	if err != nil {
		t.Error(err)
	}
}
