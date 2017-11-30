package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"evalgo.org/eve"
)

type uses []string

func (c *uses) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *uses) Set(value string) error {
	*c = append(*c, value)
	return nil
}

var (
	flags   *flag.FlagSet
	version = eve.VERSION
	command = ""
	service = ""
	use     uses
	debug   = false
	config  = ""
	target  = ""
)

// GenUsage displays the help/usage instructions
func GenUsage() {
	fmt.Println(`
usage: eve-gen {command} {type} {service-name} -use {use_flag} -use {use_flag} -target /path/to/main.go...

commands:
    help
	generate
	version

types:
	service

arguments:
	config
	target

services:
	evauth
	evbolt
	evlog
	evschedule

use flags:
	debug
	`)
}

func init() {
	if len(os.Args) <= 1 {
		GenUsage()
		os.Exit(2)
	}
	command = os.Args[1]
	flags = flag.NewFlagSet(command, flag.ExitOnError)
	flags.BoolVar(&debug, "debug", false, "turns debug mode on or off")
	flags.StringVar(&service, "service", "", "service to be generated")
	flags.StringVar(&config, "config", "", "file with a given configuration to be used to build the service")
	flags.StringVar(&target, "target", ".", "path to the directory where the fiel should be generated")
	flags.Var(&use, "use", "use a specific module/feature for the given service to be created")
	flags.Usage = GenUsage
}

func main() {
	flags.Parse(os.Args[2:])
	if debug {
		eve.SetDebug(true)
	}

	switch command {

	case "help":
		GenUsage()
		os.Exit(0)
	case "version":
		fmt.Println("eve-gen version:", version)
	default:
		srv := &eve.EVServiceConfigObj{}
		// todo implement change of the type to schedule
		//  in case evschedule is built with a config file
		eve.SetDefaultCType("rest")

		var srvConfig *eve.EVServiceConfig
		fmt.Println("eve-gen :: check if config file was provided...")
		if config != "" {
			fmt.Println("eve-gen :: found config file <" + config + ">")
			srvJSON, err := ioutil.ReadFile(config)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("eve-gen :: reading config file from:", config+"...")
			err = json.Unmarshal(srvJSON, srv)
			if err != nil {
				log.Fatal(err)
			}
			srvConfig = srv.EVServiceConfiguration()
		} else {
			fmt.Println("eve-gen :: check if service name was provided...")
			if service == "" {
				fmt.Println("error: plese specify a service name")
				GenUsage()
				os.Exit(2)
			}
			fmt.Println("eve-gen :: found service name <" + service + ">")
			var vars map[string]interface{}
			var imports []string
			switch service {
			case "evauth":
				imports = []string{
					"fmt",
					"flag",
					"os",
					"strings",
					"log",
					"net/http",
					"github.com/prometheus/client_golang/prometheus",
					"github.com/prometheus/client_golang/prometheus/promhttp",
					"github.com/dchest/uniuri",
					"github.com/gorilla/mux",
					"evalgo.org/eve",
					"net/url",
					"encoding/base64",
					"errors",
					"time",
					"path",
				}
				vars = map[string]interface{}{"Package": "main",
					"DefaultAddress":         "127.0.0.1:9093",
					"UsageFunc":              "EVUsage",
					"Version":                eve.VERSION,
					"Name":                   "EVAuth",
					"Description":            "EVAuth is a rest micro service which can be used to authenticate agains users or hosts/services",
					"Src":                    "https://git.evalgo.de:8443/",
					"DEBUG":                  true,
					"ENABLE_CROSS_ORIGIN":    true,
					"USE_PROMETHEUS":         true,
					"USE_EVTOKEN":            true,
					"TOKEN_STORAGE_URL":      "http://localhost:9092/" + eve.VERSION + "/eve/evbolt",
					"TOKEN_STORAGE_DB":       "tokens.db",
					"TOKEN_STORAGE_BUCKET":   "tokens",
					"USE_EVLOG":              false,
					"EVLOG_URL":              "",
					"USE_EVLOG_API":          false,
					"USE_EVSESSION":          true,
					"SESSION_STORAGE_URL":    "http://localhost:9092/" + eve.VERSION + "/eve/evbolt",
					"SESSION_STORAGE_DB":     "sessions.db",
					"SESSION_STORAGE_BUCKET": "sessions",
					"USE_EVSECRET":           true,
					"SECRET_STORAGE_URL":     "http://localhost:9092/" + eve.VERSION + "/eve/evbolt",
					"SECRET_STORAGE_DB":      "secrets.db",
					"SECRET_STORAGE_BUCKET":  "secrets",
					"SECRET_ENC_KEY":         "TokenKeyEnc",
					"SECRET_SIG_KEY":         "TokenKeySig",
					"USE_EVUSER":             true,
					"USER_STORAGE_URL":       "http://localhost:9092/" + eve.VERSION + "/eve/evbolt",
					"USER_STORAGE_DB":        "users.db",
					"USER_STORAGE_BUCKET":    "users",
					"USE_EVBOLT_API":         true,
					"USE_EVBOLT_AUTH":        true,
					"USE_LOGIN_API":          true,
					"SECRET_KEY_FOR_TOKEN":   "123456789012345678901234567890ab",
					"SECRET_SIG_FOR_TOKEN":   "sig.key.secret",
					"COOKIE_EXP_MINUTES":     15,
					"TIME_ZONE_LOCATION":     "Europe/Berlin",
					"TOKEN_EXP_DAYS":         7,
					"USE_EVSCHEDULE":         false,
					"USE_EVSCHEDULE_API":     false,
					"URLS": []string{
						"/help",
						"/metrics",
						"/evbolt",
						"/login",
						"/access",
						"/logout",
						"/evbolt.json",
						"/login.json",
						"/access.json",
						"/logout.json",
						"/evbolt.html",
						"/login.html",
						"/access.html",
						"/logout.html",
					},
					"ROUTE_PATH_PREFIX": "/" + eve.VERSION + "/eve/",
				}
			case "evbolt":
				imports = []string{
					"fmt",
					"flag",
					"os",
					"log",
					"net/http",
					"github.com/prometheus/client_golang/prometheus",
					"github.com/prometheus/client_golang/prometheus/promhttp",
					"github.com/gorilla/mux",
					"evalgo.org/eve",
					"errors",
					"strings",
					"path",
				}
				vars = map[string]interface{}{"Package": "main",
					"DefaultAddress":         "127.0.0.1:9092",
					"UsageFunc":              "EVUsage",
					"Version":                eve.VERSION,
					"Name":                   "EVBolt",
					"Description":            "EVBolt is a rest micro service which wrapps the golang bolt database",
					"Src":                    "https://git.evalgo.de:8443/",
					"DEBUG":                  false,
					"ENABLE_CROSS_ORIGIN":    true,
					"USE_PROMETHEUS":         true,
					"USE_EVTOKEN":            false,
					"TOKEN_STORAGE_URL":      "",
					"TOKEN_STORAGE_DB":       "",
					"TOKEN_STORAGE_BUCKET":   "",
					"USE_EVLOG":              true,
					"EVLOG_URL":              "http://localhost:9091/" + eve.VERSION + "/eve/evlog",
					"USE_EVLOG_API":          false,
					"USE_EVSESSION":          false,
					"SESSION_STORAGE_URL":    "",
					"SESSION_STORAGE_DB":     "",
					"SESSION_STORAGE_BUCKET": "",
					"USE_EVSECRET":           false,
					"SECRET_STORAGE_URL":     "",
					"SECRET_STORAGE_DB":      "",
					"SECRET_STORAGE_BUCKET":  "",
					"SECRET_ENC_KEY":         "",
					"SECRET_SIG_KEY":         "",
					"USE_EVUSER":             false,
					"USER_STORAGE_URL":       "",
					"USER_STORAGE_DB":        "",
					"USER_STORAGE_BUCKET":    "",
					"USE_evBoltRoot":         ".",
					"USE_EVBOLT_API":         true,
					"USE_EVBOLT_AUTH":        false,
					"USE_LOGIN_API":          false,
					"SECRET_KEY_FOR_TOKEN":   "",
					"SECRET_SIG_FOR_TOKEN":   "",
					"COOKIE_EXP_MINUTES":     0,
					"TOKEN_EXP_DAYS":         0,
					"USE_EVSCHEDULE":         false,
					"USE_EVSCHEDULE_API":     false,
					"URLS": []string{
						"/help",
						"/evbolt",
						"/evbolt.json",
						"/evbolt.html",
						"/metrics",
					},
					"ROUTE_PATH_PREFIX": "/" + eve.VERSION + "/eve/",
				}
			case "evlog":
				imports = []string{
					"fmt",
					"flag",
					"os",
					"log",
					"net/http",
					"github.com/prometheus/client_golang/prometheus",
					"github.com/prometheus/client_golang/prometheus/promhttp",
					"github.com/gorilla/mux",
					"evalgo.org/eve",
					"path",
					"strings",
				}
				vars = map[string]interface{}{
					"Package":                "main",
					"DefaultAddress":         "127.0.0.1:9091",
					"UsageFunc":              "EVUsage",
					"Version":                eve.VERSION,
					"Name":                   "EVLog",
					"Description":            "EVLog is a rest micro service to be used for logging messages from the other microservices",
					"Src":                    "https://git.evalgo.de:8443/",
					"DEBUG":                  false,
					"ENABLE_CROSS_ORIGIN":    true,
					"USE_PROMETHEUS":         true,
					"USE_EVTOKEN":            false,
					"TOKEN_STORAGE_URL":      "",
					"TOKEN_STORAGE_DB":       "",
					"TOKEN_STORAGE_BUCKET":   "",
					"USE_EVLOG":              false,
					"EVLOG_URL":              "",
					"USE_EVLOG_API":          true,
					"USE_EVSESSION":          false,
					"SESSION_STORAGE_URL":    "",
					"SESSION_STORAGE_DB":     "",
					"SESSION_STORAGE_BUCKET": "",
					"USE_EVSECRET":           false,
					"SECRET_STORAGE_URL":     "",
					"SECRET_STORAGE_DB":      "",
					"SECRET_STORAGE_BUCKET":  "",
					"SECRET_ENC_KEY":         "",
					"SECRET_SIG_KEY":         "",
					"USE_EVUSER":             false,
					"USER_STORAGE_URL":       "",
					"USER_STORAGE_DB":        "",
					"USER_STORAGE_BUCKET":    "",
					"USE_EVBOLT_API":         false,
					"USE_EVBOLT_AUTH":        false,
					"USE_LOGIN_API":          false,
					"SECRET_KEY_FOR_TOKEN":   "",
					"SECRET_SIG_FOR_TOKEN":   "",
					"COOKIE_EXP_MINUTES":     0,
					"TOKEN_EXP_DAYS":         0,
					"USE_EVSCHEDULE":         false,
					"USE_EVSCHEDULE_API":     false,
					"URLS": []string{
						"/help",
						"/evlog",
						"/metrics",
					},
					"ROUTE_PATH_PREFIX": "/" + eve.VERSION + "/eve/",
				}
			case "evschedule":
				imports = []string{
					"fmt",
					"flag",
					"os",
					"log",
					"net/http",
					"github.com/prometheus/client_golang/prometheus",
					"github.com/prometheus/client_golang/prometheus/promhttp",
					"github.com/gorilla/mux",
					"evalgo.org/eve",
					"errors",
					"strings",
					"sync",
					"time",
					"path",
				}
				vars = map[string]interface{}{
					"Package":                "main",
					"DefaultAddress":         "127.0.0.1:9091",
					"UsageFunc":              "EVUsage",
					"Version":                eve.VERSION,
					"Name":                   "EVSchedule",
					"Description":            "EVSchedule is a rest micro service to manage processes",
					"Src":                    "https://git.evalgo.de:8443/",
					"DEBUG":                  true,
					"ENABLE_CROSS_ORIGIN":    true,
					"USE_PROMETHEUS":         true,
					"USE_EVTOKEN":            false,
					"TOKEN_STORAGE_URL":      "",
					"TOKEN_STORAGE_DB":       "",
					"TOKEN_STORAGE_BUCKET":   "",
					"USE_EVLOG":              true,
					"EVLOG_URL":              "http://localhost:9091/" + eve.VERSION + "/eve/evlog",
					"USE_EVLOG_API":          false,
					"USE_EVSESSION":          false,
					"SESSION_STORAGE_URL":    "",
					"SESSION_STORAGE_DB":     "",
					"SESSION_STORAGE_BUCKET": "",
					"USE_EVSECRET":           false,
					"SECRET_STORAGE_URL":     "",
					"SECRET_STORAGE_DB":      "",
					"SECRET_STORAGE_BUCKET":  "",
					"SECRET_ENC_KEY":         "",
					"SECRET_SIG_KEY":         "",
					"USE_EVUSER":             false,
					"USER_STORAGE_URL":       "",
					"USER_STORAGE_DB":        "",
					"USER_STORAGE_BUCKET":    "",
					"USE_EVBOLT_API":         false,
					"USE_EVBOLT_AUTH":        false,
					"USE_LOGIN_API":          false,
					"SECRET_KEY_FOR_TOKEN":   "",
					"SECRET_SIG_FOR_TOKEN":   "",
					"COOKIE_EXP_MINUTES":     0,
					"TOKEN_EXP_DAYS":         0,
					"USE_EVSCHEDULE":         true,
					"USE_EVSCHEDULE_API":     true,
					"URLS": []string{
						"/help",
						"/evschedule",
						"/evschedule/commands/{command}",
						"/evschedule.json",
						"/evschedule/commands/{command}.json",
						"/evschedule.html",
						"/evschedule/commands/{command}.html",
						"/metrics",
					},
					"ROUTE_PATH_PREFIX": "/" + eve.VERSION + "/eve/",
				}
			default:
				fmt.Println("error: the given service name <" + service + "> is not supported yet")
				GenUsage()
				os.Exit(2)
			}
			srvConfig = &eve.EVServiceConfig{
				Main:      eve.SrvConfigMain(),
				Imports:   imports,
				Templates: eve.SrvConfigTemplates(),
				Commands:  eve.SrvConfigCommands(),
				Vars:      vars,
			}
		}
		srv.Config = srvConfig
		fmt.Println("eve-gen :: check use flags...")
		if len(use) == 0 {
			fmt.Println("eve-gen :: no use flags found")
		}
		for _, u := range use {
			fmt.Println("eve-gen :: enabling use flag <" + u + ">...")
			switch u {
			case "debug":
				srvConfig.Vars["DEBUG"] = true
			default:
				useCaseSplit := strings.Split(u, "=")
				if len(useCaseSplit) == 2 {
					switch useCaseSplit[0] {
					case "evBoltRoot":
						srvConfig.Vars["USE_"+useCaseSplit[0]] = useCaseSplit[1]
					default:
						fmt.Println("")
						fmt.Println("error: the given use flag <" + useCaseSplit[0] + "> does not exist!")
						GenUsage()
						os.Exit(2)
					}
				} else {
					fmt.Println("")
					fmt.Println("error: the given use flag <" + u + "> does not exist!")
					GenUsage()
					os.Exit(2)
				}
			}
		}
		res, err := eve.EVGenMain(srv)
		if err != nil {
			log.Fatal(err)
		}
		if target == "." {
			target = "." + string(os.PathSeparator) + "main.go"
		}
		fmt.Println("eve-gen :: generating file " + target + "...")
		err = ioutil.WriteFile(target, res, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}
