package eve

import (
	"errors"
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
	CmdQuitChannel    chan bool
	StdoutQuitChannel chan bool
	StderrQuitChannel chan bool
	Enabled           bool
	mux               sync.Mutex
}

// NewSCmd will create a new SchedulerCommand object
func NewSCmd() *SCmd {
	c := new(SCmd)
	c.Running = false
	c.PID = 0
	c.ID = "N/A"
	c.Cmd = nil
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
	mux           sync.Mutex
}

// NewCronCmd creates a new cron command object
func NewCronCmd() *CronCmd {
	cc := new(CronCmd)
	cc.Finished = false
	cc.Running = false
	cc.Cmd = nil
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
	srv.Running = false
	srv.Enabled = true
}

// Disable the service to be stopped or skiped
func (srv *SCmd) Disable() {
	srv.Running = true
	srv.Enabled = false
}

// ResetCmd the service to be stopped or skiped
func (srv *SCmd) ResetCmd(cmdPath string, args []string) error {
	srv.Enabled = false
	srv.Running = false
	srv.ID = path.Base(cmdPath)
	srv.Cmd = exec.Command(cmdPath, args...)
	srv.Cmd.Stderr = os.Stderr
	srv.Cmd.Stdout = os.Stdout
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
	mux                   sync.Mutex
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
	s.mux.Lock()
	ID := len(s.Cmds)
	s.Cmds[ID] = cmd
	s.mux.Unlock()
	return nil
}

// AppendServiceCmd creates a service object from the given command arguments
func (s *Scheduler) AppendServiceCmd(cmdPath string, cmdArgs []string) error {
	sCmd := NewSCmd()
	sCmd.Enabled = true
	sCmd.ID = path.Base(cmdPath)
	sCmd.Cmd = exec.Command(cmdPath, cmdArgs...)
	sCmd.Cmd.Stderr = os.Stderr
	sCmd.Cmd.Stdout = os.Stdout
	ID := len(s.Cmds)
	s.mux.Lock()
	s.Cmds[ID] = sCmd
	s.mux.Unlock()
	return nil
}

// ReplaceServiceCmd replaces a given service command
func (s *Scheduler) ReplaceServiceCmd(cmdPath string, args []string) error {
	sCmd := NewSCmd()
	sCmd.Enabled = true
	sCmd.ID = path.Base(cmdPath)
	sCmd.Cmd = exec.Command(cmdPath, args...)
	sCmd.Cmd.Stderr = os.Stderr
	sCmd.Cmd.Stdout = os.Stdout
	for key, cmd := range s.Cmds {
		if cmd.ID == path.Base(cmdPath) {
			s.mux.Lock()
			s.Cmds[key] = sCmd
			s.mux.Unlock()
			return nil
		}
	}
	return errors.New("could not find the given command")
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
				if debug {
					log.Println("---> kill", srvID)
				}
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
				if debug {
					fmt.Println("quit process killer")
				}
				return
			default:
				if debug {
					fmt.Println("process killer default state reached waiting 1 seconds")
				}
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
	s.Cmds[ID].mux.Lock()
	if !s.Cmds[ID].Enabled {
		s.Cmds[ID].mux.Unlock()
		return nil
	}
	fmt.Println("stopping service", s.Cmds[ID].ID, "with PID", s.Cmds[ID].PID)
	s.Cmds[ID].Running = false
	s.Cmds[ID].Enabled = false
	s.Cmds[ID].CmdQuitChannel <- true
	s.Cmds[ID].StderrQuitChannel <- true
	s.Cmds[ID].StdoutQuitChannel <- true
	s.CmdKillerChannel <- s.Cmds[ID].PID
	s.Cmds[ID].mux.Unlock()
	// fmt.Println("wait 5 sec for the process to stop")
	time.Sleep(time.Second * 5)
	return nil
}

// ServiceStart starts a running services
func (s *Scheduler) ServiceStart(ID int, syncGroup *sync.WaitGroup) error {
	s.Cmds[ID].mux.Lock()
	if !s.Cmds[ID].Enabled || s.Cmds[ID].Running {
		fmt.Println("skip disabled or already running service", s.Cmds[ID].ID)
		s.Cmds[ID].mux.Unlock()
		return nil
	}
	fmt.Println("running service", s.Cmds[ID].ID, "...")
	s.Cmds[ID].mux.Unlock()
	syncGroup.Add(1)
	go func(srv *SCmd, wg *sync.WaitGroup) {
		defer syncGroup.Done()
		for {
			if debug {
				srv.mux.Lock()
				fmt.Println("wait 1 seconds before (re)start the service", srv.ID)
				srv.mux.Unlock()
			}
			time.Sleep(time.Second * 1)
			select {
			case <-srv.CmdQuitChannel:
				if debug {
					srv.mux.Lock()
					fmt.Println("quit cmd", srv.ID)
					srv.mux.Unlock()
				}
				return
			default:
				srv.mux.Lock()
				// if the services was disabled do not try to rerun
				if !srv.Enabled {
					if debug {
						fmt.Println("the service", srv.ID, "was disabled return from cmd routine")
					}
					srv.mux.Unlock()
					return
				}
				err := srv.Cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				srv.PID = srv.Cmd.Process.Pid
				srv.mux.Unlock()
				// check if mux lock is required here
				err = srv.Cmd.Wait()
				if err != nil {
					if err.Error() == "signal: killed" {
						// do nothing this was probably done by the killer or the user
					} else {
						// this unknown error should be displayed
						fmt.Println(err)
					}
				}
			}
		}
	}(s.Cmds[ID], syncGroup)
	return nil
}

// AppendCronCmd will append a new cron command to the Scheduler
func (s *Scheduler) AppendCronCmd(cmd string, args []string, minutes, hours, monthday, month, weekday string) error {
	eCmd := exec.Command(cmd, args...)
	evCmd := NewCronCmd()
	evCmd.Cmd = eCmd
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
					// if cronjob.Stderr.Scan() {
					// 	if debug {
					// 		fmt.Println(cronjob.Stderr.Text())
					// 	}
					// }
					cronjob.mux.Lock()
					if cronjob.Finished {
						cronjob.mux.Unlock()
						return
					}
					cronjob.mux.Unlock()
				}
			}(cronjob, &wg)
			// start stdout scanner routine
			wg.Add(1)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				defer wg.Done()
				for {
					// if cronjob.Stdout.Scan() {
					// 	if debug {
					// 		fmt.Println(cronjob.Stdout.Text())
					// 	}
					// }
					cronjob.mux.Lock()
					if cronjob.Finished {
						cronjob.mux.Unlock()
						return
					}
					cronjob.mux.Unlock()
				}
			}(cronjob, &wg)
			// start cronjob routine
			wg.Add(1)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				defer wg.Done()
				if debug {
					fmt.Println("running: " + cronjob.Cmd.Path)
				}
				cronjob.mux.Lock()
				err := cronjob.Cmd.Run()
				cronjob.mux.Unlock()
				if err != nil {
					log.Println(err.Error())
				}
				cronjob.mux.Lock()
				cronjob.Finished = true
				cronjob.mux.Unlock()
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
				cronjob.mux.Lock()
				if cronjob.Running {
					if cronjob.Finished {
						finished = append(finished, true)
					}
				}
				cronjob.mux.Unlock()
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
