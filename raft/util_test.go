package raft

import (
	"fmt"
	"testing"
	"time"
)

func TestRandMs(t *testing.T) {
	tests := []struct {
		minDuration time.Duration
		maxDuration time.Duration
	}{
		{
			minDuration: 1000 * time.Millisecond,
			maxDuration: 2000 * time.Millisecond,
		},
		{
			minDuration: 10 * time.Millisecond,
			maxDuration: 50 * time.Millisecond,
		},
		{
			minDuration: 500 * time.Millisecond,
			maxDuration: 1000 * time.Millisecond,
		},
		{
			minDuration: 100 * time.Millisecond,
			maxDuration: 1000 * time.Millisecond,
		},
		{
			minDuration: 2000 * time.Millisecond,
			maxDuration: 3000 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run("TestRandMs", func(t *testing.T) {
			randDuration := randMs(tt.minDuration, tt.maxDuration)
			if randDuration < tt.minDuration || randDuration > tt.maxDuration {
				t.Errorf("Expected duration between %d and %d, got %d", tt.minDuration, tt.maxDuration, randDuration)
			}
		})
	}
}

func TestRedirectStdout(t *testing.T) {
	closeBuf := redirectStdout()

	fmt.Print("Hello, World!")

	out := closeBuf()

	string_out := string(out)

	if string_out != "Hello, World!" {
		t.Errorf("Expected output: %s, got %s", "Hello, World!", string_out)
	}

	fmt.Println("Test passed!")
	
	if string_out != "Hello, World!" {
		t.Errorf("Expected output: %s, got %s", "Hello, World!", string_out)
	}
}
