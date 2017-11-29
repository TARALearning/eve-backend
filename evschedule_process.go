package eve

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// wg is that syncing waiting group variable for the goroutines
var wg sync.WaitGroup

// mutex is the variable which is needed for mutexing the running commands map
var mutex = &sync.Mutex{}

// SCmd is the struct which is used to store the executed command and the
// stdin stderr output related to the given command
// please be careful with the Id definition because it should be the exact name of
// the given command because it is used to check for the command if it is already running
// if there is another command with the same name it will kill it during bootstrap
type SCmd struct {
	Running           bool
	PID               int
	ID                string
	Cmd               *exec.Cmd
	Stdout            *bufio.Scanner
	Stderr            *bufio.Scanner
	CmdQuitChannel    chan bool
	StdoutQuitChannel chan bool
	StderrQuitChannel chan bool
	Enabled           bool
}

// NewSCmd will create a new SchedulerCommand object
func NewSCmd() *SCmd {
	c := new(SCmd)
	c.Running = false
	c.PID = 0
	c.ID = "N/A"
	c.Cmd = nil
	c.Stdout = nil
	c.Stderr = nil
	c.CmdQuitChannel = make(chan bool, 100)
	c.StdoutQuitChannel = make(chan bool, 100)
	c.StderrQuitChannel = make(chan bool, 100)
	c.Enabled = false
	return c
}

// CronCmd is the command struct to be used for cronjobs
type CronCmd struct {
	SCmd
	Finished      bool
	Running       bool
	LastMinutes   string
	Minutes       string
	LastHours     string
	Hours         string
	LastMonthDays string
	MonthDays     string
	LastMonths    string
	Months        string
	LastWeekDays  string
	WeekDays      string
}

// NewCronCmd creates a new cron command object
func NewCronCmd() *CronCmd {
	cc := new(CronCmd)
	cc.Finished = false
	cc.Running = false
	cc.Cmd = nil
	cc.Stdout = nil
	cc.Stderr = nil
	cc.LastMinutes = ""
	cc.Minutes = "*"
	cc.LastHours = ""
	cc.Hours = "*"
	cc.LastMonthDays = ""
	cc.MonthDays = "*"
	cc.LastMonths = ""
	cc.Months = "*"
	cc.LastWeekDays = ""
	cc.WeekDays = "*"
	return cc
}

// Enable the service to be run
func (srv *SCmd) Enable() {
	mutex.Lock()
	srv.Running = false
	srv.Enabled = true
	mutex.Unlock()
}

// Disable the service to be stopped or skiped
func (srv *SCmd) Disable() {
	mutex.Lock()
	srv.Running = true
	srv.Enabled = false
	mutex.Unlock()
}

// ResetCmd the service to be stopped or skiped
func (srv *SCmd) ResetCmd(cmdPath string, args []string) error {
	fmt.Println("reset", cmdPath, args)
	mutex.Lock()
	srv.Enabled = false
	srv.Running = false
	srv.ID = path.Base(cmdPath)
	srv.Cmd = exec.Command(cmdPath, args...)
	cmdStdoutReader, err := srv.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	srv.Stdout = bufio.NewScanner(cmdStdoutReader)
	cmdStderrReader, err := srv.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	srv.Stderr = bufio.NewScanner(cmdStderrReader)
	mutex.Unlock()
	return nil
}

// Scheduler is the scheduler main struct
type Scheduler struct {
	Cmds                  map[int]*SCmd
	CronCmds              []*CronCmd
	CmdQuitChannel        chan bool
	CmdKillerQuitChannel  chan bool
	CmdKillerChannel      chan int
	MonitoringQuitChannel chan bool
	MainQuitChannel       chan bool
}

// NewScheduler will create a new Scheduler object
func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.Cmds = make(map[int]*SCmd, 0)
	s.CronCmds = make([]*CronCmd, 0)
	s.CmdQuitChannel = make(chan bool, 100)
	s.CmdKillerQuitChannel = make(chan bool, 100)
	s.CmdKillerChannel = make(chan int, 100)
	s.MonitoringQuitChannel = make(chan bool, 100)
	s.MainQuitChannel = make(chan bool, 100)
	return s
}

// AppendService appends a service to the scheduler
func (s *Scheduler) AppendService(cmd *SCmd) error {
	ID := len(s.Cmds)
	s.Cmds[ID] = cmd
	return nil
}

// AppendServiceCmd creates a service object from the given command arguments
func (s *Scheduler) AppendServiceCmd(cmdPath string, cmdArgs []string) error {
	sCmd := NewSCmd()
	sCmd.Enabled = true
	sCmd.ID = path.Base(cmdPath)
	sCmd.Cmd = exec.Command(cmdPath, cmdArgs...)
	cmdStdoutReader, err := sCmd.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	sCmd.Stdout = bufio.NewScanner(cmdStdoutReader)
	cmdStderrReader, err := sCmd.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	sCmd.Stderr = bufio.NewScanner(cmdStderrReader)
	ID := len(s.Cmds)
	s.Cmds[ID] = sCmd
	// fmt.Println(sCmd.Cmd)
	return nil
}

// ServicesRun is starting all the services which have the flag srv.Enabled set to true
func (s *Scheduler) ServicesRun(syncGroup *sync.WaitGroup) error {
	// start the process killer
	syncGroup.Add(1)
	go func(s *Scheduler, wg *sync.WaitGroup) {
		defer syncGroup.Done()
		for {
			select {
			case srvID := <-s.CmdKillerChannel:
				// log.Println("---> kill", srvID)
				proc, err := os.FindProcess(srvID)
				if err != nil {
					fmt.Println("process killer find error", err)
				}
				err = proc.Kill()
				if err != nil {
					// seems like the process is already killed
					if err.Error() == "os: process already finished" {
						continue
					}
					fmt.Println("process killer kill error", err)
				}
			case <-s.CmdKillerQuitChannel:
				// fmt.Println("quit process killer")
				return
			default:
				// fmt.Println("process killer default state reached waiting 1 seconds")
				time.Sleep(time.Second * 1)
			}
		}
	}(s, syncGroup)
	for ID := range s.Cmds {
		err := s.ServiceStart(ID, syncGroup)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown is stopping all the enabled services and sets the srv.Enabled to false
func (s *Scheduler) Shutdown() error {
	for ID := range s.Cmds {
		err := s.ServiceStop(ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ServicesStopChannels closes all the services channels
func (s *Scheduler) ServicesStopChannels() {
	for _, srv := range s.Cmds {
		// fmt.Println("close all srv channels")
		// if srv.Enabled {
		close(srv.CmdQuitChannel)
		close(srv.StderrQuitChannel)
		close(srv.StdoutQuitChannel)
		// }
	}
}

// ServiceStop stops a running services
func (s *Scheduler) ServiceStop(ID int) error {
	mutex.Lock()
	srv := s.Cmds[ID]
	if !srv.Enabled {
		mutex.Unlock()
		return nil
	}
	fmt.Println("stopping service", srv.ID, "with PID", srv.PID)
	srv.Running = false
	srv.Enabled = false
	mutex.Unlock()
	srv.CmdQuitChannel <- true
	srv.StderrQuitChannel <- true
	srv.StdoutQuitChannel <- true
	mutex.Lock()
	s.CmdKillerChannel <- srv.PID
	mutex.Unlock()
	// fmt.Println("wait 5 sec for the process to stop")
	time.Sleep(time.Second * 5)
	return nil
}

// ServiceStart starts a running services
func (s *Scheduler) ServiceStart(ID int, syncGroup *sync.WaitGroup) error {
	mutex.Lock()
	srv := s.Cmds[ID]
	fmt.Println("running service", srv.ID, "...")
	if !srv.Enabled || srv.Running {
		fmt.Println("skip disabled or already running service", srv.ID)
		mutex.Unlock()
		return nil
	}
	mutex.Unlock()
	// start the error scanner
	syncGroup.Add(1)
	go func(srv *SCmd, wg *sync.WaitGroup) {
		// fmt.Println("start stderr scanner for", srv.ID)
		defer syncGroup.Done()
		for {
			select {
			case <-srv.StderrQuitChannel:
				mutex.Lock()
				fmt.Println("quit stderr for", srv.ID)
				mutex.Unlock()
				return
			default:
				if srv.Stderr.Scan() {
					mutex.Lock()
					fmt.Println(srv.ID, ":: stderr", srv.Stderr.Text())
					mutex.Unlock()
				}
				mutex.Lock()
				err := srv.Stdout.Err()
				mutex.Unlock()
				if err != nil {
					fmt.Println("stderr ERROR :: ", err)
				}
			}
		}
	}(srv, syncGroup)
	// start the output scanner
	syncGroup.Add(1)
	go func(srv *SCmd, wg *sync.WaitGroup) {
		// fmt.Println("start stdout scanner for", srv.ID)
		defer syncGroup.Done()
		for {
			select {
			case <-srv.StdoutQuitChannel:
				mutex.Lock()
				fmt.Println("quit stdout for", srv.ID)
				mutex.Unlock()
				return
			default:
				if srv.Stdout.Scan() {
					mutex.Lock()
					fmt.Println(srv.ID, ":: stdout", srv.Stdout.Text())
					mutex.Unlock()
				}
				mutex.Lock()
				err := srv.Stdout.Err()
				mutex.Unlock()
				if err != nil {
					fmt.Println("stdout :: ERROR ", err)
				}
			}
		}
	}(srv, syncGroup)
	// start the process
	syncGroup.Add(1)
	go func(srv *SCmd, wg *sync.WaitGroup) {
		defer syncGroup.Done()
		for {
			// fmt.Println("wait 1 seconds before (re)start the service", srv.ID)
			time.Sleep(time.Second * 1)
			select {
			case <-srv.CmdQuitChannel:
				mutex.Lock()
				fmt.Println("quit cmd", srv.ID)
				mutex.Unlock()
				return
			default:
				mutex.Lock()
				// if the services was disabled do not try to rerun
				if !srv.Enabled {
					fmt.Println("the service", srv.ID, "was disabled return from cmd routine")
					return
				}
				err := srv.Cmd.Start()
				if err != nil {
					// seems like the process is already running so do nothing
					if err.Error() == "exec: already started" {
						mutex.Unlock()
						continue
					}
					fmt.Println("cmd", srv.ID, "could not be started because of an error", err)
				} else {
					srv.PID = srv.Cmd.Process.Pid
				}
				mutex.Unlock()
				err = srv.Cmd.Wait()
				if err != nil {
					// seems like somebody, probably the killer killed the process
					if err.Error() == "signal: killed" {
						continue
					}
					mutex.Lock()
					fmt.Println("cmd", srv.ID, "Wait created an error", err)
					mutex.Unlock()
				} else {
					mutex.Lock()
					fmt.Println("cmd", srv.ID, "stopped by itself", err)
					mutex.Unlock()
				}
			}
		}
	}(srv, syncGroup)
	return nil
}

// AppendCronCmd will append a new cron command to the Scheduler
func (s *Scheduler) AppendCronCmd(cmd string, args []string, minutes, hours, monthday, month, weekday string) error {
	eCmd := exec.Command(cmd, args...)
	evCmd := NewCronCmd()
	evCmd.Cmd = eCmd
	cmdStdoutReader, err := eCmd.StdoutPipe()
	if err != nil {
		return err
	}
	evCmd.Stdout = bufio.NewScanner(cmdStdoutReader)
	cmdStderrReader, err := eCmd.StderrPipe()
	if err != nil {
		return err
	}
	evCmd.Stderr = bufio.NewScanner(cmdStderrReader)
	evCmd.LastMinutes = ""
	evCmd.Minutes = minutes
	evCmd.LastHours = ""
	evCmd.Hours = hours
	evCmd.LastMonthDays = ""
	evCmd.MonthDays = monthday
	evCmd.LastMonths = ""
	evCmd.Months = month
	evCmd.LastWeekDays = ""
	evCmd.WeekDays = weekday
	s.CronCmds = append(s.CronCmds, evCmd)
	return nil
}

// RunCron will start/check all cron jobs
func (s *Scheduler) RunCron() error {
	goroutines := runtime.NumGoroutine()
	for _, cronjob := range s.CronCmds {
		now := time.Now()

		runMonth := false
		runDay := false
		runHour := false
		runMinute := false

		// month 1-12
		if strconv.Itoa(int(now.Month())) == cronjob.Months || cronjob.Months == "*" {
			runMonth = true
		}

		// weekday 0-6 ==> 1-7
		weekday := int(now.Weekday()) + 1
		if strconv.Itoa(weekday) == cronjob.WeekDays || cronjob.WeekDays == "*" {
			runDay = true
		}

		// monthday 1-31
		if strconv.Itoa(now.Day()) == cronjob.MonthDays || cronjob.MonthDays == "*" {
			runDay = true
		}

		// hours 0-23
		if strconv.Itoa(now.Hour()) == cronjob.Hours || cronjob.Hours == "*" {
			runHour = true
		}

		// minutes 0-59
		if strconv.Itoa(now.Minute()) == cronjob.Minutes || cronjob.Minutes == "*" {
			runMinute = true
		}

		if runMonth && runDay && runHour && runMinute {
			cronjob.Running = true
			// start stderror scanner routine
			wg.Add(1)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				defer wg.Done()
				for {
					if cronjob.Stderr.Scan() {
						if debug {
							fmt.Println(cronjob.Stderr.Text())
						}
					}
					mutex.Lock()
					if cronjob.Finished {
						mutex.Unlock()
						return
					}
					mutex.Unlock()
				}
			}(cronjob, &wg)
			// start stdout scanner routine
			wg.Add(1)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				defer wg.Done()
				for {
					if cronjob.Stdout.Scan() {
						if debug {
							fmt.Println(cronjob.Stdout.Text())
						}
					}
					mutex.Lock()
					if cronjob.Finished {
						mutex.Unlock()
						return
					}
					mutex.Unlock()
				}
			}(cronjob, &wg)
			// start cronjob routine
			wg.Add(1)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				defer wg.Done()
				if debug {
					fmt.Println("running: " + cronjob.Cmd.Path)
				}
				mutex.Lock()
				err := cronjob.Cmd.Run()
				mutex.Unlock()
				if err != nil {
					log.Println(err.Error())
				}
				mutex.Lock()
				cronjob.Finished = true
				mutex.Unlock()
				return
			}(cronjob, &wg)
		}

	}
	// todo we need to refactor this one
	wg.Add(1)
	go func(cronjobs []*CronCmd, wg *sync.WaitGroup, goroutines int) {
		defer wg.Done()
		for {
			finished := make([]bool, 0)
			for _, cronjob := range cronjobs {
				mutex.Lock()
				if cronjob.Running {
					if cronjob.Finished {
						finished = append(finished, true)
					}
				}
				mutex.Unlock()
			}
			closemsgs := true
			for _, closemsg := range finished {
				if !closemsg {
					closemsgs = closemsg
				}
			}
			if closemsgs {
				for {
					// wait for all routines to finish
					if (goroutines + 1) == runtime.NumGoroutine() {
						// exit goroutine
						return
					}
				}
			}
		}
		// unreachable code hopefully
		// return
	}(s.CronCmds, &wg, goroutines)
	return nil
}
