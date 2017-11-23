package eve

var (
	// DEFAULT_ADDRESS for the generated service to be used
	DEFAULT_ADDRESS = "127.0.0.1:9090"
)

// EVServiceFlag the struct which describes the flag used in the service
type EVServiceFlag struct {
	FName  string
	FType  string
	FValue string
	FDesc  string
}

// NewEVServiceFlagDebug returns the debug flag to be used
func NewEVServiceFlagDebug() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "debug",
		FType:  "bool",
		FValue: "false",
		FDesc:  "display debug information for the given command",
	}
}

// NewEVServiceFlagHelpHttp returns the help http default flag to be used
func NewEVServiceFlagHelpHttp() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "hhttp",
		FType:  "string",
		FValue: "",
		FDesc:  "display the help menu as a html website for the given command",
	}
}

// NewEVServiceFlagHttpAddress returns the address flag which should be used
func NewEVServiceFlagHttpAddress() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "address",
		FType:  "string",
		FValue: DEFAULT_ADDRESS,
		FDesc:  "address for the http service to run on the given command",
	}
}

// NewEVServiceFlagVersion returns the version flag to be used
func NewEVServiceFlagVersion() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "version",
		FType:  "string",
		FValue: VERSION,
		FDesc:  "version of the running command",
	}
}

// NewEVServiceFlagSslCrt returns the certificate flag to be used
func NewEVServiceFlagSslCrt() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "crt",
		FType:  "string",
		FValue: "",
		FDesc:  "path to the ssl certificate",
	}
}

// NewEVServiceFlagSslKey returns the certificate key flag to be used
func NewEVServiceFlagSslKey() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "key",
		FType:  "string",
		FValue: "",
		FDesc:  "path to the ssl private key",
	}
}

// NewEVServiceFlagWebRoot returns the webroot flag to be used for the assets
func NewEVServiceFlagWebRoot() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "webroot",
		FType:  "string",
		FValue: ".",
		FDesc:  "path to the webroot where all root assets [*.css,*.js,...] are stored",
	}
}

// NewEVServiceFlagEvWebRoot returns the webroot flat for the service specific assets to be used
func NewEVServiceFlagEvWebRoot() *EVServiceFlag {
	return &EVServiceFlag{
		FName:  "evwebroot",
		FType:  "string",
		FValue: ".",
		FDesc:  "path to the webroot where all service specific assets [*.css,*.js,...] are stored",
	}
}
