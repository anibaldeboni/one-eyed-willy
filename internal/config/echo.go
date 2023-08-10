package config

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
)

func GetEchoLogConfig() middleware.LoggerConfig {
	echoLogConf := middleware.DefaultLoggerConfig
	echoLogConf.CustomTimeFormat = time.RFC3339
	// echoLogConf.Format = `{ "@timestamp": "${time_rfc3339}" "request_id": "${id}", "latency": ${latency}, "status": "${status}", "method": "${method}", "uri": "${uri}", "error": "${error}" }`
	return echoLogConf
}
