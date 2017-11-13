package eve

type EVServiceConfigWS struct {
	EVServiceConfig
	Config *EVServiceConfig
}

func (ws *EVServiceConfigWS) EVServiceConfiguration() *EVServiceConfig {
	return ws.Config
}
