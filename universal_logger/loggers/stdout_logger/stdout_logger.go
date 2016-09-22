package stdout_logger

import (
	"fmt"

	"github.com/jakubkulhan/go-workshop/universal_logger"
)

var printLevelMap map[universal_logger.Level]string = map[universal_logger.Level]string{
	universal_logger.Info:    "INFO",
	universal_logger.Warning: "WARN",
	universal_logger.Error:   "ERROR",
}

type stdoutLogger struct {
}

func (l *stdoutLogger) Log(level universal_logger.Level, message string) {
	fmt.Printf("[%s] %s\n", printLevelMap[level], message)
}

func NewStdoutLogger() universal_logger.Logger {
	return &stdoutLogger{}
}
