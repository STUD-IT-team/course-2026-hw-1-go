package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestShiftProgram(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "shift")
	if output, err := exec.Command("go", "build", "-buildvcs=false", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, output)
	}

	cases := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{"several orders", "4\n700 1000 1300 0\n", "3000 1300 2\n", false},
		{"no orders", "0\n", "0 0 0\n", false},
		{"surrounding whitespace", " \t2\n0009\t100000 \n", "100009 100000 1\n", false},
		{"empty input", "", "", true},
		{"non-integer count", "two\n", "", true},
		{"negative count", "-1\n", "", true},
		{"count above limit", "1001\n", "", true},
		{"missing amount", "2\n800\n", "", true},
		{"extra amount", "0\n800\n", "", true},
		{"malformed amount", "1\n1000rub\n", "", true},
		{"negative amount", "1\n-10\n", "", true},
		{"amount above limit", "1\n100001\n", "", true},
		{"integer overflow", "1\n9999999999999999999999999\n", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, binary)
			cmd.Stdin = strings.NewReader(tc.input)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			output, err := cmd.Output()
			if ctx.Err() == context.DeadlineExceeded {
				t.Fatalf("program did not exit within 5 seconds; stderr: %q", stderr.String())
			}
			if tc.wantError {
				checkProgramError(t, err, stderr.String())
				if len(output) != 0 {
					t.Errorf("stdout = %q, want empty output", output)
				}
				return
			}
			if err != nil {
				t.Fatalf("program failed: %v; stderr: %q", err, stderr.String())
			}
			if string(output) != tc.want {
				t.Errorf("stdout = %q, want %q", output, tc.want)
			}
		})
	}

	t.Run("stdout write error", func(t *testing.T) {
		full, err := os.OpenFile("/dev/full", os.O_WRONLY, 0)
		if err != nil {
			t.Fatalf("open /dev/full: %v", err)
		}
		defer full.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary)
		cmd.Stdin = strings.NewReader("1\n1000\n")
		cmd.Stdout = full
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		checkProgramError(t, cmd.Run(), stderr.String())
	})
}

func checkProgramError(t *testing.T, err error, stderr string) {
	t.Helper()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 1 {
		t.Errorf("want exit code 1, got %v; stderr: %q", err, stderr)
	}
	if strings.TrimSpace(stderr) == "" {
		t.Error("stderr is empty, want an error message")
	}
}
