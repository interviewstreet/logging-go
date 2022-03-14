package request

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func assertFields(t *testing.T, output, key, value string) {
	compareWith := fmt.Sprintf(`"%s":"%s"`, key, value)
	if !strings.Contains(output, compareWith) {
		t.Errorf("Test failed for `%s` key", key)
	}
}

// TODO: Below test is commented due to logger memory sink conflicting with TestGinMiddleware
//func TestLogger(t *testing.T) {
//	// Create a sink instance, and register it with zap for the "memory"
//	// protocol.
//	sink := &MemorySink{new(bytes.Buffer)}
//	_ = zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
//		return sink, nil
//	})
//
//	log := createNewLogger("test", &core.LoggerOptions{Env: "test", OutputPath: "memory://"})
//
//	// Test with default generated context-id
//	log.Infow("Testing Logs", "test_key", "test_value")
//
//	// Assert sink contents
//	output := sink.String()
//	t.Logf("output = %s", output)
//
//	assertFields(t, output, "test_key", "test_value")
//	assertFields(t, output, "environment", "test")
//	assertFields(t, output, "namespace", "test")
//}
