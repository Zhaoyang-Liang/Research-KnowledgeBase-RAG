// anyu25.go — AnyuWang 2025 subring-secret bootstrapping runner.
//
// Calls the pre-compiled binary at subkey25BinPath with the I1 parameter set
// (eps=-16, scale=35, degree=127, r=3, scalestc=33, scalects=52) used in their
// published results. Output is streamed to stdout in real time; the last-line
// "Avg L1 error" and "Total:" fields are parsed and returned for the cross-method
// Summary in main().
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// subkey25BinPath is the pre-compiled AnyuWang binary (uses their private Lattigo fork).
// Path is resolved relative to this source file's directory at runtime.
const subkey25RelPath = "../../下一步实验/别人的实验对比main-25AnyuWang/subkey25/main"

func subkey25BinPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	return filepath.Join(dir, subkey25RelPath)
}

// RunAnyu25 executes the AnyuWang 2025 I1-parameter binary, streams its output,
// and returns the avg L1 precision (bits) and avg total wall-clock time.
// nRepeat is passed as -repeat=N to the binary.
func RunAnyu25(nRepeat int) (avgL1Err float64, avgTotal time.Duration) {
	if nRepeat <= 0 {
		nRepeat = 1
	}

	binPath := subkey25BinPath()
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "[subKey25] binary not found at: %s\n", binPath)
		fmt.Fprintf(os.Stderr, "           Please compile it first:\n")
		fmt.Fprintf(os.Stderr, "           cd %s && go build -o main .\n",
			filepath.Dir(binPath))
		return 0, 0
	}

	args := []string{
		"-eps=-16",
		"-scale=35",
		"-degree=127",
		"-r=3",
		"-scalestc=33",
		"-scalects=52",
		fmt.Sprintf("-repeat=%d", nRepeat),
	}
	fmt.Printf("[subKey25] running: %s %s\n\n", filepath.Base(binPath), strings.Join(args, " "))

	cmd := exec.Command(binPath, args...)
	cmd.Dir = filepath.Dir(binPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[subKey25] pipe error: %v\n", err)
		return 0, 0
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "[subKey25] start error: %v\n", err)
		return 0, 0
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if strings.HasPrefix(line, "Avg L1 error is ") {
			if v, err := strconv.ParseFloat(strings.TrimPrefix(line, "Avg L1 error is "), 64); err == nil {
				avgL1Err = v
			}
		}

		if strings.HasPrefix(line, "Avg time...") {
			for _, part := range strings.Split(line, ",") {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "Total: ") {
					durStr := strings.TrimPrefix(part, "Total: ")
					if d, err := time.ParseDuration(durStr); err == nil {
						avgTotal = d
					}
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "[subKey25] exited with error: %v\n", err)
	}

	return avgL1Err, avgTotal
}

// runAnyu25Capturing is the shim called from main(); same as RunAnyu25.
func runAnyu25Capturing(nRepeat int) (float64, time.Duration) {
	return RunAnyu25(nRepeat)
}

// Silence "imported and not used" for packages only needed in the old stub.
var _ = math.Log2
