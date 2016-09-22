package universal_logger

type Level int

const (
	Info Level = iota
	Warning
	Error
)

type Logger interface {
	Log(level Level, message string)
}
