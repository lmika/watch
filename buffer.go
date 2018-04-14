package main

import (
	"strings"
	"time"
)

// Snapshot represents a snapshot of the watched command
type Snapshot struct {
	// The command that is being watched
	Command string

	// When the command was executed
	Started time.Time

	// How long the command took to run (wall-clock)
	Duration time.Duration

	// The output lines of the snapshot
	Lines []Line

	// Any error that resulted from running the command
	Err error
}

// Line is an output line of a command
type Line struct {
	Line string
}

// StringToLines takes the output of the command and returns them
// as separated lines.
func StringToLines(s string) []Line {
	strLines := strings.Split(s, "\n")
	lines := make([]Line, len(strLines))
	for i, s := range strLines {
		lines[i] = Line{s}
	}
	return lines
}
