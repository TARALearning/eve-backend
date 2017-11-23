package eve

// EVServiceCommand contains a command to be used in the REST service main files
type EVServiceCommand struct {
	Name  string
	Desc  string
	Flags []*EVServiceFlag
}

// NewEVServiceDefaultCommandFlags returns the default service commands and flags
func NewEVServiceDefaultCommandFlags() *EVServiceCommand {
	return &EVServiceCommand{
		Name: "help",
		Desc: "help command to be used for detailed information",
		Flags: []*EVServiceFlag{
			NewEVServiceFlagHelpHttp(),
			NewEVServiceFlagWebRoot(),
			NewEVServiceFlagEvWebRoot(),
			NewEVServiceFlagDebug(),
			NewEVServiceFlagVersion(),
			NewEVServiceFlagHttpAddress(),
			NewEVServiceFlagSslCrt(),
			NewEVServiceFlagSslKey(),
		},
	}
}
