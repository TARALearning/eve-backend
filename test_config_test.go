package eve

// ConfigFotTest configures the service configuration for the unit tests
func ConfigForTest(cType string) *EVServiceConfigObj {
	defaultCType = cType
	tco := &EVServiceConfigObj{}
	tco.EVServiceConfiguration()
	return tco
}
