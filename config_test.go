package eve

func ConfigForTest(cType string) *EVServiceConfigObj {
	DEFAULT_CTYPE = cType
	tco := &EVServiceConfigObj{}
	tco.EVServiceConfiguration()
	return tco
}
