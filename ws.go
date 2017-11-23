package eve

// EVServiceConfigWS struct configuration for the web service
type EVServiceConfigWS struct {
	EVServiceConfig
	Config *EVServiceConfig
}

// EVServiceConfiguration returns the service configuration
func (ws *EVServiceConfigWS) EVServiceConfiguration() *EVServiceConfig {
	return ws.Config
}
