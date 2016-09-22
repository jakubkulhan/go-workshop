package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jakubkulhan/go-workshop/universal_logger"
	"github.com/jakubkulhan/go-workshop/universal_logger/loggers/stdout_logger"
)

var levelMap map[string]universal_logger.Level = map[string]universal_logger.Level{
	"info":    universal_logger.Info,
	"i":       universal_logger.Info,
	"warning": universal_logger.Warning,
	"warn":    universal_logger.Warning,
	"w":       universal_logger.Warning,
	"error":   universal_logger.Error,
	"err":     universal_logger.Error,
	"e":       universal_logger.Error,
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	logger := stdout_logger.NewStdoutLogger()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		parts := strings.SplitN(scanner.Text(), " ", 2)

		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, "not enough arguments")
			continue
		}

		level, ok := levelMap[parts[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknonw level [%s]\n", parts[0])
			continue
		}

		logger.Log(level, parts[1])
	}
}
