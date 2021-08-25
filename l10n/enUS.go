package l10n

var EN_US_STRING_MAP = Localzation{
	ClientTitle:     "AKASH NETWORK LITE CLIENT",
	AppDescription:  "Akash Command Center.\n\nAkash is a peer-to-peer marketplace for computing resources and \na deployment platform for heavily distributed applications. \nFind out more at https://akash.network\n\n",
	WelcomeFirstRun: "Welcome, this is the first run, prepping...",
	Command: map[string]string{
		"welcome":                         "Runs the welcome flow of the client",
		"config":                          "Modify configuration options for the client",
		"config.network":                  "Get the current information about the network configuration for the client",
		"config.network.update":           "Retrieve the latet network settings from the Akash repository",
		"config.network.defaultSelection": "Please select an available (default) network",
	},
}
