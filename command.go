package eve

type EVServiceCommand struct {
	Name  string
	Desc  string
	Flags []*EVServiceFlag
}

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
