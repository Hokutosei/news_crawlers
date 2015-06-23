package main

import "github.com/yvasiyarov/gorelic"

// InitNewRelic start newrelic
func InitNewRelic() {
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = "0f9d8145340e20588fe6fe3df67349f6bcd806d8"
	agent.Run()
}
