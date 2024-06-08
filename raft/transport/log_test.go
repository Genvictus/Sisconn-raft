package transport

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"google.golang.org/grpc"
)

func TestLogPrint(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "[Raft] Client: ", 0)

	message := "Test message"
	LogPrint(logger, message)

	expectedOutput := "[Raft] Client: Test message\n"

	if buf.String() != expectedOutput {
		t.Errorf("LogPrint failed, Expected: %s, but got: %s", expectedOutput, buf.String())
	}
}

func TestLoggingInterceptorServer(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "[Raft] Client: ", 0)

	ctx := context.Background()
	req := "test request"
	resp := "test response"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.TestService/TestMethod",
	}

	// Mock the handler
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return resp, nil
	}

	tests := []struct {
		testName    string
		interceptor grpc.UnaryServerInterceptor
		logPrefix   string
	}{
		{"WithLogger", LoggingInterceptorServer(logger), "[Raft] Client: "},
		{"WithoutLogger", LoggingInterceptorServer(nil), ""},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Clear buffer
			buf.Reset()

			// Call the interceptor
			result, err := tt.interceptor(ctx, req, info, handler)
			if err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if resp != result {
				t.Errorf("Expected response: %v, but got: %v", resp, result)
			}

			// Checking log output
			logOutput := buf.String()
			if !strings.HasPrefix(logOutput, tt.logPrefix) {
				t.Errorf("Expected log output to has prefix: %s, but got: %s", tt.logPrefix, logOutput)
			}

			if tt.testName == "WithLogger" {
				if !strings.Contains(logOutput, "RPC response: /test.TestService/TestMethod") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "RPC response: /test.TestService/TestMethod", logOutput)
				}

				if !strings.Contains(logOutput, "Duration:") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "Duration:", logOutput)
				}

				if !strings.Contains(logOutput, "Error: <nil>") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "Error: <nil>", logOutput)
				}
			}
		})
	}
}

func TestLoggingInterceptorClient(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "[Raft] Client: ", 0)

	ctx := context.Background()
	method := "/example.Method"
	req := "test request"
	resp := "test response"
	clientConn := &grpc.ClientConn{}

	tests := []struct {
		testName    string
		interceptor grpc.UnaryClientInterceptor
		logPrefix   string
	}{
		{"WithLogger", LoggingInterceptorClient(logger), "[Raft] Client: "},
		{"WithoutLogger", LoggingInterceptorClient(nil), ""},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Clear buffer
			buf.Reset()

			// Call the interceptor
			err := tt.interceptor(ctx, method, req, resp, clientConn, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			})

			if err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		
			// Checking log output
			logOutput := buf.String()
			if !strings.HasPrefix(logOutput, tt.logPrefix) {
				t.Errorf("Expected log output to has prefix: %s, but got: %s", tt.logPrefix, logOutput)
			}

			if tt.testName == "WithLogger" {
				if !strings.Contains(logOutput, "RPC request: /example.Method") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "RPC request: /example.Method", logOutput)
				}

				if !strings.Contains(logOutput, "Duration:") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "Duration:", logOutput)
				}

				if !strings.Contains(logOutput, "Error: <nil>") {
					t.Errorf("Expected log output to contain: %s, but got: %s", "Error: <nil>", logOutput)
				}
			}
		})
	}
}
