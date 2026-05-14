package monitorsvc

import (
	"math"
	"path/filepath"
	"strings"
	"time"
)

func timeString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func nowString() string {
	return time.Now().Format(time.RFC3339)
}

func millisSince(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}

func roundPercent(v float64) float64 {
	return math.Round(v*100) / 100
}

func safeError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	if len([]rune(msg)) <= 200 {
		return msg
	}
	runes := []rune(msg)
	return string(runes[:200])
}

func binaryName(path string) string {
	name := filepath.Base(path)
	if name == "." || name == string(filepath.Separator) {
		return ""
	}
	return name
}
