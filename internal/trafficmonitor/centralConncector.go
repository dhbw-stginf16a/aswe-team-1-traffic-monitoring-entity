package trafficmonitor

// CentralConnector abstracts connection to central node
type CentralConnector struct {
	url string
}

// NewCentralConnector create a new CentralConnector
func NewCentralConnector(url string) *CentralConnector {
	return &CentralConnector{url}
}

// Register to central node
func (con CentralConnector) Register(hostname string) error {
	return nil
}

// GetGlobalPrefs from central node
func (con CentralConnector) GetGlobalPrefs() {

}

// GetUserPrefs from central node
func (con CentralConnector) GetUserPrefs(userid string) {

}
