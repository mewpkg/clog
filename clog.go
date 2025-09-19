// Package clog provides coloured logging.
package clog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/mewpkg/term"
)

// --- [ log levels ] ----------------------------------------------------------

// Level specifies a log level which denotes the importance or severity of a log
// event.
type Level int

// Common log levels.
const (
	// LevelDebug is used for debug messages (Magenta).
	LevelDebug Level = -4
	// LevelInfo is used for informational messages (Cyan).
	LevelInfo Level = 0
	// LevelWarn is used for non-fatal warnings (Red).
	LevelWarn Level = 4
	// LevelError is used for fatal errors (Red).
	LevelError Level = 8
)

var (
	// mu is a mutex for concurrent access to activeLevel.
	mu sync.Mutex
	// activeLevel specifies the active log level at package and function
	// granularity.
	activeLevel = make(map[string]Level)
)

// SetPathLevel sets the log level of the given path at package
// (e.g. "github.com/user/repo/pkg") or function
// (e.g. "github.com/user/repo/pkg.Func") granularity.
//
// For function ganularity of leaf node functions, function inlining may have to
// be disabled (use the `//go:noinline` build tag).
func SetPathLevel(path string, level Level) {
	mu.Lock()
	defer mu.Unlock()
	activeLevel[path] = level
}

// PathLevel returns the current log level of the given path at package or
// function granularity (or a parent directory thereof), and a boolean
// indicating whether the log level was set.
func PathLevel(path string) (Level, bool) {
	mu.Lock()
	defer mu.Unlock()
	for {
		if level, ok := activeLevel[path]; ok {
			return level, true
		}
		// check if path level is specified for parent directory.
		dir := filepath.Dir(path)
		if len(dir) == 0 || dir == path {
			break
		}
		path = dir
	}
	return 0, false
}

// skip reports whether to skip log output of the given log level for the
// package path and function path of the caller.
func skip(cur Level) bool {
	pkgPath, funcPath := getQualifiedPaths()
	if funcLevel, ok := PathLevel(funcPath); ok {
		return funcLevel > cur
	}
	if pkgLevel, ok := PathLevel(pkgPath); ok {
		return pkgLevel > cur
	}
	return false
}

// mainPrefixName specifies the prefix name used for the 'main' package.
var mainPrefixName string

// SetMainPrefixName sets the prefix name used for the 'main' package.
func SetMainPrefixName(name string) {
	mainPrefixName = name
}

// --- [ debug ] ---------------------------------------------------------------

// outputMutex is a mutex for concurrent writes to output writers.
var outputMutex sync.Mutex

var (
	// debugOutput specifies the output writer of debug messages.
	debugOutput io.Writer = os.Stderr

	// debugUsePrefix specifies whether to use a prefix for debug messages.
	debugUsePrefix = true
)

// SetDebugOutput sets the output writer of debug messages.
func SetDebugOutput(w io.Writer) {
	outputMutex.Lock()
	defer outputMutex.Unlock()
	debugOutput = w
}

// SetDebugPrefix sets whether to use a prefix for debug messages.
func SetDebugPrefix(usePrefix bool) {
	debugUsePrefix = usePrefix
}

// Debug outputs the given debug message to standard error.
func Debug(args ...any) {
	if skip(LevelDebug) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if debugUsePrefix {
		prefix := getPrefix(term.MagentaBold)
		fmt.Fprint(debugOutput, prefix)
	}
	fmt.Fprint(debugOutput, args...)
	fmt.Fprintln(debugOutput)
}

// Debugf outputs the given debug message to standard error.
func Debugf(format string, args ...any) {
	if skip(LevelDebug) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if debugUsePrefix {
		prefix := getPrefix(term.MagentaBold)
		fmt.Fprint(debugOutput, prefix)
	}
	fmt.Fprintf(debugOutput, format, args...)
	fmt.Fprintln(debugOutput)
}

// Debugln outputs the given debug message to standard error.
func Debugln(args ...any) {
	if skip(LevelDebug) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if debugUsePrefix {
		prefix := getPrefix(term.MagentaBold)
		fmt.Fprint(debugOutput, prefix)
	}
	fmt.Fprintln(debugOutput, args...)
}

// --- [ info ] ----------------------------------------------------------------

var (
	// infoOutput specifies the output writer of info messages.
	infoOutput io.Writer = os.Stderr

	// infoUsePrefix specifies whether to use a prefix for info messages.
	infoUsePrefix = true
)

// SetInfoOutput sets the output writer of info messages.
func SetInfoOutput(w io.Writer) {
	outputMutex.Lock()
	defer outputMutex.Unlock()
	infoOutput = w
}

// SetInfoPrefix sets whether to use a prefix for info messages.
func SetInfoPrefix(usePrefix bool) {
	infoUsePrefix = usePrefix
}

// Info outputs the given info message to standard error.
func Info(args ...any) {
	if skip(LevelInfo) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if infoUsePrefix {
		prefix := getPrefix(term.CyanBold)
		fmt.Fprint(infoOutput, prefix)
	}
	fmt.Fprint(infoOutput, args...)
	fmt.Fprintln(infoOutput)
}

// Infof outputs the given info message to standard error.
func Infof(format string, args ...any) {
	if skip(LevelInfo) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if infoUsePrefix {
		prefix := getPrefix(term.CyanBold)
		fmt.Fprint(infoOutput, prefix)
	}
	fmt.Fprintf(infoOutput, format, args...)
	fmt.Fprintln(infoOutput)
}

// Infoln outputs the given info message to standard error.
func Infoln(args ...any) {
	if skip(LevelInfo) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if infoUsePrefix {
		prefix := getPrefix(term.CyanBold)
		fmt.Fprint(infoOutput, prefix)
	}
	fmt.Fprintln(infoOutput, args...)
}

// --- [ warning ] -------------------------------------------------------------

var (
	// warnOutput specifies the output writer of non-fatal warning messages.
	warnOutput io.Writer = os.Stderr

	// warnUsePrefix specifies whether to use a prefix for warning messages.
	warnUsePrefix = true
)

// SetWarnOutput sets the output writer of non-fatal warning messages.
func SetWarnOutput(w io.Writer) {
	outputMutex.Lock()
	defer outputMutex.Unlock()
	warnOutput = w
}

// SetWarnPrefix sets whether to use a prefix for warning messages.
func SetWarnPrefix(usePrefix bool) {
	warnUsePrefix = usePrefix
}

// Warn outputs the given non-fatal warning message to standard error.
func Warn(args ...any) {
	if skip(LevelWarn) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if warnUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(warnOutput, prefix)
	}
	fmt.Fprint(warnOutput, args...)
	fmt.Fprintln(warnOutput)
}

// Warnf outputs the given non-fatal warning message to standard error.
func Warnf(format string, args ...any) {
	if skip(LevelWarn) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if warnUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(warnOutput, prefix)
	}
	fmt.Fprintf(warnOutput, format, args...)
	fmt.Fprintln(warnOutput)
}

// Warnln outputs the given non-fatal warning message to standard error.
func Warnln(args ...any) {
	if skip(LevelWarn) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if warnUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(warnOutput, prefix)
	}
	fmt.Fprintln(warnOutput, args...)
}

// --- [ error ] ---------------------------------------------------------------

var (
	// errorOutput specifies the output writer of fatal error messages.
	errorOutput io.Writer = os.Stderr

	// errorUsePrefix specifies whether to use a prefix for error messages.
	errorUsePrefix = true
)

// SetErrorOutput sets the output writer of fatal error messages.
func SetErrorOutput(w io.Writer) {
	outputMutex.Lock()
	defer outputMutex.Unlock()
	errorOutput = w
}

// SetErrorPrefix sets whether to use a prefix for error messages.
func SetErrorPrefix(usePrefix bool) {
	errorUsePrefix = usePrefix
}

// Fatal outputs the given fatal error message to standard error and terminates
// the application.
func Fatal(args ...any) {
	if skip(LevelError) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if errorUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(errorOutput, prefix)
	}
	fmt.Fprint(errorOutput, args...)
	fmt.Fprintln(errorOutput)
	os.Exit(1)
}

// Fatalf outputs the given fatal error message to standard error and terminates
// the application.
func Fatalf(format string, args ...any) {
	if skip(LevelError) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if errorUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(errorOutput, prefix)
	}
	fmt.Fprintf(errorOutput, format, args...)
	fmt.Fprintln(errorOutput)
	os.Exit(1)
}

// Fatalln outputs the given fatal error message to standard error and
// terminates the application.
func Fatalln(args ...any) {
	if skip(LevelError) {
		return
	}
	outputMutex.Lock()
	defer outputMutex.Unlock()
	if errorUsePrefix {
		prefix := getPrefix(term.RedBold)
		prefix += getFileLine()
		fmt.Fprint(errorOutput, prefix)
	}
	fmt.Fprintln(errorOutput, args...)
	os.Exit(1)
}

// ### [ Helper functions ] ####################################################

// getQualifiedPaths returns the qualified package and and qualified function
// paths of the caller.
func getQualifiedPaths() (pkgPath, funcPath string) {
	const skip = 3 // skip 3 call frames: {Debugf,Warnf}, skip and getQualifiedPaths.
	pathQualifiedName, _, _, ok := callerName(skip)
	if !ok {
		return "", ""
	}
	funcPath = pathQualifiedName
	pkgPath = getPkgPath(funcPath)
	return pkgPath, funcPath
}

// getPrefix returns the prefix used for logging based on the function name of
// the caller and the given terminal color.
func getPrefix(colorFunc func(string) string) string {
	const skip = 2 // skip 2 call frames: {Debugf,Warnf} and getPrefix.
	pathQualifiedName, _, _, ok := callerName(skip)
	if !ok {
		return ""
	}
	pkgName := getPkgName(pathQualifiedName)
	prefix := colorFunc(pkgName+":") + " "
	return prefix
}

// getFileLine returns the file name and line number of the caller.
func getFileLine() string {
	const skip = 2 // skip 2 call frames: {Debugf,Warnf} and getFileLine.
	_, file, line, ok := callerName(skip)
	if !ok {
		return ""
	}
	// TODO: use getFuncName?
	s := fmt.Sprintf("%s:%d", file, line)
	fileLine := term.WhiteBold(s+":") + " "
	return fileLine
}

// callerName returns the path-qualified function name of the caller.
func callerName(skip int) (pathQualifiedName string, fileName string, lineNum int, ok bool) {
	var pcs [1]uintptr
	n := runtime.Callers(skip+2, pcs[:]) // always skip the 2 deepest call frames: callerName and runtime.Callers
	if n != len(pcs) {
		// unable to get program counter of callers
		return "", "", 0, false
	}
	fn := runtime.FuncForPC(pcs[0])
	if fn == nil {
		// unable to get function with program counter pcs[0]
		return "", "", 0, false
	}
	pathQualifiedName = fn.Name()
	fileName, lineNum = fn.FileLine(pcs[0])
	return pathQualifiedName, fileName, lineNum, true
}

// getPkgPath returns the package path of the path-qualified function name.
//
// Example input:
//
//	github.com/mewpkg/clog.getPrefix
//	github.com/mewpkg/clog.Debugf
//	main.main
//
// Example output:
//
//	github.com/mewpkg/clog
//	github.com/mewpkg/clog
//	main
func getPkgPath(name string) string {
	// find last slash of package path.
	end := 0
	pos := strings.LastIndex(name, "/")
	if pos != -1 {
		end = pos + 1
	}
	// strip function name.
	pos = strings.Index(name[end:], ".")
	if pos != -1 {
		end += pos
	}
	return name[:end]
}

// getPkgName returns the package name of the path-qualified function name.
//
// Example input:
//
//	github.com/mewpkg/clog.getPrefix
//	github.com/mewpkg/clog.Debugf
//	main.main
//
// Example output:
//
//	clog
//	clog
//	main
func getPkgName(name string) string {
	// strip package path; keep package name and function name.
	pos := strings.LastIndex(name, "/")
	if pos != -1 {
		name = name[pos+1:]
	}
	// get package name.
	pos = strings.Index(name, ".")
	if pos != -1 {
		name = name[:pos]
	}
	if name == "main" && len(mainPrefixName) > 0 {
		return mainPrefixName
	}
	return name
}

// getFuncName returns the function name of the path-qualified function name.
//
// Example input:
//
//	github.com/mewpkg/clog.getPrefix
//	github.com/mewpkg/clog.Debugf
//	main.main
//
// Example output:
//
//	getPrefix
//	Debugf
//	main
func getFuncName(name string) string {
	// strip package path; keep package name and function name.
	pos := strings.LastIndex(name, "/")
	if pos != -1 {
		name = name[pos+1:]
	}
	// get function name.
	pos = strings.Index(name, ".")
	if pos != -1 {
		name = name[pos+1:]
	}
	return name
}
