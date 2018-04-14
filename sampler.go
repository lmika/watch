package main

import (
	"context"
	"os/exec"
	"sync"
	"time"
)

// Sampler is responsible for sampling the output of a command
type Sampler struct {
	// Command is the command to sample
	Command string

	// Interval is how often to sample the command
	Interval time.Duration

	// OnSampled is called when a new sample is available
	OnSampled func()

	lastSnapshot  *Snapshot
	snapshotMutex *sync.RWMutex
}

// Init initialises the sampler
func (s *Sampler) Init() {
	s.snapshotMutex = new(sync.RWMutex)
}

// LastSnapshot returns the last sampled snapshot.  If no snapshot has been
// collected yet, returns nil.
func (s *Sampler) LastSnapshot() *Snapshot {
	s.snapshotMutex.RLock()
	defer s.snapshotMutex.RUnlock()

	return s.lastSnapshot
}

// Start starts sampling commands.  This should be launched in a go routine.
func (s *Sampler) Start() {

	takeSnapshot := make(chan struct{})
	defer close(takeSnapshot)

	go func() {
		for range takeSnapshot {
			s.SampleNow()
		}
	}()

	// Take a snapshot now
	takeSnapshot <- struct{}{}

	// Once taken, start ticking
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for range ticker.C {
		// Attempt to send a request to take a snapshot
		// If it missfires, skip and wait for the next interval
		select {
		case takeSnapshot <- struct{}{}:
		default:
		}
	}

}

// SampleNow executes the command and saves the snapshot
func (s *Sampler) SampleNow() {
	snaphot := s.sample(s.Command)
	s.pushSnapshot(snaphot)

	if s.OnSampled != nil {
		s.OnSampled()
	}
}

func (s *Sampler) pushSnapshot(newSnapshot *Snapshot) {
	s.snapshotMutex.Lock()
	defer s.snapshotMutex.Unlock()

	// TODO: Allow some historical analysis of snapshots
	s.lastSnapshot = newSnapshot
}

func (s *Sampler) sample(command string) *Snapshot {
	snapshot := &Snapshot{}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	snapshot.Command = command
	snapshot.Started = time.Now()

	output, err := cmd.CombinedOutput()

	snapshot.Duration = time.Now().Sub(snapshot.Started)
	snapshot.Err = err
	snapshot.Lines = StringToLines(string(output))

	return snapshot
}
