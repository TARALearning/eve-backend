package eve

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

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
	go func(s *Scheduler) {
		log.Println("wait 5 sec and then quit the scheduler")
		time.Sleep(5 * time.Second)
		log.Println(s.Shutdown())
		return
	}(s)
	err = s.Run()
	if err != nil {
		t.Error(err)
	}
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
