package eve

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
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
	cmdStdoutReader, err := sleep.Cmd.StdoutPipe()
	if err != nil {
		t.Error(err)
	}
	sleep.Stdout = bufio.NewScanner(cmdStdoutReader)
	cmdStderrReader, err := sleep.Cmd.StderrPipe()
	if err != nil {
		t.Error(err)
	}
	sleep.Stderr = bufio.NewScanner(cmdStderrReader)
	s.AppendService(sleep)
	s.AppendServiceCmd("./tests/tmp/eve-test.exe", []string{"1m"})
	err = s.ServicesRun(&wg)
	if err != nil {
		t.Error(err)
	}
	// wait 10 seconds for the commands to start
	time.Sleep(time.Second * 10)
	// stop the command and all the routines running for the command
	err = s.Shutdown()
	if err != nil {
		t.Error(err)
	}
	// wait 5 sec for the commands to stop
	time.Sleep(time.Second * 5)
	s.CmdKillerQuitChannel <- true
	// wait 3 seconds for the killer routine to finish
	time.Sleep(time.Second * 3)
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

func Test_SchedulerServicesRestart(t *testing.T) {
	SetDebug(true)
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
	s.AppendServiceCmd("./tests/tmp/"+srv1, []string{"1m"})
	s.AppendServiceCmd("./tests/tmp/"+srv2, []string{"1m"})
	err = s.ServicesRun(&wg)
	if err != nil {
		t.Error(err)
	}
	// wait 10 seconds for the commands to start
	time.Sleep(time.Second * 10)
	// restart the srv1 executable
	err = s.ServiceStop(0)
	if err != nil {
		t.Error(err)
	}
	// wait 5 sec for the command to stops
	// time.Sleep(time.Second * 15)
	s.Cmds[0].ResetCmd("./tests/tmp/"+srv1, []string{"1m"})
	s.Cmds[0].Enable()
	//  give the command 10 sec to reinitialize
	// time.Sleep(time.Second * 10)
	err = s.ServiceStart(0, &wg)
	if err != nil {
		t.Error(err)
	}
	// wait 10 sec for the command to start again
	time.Sleep(time.Second * 10)
	// stop the command and all the routines running for the command
	err = s.Shutdown()
	if err != nil {
		t.Error(err)
	}
	// wait 5 sec for the commands to stop
	time.Sleep(time.Second * 5)
	s.CmdKillerQuitChannel <- true
	// wait 3 seconds for the killer routine to finish
	time.Sleep(time.Second * 3)
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

func Test_SchedulerCronJobs(t *testing.T) {
	startRoutines := runtime.NumGoroutine()
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
	wg.Wait()
	current := runtime.NumGoroutine()
	if startRoutines != current {
		t.Log(current)
		t.Error("SchedulerCronJobs does not work as expected")
	}
}
