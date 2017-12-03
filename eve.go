package eve

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Uses []string

func (c *Uses) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *Uses) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func GenEvSchedule(config string, use Uses, filepath string) (string, error) {
	srv := &EVServiceConfigObj{}
	SetDefaultCType("rest")
	var srvConfig *EVServiceConfig
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
		var vars map[string]interface{}
		var imports []string
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
			"strconv",
		}
		vars = map[string]interface{}{
			"Package":                "main",
			"DefaultAddress":         "127.0.0.1:9091",
			"UsageFunc":              "EVUsage",
			"Version":                VERSION,
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
			"EVLOG_URL":              "http://127.0.0.1:9091/" + VERSION + "/eve/evlog",
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
			"ROUTE_PATH_PREFIX": "/" + VERSION + "/eve/",
		}
		srvConfig = &EVServiceConfig{
			Main:      SrvConfigMain(),
			Imports:   imports,
			Templates: SrvConfigTemplates(),
			Commands:  SrvConfigCommands(),
			Vars:      vars,
		}
	}
	return genEvService(srv, srvConfig, use, filepath)
}

func GenEvLog(config string, use Uses, filepath string) (string, error) {
	srv := &EVServiceConfigObj{}
	SetDefaultCType("rest")
	var srvConfig *EVServiceConfig
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
		var vars map[string]interface{}
		var imports []string
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
			"Version":                VERSION,
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
			"ROUTE_PATH_PREFIX": "/" + VERSION + "/eve/",
		}
		srvConfig = &EVServiceConfig{
			Main:      SrvConfigMain(),
			Imports:   imports,
			Templates: SrvConfigTemplates(),
			Commands:  SrvConfigCommands(),
			Vars:      vars,
		}
	}
	return genEvService(srv, srvConfig, use, filepath)
}

func GenEvBolt(config string, use Uses, filepath string) (string, error) {
	srv := &EVServiceConfigObj{}
	SetDefaultCType("rest")
	var srvConfig *EVServiceConfig
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
		var vars map[string]interface{}
		var imports []string
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
			"Version":                VERSION,
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
			"EVLOG_URL":              "http://127.0.0.1:9091/" + VERSION + "/eve/evlog",
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
			"ROUTE_PATH_PREFIX": "/" + VERSION + "/eve/",
		}
		srvConfig = &EVServiceConfig{
			Main:      SrvConfigMain(),
			Imports:   imports,
			Templates: SrvConfigTemplates(),
			Commands:  SrvConfigCommands(),
			Vars:      vars,
		}
	}
	return genEvService(srv, srvConfig, use, filepath)
}

func GenEvAuth(config string, use Uses, filepath string) (string, error) {
	srv := &EVServiceConfigObj{}
	SetDefaultCType("rest")
	var srvConfig *EVServiceConfig
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
		var vars map[string]interface{}
		var imports []string
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
			"Version":                VERSION,
			"Name":                   "EVAuth",
			"Description":            "EVAuth is a rest micro service which can be used to authenticate agains users or hosts/services",
			"Src":                    "https://git.evalgo.de:8443/",
			"DEBUG":                  true,
			"ENABLE_CROSS_ORIGIN":    true,
			"USE_PROMETHEUS":         true,
			"USE_EVTOKEN":            true,
			"TOKEN_STORAGE_URL":      "http://127.0.0.1:9092/" + VERSION + "/eve/evbolt",
			"TOKEN_STORAGE_DB":       "tokens.db",
			"TOKEN_STORAGE_BUCKET":   "tokens",
			"USE_EVLOG":              false,
			"EVLOG_URL":              "",
			"USE_EVLOG_API":          false,
			"USE_EVSESSION":          true,
			"SESSION_STORAGE_URL":    "http://127.0.0.1:9092/" + VERSION + "/eve/evbolt",
			"SESSION_STORAGE_DB":     "sessions.db",
			"SESSION_STORAGE_BUCKET": "sessions",
			"USE_EVSECRET":           true,
			"SECRET_STORAGE_URL":     "http://127.0.0.1:9092/" + VERSION + "/eve/evbolt",
			"SECRET_STORAGE_DB":      "secrets.db",
			"SECRET_STORAGE_BUCKET":  "secrets",
			"SECRET_ENC_KEY":         "TokenKeyEnc",
			"SECRET_SIG_KEY":         "TokenKeySig",
			"USE_EVUSER":             true,
			"USER_STORAGE_URL":       "http://127.0.0.1:9092/" + VERSION + "/eve/evbolt",
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
			"ROUTE_PATH_PREFIX": "/" + VERSION + "/eve/",
		}
		srvConfig = &EVServiceConfig{
			Main:      SrvConfigMain(),
			Imports:   imports,
			Templates: SrvConfigTemplates(),
			Commands:  SrvConfigCommands(),
			Vars:      vars,
		}
	}
	return genEvService(srv, srvConfig, use, filepath)
}
func genEvService(srv *EVServiceConfigObj, srvConfig *EVServiceConfig, use Uses, filepath string) (string, error) {
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
					return "", errors.New("error: the given use flag <" + useCaseSplit[0] + "> does not exist!")
				}
			} else {
				return "", errors.New("error: the given use flag <" + u + "> does not exist!")
			}
		}
	}
	res, err := EVGenMain(srv)
	if err != nil {
		return "", err
	}
	if filepath == "." {
		filepath = "." + string(os.PathSeparator) + "main.go"
	}
	fmt.Println("eve-gen :: generating file " + filepath + "...")
	err = ioutil.WriteFile(filepath, res, 0777)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
