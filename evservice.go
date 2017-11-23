package eve

var (
	// defaultCType defines the default config type to be created
	defaultCType = "default"
	// urls defines the default urls to be used
	urls = []string{
		"/help",
		"/bolt",
		"/metrics",
	}
	// imports defines the default import packages to be used in the generated service
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
	}
	// TEMPLATES defines the templates that should be used for code generation
	templates = []string{
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
	// VARS defines the variables that should be used during the code generation
	vars = map[string]interface{}{
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
		"URLS":                   urls,
		"ROUTE_PATH_PREFIX":      "/0.0.1/eve/",
	}
	// COMMANDS defines the default commands that should be used in the generated service code
	commands = []*EVServiceCommand{
		NewEVServiceDefaultCommandFlags(),
	}

	// MAIN i don't know if we need it any longer todo check if this variable is needed
	main = "EVREST"
)

// EVServiceConfig defines the service configuration struct
type EVServiceConfig struct {
	Main      string
	Templates []string
	Imports   []string
	Vars      map[string]interface{}
	Commands  []*EVServiceCommand
}

// EVService defines the EVService interface to be implemented
type EVService interface {
	EVServiceConfiguration() *EVServiceConfig
}

// EVServiceConfigObj contents the service configuration object
type EVServiceConfigObj struct {
	Config *EVServiceConfig
}

// NewEVServiceConfig creates a new service configuration with the default values
func (tco *EVServiceConfigObj) NewEVServiceConfig(cType string) *EVServiceConfigObj {
	switch cType {
	case "rest":
		tco.Config = &EVServiceConfig{
			Main:      main,
			Imports:   imports,
			Templates: templates,
			Commands:  commands,
			Vars:      vars,
		}
	case "rest_all":
		vars["USE_EVLOG"] = true
		vars["USE_EVLOG_API"] = true
		vars["EVLOG_URL"] = "http://localhost:9090/0.0.1/eve/evlog"
		urls = append(urls, "login")
		urls = append(urls, "setup")
		urls = append(urls, "access")
		urls = append(urls, "logout")
		vars["URLS"] = urls
		tco.Config = &EVServiceConfig{
			Main:      "EVREST",
			Imports:   imports,
			Templates: templates,
			Commands:  commands,
			Vars:      vars,
		}
	case "schedule":
		vars["USE_EVTOKEN"] = false
		vars["TOKEN_STORAGE_URL"] = ""
		vars["TOKEN_STORAGE_DB"] = ""
		vars["TOKEN_STORAGE_BUCKET"] = ""
		vars["USE_EVLOG"] = false
		vars["EVLOG_URL"] = ""
		vars["USE_EVLOG_API"] = false
		vars["USE_EVSESSION"] = false
		vars["SESSION_STORAGE_URL"] = ""
		vars["SESSION_STORAGE_DB"] = ""
		vars["SESSION_STORAGE_BUCKET"] = ""
		vars["USE_EVSECRET"] = false
		vars["SECRET_STORAGE_URL"] = ""
		vars["SECRET_STORAGE_DB"] = ""
		vars["SECRET_STORAGE_BUCKET"] = ""
		vars["SECRET_ENC_KEY"] = ""
		vars["SECRET_SIG_KEY"] = ""
		vars["USE_EVUSER"] = false
		vars["USER_STORAGE_URL"] = ""
		vars["USER_STORAGE_DB"] = ""
		vars["USER_STORAGE_BUCKET"] = ""
		vars["USE_EVBOLT_API"] = false
		vars["USE_EVBOLT_AUTH"] = false
		vars["USE_LOGIN_API"] = false
		vars["SECRET_KEY_FOR_TOKEN"] = ""
		vars["SECRET_SIG_FOR_TOKEN"] = ""
		vars["COOKIE_EXP_MINUTES"] = 0
		vars["TOKEN_EXP_DAYS"] = 0
		vars["USE_EVSCHEDULE"] = true
		vars["USE_EVSCHEDULE_API"] = true
		vars["URLS"] = []string{
			"/help",
			"/evschedule",
			"/metrics",
		}
		vars["ROUTE_PATH_PREFIX"] = "/0.0.1/eve/"
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
			"time",
			"sync",
		}
		tco.Config = &EVServiceConfig{
			Main:      main,
			Imports:   imports,
			Templates: templates,
			Commands:  commands,
			Vars:      vars,
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

// EVServiceConfiguration returns the configuration of the Service
func (tco *EVServiceConfigObj) EVServiceConfiguration() *EVServiceConfig {
	if tco.Config == nil {
		return tco.NewEVServiceConfig(defaultCType).Config
	}
	return tco.Config
}

// SrvConfigMain returns the service config main value
func SrvConfigMain() string {
	return main
}

// SrvConfigTemplates returns the service config templates slice of strings
func SrvConfigTemplates() []string {
	return templates
}

// SrvConfigCommands returns the default service commands
func SrvConfigCommands() []*EVServiceCommand {
	return commands
}

// SetDefaultCType sets the default cType
func SetDefaultCType(cType string) {
	defaultCType = cType
}
