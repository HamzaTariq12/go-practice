package types

// Note Receiver funcions should be in same struct package
type LogEntry struct {
	Message string
	Bytes   int64
}

func (l LogEntry) SizeInBytes() int64 {
	return l.Bytes
}
