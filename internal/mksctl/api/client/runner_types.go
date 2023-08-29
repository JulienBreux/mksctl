package client

import api "github.com/JulienBreux/microcks-client-go"

var (
	RunnerTypes = []api.TestRunnerType{
		api.TestRunnerTypeASYNCAPISCHEMA,
		api.TestRunnerTypeGRAPHQLSCHEMA,
		api.TestRunnerTypeGRPCPROTOBUF,
		api.TestRunnerTypeHTTP,
		api.TestRunnerTypeOPENAPISCHEMA,
		api.TestRunnerTypePOSTMAN,
		api.TestRunnerTypeSOAPHTTP,
		api.TestRunnerTypeSOAPUI,
	}
)
