package eve

import (
	"bytes"
	"strings"
	"text/template"
)

var (
	// EVServiceUsageLines defines the number of usage lines to be displayed
	EVServiceUsageLines = 50
)

// EVGenImports specifies the import packages for code generation
// it transforms the given slice of strings into import commands in go
func EVGenImports(importsPkgs []string) string {
	newImports := []string{}
	// append prefix and postfix " to the import names
	for k := range importsPkgs {
		switch importsPkgs[k] {
		case "github.com/mitchellh/go-ps":
			newImports = append(newImports, "gops \""+importsPkgs[k]+"\"")
		default:
			newImports = append(newImports, "\""+importsPkgs[k]+"\"")
		}
	}
	return strings.Join(newImports, "\n")
}

// EVGenCommandFlags generates the from the given commands slice the flags go code strings
func EVGenCommandFlags(cmds []*EVServiceCommand) string {
	cmdsString := []string{}
	for k := range cmds {
		for kk := range cmds[k].Flags {
			switch cmds[k].Flags[kk].FType {
			case "bool":
				cmdsString = append(cmdsString, "flags.BoolVar(&"+cmds[k].Flags[kk].FName+",\""+cmds[k].Flags[kk].FName+"\","+cmds[k].Flags[kk].FValue+",\""+cmds[k].Flags[kk].FDesc+"\")")
			case "string":
				cmdsString = append(cmdsString, "flags.StringVar(&"+cmds[k].Flags[kk].FName+",\""+cmds[k].Flags[kk].FName+"\",\""+cmds[k].Flags[kk].FValue+"\",\""+cmds[k].Flags[kk].FDesc+"\")")
			}
		}
	}
	return strings.Join(cmdsString, "\n")
}

// EVGenCommandFlagsVars generates the var code strings from the given commands slice
func EVGenCommandFlagsVars(cmds []*EVServiceCommand) string {
	cmdsVars := []string{}
	for k := range cmds {
		for kk := range cmds[k].Flags {
			switch cmds[k].Flags[kk].FType {
			case "bool":
				cmdsVars = append(cmdsVars, "var "+cmds[k].Flags[kk].FName+" "+cmds[k].Flags[kk].FType+" = "+cmds[k].Flags[kk].FValue)
			case "string":
				cmdsVars = append(cmdsVars, "var "+cmds[k].Flags[kk].FName+" "+cmds[k].Flags[kk].FType+" = \""+cmds[k].Flags[kk].FValue+"\"")
			}
		}
	}
	return strings.Join(cmdsVars, "\n")
}

// EVGenUsageAppendMoreLines appends more usage lines 2 times more default which is EVServiceUsageLines
func EVGenUsageAppendMoreLines(lines []string) []string {
	newLines := make([]string, 0)
	lines = EVGenUsageRemoveEmptyLines(lines)
	if len(lines) >= EVServiceUsageLines {
		EVServiceUsageLines = EVServiceUsageLines * 2
		newLines = make([]string, EVServiceUsageLines)
		for k := range lines {
			newLines[k] = lines[k]
		}
		return newLines
	}
	return lines
}

// EVGenUsageRemoveEmptyLines removes empty usage lines
func EVGenUsageRemoveEmptyLines(lines []string) []string {
	last := EVServiceUsageLines
	for k := range lines {
		if lines[k] == "" {
			last = k - 1
			break
		}
	}
	if last < EVServiceUsageLines {
		newLines := make([]string, last)
		for k := range lines {
			if k < last {
				newLines[k] = lines[k]
			}
		}
		return newLines
	}
	return lines
}

// EVGenCommandUsage generates the default usage lines to be displayed in the usage/help text
func EVGenCommandUsage(cmds []*EVServiceCommand, vars map[string]interface{}) string {
	usage := make([]string, EVServiceUsageLines)
	usage[0] = vars["Name"].(string) + "\n"
	usage[1] = "Version: " + vars["Version"].(string) + "\n"
	usage[2] = "Description: " + vars["Description"].(string) + "\n"
	usage[3] = "Src: " + vars["Src"].(string) + "\n"
	usage[4] = "Usage: [command] -[argument_label] [argument_value]...\n"
	usage[5] = "Commands: \n"
	sk := 5
	skk := 0
	for k := range cmds {
		if skk == 0 {
			sk = sk + k
		} else {
			sk = skk + 1
		}
		usage[sk] = "command: " + cmds[k].Name + "\n"
		usage[(sk + 1)] = "attributes:"
		for kk := range cmds[k].Flags {
			skk = sk + kk + 2
			usage[skk] = cmds[k].Flags[kk].FType + ": -" + cmds[k].Flags[kk].FName + " \"" + cmds[k].Flags[kk].FValue + "\" " + cmds[k].Flags[kk].FDesc
		}
	}
	usage = EVGenUsageAppendMoreLines(usage)
	return strings.Join(usage, "\n")
}

// EVGenMain generates the main go code file with the default settings
func EVGenMain(srv EVService) ([]byte, error) {
	funcMap := template.FuncMap{
		"imports":          EVGenImports,
		"commandFlags":     EVGenCommandFlags,
		"commandFlagsVars": EVGenCommandFlagsVars,
		"commandUsage":     EVGenCommandUsage,
	}
	config := srv.EVServiceConfiguration()
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(config.Templates...)
	if err != nil {
		return nil, err
	}
	// check if vars is set if not set it
	if config.Vars == nil {
		config.Vars = map[string]interface{}{}
	}
	// check if Commands is set
	if config.Commands != nil {
		config.Vars["Commands"] = config.Commands
	}
	// check if imports is set
	if config.Imports != nil {
		config.Vars["Imports"] = config.Imports
	}
	// check if the default address is in the commands and in the vars the same
	// if not set it to the value from the vars which is set by the user
	if cmds, ok := config.Vars["Commands"].([]*EVServiceCommand); ok {
		if addr, ok := config.Vars["DefaultAddress"]; ok {
			for ck := range cmds {
				cmd := cmds[ck]
				for k := range cmd.Flags {
					if cmd.Flags[k].FName == "address" {
						if cmd.Flags[k].FValue != addr {
							cmd.Flags[k].FValue = addr.(string)
						}
					}
				}
			}
		}
	}

	buff := bytes.NewBuffer(nil)
	err = tmpl.ExecuteTemplate(buff, config.Main, config.Vars)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
