package eve

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	gops "github.com/mitchellh/go-ps"
)

// wg is that syncing waiting group variable for the goroutines
var wg sync.WaitGroup

// Running contents the commands which are executed by the scheduler
var Running = make(map[string]*SCmd, 0)

// mutex is the variable which is needed for mutexing the Running commands map
var mutex = &sync.Mutex{}

// serviceEnabled describes the service status enabled
var serviceEnabled = "enabled"

// serviceDisabled describes the service status disabled
var serviceDisabled = "disabled"

// SCmd is the struct which is used to store the executed command and the
// stdin stderr output related to the given command
// please be careful with the Id definition because it should be the exact name of
// the given command because it is used to check for the command if it is already running
// if there is another command with the same name it will kill it during bootstrap
type SCmd struct {
	ID                string
	Cmd               *exec.Cmd
	Owner             string
	Stdout            *bufio.Scanner
	Stderr            *bufio.Scanner
	StdoutQuitChannel chan bool
	StderrQuitChannel chan bool
	ServiceType       string
}

// NewSCmd will create a new SchedulerCommand object
func NewSCmd() *SCmd {
	c := new(SCmd)
	c.ID = "N/A"
	c.Cmd = nil
	c.Owner = ""
	c.Stdout = nil
	c.Stderr = nil
	c.StdoutQuitChannel = make(chan bool)
	c.StderrQuitChannel = make(chan bool)
	c.ServiceType = serviceEnabled
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

// Scheduler is the scheduler main struct
type Scheduler struct {
	Cmds                  map[int]*SCmd
	CronCmds              []*CronCmd
	MonitoringQuitChannel chan bool
	MainQuitChannel       chan bool
}

// NewScheduler will create a new Scheduler object
func NewScheduler() *Scheduler {
	s := new(Scheduler)
	s.Cmds = make(map[int]*SCmd, 0)
	s.CronCmds = make([]*CronCmd, 0)
	s.MonitoringQuitChannel = make(chan bool)
	s.MainQuitChannel = make(chan bool)
	return s
}

// RestartCmd does restart a command which was already started before
func (s *Scheduler) RestartCmd(cmdID string) error {
	found := false
	mutex.Lock()
	defer mutex.Unlock()
	for id, oCmd := range s.Cmds {
		if oCmd.ID == cmdID {
			found = true
			eCmd := exec.Command(oCmd.Cmd.Args[0], oCmd.Cmd.Args[1:]...)
			evCmd := NewSCmd()
			evCmd.ID = cmdID
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
			s.Cmds[id] = evCmd
			err = s.RunCmd(evCmd)
			if err != nil {
				return err
			}
		}
	}
	if !found {
		return errors.New("the given process with the id:" + cmdID + " could not be found!")
	}
	return nil
}

// appendCmd appends a command
func (s *Scheduler) appendCmd(cmdID, cmd, owner string, args []string) error {
	eCmd := exec.Command(cmd, args...)
	evCmd := NewSCmd()
	evCmd.ID = cmdID
	evCmd.Cmd = eCmd
	evCmd.Owner = owner
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
	order := len(s.Cmds)
	s.Cmds[order] = evCmd
	return nil
}

// AppendCmd will append a command to the Scheduler Running map
func (s *Scheduler) AppendCmd(cmdID, cmd, owner string, args []string) error {
	mutex.Lock()
	defer mutex.Unlock()
	return s.appendCmd(cmdID, cmd, owner, args)
}

// DeleteCmd removes a cmd from the Scheduler Running map
func (s *Scheduler) DeleteCmd(cmdID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for idx, cmd := range s.Cmds {
		if cmd.ID == cmdID {
			delete(s.Cmds, idx)
		}
	}
	return nil
}

// ReplaceCmd removes/appends a cmd from/to the Scheduler Running map
func (s *Scheduler) ReplaceCmd(cmdID, cmd, owner string, args []string) error {
	err := s.DeleteCmd(cmdID)
	if err != nil {
		return err
	}
	return s.AppendCmd(cmdID, cmd, owner, args)
}

// killCmd kills the given command
func killCmd(cmd *SCmd) error {
	if DEBUG {
		fmt.Println("killing process with ID", cmd.ID)
	}
	cmd.ServiceType = serviceDisabled
	delete(Running, cmd.ID)
	close(cmd.StderrQuitChannel)
	close(cmd.StdoutQuitChannel)
	return EnforceProcessKill(cmd.ID)
}

// KillCmd kills a cmd from the Scheduler Running map
func (s *Scheduler) KillCmd(cmdID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for _, cmd := range s.Cmds {
		if cmd.ID == cmdID {
			return killCmd(cmd)
		}
	}
	return nil
}

// KillAllCmds kills all cmds from the Scheduler Running map
func (s *Scheduler) KillAllCmds() error {
	mutex.Lock()
	defer mutex.Unlock()
	for _, cmd := range s.Cmds {
		if cmd.ServiceType != serviceDisabled {
			err := killCmd(cmd)
			if err != nil {
				log.Println("Scheduler::KillAllCmds cannot kill command with Id", cmd.ID, "ERROR::", err)
				return err
			}
		}
	}
	return nil
}

// RunCmd will create a stdout and stderr scanner and will run the given command
func (s *Scheduler) RunCmd(cmd *SCmd) error {
	// todo implement more service types like enabled,disabled,stopped,...
	switch cmd.ServiceType {
	case serviceDisabled:
		if DEBUG {
			fmt.Println("Scheduler::RunCmd found ServiceType disabled do not start cmd", cmd.ID)
		}
		return nil
	default:
		if DEBUG {
			fmt.Println("Scheduler::RunCmd starting found ServiceType", cmd.ServiceType, " with id", cmd.ID)
		}
	}
	go func(cmd *SCmd, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-cmd.StdoutQuitChannel:
				return
			default:
				if cmd.Stdout.Scan() {
					if DEBUG {
						fmt.Println(cmd.Stdout.Text())
					}
				}
			}
		}
	}(cmd, &wg)
	go func(cmd *SCmd, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-cmd.StderrQuitChannel:
				return
			default:
				if cmd.Stderr.Scan() {
					if DEBUG {
						fmt.Println("ERROR MSG:::" + cmd.Stderr.Text())
					}
				}
			}
		}
	}(cmd, &wg)
	go func(cmd *SCmd, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		if cmd.Owner != "" {
			fmt.Println("EVSchedule", "running", cmd.Cmd.Path, "as user", cmd.Owner)
		} else {
			fmt.Println("EVSchedule", "running", cmd.Cmd.Path)
		}
		/*
			// todo try to find a way to fix that problem
			// till go version 1.7.x there is no official support
			// for syscall.Setuid syscall.Setgid :-(
			if cmd.Owner != "" {
				owner := strings.Split(cmd.Owner, ":")
				usr, err := user.Lookup(owner[0])
				if err != nil {
					if DEBUG {
						log.Fatal(err)
					}
				}
				gid, err := user.LookupGroupID(owner[1])
				if err != nil {
					if DEBUG {
						log.Println("WARNING:", err)
					}
				}
				iuid, err := strconv.Atoi(usr.Uid)
				if err != nil {
					if DEBUG {
						log.Fatal(err)
					}
				}
				err = syscall.Setuid(iuid)
				if err != nil {
					if DEBUG {
						log.Fatal(err)
					}
				}
				// if LookupGroupID does not return a value
				// fallback to the usr.Gid group
				var igid int
				if gid != nil {
					igid, err = strconv.Atoi(gid.Gid)
					if err != nil {
						if DEBUG {
							log.Fatal(err)
						}
					}
				} else {
					igid, err = strconv.Atoi(usr.Gid)
					if err != nil {
						if DEBUG {
							log.Fatal(err)
						}
					}
				}
				err = syscall.Setgid(igid)
				if err != nil {
					if DEBUG {
						log.Fatal(err)
					}
				}
			}
		*/
		err := cmd.Cmd.Run()
		if err != nil {
			if DEBUG {
				log.Fatal(err)
			}
		}
	}(cmd, &wg)
	go func(cmd *SCmd, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		// wait 3 seconds for the cmd to be initialized
		time.Sleep(3 * time.Second)
		mutex.Lock()
		Running[cmd.ID] = cmd
		mutex.Unlock()
	}(cmd, &wg)
	// time.Sleep(6 * time.Second)
	return nil
}

// CmdFind returns the map int key of the command with the given cmd id
func (s *Scheduler) CmdFind(pID string) int {
	// mutex.Lock()
	// defer mutex.Unlock()
	for id, cmd := range s.Cmds {
		if pID == cmd.ID {
			return id
		}
	}
	return -1
}

// Monitoring will start a monitoring goroutine for all started commands
func (s *Scheduler) Monitoring(wg *sync.WaitGroup) {
	fmt.Println("EVSchedule :: starting monitoring for the scheduler...")
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-s.MonitoringQuitChannel:
			return
		default:
			if DEBUG {
				fmt.Println("waiting 5 seconds before check if the command is still running")
			}
			time.Sleep(5 * time.Second)
			mutex.Lock()
			if DEBUG {
				fmt.Println(Running)
			}
			for id, cmd := range Running {
				cont := false
				if DEBUG {
					fmt.Println(cmd.Cmd.Process)
				}
				if cmd.Cmd.Process == nil {
					if DEBUG {
						fmt.Println(id + ": the process is not running or it is still booting waiting 5 seconds for the process to boot!")
					}
					if cmd.Cmd.Process == nil {
						if DEBUG {
							fmt.Println(id + ": the process seems not to run try to restart it!")
						}
					}
				}
				proc, err := gops.FindProcess(cmd.Cmd.Process.Pid)
				if DEBUG {
					fmt.Println("proc == nil && err == nil")
				}
				if proc == nil && err == nil {
					if DEBUG {
						fmt.Println("the process with id", id, "was finished successfully it will be now restarted")
					}
					err = killCmd(cmd)
					if DEBUG {
						fmt.Println("kill command", err)
					}
					if err != nil {
						log.Println(err.Error())
						log.Println("EVSchedule", err)
					}
					// err = nil
					cmdID := s.CmdFind(id)
					if DEBUG {
						fmt.Println("command ID::", cmdID)
					}
					if cmdID == -1 {
						log.Println("error:: can not find the given command with id::", id, "in the Scheduler.Cmds map")
					} else {
						cmdArgs := cmd.Cmd.Args
						if DEBUG {
							log.Println("delete cmd with id", cmdID, "from s.Cmds")
						}
						delete(s.Cmds, cmdID)
						if DEBUG {
							fmt.Println(cmdArgs)
						}
						err = s.appendCmd(id, cmdArgs[0], "", cmdArgs[1:])
						if err != nil {
							log.Println(err.Error())
							log.Println("EVSchedule", err)
						}
						cmdID = s.CmdFind(id)
						if cmdID == -1 {
							log.Println("error:: can not find the given command with id::", id, "in the Scheduler.Cmds map")
						} else {
							*cmd = *s.Cmds[cmdID]
						}
					}
					err = startCmd(cmd)
					if DEBUG {
						fmt.Println("start command", err)
					}
					if err != nil {
						log.Println(err.Error())
						log.Println("EVSchedule", err)
					}
					err = s.RunCmd(cmd)
					if DEBUG {
						fmt.Println("run command", err)
					}
					if err != nil {
						log.Println(err.Error())
						log.Println("EVSchedule", err)
					}
					cont = true
				}
				if !cont {
					if DEBUG {
						fmt.Println("err != nil || proc == nil")
					}
					if err != nil || proc == nil {
						if DEBUG {
							fmt.Println("EVSchedule", err)
						}
						err = s.RestartCmd(id)
						if err != nil {
							log.Println(err.Error())
							log.Println("EVSchedule", err)
						}
						if proc == nil {
							if DEBUG {
								fmt.Println("process could not be found, suggested process is probably death!")
							}
						}
					}
					if DEBUG {
						fmt.Println("err != nil && proc != nil")
					}
					if err != nil && proc != nil {
						log.Println("service: " + id + "(" + strconv.Itoa(cmd.Cmd.Process.Pid) + ")" + " seems to be OK")
					}
					if DEBUG {
						fmt.Println("proc == nil || proc.pID() != cmd.Cmd.Process.pID")
					}
					if proc == nil || proc.Pid() != cmd.Cmd.Process.Pid {
						log.Println("the process seems to be restarted!")
					}
				}
			}
			mutex.Unlock()
		}
	}
}

// EnforceProcessKill will check for a command based on it's Id if it is running
// and take care to kill it anyway
func EnforceProcessKill(cmdID string) error {
	allProcs, err := gops.Processes()
	if err != nil {
		return err
	}
	for _, wProc := range allProcs {
		if DEBUG {
			fmt.Println("found::", wProc.Executable())
		}
		if wProc.Executable() == cmdID {
			if DEBUG {
				fmt.Println("EnforceProcessKill found executable with id, ", wProc.Executable())
			}
			proc, err := os.FindProcess(wProc.Pid())
			if err != nil {
				return err
			}
			err = proc.Kill()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("EnforceProcessKill::ERROR could not find process with id::" + cmdID)
}

// Run will start all available commands in the Running map
// and it will also start the monitoring goroutine to take care
// of all started/running commands
func (s *Scheduler) Run() error {
	for _, cmd := range s.Cmds {
		if cmd.ServiceType != "disabled" {
			err := EnforceProcessKill(cmd.ID)
			if err != nil {
				// do nothing because if the executable
				// is not available it is the best it could happen
				// return err
			}
			err = s.RunCmd(cmd)
			if err != nil {
				return err
			}
		}
	}
	go s.Monitoring(&wg)
	wg.Wait()
	return nil
}

// Shutdown shuts all goroutines and commands gracefully down
func (s *Scheduler) Shutdown() error {
	s.MonitoringQuitChannel <- true
	err := s.KillAllCmds()
	if err != nil {
		return err
	}
	wg.Add(1)
	go func(s *Scheduler, wg *sync.WaitGroup) {
		defer wg.Done()
		if DEBUG {
			fmt.Println("sending quit message to scheduler main routine in 3 sec")
		}
		time.Sleep(3 * time.Second)
		s.MainQuitChannel <- true
	}(s, &wg)
	return nil
}

// starts the given command
func startCmd(cmd *SCmd) error {
	cmd.ServiceType = "enabled"
	cmd.StdoutQuitChannel = make(chan bool)
	cmd.StderrQuitChannel = make(chan bool)
	return nil
}

// EnableProcess will enable the process with the given id
func (s *Scheduler) EnableProcess(pID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for _, cmd := range s.Cmds {
		if cmd.ID == pID {
			cmd.ServiceType = "enabled"
			return nil
		}
	}
	return errors.New("EnableProcess could not found the process with the given id::" + pID)
}

// DisableProcess will disable the process with the given id
func (s *Scheduler) DisableProcess(pID string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for _, cmd := range s.Cmds {
		if cmd.ID == pID {
			cmd.ServiceType = "disabled"
			return nil
		}
	}
	return errors.New("DisableProcess could not found the process with the given id::" + pID)
}

// StartProcess will start the process with the given id
func (s *Scheduler) StartProcess(pID string) error {
	if DEBUG {
		fmt.Println("Starting Processsssss :::: ", pID)
	}
	mutex.Lock()
	defer mutex.Unlock()
	if DEBUG {
		fmt.Println(s.Cmds)
	}
	for _, cmd := range s.Cmds {
		if DEBUG {
			fmt.Println(cmd)
		}
		if cmd.ID == pID {
			err := startCmd(cmd)
			if err != nil {
				log.Println("Scheduler::StartProcesses cannot start process with Id", cmd.ID, "ERROR::", err)
			}
			return s.RunCmd(cmd)
		}
	}
	return errors.New("StartProcess could not found the process with the given id::" + pID)
}

// StartAllProcesses will start all processes
func (s *Scheduler) StartAllProcesses() error {
	for _, cmd := range s.Cmds {
		log.Println(cmd)
		err := s.StartProcess(cmd.ID)
		if err != nil {
			log.Println("Scheduler::StartAllProcesses cannot start process with Id", cmd.ID, "ERROR::", err)
			return err
		}
	}
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
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				wg.Add(1)
				defer wg.Done()
				for {
					if cronjob.Stdout.Scan() {
						if DEBUG {
							fmt.Println(cronjob.Stdout.Text())
						}
					}
					if cronjob.Finished {
						return
					}
				}
			}(cronjob, &wg)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				wg.Add(1)
				defer wg.Done()
				for {
					if cronjob.Stderr.Scan() {
						if DEBUG {
							fmt.Println(cronjob.Stderr.Text())
						}
					}
					if cronjob.Finished {
						return
					}
				}
			}(cronjob, &wg)
			go func(cronjob *CronCmd, wg *sync.WaitGroup) {
				wg.Add(1)
				defer wg.Done()
				if DEBUG {
					fmt.Println("running: " + cronjob.Cmd.Path)
				}
				err := cronjob.Cmd.Run()
				if err != nil {
					log.Println(err.Error())
				}
				cronjob.Finished = true
				return
			}(cronjob, &wg)
		}

	}
	go func(cronjobs []*CronCmd, wg *sync.WaitGroup, goroutines int) {
		wg.Add(1)
		defer wg.Done()
		for {
			finished := make([]bool, 0)
			for _, cronjob := range cronjobs {
				if cronjob.Running {
					if cronjob.Finished {
						finished = append(finished, true)
					}
				}
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
