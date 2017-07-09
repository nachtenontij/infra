package server

var Settings struct {
	DatabaseAddress string
	DatabaseName    string
	BindAddress     string

	// Session key to be used to authorize if there are no users yet.
	GenesisSessionKey string
}
