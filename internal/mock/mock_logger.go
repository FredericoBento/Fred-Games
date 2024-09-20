package mock

type MockIoWriter struct {
}

func (w MockIoWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
