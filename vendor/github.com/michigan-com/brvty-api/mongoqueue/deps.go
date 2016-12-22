package mongoqueue

type Logger interface {
	Printf(format string, v ...interface{})
}
