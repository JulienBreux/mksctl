package client

import api "github.com/JulienBreux/mksctl/internal/mksctl/api/gen"

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
