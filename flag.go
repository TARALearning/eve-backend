package eve

var (
	DEFAULT_ADDRESS = "127.0.0.1:9090"
)

type EVServiceFlag struct {
	FName  string
	FType  string
	FValue string
	FDesc  string
}

func NewEVServiceFlagDebug() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "debug",
		FType:  "bool",
		FValue: "false",
		FDesc:  "display debug information for the given command",
	}
}

func NewEVServiceFlagHelpHttp() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "hhttp",
		FType:  "string",
		FValue: "",
		FDesc:  "display the help menu as a html website for the given command",
	}
}

func NewEVServiceFlagHttpAddress() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "address",
		FType:  "string",
		FValue: DEFAULT_ADDRESS,
		FDesc:  "address for the http service to run on the given command",
	}
}

func NewEVServiceFlagVersion() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "version",
		FType:  "string",
		FValue: VERSION,
		FDesc:  "version of the running command",
	}
}

func NewEVServiceFlagSslCrt() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "crt",
		FType:  "string",
		FValue: "",
		FDesc:  "path to the ssl certificate",
	}
}

func NewEVServiceFlagSslKey() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "key",
		FType:  "string",
		FValue: "",
		FDesc:  "path to the ssl private key",
	}
}

func NewEVServiceFlagWebRoot() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "webroot",
		FType:  "string",
		FValue: ".",
		FDesc:  "path to the webroot where all root assets [*.css,*.js,...] are stored",
	}
}

func NewEVServiceFlagEvWebRoot() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "evwebroot",
		FType:  "string",
		FValue: ".",
		FDesc:  "path to the webroot where all service specific assets [*.css,*.js,...] are stored",
	}
}
