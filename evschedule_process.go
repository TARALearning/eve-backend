package eve

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
	Path              string
	Args              []string
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
	c.Path = ""
	c.Args = make([]string, 0)
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

// WaitGroup returns the sync WaitGroup
func (s *Scheduler) WaitGroup() *sync.WaitGroup {
	return &wg
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
	srv.Path = cmdPath
	srv.Args = args
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
	sCmd.Args = cmdArgs
	sCmd.Path = cmdPath
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
	sCmd.Args = args
	sCmd.Path = cmdPath
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

// ServicesStop stops all running services
func (s *Scheduler) ServicesStop() error {
	for ID := range s.Cmds {
		err := s.ServiceStop(ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ServiceRestart stops all running services
func (s *Scheduler) ServiceRestart(ID int, cmdPath string, args []string) error {
	err := s.ServiceStop(ID)
	if err != nil {
		return err
	}
	s.ReplaceServiceCmd(cmdPath, args)
	s.Cmds[ID].mux.Lock()
	s.Cmds[ID].Enable()
	s.Cmds[ID].mux.Unlock()
	err = s.ServiceStart(0, &wg)
	if err != nil {
		return err
	}
	return nil
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
	if debug {
		fmt.Println("wait 5 sec for the service " + s.Cmds[ID].ID + " to stop...")
	}
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
					// check if the command exited from it self
					// in this case just restart the command
					if err.Error() == "exec: already started" {
						srv.Cmd = exec.Command(srv.Path, srv.Args...)
						srv.Cmd.Stderr = os.Stderr
						srv.Cmd.Stdout = os.Stdout
						err = srv.Cmd.Start()
					}
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
	if debug {
		s.Cmds[ID].mux.Lock()
		fmt.Println("waiting 3 sec for the service " + s.Cmds[ID].ID + " to start...")
		s.Cmds[ID].mux.Unlock()
	}
	// wait 3 sec for the service to initialize
	time.Sleep(time.Second * 3)
	return nil
}
