package multiresolver

import "github.com/redesblock/hop/core/logging"

func GetLogger(mr *MultiResolver) logging.Logger {
	return mr.logger
}

func GetCfgs(mr *MultiResolver) []ConnectionConfig {
	return mr.cfgs
}
