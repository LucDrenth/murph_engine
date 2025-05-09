package log

type noOpLogger struct{}

var _ Logger = &noOpLogger{}

// NoOp returns a no-op logger, meaning it does nothing when its methods are called.
func NoOp() noOpLogger {
	return noOpLogger{}
}

func (*noOpLogger) Debug(string)     {}
func (*noOpLogger) DebugOnce(string) {}
func (*noOpLogger) Info(string)      {}
func (*noOpLogger) InfoOnce(string)  {}
func (*noOpLogger) Warn(string)      {}
func (*noOpLogger) WarnOnce(string)  {}
func (*noOpLogger) Error(string)     {}
func (*noOpLogger) ErrorOnce(string) {}
func (*noOpLogger) Trace(string)     {}
func (*noOpLogger) TraceOnce(string) {}
func (*noOpLogger) ClearStorage()    {}
