package eve

var (
	DEFAULT_CTYPE = "default"
	IMPORTS       = []string{
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
	}
	TEMPLATES = []string{
		"templates/connector/bolt.tmpl",
		"templates/connector/log.tmpl",
		"templates/connector/secret.tmpl",
		"templates/connector/session.tmpl",
		"templates/connector/token.tmpl",
		"templates/connector/user.tmpl",
		"templates/rest/assets.tmpl",
		"templates/rest/help.tmpl",
		"templates/rest/main.tmpl",
		"templates/rest/prometheus.init.tmpl",
		"templates/rest/prometheus.tmpl",
		"templates/rest/schedule.tmpl",
		"templates/rest/service.tmpl",
		"templates/rest/time_zone_location.tmpl",
		"templates/service/evauth.tmpl",
		"templates/service/evbolt.tmpl",
		"templates/service/evlog.tmpl",
		"templates/service/evschedule.tmpl",
	}
	VARS = map[string]interface{}{
		"Package":                "main",
		"DefaultAddress":         "127.0.0.1:9090",
		"UsageFunc":              "EVUsage",
		"Version":                "0.0.1",
		"Name":                   "EVBolt",
		"Description":            "EVBolt is a rest micro service which wrapps the golang bolt database",
		"Src":                    "https://git.evalgo.de:8443/",
		"DEBUG":                  true,
		"ENABLE_CROSS_ORIGIN":    true,
		"USE_PROMETHEUS":         true,
		"USE_EVTOKEN":            true,
		"TOKEN_STORAGE_URL":      "http://localhost:9090/0.0.1/eve/bolt",
		"TOKEN_STORAGE_DB":       "tokens.db",
		"TOKEN_STORAGE_BUCKET":   "tokens",
		"USE_EVLOG":              false,
		"EVLOG_URL":              "",
		"USE_EVLOG_API":          false,
		"USE_EVSESSION":          true,
		"SESSION_STORAGE_URL":    "http://localhost:9090/0.0.1/eve/bolt",
		"SESSION_STORAGE_DB":     "sessions.db",
		"SESSION_STORAGE_BUCKET": "sessions",
		"USE_EVSECRET":           true,
		"SECRET_STORAGE_URL":     "http://localhost:9090/0.0.1/eve/bolt",
		"SECRET_STORAGE_DB":      "secrets.db",
		"SECRET_STORAGE_BUCKET":  "secrets",
		"SECRET_ENC_KEY":         "TokenKeyEnc",
		"SECRET_SIG_KEY":         "TokenKeySig",
		"USE_EVUSER":             true,
		"USER_STORAGE_URL":       "http://localhost:9090/0.0.1/eve/bolt",
		"USER_STORAGE_DB":        "users.db",
		"USER_STORAGE_BUCKET":    "users",
		"USE_EVBOLT_API":         true,
		"USE_EVBOLT_AUTH":        true,
		"USE_LOGIN_API":          false,
		"SECRET_KEY_FOR_TOKEN":   "123456789012345678901234567890ab",
		"SECRET_SIG_FOR_TOKEN":   "sig.key.secret",
		"COOKIE_EXP_MINUTES":     15,
		"TOKEN_EXP_DAYS":         7,
		"USE_EVSCHEDULE":         false,
		"USE_EVSCHEDULE_API":     false,
		"URLS": []string{
			"/help",
			"/bolt",
			"/metrics",
		},
		"ROUTE_PATH_PREFIX": "/0.0.1/eve/",
	}
	COMMANDS = []*EVServiceCommand{
		NewEVServiceDefaultCommandFlags(),
	}
	MAIN = "EVREST"
)

type EVServiceConfig struct {
	Main      string
	Templates []string
	Imports   []string
	Vars      map[string]interface{}
	Commands  []*EVServiceCommand
}

type EVService interface {
	EVServiceConfiguration() *EVServiceConfig
}

type EVServiceConfigObj struct {
	Config *EVServiceConfig
}

func (tco *EVServiceConfigObj) NewEVServiceConfig(cType string) *EVServiceConfigObj {
	switch cType {
	case "rest":
		tco.Config = &EVServiceConfig{
			Main:      MAIN,
			Imports:   IMPORTS,
			Templates: TEMPLATES,
			Commands:  COMMANDS,
			Vars:      VARS,
		}
	case "rest_all":
		VARS["USE_EVLOG"] = true
		VARS["USE_EVLOG_API"] = true
		VARS["EVLOG_URL"] = "http://localhost:9090/0.0.1/eve/evlog"
		cVars := VARS["URLS"].([]string)
		cVars = append(cVars, "login")
		cVars = append(cVars, "setup")
		cVars = append(cVars, "access")
		cVars = append(cVars, "logout")
		tco.Config = &EVServiceConfig{
			Main:      "EVREST",
			Imports:   IMPORTS,
			Templates: TEMPLATES,
			Commands:  COMMANDS,
			Vars:      VARS,
		}
	case "schedule":
		VARS["USE_EVTOKEN"] = false
		VARS["TOKEN_STORAGE_URL"] = ""
		VARS["TOKEN_STORAGE_DB"] = ""
		VARS["TOKEN_STORAGE_BUCKET"] = ""
		VARS["USE_EVLOG"] = false
		VARS["EVLOG_URL"] = ""
		VARS["USE_EVLOG_API"] = false
		VARS["USE_EVSESSION"] = false
		VARS["SESSION_STORAGE_URL"] = ""
		VARS["SESSION_STORAGE_DB"] = ""
		VARS["SESSION_STORAGE_BUCKET"] = ""
		VARS["USE_EVSECRET"] = false
		VARS["SECRET_STORAGE_URL"] = ""
		VARS["SECRET_STORAGE_DB"] = ""
		VARS["SECRET_STORAGE_BUCKET"] = ""
		VARS["SECRET_ENC_KEY"] = ""
		VARS["SECRET_SIG_KEY"] = ""
		VARS["USE_EVUSER"] = false
		VARS["USER_STORAGE_URL"] = ""
		VARS["USER_STORAGE_DB"] = ""
		VARS["USER_STORAGE_BUCKET"] = ""
		VARS["USE_EVBOLT_API"] = false
		VARS["USE_EVBOLT_AUTH"] = false
		VARS["USE_LOGIN_API"] = false
		VARS["SECRET_KEY_FOR_TOKEN"] = ""
		VARS["SECRET_SIG_FOR_TOKEN"] = ""
		VARS["COOKIE_EXP_MINUTES"] = 0
		VARS["TOKEN_EXP_DAYS"] = 0
		VARS["USE_EVSCHEDULE"] = true
		VARS["USE_EVSCHEDULE_API"] = true
		VARS["URLS"] = []string{
			"/help",
			"/evschedule",
			"/metrics",
		}
		VARS["ROUTE_PATH_PREFIX"] = "/0.0.1/eve/"
		IMPORTS = []string{
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
			"time",
			"sync",
		}
		tco.Config = &EVServiceConfig{
			Main:      MAIN,
			Imports:   IMPORTS,
			Templates: TEMPLATES,
			Commands:  COMMANDS,
			Vars:      VARS,
		}
	default:
		tco.Config = &EVServiceConfig{
			Main: "TestMain",
			Imports: []string{
				"evalgo.org/eve",
			},
			Templates: []string{
				"tests/test.main.tmpl",
			},
			Vars: map[string]interface{}{
				"TestVar": "TestValue",
			},
		}
	}
	return tco
}

func (tco *EVServiceConfigObj) EVServiceConfiguration() *EVServiceConfig {
	if tco.Config == nil {
		return tco.NewEVServiceConfig(DEFAULT_CTYPE).Config
	}
	return tco.Config
}
