package service

import (
	"runtime"
	"strconv"
	"strings"
)

// GetCaller returns a "file:line" string for the caller at the given skip level.
// skip=0 returns the caller of GetCaller, skip=1 returns the caller's caller, etc.
func GetCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "???:0"
	}
	// Trim to last 2 path segments for readability
	parts := strings.Split(file, "/")
	if len(parts) > 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	return file + ":" + strconv.Itoa(line)
}

// getCaller is an unexported alias for internal use within the service package.
func getCaller(skip int) string {
	return GetCaller(skip + 1)
}
